#!/bin/bash
# Author: XORbit01
# Update the README.md file with the new version
# this script is called from the deployer script
# a new release is created to replace the version 
# regex in the README.md to the new version

version=$1
if [ -z "$version" ]; then
  echo "version is required"
  exit 1
fi

sed -i "s/v[0-9]\+\.[0-9]\+\.[0-9]\+/$version/" README.md
