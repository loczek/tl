ssh "ubuntu@$(terraform -chdir=deployment/infra/staging output -raw ip)"
