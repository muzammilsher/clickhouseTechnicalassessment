package models

type Finding struct {
    Protocol string `json:"protocol"`
    Port     any    `json:"port"`
    Source   string `json:"source"`
}

type SecurityGroupFinding struct {
    ID            string    `json:"id"`
    Name          string    `json:"name"`
    InsecureRules []Finding `json:"insecure_rules"`
}

type GCPFinding struct {
    Name          string    `json:"name"`
    Network       string    `json:"network"`
    InsecureRules []Finding `json:"insecure_rules"`
}

type AzureFinding struct {
    ID            string    `json:"id"`
    Name          string    `json:"name"`
    InsecureRules []Finding `json:"insecure_rules"`
}

type ScanResult struct {
    AWS struct {
        SecurityGroups []SecurityGroupFinding `json:"security_groups"`
    } `json:"aws"`

    GCP struct {
        FirewallRules []GCPFinding `json:"firewall_rules"`
    } `json:"gcp"`

    Azure struct {
        NSGs []AzureFinding `json:"nsgs"`
    } `json:"azure"`
}
