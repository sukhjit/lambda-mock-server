#!/bin/bash

set -euo pipefail

docker run --rm -it \
    -p 8000:8000 \
    -e AWS_ACCESS_KEY_ID=$(pass my-aws-id) \
    -e AWS_SECRET_ACCESS_KEY=$(pass my-aws-secret) \
    -v $(pwd):/app \
    -w /app golang:1.26-trixie \
    /bin/bash
