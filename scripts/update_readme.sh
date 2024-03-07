#!/bin/bash

# Update the README.md file with the new version
# this script is called from github actions when 
# a new release is created to replace the version 
# regex in the README.md to the new version

version=$1
sed -i "s/v[0-9]\.[0-9]\.[0-9]/v$version/" README.md


