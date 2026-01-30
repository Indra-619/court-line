Write-Host "Starting Build Safety Verification..." -ForegroundColor Cyan

# 1. Backend Verification
Write-Host "Verifying Backend (Go)..." -ForegroundColor Yellow
Set-Location "$PSScriptRoot\..\backend"

# Run Unit Tests
Write-Host "   Running Unit Tests..." -NoNewline
$testOutput = go test ./... 2>&1
if ($LASTEXITCODE -ne 0) {
    Write-Host " FAILED" -ForegroundColor Red
    Write-Host $testOutput
    exit 1
}
Write-Host " PASSED" -ForegroundColor Green

# 2. Frontend Verification
Write-Host "Verifying Frontend (Nuxt/Vue)..." -ForegroundColor Yellow
Set-Location "$PSScriptRoot\..\frontend"

# Run Type Check
Write-Host "   Running Type Check..." -NoNewline
# We use cmd /c to verify typecheck to bypass potential npm.ps1 restrictions
$typeCheck = cmd /c "npm run typecheck" 2>&1
if ($LASTEXITCODE -ne 0) {
    Write-Host " FAILED" -ForegroundColor Red
    Write-Host $typeCheck
    exit 1
}
Write-Host " PASSED" -ForegroundColor Green

# 3. Infrastructure Verification
Write-Host "Verifying Infrastructure (Docker)..." -ForegroundColor Yellow
Set-Location "$PSScriptRoot\.."

Write-Host "   Checking Docker Compose Config..." -NoNewline
$composeCheck = docker-compose config -q 2>&1
if ($LASTEXITCODE -ne 0) {
    Write-Host " FAILED" -ForegroundColor Red
    Write-Host "   Error: Docker Compose configuration is invalid."
    Write-Host $composeCheck
    exit 1
}
Write-Host " PASSED" -ForegroundColor Green

Write-Host "All Systems Go! Build is strictly safe." -ForegroundColor Cyan
