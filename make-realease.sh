#!/bin/bash

echo "Build CHANGELOG.md"
git-chglog --output CHANGELOG.md 

echo "Tagging"
git tag -a "$1" -m "WALLETS $1"

echo "Pushing release"
git push --tags origin main
