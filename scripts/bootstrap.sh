#!/usr/bin/env bash
set -euo pipefail

SECRET_PATH="projects/tl/nomad/secret"

SECRET=$(nomad acl bootstrap | head -n 2 | tail -n 1 | cut -d '=' -f 2 | xargs)

echo $SECRET

pass insert -e $SECRET_PATH $SECRET
