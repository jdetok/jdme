#!/bin/bash
# replace all instances of local host url with production

set -e

find . -type f \( -name "*.js" -o -name "*.html" -o -name "*.go" \) \
    -exec sed -i -E 's|http://localhost:[0-9]+|https://jdeko.me|g' {} +