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
sudo wget --quiet "${release_url}/${version}/fspop-${version}-linux_${arch}.zip" -O /var/tmp/fspop.zip
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
