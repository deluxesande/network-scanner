$ErrorActionPreference = "Stop"

$repo = "deluxesande/network-scanner"
$installPath = "$env:USERPROFILE\netscanner"

# Architecture detection. PROCESSOR_ARCHITEW6432 is set when a 32-bit
# PowerShell runs on 64-bit Windows, so check it first.
$rawArch = if ($env:PROCESSOR_ARCHITEW6432) { $env:PROCESSOR_ARCHITEW6432 } else { $env:PROCESSOR_ARCHITECTURE }
switch ($rawArch) {
    "AMD64" { $arch = "x86_64" }
    "ARM64" { $arch = "arm64" }
    default { Write-Host "Unsupported architecture: $rawArch"; exit 1 }
}

$url = "https://github.com/$repo/releases/latest/download/netscanner_Windows_$arch.zip"
$tmp = Join-Path $env:TEMP "netscanner_install"

if (Test-Path $tmp) { Remove-Item -Recurse -Force $tmp }
New-Item -ItemType Directory -Force -Path $tmp, $installPath | Out-Null

Write-Host "Downloading netscanner (windows/$arch)..."
Invoke-WebRequest -Uri $url -OutFile "$tmp\netscanner.zip"
Expand-Archive -Path "$tmp\netscanner.zip" -DestinationPath $tmp -Force

# Copy -Force overwrites an existing install, so reinstalling and upgrading work.
# (Interim, before the new release: use "$tmp\network-scanner.exe" here.)
Copy-Item "$tmp\netscanner.exe" "$installPath\netscanner.exe" -Force
Remove-Item -Recurse -Force $tmp

# PATH: read ONLY the user-scope PATH, so the system PATH is not copied into it
$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
if (($userPath -split ";") -notcontains $installPath) {
    $newPath = if ([string]::IsNullOrEmpty($userPath)) { $installPath } else { "$userPath;$installPath" }
    [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
}
# Current session only, so it works right away
if (($env:Path -split ";") -notcontains $installPath) { $env:Path += ";$installPath" }

Write-Host "netscanner installed successfully! Run: netscanner -h"