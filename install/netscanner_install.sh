#!/usr/bin/env bash
set -euo pipefail

repo="deluxesande/network-scanner"
install_path="$HOME/.local/bin"

# OS detection (matches the archive names: Linux, Darwin)
case "$(uname -s)" in
  Linux)  os="Linux" ;;
  Darwin) os="Darwin" ;;
  *) echo "Unsupported OS: $(uname -s)"; exit 1 ;;
esac

# Architecture detection (matches the archive names: x86_64, arm64)
case "$(uname -m)" in
  x86_64|amd64)  arch="x86_64" ;;
  aarch64|arm64) arch="arm64" ;;
  *) echo "Unsupported architecture: $(uname -m)"; exit 1 ;;
esac

# "latest/download" always points at the newest published release
url="https://github.com/$repo/releases/latest/download/netscanner_${os}_${arch}.tar.gz"

tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

echo "Downloading netscanner ($os/$arch)..."
# -f makes curl fail on a 404 instead of saving the error page
curl -fL "$url" -o "$tmp/netscanner.tar.gz"

# Extract to a temp dir so LICENSE/README/CHANGELOG don't land in ~/.local/bin
tar -xzf "$tmp/netscanner.tar.gz" -C "$tmp"
mkdir -p "$install_path"
install -m 755 "$tmp/netscanner" "$install_path/netscanner"

# Add to PATH using the user's actual shell
case "$(basename "${SHELL:-}")" in
  zsh)  rc="$HOME/.zshrc" ;;
  bash) rc="$HOME/.bashrc" ;;
  *)    rc="$HOME/.profile" ;;
esac

case ":$PATH:" in
  *":$install_path:"*) ;;  # already on PATH
  *)
    # The marker comment lets the uninstaller remove only this line
    echo "export PATH=\"$install_path:\$PATH\" # added by netscanner installer" >> "$rc"
    echo "Added $install_path to PATH in $rc"
    echo "Restart your terminal or run: source $rc"
    ;;
esac

echo "netscanner installed successfully! Run: netscanner -h"