SECRET_PATH="projects/tl/nomad/secret"

export NOMAD_ADDR="http://$(terraform -chdir=deployment/infra/staging output -raw ip):4646"

if ! pass $SECRET_PATH >/dev/null 2>&1; then
    nomad acl bootstrap | head -n 2 | tail -n 1 | cut -d '=' -f 2 | xargs | pass insert -m $SECRET_PATH
fi

export NOMAD_TOKEN="$(pass $SECRET_PATH)"
