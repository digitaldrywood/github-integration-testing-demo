# Setup Instructions for GitHub Integration Testing Demo

## 🚀 Quick Start

### 1. Initialize Git Repository

```bash
cd ~/projects/digitaldrywood/github-integration-testing-demo
git init
git add .
git commit -m "Initial commit - GitHub Actions integration testing demo"
```

### 2. Create GitHub Repository

1. Go to https://github.com/new
2. Repository name: `github-integration-testing-demo`
3. Description: "Demo repository showing GitHub Actions integration testing with manual triggers"
4. Choose: Public or Private
5. **DO NOT** initialize with README, .gitignore, or license (we already have them)
6. Click "Create repository"

### 3. Push to GitHub

```bash
# Add your remote (replace with your actual repository URL)
git remote add origin https://github.com/digitaldrywood/github-integration-testing-demo.git

# Push to GitHub
git branch -M main
git push -u origin main
```

### 4. Configure GitHub Environments

#### Create Integration Testing Environment:

1. Go to your repo on GitHub
2. Click **Settings** → **Environments**
3. Click **New environment**
4. Name: `integration-testing`
5. Click **Configure environment**
6. Add Protection Rules (optional):
   - ✅ Required reviewers: Add your GitHub username
   - 🕐 Wait timer: 0-5 minutes (optional)
   - Click **Save protection rules**

#### Create Production Integration Environment (optional):

1. Click **New environment**
2. Name: `production-integration`
3. Configure with stricter rules:
   - ✅ Required reviewers: Add multiple maintainers
   - 🕐 Wait timer: 5-10 minutes
   - 🌿 Deployment branches: Only from `main`

### 5. Add Repository Secrets

Go to **Settings** → **Secrets and variables** → **Actions**

#### Add these secrets:

```yaml
# S3 Testing (use your AWS credentials or create IAM user)
AWS_ACCESS_KEY_ID: your-access-key
AWS_SECRET_ACCESS_KEY: your-secret-key
AWS_REGION: us-east-1
TEST_S3_BUCKET: your-test-bucket-name

# API Testing (optional, can use dummy values)
API_KEY: your-api-key
API_ENDPOINT: https://api.example.com

# Database is using GitHub Actions service container
# No secrets needed for basic testing
```

### 6. Test Your Setup

#### Test CI Workflow:
1. Make a small change to README.md
2. Commit and push
3. Check Actions tab - CI workflow should run automatically

#### Test Manual Integration Workflow:
1. Go to **Actions** tab
2. Click **Manual Integration Tests**
3. Click **Run workflow**
4. Select options:
   - Leave PR number empty
   - Check "Run S3 integration tests"
   - Click **Run workflow**
5. If you configured approval, approve when prompted

### 7. Local Testing

```bash
# Install dependencies
make deps

# Run unit tests
make test

# Setup local environment (requires Docker)
make setup-local
source .env.test

# Run integration tests locally
make test-integration

# Cleanup
make teardown-local
```

## 📝 Checklist

- [ ] Repository created on GitHub
- [ ] Code pushed to repository
- [ ] GitHub Actions enabled (should be by default)
- [ ] `integration-testing` environment created
- [ ] Protection rules configured (optional)
- [ ] Secrets added
- [ ] CI workflow runs on push
- [ ] Manual workflow appears in Actions tab
- [ ] Manual workflow can be triggered
- [ ] Approval flow works (if configured)

## 🎯 What You Can Demo

1. **Automatic CI** - Push code, see tests run
2. **Manual Triggers** - Run integration tests on demand
3. **PR Testing** - Test specific PRs by number
4. **Environment Protection** - Require approval for sensitive tests
5. **Selective Testing** - Choose which integration tests to run
6. **Secret Management** - Secure credential handling
7. **Local Testing** - Run same tests locally with Docker

## 🆘 Troubleshooting

### Workflow not appearing in Actions tab
- Make sure workflow files are in `.github/workflows/`
- Check YAML syntax (use yamllint.com)
- Ensure Actions are enabled in repository settings

### Environment protection not working
- Environment name must match exactly (case-sensitive)
- User must have write access to repository
- Required reviewers must have repository access

### Secrets not working
- Secret names are case-sensitive
- Check they're added to correct repository
- Environment secrets override repository secrets

### Tests failing
- Check AWS credentials are valid
- Ensure S3 bucket exists and is accessible
- For local testing, ensure Docker is running

## 📚 Next Steps

1. **Customize tests** - Add your own integration tests
2. **Add badges** - Show test status in README
3. **Branch protection** - Require tests to pass before merge
4. **Notifications** - Add Slack/Discord notifications
5. **Matrix testing** - Test multiple Go versions
6. **Caching** - Speed up workflows with dependency caching

## 🔗 Resources

- [GitHub Actions Docs](https://docs.github.com/actions)
- [Environment Protection](https://docs.github.com/actions/deployment/environments)
- [workflow_dispatch](https://docs.github.com/actions/using-workflows/events-that-trigger-workflows#workflow_dispatch)
- [Using Secrets](https://docs.github.com/actions/security-guides/encrypted-secrets)