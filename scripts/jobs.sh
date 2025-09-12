nomad job run ./deployment/jobs/prometheus.nomad.hcl
nomad job run ./deployment/jobs/loki.nomad.hcl
nomad job run ./deployment/jobs/tempo.nomad.hcl
nomad job run -var="domain=grafana.short.com" ./deployment/jobs/grafana.nomad.hcl
nomad job run -var="image=tl-website:local" -var="domain=short.com" ./deployment/jobs/website.nomad.hcl
nomad job run ./deployment/jobs/redis.nomad.hcl
nomad job run ./deployment/jobs/postgres.nomad.hcl
nomad job run -var="image=tl-server:local" -var="short_domain=sho.rt" ./deployment/jobs/server.nomad.hcl
nomad job run ./deployment/jobs/traefik.nomad.hcl
