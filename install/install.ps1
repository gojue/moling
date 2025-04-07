#!/usr/bin/env pwsh

Set-StrictMode -Version Latest

Write-Output "Welcome to MoLing MCP Server initialization script."
Write-Output "Home page: https://gojue.cc/moling"
Write-Output "Github: https://github.com/gojue/moling"

# Determine the OS and architecture
$OS = (Get-CimInstance Win32_OperatingSystem).Caption
$ARCH = (Get-CimInstance Win32_Processor).Architecture

switch ($ARCH) {
    9 { $ARCH = "amd64" }
    5 { $ARCH = "arm64" }
    default {
        Write-Error "Unsupported architecture: $ARCH"
        exit 1
    }
}

# Determine the download URL
$VERSION = "v0.0.1"
$BASE_URL = "https://github.com/gojue/moling/releases/download/$VERSION"
$FILE_NAME = "moling-$VERSION-windows-$ARCH.zip"
$DOWNLOAD_URL = "$BASE_URL/$FILE_NAME"

# Download the installation package
Write-Output "Downloading $DOWNLOAD_URL..."
Invoke-WebRequest -Uri $DOWNLOAD_URL -OutFile $FILE_NAME

# Extract the package
Write-Output "Extracting $FILE_NAME..."
Expand-Archive -Path $FILE_NAME -DestinationPath "moling"

# Move the binary to C:\Program Files
$destination = "C:\Program Files\moling"
if (-Not (Test-Path -Path $destination)) {
    New-Item -ItemType Directory -Path $destination
}
Move-Item -Path "moling\moling.exe" -Destination "$destination\moling.exe"

# Add to PATH
$env:Path += ";$destination"
[System.Environment]::SetEnvironmentVariable("Path", $env:Path, [System.EnvironmentVariableTarget]::Machine)

# MCP Client configuration
& "$destination\moling.exe" client --install

# Clean up
Remove-Item -Recurse -Force "moling"
Remove-Item -Force $FILE_NAME

Write-Output "MoLing has been installed successfully!"