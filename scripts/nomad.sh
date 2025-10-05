SECRET_PATH="projects/tl/nomad/secret"

export NOMAD_ADDR="http://$(terraform -chdir=deployment/infra/staging output -raw ip):4646"

export NOMAD_TOKEN="$(pass $SECRET_PATH)"
