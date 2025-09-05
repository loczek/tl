# yaml-language-server: $schema=https://www.schemastore.org/traefik-v2.json

entryPoints:
  web:
    address: :80
  websecure:
    address: :433

api:
  dashboard: true
  insecure: true

providers:
  nomad:
    exposedByDefault: false
