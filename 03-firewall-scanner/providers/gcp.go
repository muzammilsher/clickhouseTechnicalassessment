package providers

import (
    "context"
    "os"
    "path"
    "strings"

    "cloud-security-firewall-scanner/models"
    compute "cloud.google.com/go/compute/apiv1"
    "cloud.google.com/go/compute/apiv1/computepb"
    "google.golang.org/api/iterator"
)

// ScanGCP authenticates via Google Cloud Compute REST API, lists project firewall rules,
// and identifies any ingress rule with sourceRange 0.0.0.0/0
func ScanGCP(ctx context.Context) ([]models.GCPFinding, error) {
    // Resolve project ID from standard GCP environment variables
    projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
    if projectID == "" {
        projectID = os.Getenv("GCP_PROJECT")
    }
    if projectID == "" {
        projectID = os.Getenv("CLOUDSDK_CORE_PROJECT")
    }
    if projectID == "" {
        return nil, &missingEnvironmentError{"GOOGLE_CLOUD_PROJECT"}
    }

    client, err := compute.NewFirewallsRESTClient(ctx)
    if err != nil {
        return nil, err
    }
    defer client.Close()

    // Query project-scoped firewall rules
    it := client.List(ctx, &computepb.ListFirewallsRequest{
        Project: projectID,
    })

    var findings []models.GCPFinding

    for {
        firewall, err := it.Next()
        if err == iterator.Done {
            break
        }
        if err != nil {
            return nil, err
        }

        var insecure []models.Finding

        for _, source := range firewall.GetSourceRanges() {
            if source != "0.0.0.0/0" {
                continue
            }

            for _, allowed := range firewall.GetAllowed() {
                insecure = append(insecure, models.Finding{
                    Protocol: allowed.GetIPProtocol(),
                    Port:     gcpPorts(allowed.GetPorts()),
                    Source:   source,
                })
            }

            if len(firewall.GetAllowed()) == 0 {
                insecure = append(insecure, models.Finding{
                    Protocol: "all",
                    Port:     "all",
                    Source:   source,
                })
            }
        }

        if len(insecure) > 0 {
            networkName := path.Base(firewall.GetNetwork())
            findings = append(findings, models.GCPFinding{
                Name:          firewall.GetName(),
                Network:       networkName,
                InsecureRules: insecure,
            })
        }
    }

    return findings, nil
}

func gcpPorts(ports []string) string {
    if len(ports) == 0 {
        return "all"
    }
    return strings.Join(ports, ",")
}
