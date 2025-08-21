# yaml-language-server: $schema=https://www.schemastore.org/traefik-v2.json

entryPoints:
  http:
    address: :8080
  traefik:
    address: :1936
  web:
    address: :80

api:
  dashboard: true
  insecure: true

providers:
  nomad:
    exposedByDefault: false
