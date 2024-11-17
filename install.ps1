$INFO = "::"
$WARN = "::"

Write-Host "$INFO Getting ready"

$REPO = "Modern-Modpacks/kjspkg"
$LATEST_RELEASE_URL = (Invoke-RestMethod -Uri "https://api.github.com/repos/$REPO/releases/latest").assets | Where-Object { $_.browser_download_url -like "*kjspkg_windows_amd64.exe" } | Select-Object -ExpandProperty browser_download_url

if (-not $LATEST_RELEASE_URL) {
    Write-Host "$WARN Failed to fetch the latest release. Please check your internet connection!"
    exit 1
}

Write-Host "$INFO Downloading KJSPKG"
$installPath = "$env:LOCALAPPDATA\Microsoft\WindowsApps"

try {
    Invoke-WebRequest -Uri $LATEST_RELEASE_URL -OutFile "$installPath\kjspkg.exe"
} catch {
    Write-Host "$WARN Download failed. Please try again."
    exit 1
}

Write-Host "$INFO Done! Run 'kjspkg' to get started"
