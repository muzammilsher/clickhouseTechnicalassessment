# Configuration variables for GCP infrastructure

# Target GCP Project where infrastructure will be deployed
variable "project_id" {
  description = "GCP project ID"
  type        = string
}

# GCP Deployment Region (Default: London / europe-west2)
variable "region" {
  description = "GCP region"
  type        = string
  default     = "europe-west2"
}

# GCP Deployment Zone for Compute Engine instance
variable "zone" {
  description = "GCP zone"
  type        = string
  default     = "europe-west2-a"
}

# Trusted source CIDR for SSH access (must not be open to internet)
variable "trusted_ssh_cidr" {
  description = "Trusted source CIDR permitted to SSH to the private VM"
  type        = string

  validation {
    condition     = var.trusted_ssh_cidr != "0.0.0.0/0"
    error_message = "Security Violation: SSH access must not be allowed from the entire internet (0.0.0.0/0)."
  }
}

# Domain name for Google-managed SSL Certificate termination on the Load Balancer
variable "domain_name" {
  description = "Domain name for the managed SSL certificate"
  type        = string
  default     = "app.example.com"
}
