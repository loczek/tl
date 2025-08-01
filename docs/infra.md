# Infrastructure as Code (IaC)

- Terraform
- OpenTofu (Terraform fork)
- AWS CloudFormation
- Ansible

# Deployment Options

- bare metal (physical server)
  - not a vm
  - only used by one consumer/tenant


# Definitions

- On-premises
  - On-premises owns everything like
    - software
    - hardware
    - maintanance
- Cloud


# Providers

| Platform      | Free Tier        | Cold Starts |
| ------------- | ---------------- | ----------- |
| Digital Ocean | No               | No          |
| Fly.io        | Generous         | No          |
| Render        | Yes Limited      | Yes         |
| Railway       | Yes ($5 credits) | Yes         |
| Hetzner VPS   | (~€4/month)      | No          |
| Oracle Cloud  | Yes (ARM only)   | No          |

# Providers

- AWS
  - Render
  - Vercel
- Digital Ocean
- Fly.io
- GCP
- Hetzner
- Oracle Cloud
- Railway (was on GCP, but now they have their own datacenters)

## AWS Services

- EC2
- Elastic Kubernetes Service (seems expensive)
- Elastic Container Service (containers on EC2 or Fargate)
- AWS RDS (managed postgresql)
- Aurora DSQL
- Route 53
- CloudFront
- CloudWatch

## Uptime

- Uptime Robot
- Uptime kuma
