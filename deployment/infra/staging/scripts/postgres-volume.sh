#!/usr/bin/env bash
set -euo pipefail

sudo cat <<EOF | nomad volume register -
id          = "postgres-volume"
name        = "postgres-volume"
plugin_id   = "aws-ebs0"
external_id = $(terraform -chdir=deployment/infra/staging output db_volume_id)
type        = "csi"

capability {
  access_mode     = "single-node-writer"
  attachment_mode = "file-system"
}

mount_options {
  fs_type = "xfs"
}
EOF
