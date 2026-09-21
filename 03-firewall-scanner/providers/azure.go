package providers

import (
    "context"
    "os"
    "strings"

    "cloud-security-firewall-scanner/models"
    "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
    "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v6"
)

// ScanAzure authenticates via Azure SDK v6 (azidentity), lists Network Security Groups (NSGs),
// and identifies inbound Allow rules permitting ingress from 0.0.0.0/0, *, or Internet
func ScanAzure(ctx context.Context) ([]models.AzureFinding, error) {
    subscriptionID := os.Getenv("AZURE_SUBSCRIPTION_ID")
    if subscriptionID == "" {
        return nil, &missingEnvironmentError{"AZURE_SUBSCRIPTION_ID"}
    }

    // Authenticate using DefaultAzureCredential (supports environment, managed identity, CLI)
    cred, err := azidentity.NewDefaultAzureCredential(nil)
    if err != nil {
        return nil, err
    }

    // Azure SDK v6 constructor: armnetwork.NewSecurityGroupsClient
    client, err := armnetwork.NewSecurityGroupsClient(subscriptionID, cred, nil)
    if err != nil {
        return nil, err
    }

    pager := client.NewListAllPager(nil)
    var findings []models.AzureFinding

    for pager.More() {
        page, err := pager.NextPage(ctx)
        if err != nil {
            return nil, err
        }

        for _, nsg := range page.Value {
            var insecure []models.Finding

            if nsg.Properties == nil {
                continue
            }

            for _, rule := range nsg.Properties.SecurityRules {
                if rule.Properties == nil {
                    continue
                }

                if !strings.EqualFold(stringValue(rule.Properties.Direction), "Inbound") ||
                    !strings.EqualFold(stringValue(rule.Properties.Access), "Allow") {
                    continue
                }

                var openSources []string
                if single := stringValue(rule.Properties.SourceAddressPrefix); isAzureOpenSource(single) {
                    openSources = append(openSources, strings.TrimSpace(single))
                }
                for _, multi := range rule.Properties.SourceAddressPrefixes {
                    if val := stringValue(multi); isAzureOpenSource(val) {
                        openSources = append(openSources, strings.TrimSpace(val))
                    }
                }

                if len(openSources) == 0 {
                    continue
                }

                port := stringValue(rule.Properties.DestinationPortRange)
                if port == "" && len(rule.Properties.DestinationPortRanges) > 0 {
                    var ranges []string
                    for _, pr := range rule.Properties.DestinationPortRanges {
                        ranges = append(ranges, stringValue(pr))
                    }
                    port = strings.Join(ranges, ",")
                }
                if port == "" {
                    port = "*"
                }

                for _, source := range openSources {
                    insecure = append(insecure, models.Finding{
                        Protocol: stringValue(rule.Properties.Protocol),
                        Port:     port,
                        Source:   source,
                    })
                }
            }

            if len(insecure) > 0 {
                findings = append(findings, models.AzureFinding{
                    ID:            stringValue(nsg.ID),
                    Name:          stringValue(nsg.Name),
                    InsecureRules: insecure,
                })
            }
        }
    }

    return findings, nil
}


func isAzureOpenSource(source string) bool {
    s := strings.TrimSpace(source)
    return s == "0.0.0.0/0" || s == "*" || strings.EqualFold(s, "Internet")
}
