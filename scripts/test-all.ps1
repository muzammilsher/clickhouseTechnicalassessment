# PowerShell Automated Assessment Verification Script
# Runs all 4 modules using local tools (Go, Python) and cached Docker containers (Terraform, Hadolint)

$ErrorActionPreference = "Continue"

Write-Host "==================================================" -ForegroundColor Cyan
Write-Host "         Security Test Suite Runner               " -ForegroundColor Cyan
Write-Host "==================================================" -ForegroundColor Cyan

$dockerAvailable = $false
try {
    $dockerInfo = docker info 2>&1
    if ($LASTEXITCODE -eq 0) { $dockerAvailable = $true }
}
catch {}

# 1. Terraform GCP
Write-Host "`n[1/4] Running Module 01: Terraform GCP Checks..." -ForegroundColor Yellow
if ($dockerAvailable) {
    docker run --rm -v "${PWD}/01-terraform-gcp:/workspace" -w /workspace hashicorp/terraform:latest fmt -check -recursive
    if ($LASTEXITCODE -eq 0) { Write-Host "  [PASS] terraform fmt" -ForegroundColor Green } else { Write-Host "  [FAIL] terraform fmt" -ForegroundColor Red }

    docker run --rm -v "${PWD}/01-terraform-gcp:/workspace" -w /workspace hashicorp/terraform:latest init -backend=false | Out-Null
    docker run --rm -v "${PWD}/01-terraform-gcp:/workspace" -w /workspace hashicorp/terraform:latest validate
    if ($LASTEXITCODE -eq 0) { Write-Host "  [PASS] terraform validate" -ForegroundColor Green } else { Write-Host "  [FAIL] terraform validate" -ForegroundColor Red }
}
else {
    Write-Host "  [SKIP] Docker daemon not running. Start Docker Desktop to validate Terraform in container." -ForegroundColor DarkGray
}

# 2. Dockerfile Linting
Write-Host "`n[2/4] Running Module 02: Hadolint Container Security..." -ForegroundColor Yellow
if ($dockerAvailable) {
    docker run --rm -v "${PWD}:/work" -w /work hadolint/hadolint hadolint 02-cloud-security-check/Dockerfile.secure
    if ($LASTEXITCODE -eq 0) { Write-Host "  [PASS] hadolint Dockerfile.secure (0 warnings, 0 errors)" -ForegroundColor Green } else { Write-Host "  [FAIL] hadolint failed" -ForegroundColor Red }
}
else {
    Write-Host "  [SKIP] Docker daemon not running. Start Docker Desktop to run Hadolint in container." -ForegroundColor DarkGray
}

# 3. Go Firewall Scanner
Write-Host "`n[3/4] Running Module 03: Go Firewall Scanner..." -ForegroundColor Yellow
Push-Location 03-firewall-scanner
go test -v ./...
if ($LASTEXITCODE -eq 0) { Write-Host "  [PASS] go test ./..." -ForegroundColor Green } else { Write-Host "  [FAIL] go test" -ForegroundColor Red }
go vet ./...
if ($LASTEXITCODE -eq 0) { Write-Host "  [PASS] go vet ./..." -ForegroundColor Green } else { Write-Host "  [FAIL] go vet" -ForegroundColor Red }
go build -o scanner.exe .
if ($LASTEXITCODE -eq 0) { Write-Host "  [PASS] go build ." -ForegroundColor Green } else { Write-Host "  [FAIL] go build" -ForegroundColor Red }
Remove-Item -Force scanner.exe -ErrorAction SilentlyContinue
Pop-Location

# 4. Dependency Security Audit
Write-Host "`n[4/4] Running Module 04: Dependency Vulnerability Audit..." -ForegroundColor Yellow
Push-Location 04-dependency-audit
python audit.py npm-audit.json vulnerability-report.json
if ($LASTEXITCODE -eq 0) { Write-Host "  [PASS] python audit.py (JSON)" -ForegroundColor Green } else { Write-Host "  [FAIL] python audit.py (JSON)" -ForegroundColor Red }
python audit.py npm-audit.json vulnerability-report.csv
if ($LASTEXITCODE -eq 0) { Write-Host "  [PASS] python audit.py (CSV)" -ForegroundColor Green } else { Write-Host "  [FAIL] python audit.py (CSV)" -ForegroundColor Red }
Pop-Location

Write-Host "`n==================================================" -ForegroundColor Cyan
Write-Host "       All Module Checks Finished                 " -ForegroundColor Cyan
Write-Host "==================================================" -ForegroundColor Cyan
