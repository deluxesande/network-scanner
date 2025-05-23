# Download and install netscanner on Windows
$version = "v0.2.0"
$os = "windows"
$arch = "amd64"
$url = "https://github.com/deluxesande/network-scanner/releases/download/$version/netscanner_Windows_x86_64.zip"
$installPath = "$env:USERPROFILE\netscanner"

# Create installation directory silently
try {
    New-Item -ItemType Directory -Force -Path $installPath | Out-Null
} catch {
    Write-Host "Failed to create directory: $installPath"
    exit 1
}

# Download the binary
Write-Host "Downloading netscanner..."
Invoke-WebRequest -Uri $url -OutFile "$installPath\netscanner.zip"

# Extract the binary
Write-Host "Extracting netscanner..."
Expand-Archive -Path "$installPath\netscanner.zip" -DestinationPath $installPath -Force

# Rename the binary to 'netscanner.exe'
Rename-Item -Path "$installPath\network-scanner.exe" -NewName "netscanner.exe"

# Add to PATH
Write-Host "Adding netscanner to PATH..."
$env:Path += ";$installPath"
[Environment]::SetEnvironmentVariable("Path", $env:Path, [EnvironmentVariableTarget]::User)

# Clean up
Remove-Item "$installPath\netscanner.zip"

Write-Host "netscanner installed successfully! You can now run 'netscanner' from anywhere."