package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// Colors
var (
	green  = color.New(color.FgGreen, color.Bold)
	red    = color.New(color.FgRed, color.Bold)
	yellow = color.New(color.FgYellow, color.Bold)
	blue   = color.New(color.FgBlue, color.Bold)
	cyan   = color.New(color.FgCyan, color.Bold)
)

// Config holds the application configuration
type Config struct {
	Profile       string
	Region        string
	Cluster       string
	Interactive   bool
	SkipSSO       bool
	DefaultRegion string
}

// EKSCluster represents an EKS cluster
type EKSCluster struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Region string `json:"region"`
}

// ListClustersResponse represents the response from eks list-clusters
type ListClustersResponse struct {
	Clusters []string `json:"clusters"`
}

// ProfileInfo holds AWS profile information
type ProfileInfo struct {
	Name   string
	Region string
}

// EKSLoginApp represents the main application
type EKSLoginApp struct {
	config *Config
}

// NewEKSLoginApp creates a new instance of the application
func NewEKSLoginApp() *EKSLoginApp {
	return &EKSLoginApp{
		config: &Config{
			DefaultRegion: "us-west-2",
			Interactive:   true,
		},
	}
}

// Execute runs a command and returns the output
func (app *EKSLoginApp) Execute(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("command failed: %s\nstderr: %s", err, exitError.Stderr)
		}
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// CheckDependencies verifies that required tools are installed
func (app *EKSLoginApp) CheckDependencies() error {
	dependencies := []string{"aws", "kubectl"}

	blue.Println("🔍 Checking dependencies...")

	for _, dep := range dependencies {
		if _, err := exec.LookPath(dep); err != nil {
			return fmt.Errorf("required dependency '%s' not found in PATH", dep)
		}
		green.Printf("  ✓ %s found\n", dep)
	}

	return nil
}

// GetAWSProfiles retrieves available AWS profiles
func (app *EKSLoginApp) GetAWSProfiles() ([]ProfileInfo, error) {
	output, err := app.Execute("aws", "configure", "list-profiles")
	if err != nil {
		return nil, fmt.Errorf("failed to list AWS profiles: %w", err)
	}

	lines := strings.Split(output, "\n")
	profiles := make([]ProfileInfo, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			// Try to get region for this profile
			region, _ := app.Execute("aws", "configure", "get", "region", "--profile", line)
			if region == "" {
				region = app.config.DefaultRegion
			}

			profiles = append(profiles, ProfileInfo{
				Name:   line,
				Region: region,
			})
		}
	}

	return profiles, nil
}

// SelectProfile allows interactive profile selection
func (app *EKSLoginApp) SelectProfile() error {
	profiles, err := app.GetAWSProfiles()
	if err != nil {
		return err
	}

	if len(profiles) == 0 {
		return fmt.Errorf("no AWS profiles found. Please configure AWS CLI first")
	}

	// If only one profile, use it
	if len(profiles) == 1 {
		app.config.Profile = profiles[0].Name
		app.config.Region = profiles[0].Region
		cyan.Printf("📋 Using profile: %s (region: %s)\n", app.config.Profile, app.config.Region)
		return nil
	}

	// Interactive selection
	blue.Println("\n📋 Available AWS Profiles:")
	for i, profile := range profiles {
		fmt.Printf("  %d. %s (region: %s)\n", i+1, profile.Name, profile.Region)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		yellow.Printf("\nSelect profile (1-%d): ", len(profiles))
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		choice, err := strconv.Atoi(strings.TrimSpace(input))
		if err != nil || choice < 1 || choice > len(profiles) {
			red.Printf("Invalid selection. Please choose a number between 1 and %d.\n", len(profiles))
			continue
		}

		selectedProfile := profiles[choice-1]
		app.config.Profile = selectedProfile.Name
		app.config.Region = selectedProfile.Region
		break
	}

	return nil
}

// CheckSSOSession verifies if the SSO session is valid
func (app *EKSLoginApp) CheckSSOSession() (bool, error) {
	_, err := app.Execute("aws", "sts", "get-caller-identity", "--profile", app.config.Profile)
	return err == nil, nil
}

// LoginSSO performs AWS SSO login
func (app *EKSLoginApp) LoginSSO() error {
	if app.config.SkipSSO {
		return nil
	}

	blue.Println("🔐 Logging in to AWS SSO...")

	cmd := exec.Command("aws", "sso", "login", "--profile", app.config.Profile)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("SSO login failed: %w", err)
	}

	green.Println("✓ SSO login successful")
	return nil
}

// ListEKSClusters retrieves available EKS clusters
func (app *EKSLoginApp) ListEKSClusters() ([]string, error) {
	blue.Println("📋 Fetching EKS clusters...")

	output, err := app.Execute("aws", "eks", "list-clusters",
		"--profile", app.config.Profile,
		"--region", app.config.Region,
		"--output", "json")

	if err != nil {
		return nil, fmt.Errorf("failed to list EKS clusters: %w", err)
	}

	var response ListClustersResponse
	if err := json.Unmarshal([]byte(output), &response); err != nil {
		return nil, fmt.Errorf("failed to parse cluster list: %w", err)
	}

	return response.Clusters, nil
}

// SelectCluster allows interactive cluster selection
func (app *EKSLoginApp) SelectCluster() error {
	clusters, err := app.ListEKSClusters()
	if err != nil {
		return err
	}

	if len(clusters) == 0 {
		return fmt.Errorf("no EKS clusters found in region %s with profile %s", app.config.Region, app.config.Profile)
	}

	// If only one cluster, use it
	if len(clusters) == 1 {
		app.config.Cluster = clusters[0]
		cyan.Printf("🎯 Using cluster: %s\n", app.config.Cluster)
		return nil
	}

	// Interactive selection
	blue.Printf("\n🎯 Available EKS Clusters in %s:\n", app.config.Region)
	for i, cluster := range clusters {
		fmt.Printf("  %d. %s\n", i+1, cluster)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		yellow.Printf("\nSelect cluster (1-%d): ", len(clusters))
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		choice, err := strconv.Atoi(strings.TrimSpace(input))
		if err != nil || choice < 1 || choice > len(clusters) {
			red.Printf("Invalid selection. Please choose a number between 1 and %d.\n", len(clusters))
			continue
		}

		app.config.Cluster = clusters[choice-1]
		break
	}

	return nil
}

// UpdateKubeconfig updates the kubeconfig file
func (app *EKSLoginApp) UpdateKubeconfig() error {
	blue.Printf("⚙️  Updating kubeconfig for cluster: %s\n", app.config.Cluster)

	args := []string{
		"eks", "update-kubeconfig",
		"--region", app.config.Region,
		"--name", app.config.Cluster,
		"--profile", app.config.Profile,
	}

	cmd := exec.Command("aws", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to update kubeconfig: %w", err)
	}

	green.Println("✓ Kubeconfig updated successfully!")
	return nil
}

// VerifyConnection verifies the connection to the cluster
func (app *EKSLoginApp) VerifyConnection() error {
	blue.Println("🔍 Verifying cluster connection...")

	// Check if kubectl can connect
	output, err := app.Execute("kubectl", "cluster-info")
	if err != nil {
		yellow.Println("⚠️  Kubeconfig updated but unable to verify connection")
		return nil
	}

	green.Println("✓ Successfully connected to cluster!")

	// Show current context
	if context, err := app.Execute("kubectl", "config", "current-context"); err == nil {
		cyan.Printf("📍 Current context: %s\n", context)
	}

	// Optionally show cluster info
	fmt.Println("\n" + strings.TrimSpace(output))

	return nil
}

// ShowSummary displays a summary of the operation
func (app *EKSLoginApp) ShowSummary() {
	green.Println("\n🎉 EKS Login Complete!")
	fmt.Printf("Profile: %s\n", app.config.Profile)
	fmt.Printf("Region: %s\n", app.config.Region)
	fmt.Printf("Cluster: %s\n", app.config.Cluster)
	fmt.Println("\nYou can now use kubectl to interact with your cluster.")
}

// Run executes the main application logic
func (app *EKSLoginApp) Run() error {
	// Check dependencies
	if err := app.CheckDependencies(); err != nil {
		return err
	}

	// Select profile if not provided
	if app.config.Profile == "" {
		if err := app.SelectProfile(); err != nil {
			return err
		}
	}

	// Check SSO session
	if sessionValid, err := app.CheckSSOSession(); err != nil {
		return fmt.Errorf("failed to check SSO session: %w", err)
	} else if sessionValid {
		green.Println("✓ SSO session is valid")
	} else {
		if err := app.LoginSSO(); err != nil {
			return err
		}
	}

	// Select cluster if not provided
	if app.config.Cluster == "" {
		if err := app.SelectCluster(); err != nil {
			return err
		}
	}

	// Update kubeconfig
	if err := app.UpdateKubeconfig(); err != nil {
		return err
	}

	// Verify connection
	if err := app.VerifyConnection(); err != nil {
		return err
	}

	// Show summary
	app.ShowSummary()

	return nil
}

func main() {
	app := NewEKSLoginApp()

	var rootCmd = &cobra.Command{
		Use:   "eks-login",
		Short: "🚀 EKS Login Helper - Streamline your AWS EKS authentication",
		Long: `EKS Login Helper automates the process of logging into AWS SSO,
listing available EKS clusters, and updating your kubeconfig.

Examples:
  eks-login                           # Interactive mode
  eks-login --profile my-profile      # Use specific profile
  eks-login --profile my-profile --region us-east-1 --cluster my-cluster`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.Run()
		},
	}

	// Flags
	rootCmd.Flags().StringVarP(&app.config.Profile, "profile", "p", "", "AWS profile to use")
	rootCmd.Flags().StringVarP(&app.config.Region, "region", "r", app.config.DefaultRegion, "AWS region")
	rootCmd.Flags().StringVarP(&app.config.Cluster, "cluster", "c", "", "EKS cluster name")
	rootCmd.Flags().BoolVar(&app.config.SkipSSO, "skip-sso", false, "Skip SSO login (assume already logged in)")
	rootCmd.Flags().BoolVar(&app.config.Interactive, "interactive", true, "Enable interactive mode")

	// Version command
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("EKS Login Helper v1.0.0")
		},
	}

	rootCmd.AddCommand(versionCmd)

	// Execute
	if err := rootCmd.Execute(); err != nil {
		red.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
