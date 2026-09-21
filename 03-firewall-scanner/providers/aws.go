package providers

import (
    "context"
    "fmt"

    "cloud-security-firewall-scanner/models"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/ec2"
)

// ScanAWS authenticates via AWS SDK v2, lists EC2 Security Groups,
// and identifies any inbound rule permitting ingress from 0.0.0.0/0
func ScanAWS(ctx context.Context) ([]models.SecurityGroupFinding, error) {
    // Authenticate using standard AWS credential chain (env vars, IAM role, or ~/.aws)
    cfg, err := config.LoadDefaultConfig(ctx)
    if err != nil {
        return nil, err
    }

    client := ec2.NewFromConfig(cfg)
    paginator := ec2.NewDescribeSecurityGroupsPaginator(client, &ec2.DescribeSecurityGroupsInput{})

    var findings []models.SecurityGroupFinding

    for paginator.HasMorePages() {
        page, err := paginator.NextPage(ctx)
        if err != nil {
            return nil, err
        }

        for _, sg := range page.SecurityGroups {
            var insecure []models.Finding

            // Inspect inbound permission blocks for public 0.0.0.0/0 CIDRs
            for _, permission := range sg.IpPermissions {
                for _, ipRange := range permission.IpRanges {
                    if ipRange.CidrIp == nil || *ipRange.CidrIp != "0.0.0.0/0" {
                        continue
                    }

                    insecure = append(insecure, models.Finding{
                        Protocol: protocolValue(permission.IpProtocol),
                        Port:     awsPort(permission.FromPort, permission.ToPort),
                        Source:   "0.0.0.0/0",
                    })
                }
            }

            if len(insecure) > 0 {
                findings = append(findings, models.SecurityGroupFinding{
                    ID:            stringValue(sg.GroupId),
                    Name:          stringValue(sg.GroupName),
                    InsecureRules: insecure,
                })
            }
        }
    }

    return findings, nil
}

func protocolValue(v *string) string {
    if v == nil {
        return "*"
    }
    return *v
}

func awsPort(from, to *int32) any {
    if from == nil && to == nil {
        return "*"
    }
    if from != nil && to != nil && *from == *to {
        return int(*from)
    }
    if from != nil && to != nil {
        return fmt.Sprintf("%d-%d", *from, *to)
    }
    return "*"
}

