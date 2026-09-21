# Cloud Security Engineer Technical Assessment

Welcome to the Cloud Security Engineer technical assessment. This take-home assessment is designed to evaluate your knowledge and skills in cloud security engineering, specifically focusing on scenario-based problem-solving.

The assessment will cover the following key areas:

- Infrastructure as Code
- Automation
- Vulnerability Management

You will be presented with real-world scenarios and will be expected to provide solutions that demonstrate your understanding of cloud security best practices and your ability to implement them using the relevant tools and technologies.

You will have until your next interview to complete the assessment. Please note that while you have flexibility with time, we recommend managing your time effectively to ensure you can submit your best work.

The completed assessment will be submitted as a link to a public GitHub repository in your account.

## 1. Secure GCP Infrastructure with Terraform

Write a Terraform module that deploys the following resources securely in Google Cloud Platform (GCP).

### Requirements

- A VPC network with a public and private subnet.
- A Compute Engine instance in the private subnet.
- A Global HTTP Load Balancer in the public subnet that routes traffic to the Compute Engine instance.
- Firewall rules that:
  - Allow only HTTPS (443) from the internet to the Load Balancer.
  - Allow only SSH (22) from a trusted IP range (not `0.0.0.0/0`) to the Compute Engine instance.
  - Allow only HTTP (80) from the Load Balancer to the Compute Engine instance.

### Bonus

Implement least-privilege IAM roles for resources instead of using `roles/editor` or `roles/owner`.

## 2. Cloud Security Check

You are given the following CloudFormation template (YAML format) and Dockerfile.

### CloudFormation template

```yaml
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
  MyIAMUser:
    Type: AWS::IAM::User
    Properties:
      UserName: "UnrestrictedUser"
      Policies:
        - PolicyName: "FullAccess"
          PolicyDocument:
            Statement:
              - Effect: Allow
                Action: "*"
                Resource: "*"
```

### Dockerfile

```dockerfile
FROM ubuntu:18.04

USER root

RUN apt-get update && apt-get install -y \
    python3 \
    python3-pip \
    curl \
    vim \
    dnsutils

ENV DB_PASSWORD="SuperSecretPassword123!"
ENV API_KEY="prod-key-778899"

COPY . /app
RUN chmod -R 777 /app
EXPOSE 80

WORKDIR /app
CMD python3 app.py
```

### Tasks

#### CloudFormation template

- Identify at least three security issues with the CloudFormation template.
- Write an improved version of the CloudFormation template that follows AWS security best practices.
- Assuming the user with these permissions exists in the environment, explain how you would find misconfigurations like this and how you would fix them.

#### Dockerfile

- What are the issues with the Dockerfile?
- How do we address those issues by following security best practices?

#### Prevention

- How would you prevent the misconfigurations from being introduced or applied to the environment for both files?

## 3. Automate Security Checks with Go

### Task

Write a Go CLI tool that:

- Scans AWS, GCP, and Azure firewall/security group rules.
- Detects any rule that allows inbound traffic from `0.0.0.0/0` (open to the internet) for any protocol and any port.
- Outputs the list of insecure rules in JSON format.

### Expected implementation

Your Go program should:

1. Authenticate with AWS, GCP, and Azure.
2. Fetch all firewall/security group rules from each cloud provider.
3. Identify rules where `0.0.0.0/0` is allowed for any protocol/port.
4. Output the findings in JSON format.

### Cloud-specific details

#### AWS security check

- Use the AWS SDK for Go to list EC2 Security Groups.
- Identify rules allowing inbound traffic from `0.0.0.0/0` for any protocol/port.

#### GCP security check

- Use the Google Cloud SDK for Go to list firewall rules.
- Look for `sourceRanges: ["0.0.0.0/0"]` allowing any protocol/port.

#### Azure security check

- Use the Azure SDK for Go to list Network Security Groups (NSGs).
- Look for NSG rules allowing `Any` source (`0.0.0.0/0`) for any protocol/port.

### Example JSON output

```json
{
  "aws": {
    "security_groups": [
      {
        "id": "sg-12345678",
        "name": "open-sg",
        "insecure_rules": [
          {
            "protocol": "tcp",
            "port": 8080,
            "source": "0.0.0.0/0"
          },
          {
            "protocol": "udp",
            "port": 123,
            "source": "0.0.0.0/0"
          }
        ]
      }
    ]
  },
  "gcp": {
    "firewall_rules": [
      {
        "name": "allow-all",
        "network": "default",
        "insecure_rules": [
          {
            "protocol": "icmp",
            "port": "all",
            "source": "0.0.0.0/0"
          }
        ]
      }
    ]
  },
  "azure": {
    "nsgs": [
      {
        "id": "/subscriptions/xxxx/resourceGroups/my-rg/providers/Microsoft.Network/networkSecurityGroups/my-nsg",
        "name": "open-nsg",
        "insecure_rules": [
          {
            "protocol": "*",
            "port": "*",
            "source": "0.0.0.0/0"
          }
        ]
      }
    ]
  }
}
```

## 4. Third-Party Dependency Security Audit

### Scenario

You are given a `package.json` file containing dependencies for a Node.js application. Your task is to scan it for vulnerabilities, analyze the results, and provide security recommendations.

### Instructions

1. Analyze the given `package.json` to identify security vulnerabilities.
2. Use at least one vulnerability scanning tool, such as:
   - `npm audit`
   - `yarn audit`
   - `snyk test`
   - OWASP Dependency-Check
3. Document:
   - The list of detected vulnerabilities (high/critical priority).
   - Which dependencies are vulnerable.
   - How to fix the issues, such as updating or replacing dependencies.
   - How many issues can be fixed.
4. **Bonus:** Automate the scan with a Python or Go script that:
   - Parses the vulnerabilities.
   - Generates a CSV or JSON report.

### Example `package.json` for testing

```json
{
  "name": "vulnerable-node-app",
  "version": "1.0.0",
  "dependencies": {
    "body-parser": "^1.19.0",
    "dotenv": "^8.2.0",
    "express": "4.15.0",
    "lodash": "4.17.10",
    "jsonwebtoken": "8.1.0",
    "mongoose": "^4.2.4"
  },
  "devDependencies": {
    "mocha": "5.0.5"
  }
}
```

### Expected deliverables

- A list of vulnerabilities, for example: “Lodash prototype pollution vulnerability in version `4.17.10`.”
- Remediation steps, for example: “Upgrade lodash to `4.17.21`.”
- A script to automate vulnerability reporting (optional bonus).
