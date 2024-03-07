#!/bin/bash
# author: XORbit01@protonmail
#
# this script is used localy to deploy the app 
# to the server 
# get the latest tag then increment it by 1 
# then push the new tag to the repo 
# get the latest tag name from remote
# then update the (
#   README.md,
# )

git fetch --tags
latest_tag=$(git describe --tags `git rev-list --tags --max-count=1`)
echo "latest tag is $latest_tag"

# then increment it
new_tag=$(echo "$latest_tag" | awk -F. -v OFS=. '{$NF++;print}')
echo "new tag is $new_tag"

# update the README 
echo "updating README.md with the new version $new_tag"
./scripts/update_readme.sh $new_tag
echo "README.md updated with the new version $new_tag"
git add README.md 
git commit -m "bump version to $new_tag" -a
git push origin main

# create a tag
echo "creating new tag"
git tag "$new_tag"
echo "push tag to origin"
git push origin "$new_tag"


