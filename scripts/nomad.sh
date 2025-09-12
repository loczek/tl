export NOMAD_ADDR="http://$(terraform -chdir=deployment/infra/staging output -raw ip):4646"
