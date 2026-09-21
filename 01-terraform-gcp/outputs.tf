# Infrastructure outputs

# VPC Network identifier
output "vpc_name" {
  description = "The name of the VPC network created"
  value       = google_compute_network.secure_vpc.name
}

# Private subnet hosting the isolated Compute Engine instance
output "private_subnet" {
  description = "The name of the private subnet (no external IPs)"
  value       = google_compute_subnetwork.private.name
}

# Compute Engine VM instance identifier
output "instance_name" {
  description = "The name of the private Compute Engine VM instance"
  value       = google_compute_instance.application.name
}

# Public ingress IP for the Global External HTTPS Load Balancer
output "https_forwarding_rule_ip" {
  description = "Public IP address of the Global External HTTPS Load Balancer"
  value       = google_compute_global_forwarding_rule.https.ip_address
}

# Dedicated service account email
output "service_account_email" {
  description = "Email of the dedicated least-privilege application service account"
  value       = google_service_account.application.email
}
