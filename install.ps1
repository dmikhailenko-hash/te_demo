# te_demo installer for Windows
# Usage: irm https://raw.githubusercontent.com/dmikhailenko-hash/te_demo/main/install.ps1 | iex

$ErrorActionPreference = "Stop"
$Repo    = "dmikhailenko-hash/te_demo"
$Binary  = "te_demo.exe"
$InstDir = "$env:USERPROFILE\.te_demo\bin"
$Tag     = "1.0.2"

Write-Host ""
Write-Host "  Installing te_demo $Tag ..." -ForegroundColor Cyan

# Direct download URL (no API needed)
$url = "https://github.com/$Repo/releases/download/$Tag/te_demo-windows-amd64.exe"

# Create install dir
New-Item -ItemType Directory -Force -Path $InstDir | Out-Null

# Download
$dest = Join-Path $InstDir $Binary
Write-Host "  Downloading from $url ..."
Invoke-WebRequest -Uri $url -OutFile $dest -UseBasicParsing

# Add to PATH (user scope)
$userPath = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($userPath -notlike "*$InstDir*") {
    [Environment]::SetEnvironmentVariable("PATH", "$userPath;$InstDir", "User")
    Write-Host "  Added $InstDir to PATH" -ForegroundColor Green
    Write-Host "  Restart terminal for PATH to take effect." -ForegroundColor Yellow
}

Write-Host ""
Write-Host "  Done! Run: te_demo help" -ForegroundColor Green
Write-Host ""
