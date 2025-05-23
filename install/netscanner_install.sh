#!/bin/bash

# Download and install netscanner on Linux
version="v0.2.0"
os="linux"
arch="amd64"
url="https://github.com/deluxesande/network-scanner/releases/download/$version/netscanner_Linux_x86_64.tar.gz"
install_path="$HOME/.local/bin"

# Create installation directory
mkdir -p "$install_path"

# Download the binary
echo "Downloading netscanner..."
curl -L "$url" -o "/tmp/netscanner.tar.gz"

# Extract the binary
echo "Extracting netscanner..."
tar -xzf "/tmp/netscanner.tar.gz" -C "$install_path"

# Rename the binary from network-scanner to netscanner
mv "$install_path/network-scanner" "$install_path/netscanner"

# Add to PATH
echo "Adding netscanner to PATH..."
if ! grep -q "$install_path" <<< "$PATH"; then
  echo "export PATH=\"$install_path:\$PATH\"" >> "$HOME/.bashrc"
  source "$HOME/.bashrc"
fi

# Clean up
rm "/tmp/netscanner.tar.gz"

echo "netscanner installed successfully! You can now run 'netscanner' from anywhere."