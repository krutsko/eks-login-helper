# 🚀 EKS Login Helper

A powerful, user-friendly CLI tool written in Go that streamlines your AWS EKS authentication workflow. Say goodbye to manual SSO logins, cluster hunting, and kubeconfig management!

## ✨ Features

- 🔐 **Smart SSO Management**: Automatically detects and validates AWS SSO sessions
- 📋 **Interactive Cluster Selection**: Browse and select from available EKS clusters  
- ⚙️ **Automatic Kubeconfig Updates**: Seamlessly updates your kubectl configuration
- 🎨 **Beautiful CLI Interface**: Colorized output with clear status indicators
- 🚄 **Lightning Fast**: Built in Go for optimal performance
- 🔧 **Flexible Usage**: Support for both interactive and non-interactive modes
- 🌍 **Multi-Platform**: Works on Linux, macOS, and Windows

## 🎯 Your Workflow, Simplified

**Before** (your current manual process):
```bash
aws sso login --profile my-profile
aws eks list-clusters --profile my-profile
aws eks --region us-west-2 update-kubeconfig --name my-cluster --profile my-profile
kubectl config use-context my-cluster
```

**After** (with eks-login):
```bash
eks-login
```

That's it! The tool handles everything interactively.

## 📦 Installation

### Option 1: Build from Source
```bash
# Clone or download the source files
git clone <your-repo-url>  # or create the files manually
cd eks-login

# Build and install
make install
```

### Option 2: Manual Build
```bash
# Install dependencies
go mod tidy

# Build
go build -o eks-login .

# Install (optional)
sudo mv eks-login /usr/local/bin/
```

### Prerequisites
- Go 1.21 or later
- AWS CLI v2 configured
- kubectl installed
- Valid AWS SSO configuration

## 🚀 Usage

### Interactive Mode (Recommended)
```bash
eks-login
```
The tool will guide you through:
1. Selecting your AWS profile
2. Validating/refreshing SSO session
3. Choosing your EKS cluster
4. Updating kubeconfig automatically

### Non-Interactive Mode
```bash
# Specify all parameters
eks-login --profile my-profile --region us-west-2 --cluster my-cluster

# Skip SSO login if already authenticated
eks-login --profile my-profile --skip-sso
```

### Command Line Options
```
Flags:
  -c, --cluster string    EKS cluster name
  -h, --help             help for eks-login
      --interactive      Enable interactive mode (default true)
  -p, --profile string   AWS profile to use
  -r, --region string    AWS region (default "us-west-2")
      --skip-sso         Skip SSO login (assume already logged in)
```

## 📖 Examples

### Basic Interactive Usage
```bash
$ eks-login

🔍 Checking dependencies...
  ✓ aws found
  ✓ kubectl found

📋 Available AWS Profiles:
  1. dev-profile (region: us-west-2)
  2. prod-profile (region: us-east-1)
  3. staging-profile (region: eu-west-1)

Select profile (1-3): 1

✓ SSO session is valid

📋 Fetching EKS clusters...

🎯 Available EKS Clusters in us-west-2:
  1. dev-cluster
  2. staging-cluster
  3. prod-cluster

Select cluster (1-3): 2

⚙️  Updating kubeconfig for cluster: staging-cluster
Updated context arn:aws:eks:us-west-2:123456789:cluster/staging-cluster in /Users/you/.kube/config
✓ Kubeconfig updated successfully!

🔍 Verifying cluster connection...
✓ Successfully connected to cluster!
📍 Current context: arn:aws:eks:us-west-2:123456789:cluster/staging-cluster

🎉 EKS Login Complete!
Profile: dev-profile
Region: us-west-2
Cluster: staging-cluster

You can now use kubectl to interact with your cluster.
```

### Quick Non-Interactive Usage
```bash
# If you know exactly what you want
$ eks-login --profile prod-profile --region us-east-1 --cluster production

🔍 Checking dependencies...
  ✓ aws found
  ✓ kubectl found
🔐 Logging in to AWS SSO...
✓ SSO login successful
⚙️  Updating kubeconfig for cluster: production
✓ Kubeconfig updated successfully!
✓ Successfully connected to cluster!

🎉 EKS Login Complete!
```

## 🔧 Development

### Building
```bash
# Install dependencies
make deps

# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
make test

# Development workflow
make dev
```

### Project Structure
```
eks-login/
├── main.go           # Main application code
├── go.mod            # Go module definition
├── Makefile          # Build automation
├── README.md         # This file
└── build/            # Build artifacts (created during build)
```

## 🌟 Key Benefits

- **Time Saver**: Reduces 4-step manual process to 1 command
- **Error Reduction**: Eliminates typing mistakes in cluster names and regions
- **Session Management**: Intelligent SSO session handling
- **Multi-Cluster Friendly**: Easy switching between different clusters
- **Developer Experience**: Beautiful, intuitive CLI interface
- **Cross-Platform**: Single binary works everywhere

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🚨 Troubleshooting

### Common Issues

**"aws command not found"**
- Install AWS CLI v2: https://docs.aws.amazon.com/cli/latest/userguide/install-cliv2.html

**"kubectl command not found"**
- Install kubectl: https://kubernetes.io/docs/tasks/tools/install-kubectl/

**"No AWS profiles found"**
- Configure AWS CLI with: `aws configure sso`

**"SSO login failed"**
- Check your AWS SSO configuration
- Ensure you have internet connectivity
- Verify your SSO start URL is correct

**"No EKS clusters found"**
- Verify you're in the correct AWS region
- Check that your profile has EKS permissions
- Ensure clusters exist in the specified region

## 🙏 Acknowledgments

- AWS CLI team for the excellent `aws eks update-kubeconfig` command
- The Go community for amazing libraries
- [Cobra](https://github.com/spf13/cobra) for CLI framework
- [Color](https://github.com/fatih/color) for terminal colors