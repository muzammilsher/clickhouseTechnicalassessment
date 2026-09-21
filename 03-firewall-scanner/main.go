package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"cloud-security-firewall-scanner/models"
	"cloud-security-firewall-scanner/providers"
)

// Audits AWS Security Groups, GCP Firewall Rules, and Azure NSGs for 0.0.0.0/0 ingress.
func main() {
	// Allows targeting specific clouds (e.g., -providers=aws,gcp) or all three by default
	providerFlag := flag.String("providers", "aws,gcp,azure", "Comma-separated providers")
	flag.Parse()

	ctx := context.Background()
	wanted := map[string]bool{}

	for _, p := range strings.Split(*providerFlag, ",") {
		wanted[strings.ToLower(strings.TrimSpace(p))] = true
	}

	result := models.ScanResult{}

	if wanted["aws"] {
		findings, err := providers.ScanAWS(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "AWS scan failed: %v\n", err)
			os.Exit(1)
		}
		result.AWS.SecurityGroups = findings
	}

	if wanted["gcp"] {
		findings, err := providers.ScanGCP(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "GCP scan failed: %v\n", err)
			os.Exit(1)
		}
		result.GCP.FirewallRules = findings
	}

	if wanted["azure"] {
		findings, err := providers.ScanAzure(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Azure scan failed: %v\n", err)
			os.Exit(1)
		}
		result.Azure.NSGs = findings
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(result); err != nil {
		fmt.Fprintf(os.Stderr, "failed to encode JSON: %v\n", err)
		os.Exit(1)
	}
}
