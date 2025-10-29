SECRET_PATH="projects/tl/prod/nomad/secret"

export NOMAD_ADDR="http://$(terraform -chdir=deployment/infra/production output -raw ip):4646"

export NOMAD_TOKEN="$(pass $SECRET_PATH)"
