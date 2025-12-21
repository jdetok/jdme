#!/usr/bin/bash
if git diff --quiet; then
    echo "No changes to commit"
else
    git config user.name "github-actions[bot]"
    git config user.email "github-actions[bot]@users.noreply.github.com"
    git add .
    git commit -m "Replace localhost URLs"
    git push
fi