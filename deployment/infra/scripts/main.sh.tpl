#!/usr/bin/env bash
set -eo pipefail

${INSTALL_DOCKER_SCRIPT}

${INSTALL_NOMAD_SCRIPT}

${SETUP_NOMAD_SCRIPT}

sudo systemctl enable nomad.service
sudo systemctl start nomad.service
