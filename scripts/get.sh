#!/bin/bash


#
# This script downloads, extracts and then moves the
# 'fspop' binary into the bin directory. Requires 'sudo' access
#


# fspop latest version
version="v1.1.0"
release_url="--URL TO RELEASES DIRECTORY--"


# Detect the platform
OS="$(uname)"
case $OS in
  Linux)
    OS='linux'
    ;;
  FreeBSD)
    OS='freebsd'
    echo 'OS not supported'
    exit 2
    ;;
  NetBSD)
    OS='netbsd'
	echo 'OS not supported'
    exit 2
    ;;
  OpenBSD)
    OS='openbsd'
	echo 'OS not supported'
    exit 2
    ;;  
  Darwin)
    OS='darwin'
    ;;
  SunOS)
    OS='solaris'
    echo 'OS not supported'
    exit 2
    ;;
  *)
    echo 'OS not supported'
    exit 2
    ;;
esac


# Get device architecture
# -> amd64 : arm64
arch="$(uname -m)"
case "$arch" in
  x86_64|amd64)
    arch='amd64'
    ;;
  i?86|x86)
    arch='386'
    echo 'OS type not supported'
    exit 2
    ;;
  aarch64|arm64)
    arch='arm64'
    ;;
  arm*)
    arch='arm'
    echo 'OS type not supported'
    exit 2
    ;;
  *)
    echo 'OS type not supported'
    exit 2
    ;;
esac


# Print message when a process fails
onfail () {
	if [ "${?}" != "0" ]; then
		echo
		echo "${1}"
		exit 1
	fi
}


# Fetch install dependencies
sudo apt -qq install wget unzip


# Download latest fspop
echo
echo "Downloading fspop from ${release_url}"
sudo wget --quiet "${release_url}/${version}/fspop-${version}-${OS}_${arch}.zip" -O /var/tmp/fspop.zip
echo
onfail "ERROR: Unable to fetch fspop"


# Extract fspop release zip
sudo unzip -o /var/tmp/fspop.zip -d /var/tmp
onfail "ERROR: Unable to extract fspop. Make sure to have 'unzip' installed."

sudo rm /var/tmp/fspop.zip


# Move to user bin
sudo mv /var/tmp/fspop /usr/bin/fspop
onfail "ERROR: Unable to move fspop into /usr/bin directory"


# Update permissions, make executable
sudo chmod +x /usr/bin/fspop
onfail "ERROR: Unable to make fspop executable. chmod +x /usr/bin/fspop"

echo
echo "Success. fspop is good to go!"
