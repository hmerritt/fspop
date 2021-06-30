#!/bin/bash


#
# This script downloads, extracts and then moves the
# 'fspop' binary into the bin directory. Requires 'sudo' access
#


# fspop latest version
version="v1.1.0"
release_url="--URL TO RELEASES DIRECTORY--"


# Get device architecture
# -> amd64 : arm64
arch=$(dpkg --print-architecture)


# Download latest fspop
sudo wget "${release_url}/${version}/fspop-${version}-linux_${arch}.zip" -O /var/tmp/fspop.zip


# Extract fspop release zip
sudo unzip /var/tmp/fspop.zip
sudo rm /var/tmp/fspop.zip


# Move to user bin
sudo mv /usr/bin/fspop


# Update permissions, make executable
sudo chmod +x /usr/bin/fspop
