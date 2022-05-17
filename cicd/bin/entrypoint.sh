#!/bin/bash

curl -sL https://deb.nodesource.com/setup_16.x | bash

if ! which node > /dev/null 2>&1; then
    apt-get update && apt-get install -y nodejs
fi

echo "Container running..."
tail -f /dev/null
