# te_demo installer for Windows
# Usage: irm https://raw.githubusercontent.com/YOUR_ORG/te_demo/main/install.ps1 | iex

$ErrorActionPreference = "Stop"
$Repo    = "YOUR_ORG/te_demo"
$Binary  = "te_demo.exe"
$InstDir = "$env:USERPROFILE\.te_demo\bin"

Write-Host ""
Write-Host "  Installing te_demo..." -ForegroundColor Cyan

# Get latest release
$release  = Invoke-RestMethod "https://api.github.com/repos/$Repo/releases/latest"
$version  = $release.tag_name
$asset    = $release.assets | Where-Object { $_.name -eq "te_demo-windows-amd64.exe" }
$url      = $asset.browser_download_url

Write-Host "  Version : $version" -ForegroundColor Green
Write-Host "  From    : $url"

# Create install dir
New-Item -ItemType Directory -Force -Path $InstDir | Out-Null

# Download
$dest = Join-Path $InstDir $Binary
Write-Host "  Saving to $dest ..."
Invoke-WebRequest -Uri $url -OutFile $dest -UseBasicParsing

# Add to PATH (user scope)
$userPath = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($userPath -notlike "*$InstDir*") {
    [Environment]::SetEnvironmentVariable("PATH", "$userPath;$InstDir", "User")
    Write-Host "  Added $InstDir to PATH" -ForegroundColor Green
    Write-Host "  Restart your terminal for PATH to take effect." -ForegroundColor Yellow
}

Write-Host ""
Write-Host "  Done! Run: te_demo help" -ForegroundColor Green
Write-Host ""
