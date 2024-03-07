#!/bin/bash

# Update the README.md file with the new version
# this script is called from github actions when 
# a new release is created to replace the version 
# regex in the README.md to the new version

version=$1
if [ -z "$version" ]; then
  echo "version is required"
  exit 1
fi

git config user.email "41898282+github-actions[bot]@users.noreply.github.com"
git config user.name "github-actions[bot]"

sed -i "s/v[0-9]\.[0-9]\.[0-9]/v$version/" README.md
git add README.md
git commit -m "bump version to $version"
git push origin main 

