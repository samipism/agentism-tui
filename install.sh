#!/bin/sh
# Installs the latest agentism-tui release binary for this machine.
#
#   curl -fsSL https://raw.githubusercontent.com/samipism/agentism-tui/main/install.sh | sh
#
# Override install location with INSTALL_DIR (default /usr/local/bin).
set -eu

REPO="samipism/agentism-tui"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# asset_name prints the goreleaser archive name for the given uname -s/-m
# pair, or exits 1 if this platform has no build (see .goreleaser.yaml).
asset_name() {
	os="$1" arch="$2"
	case "$os" in
	Darwin) os=darwin ;;
	Linux) os=linux ;;
	*)
		echo "install.sh: unsupported OS '$os' (only darwin and linux have releases)" >&2
		exit 1
		;;
	esac
	case "$arch" in
	x86_64 | amd64) arch=amd64 ;;
	arm64 | aarch64) arch=arm64 ;;
	*)
		echo "install.sh: unsupported arch '$arch' (only amd64 and arm64 have releases)" >&2
		exit 1
		;;
	esac
	echo "agentism-tui_${os}_${arch}.tar.gz"
}

# --self-test checks asset_name's platform mapping without touching the
# network. Run after editing this script or .goreleaser.yaml's build matrix.
if [ "${1:-}" = "--self-test" ]; then
	got=$(asset_name Darwin arm64)
	[ "$got" = "agentism-tui_darwin_arm64.tar.gz" ] || { echo "FAIL: $got"; exit 1; }
	got=$(asset_name Linux x86_64)
	[ "$got" = "agentism-tui_linux_amd64.tar.gz" ] || { echo "FAIL: $got"; exit 1; }
	echo "ok"
	exit 0
fi

asset=$(asset_name "$(uname -s)" "$(uname -m)")
base_url="https://github.com/${REPO}/releases/latest/download"

tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

echo "install.sh: downloading $asset ..."
curl -fsSL "$base_url/$asset" -o "$tmp/$asset"
curl -fsSL "$base_url/checksums.txt" -o "$tmp/checksums.txt"

( cd "$tmp" && grep " $asset\$" checksums.txt | shasum -a 256 -c - ) ||
	{ echo "install.sh: checksum verification failed" >&2; exit 1; }

tar -xzf "$tmp/$asset" -C "$tmp" agentism-tui

if [ -w "$INSTALL_DIR" ]; then
	mv "$tmp/agentism-tui" "$INSTALL_DIR/agentism-tui"
else
	sudo mv "$tmp/agentism-tui" "$INSTALL_DIR/agentism-tui"
fi
chmod +x "$INSTALL_DIR/agentism-tui"

echo "install.sh: installed to $INSTALL_DIR/agentism-tui"
