#!/bin/bash

if [ -z "$1" ]; then
    echo "No tag supplied"
    exit
fi

echo "Build CHANGELOG.md"
git-chglog --output CHANGELOG.md 

git add CHANGELOG.md
git commit -m "chore: updated changelog"

echo "Tagging"
git tag -a "$1" -m "WALLETS $1"

echo "Pushing release"
read -r -p "Are you sure? [y/N] " response
case "$response" in
    [yY][eE][sS]|[yY]) 
        git push --tags origin main
        ;;
esac
