# Dedicated VPC network with public and private subnets
resource "google_compute_network" "secure_vpc" {
  name                    = "secure-cloud-vpc"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "public" {
  name          = "public-subnet"
  ip_cidr_range = "10.10.1.0/24"
  region        = var.region
  network       = google_compute_network.secure_vpc.id
}

resource "google_compute_subnetwork" "private" {
  name                     = "private-subnet"
  ip_cidr_range            = "10.10.2.0/24"
  region                   = var.region
  network                  = google_compute_network.secure_vpc.id
  private_ip_google_access = true
}

# Cloud NAT for private subnet outbound connectivity
resource "google_compute_router" "router" {
  name    = "secure-cloud-router"
  region  = var.region
  network = google_compute_network.secure_vpc.id
}

resource "google_compute_router_nat" "nat" {
  name                               = "private-egress-nat"
  router                             = google_compute_router.router.name
  region                             = var.region
  nat_ip_allocate_option             = "AUTO_ONLY"
  source_subnetwork_ip_ranges_to_nat = "LIST_OF_SUBNETWORKS"

  subnetwork {
    name                    = google_compute_subnetwork.private.id
    source_ip_ranges_to_nat = ["ALL_IP_RANGES"]
  }
}

# Dedicated service account with least-privilege IAM bindings
resource "google_service_account" "application" {
  account_id   = "app-workload-sa"
  display_name = "Least Privilege Application Service Account"
  description  = "Dedicated service account for the private compute workload"
}

resource "google_project_iam_member" "application_logging" {
  project = var.project_id
  role    = "roles/logging.logWriter"
  member  = "serviceAccount:${google_service_account.application.email}"
}

resource "google_project_iam_member" "application_monitoring" {
  project = var.project_id
  role    = "roles/monitoring.metricWriter"
  member  = "serviceAccount:${google_service_account.application.email}"
}

# Compute Engine instance deployed to private subnet without external IP
resource "google_compute_instance" "application" {
  name         = "private-web-server"
  machine_type = "e2-small"
  zone         = var.zone

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-12"
      size  = 20
      type  = "pd-balanced"
    }
  }

  network_interface {
    subnetwork = google_compute_subnetwork.private.id
  }

  tags = ["private-web-server"]

  metadata = {
    enable-oslogin = "TRUE"
  }

  service_account {
    email  = google_service_account.application.email
    scopes = ["cloud-platform"]
  }
}

resource "google_compute_firewall" "allow_lb_to_backend" {
  name    = "allow-load-balancer-to-backend"
  network = google_compute_network.secure_vpc.name

  direction = "INGRESS"

  allow {
    protocol = "tcp"
    ports    = ["80"]
  }

  source_ranges = [
    "35.191.0.0/16",
    "130.211.0.0/22"
  ]

  target_tags = ["private-web-server"]
}

resource "google_compute_firewall" "allow_trusted_ssh" {
  name    = "allow-trusted-ssh"
  network = google_compute_network.secure_vpc.name

  direction = "INGRESS"

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }

  source_ranges = [var.trusted_ssh_cidr]
  target_tags   = ["private-web-server"]
}

resource "google_compute_health_check" "application" {
  name = "application-health-check"

  http_health_check {
    port         = 80
    request_path = "/"
  }
}

resource "google_compute_instance_group" "application" {
  name      = "application-instance-group"
  zone      = var.zone
  instances = [google_compute_instance.application.self_link]

  named_port {
    name = "http"
    port = 80
  }
}

resource "google_compute_backend_service" "application" {
  name                  = "application-backend"
  protocol              = "HTTP"
  port_name             = "http"
  load_balancing_scheme = "EXTERNAL_MANAGED"
  health_checks         = [google_compute_health_check.application.id]

  backend {
    group = google_compute_instance_group.application.id
  }
}

resource "google_compute_url_map" "application" {
  name            = "application-url-map"
  default_service = google_compute_backend_service.application.id
}

resource "google_compute_managed_ssl_certificate" "application" {
  name = "application-certificate"

  managed {
    domains = [var.domain_name]
  }
}

resource "google_compute_target_https_proxy" "application" {
  name             = "application-https-proxy"
  url_map          = google_compute_url_map.application.id
  ssl_certificates = [google_compute_managed_ssl_certificate.application.id]
}

resource "google_compute_global_forwarding_rule" "https" {
  name                  = "application-https-forwarding-rule"
  target                = google_compute_target_https_proxy.application.id
  port_range            = "443"
  load_balancing_scheme = "EXTERNAL_MANAGED"
}

# HTTP-to-HTTPS redirect
resource "google_compute_url_map" "http_redirect" {
  name = "application-http-redirect"

  default_url_redirect {
    https_redirect         = true
    strip_query            = false
    redirect_response_code = "MOVED_PERMANENTLY_DEFAULT"
  }
}

resource "google_compute_target_http_proxy" "http_redirect" {
  name    = "application-http-proxy"
  url_map = google_compute_url_map.http_redirect.id
}

resource "google_compute_global_forwarding_rule" "http" {
  name                  = "application-http-forwarding-rule"
  target                = google_compute_target_http_proxy.http_redirect.id
  port_range            = "80"
  load_balancing_scheme = "EXTERNAL_MANAGED"
}
