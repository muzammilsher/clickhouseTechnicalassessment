"""
Parses 'npm audit --json' output, extracts direct and transitive vulnerabilities,
ranks findings by severity (Critical -> High -> Moderate -> Low),
and exports structured reports in JSON or CSV format.
"""

import csv
import json
import sys
from pathlib import Path


def load_json(path: str):
    """Load and parse raw npm audit JSON file."""
    with open(path, "r", encoding="utf-8") as handle:
        return json.load(handle)


def normalise_via(via):
    result = []
    for item in via or []:
        if isinstance(item, str):
            result.append(item)
        elif isinstance(item, dict):
            result.append({
                "source": item.get("source"),
                "severity": item.get("severity"),
                "title": item.get("title"),
                "url": item.get("url"),
                "range": item.get("range"),
            })
    return result


def build_report(audit):
    """Normalize npm audit entries and sort by severity (Critical -> High -> Moderate -> Low)."""
    findings = []

    for package, details in audit.get("vulnerabilities", {}).items():
        findings.append({
            "package": package,
            "severity": details.get("severity"),
            "is_direct": details.get("isDirect"),
            "range": details.get("range"),
            "via": normalise_via(details.get("via")),
            "fix_available": details.get("fixAvailable"),
        })

    # Strict priority ranking: critical vulnerabilities surfaced first
    severity_order = {"critical": 0, "high": 1, "moderate": 2, "low": 3}
    findings.sort(key=lambda x: severity_order.get(x["severity"], 99))

    return {
        "audit_version": audit.get("auditReportVersion"),
        "metadata": audit.get("metadata", {}),
        "findings": findings,
    }


def write_csv(report, out_path):
    """Export normalized findings into CSV spreadsheet format for audits and ticketing."""
    fieldnames = ["package", "severity", "is_direct", "range", "fix_available", "advisories"]
    with open(out_path, "w", encoding="utf-8", newline="") as f:
        writer = csv.DictWriter(f, fieldnames=fieldnames)
        writer.writeheader()
        for item in report["findings"]:
            advisories = []
            for v in item["via"]:
                if isinstance(v, dict):
                    advisories.append(f"{v.get('title', '')} ({v.get('url', '')})")
                elif isinstance(v, str):
                    advisories.append(v)
            fix = item["fix_available"]
            fix_str = f"{fix.get('name')}@{fix.get('version')}" if isinstance(fix, dict) else str(fix)
            writer.writerow({
                "package": item["package"],
                "severity": item["severity"],
                "is_direct": item["is_direct"],
                "range": item["range"],
                "fix_available": fix_str,
                "advisories": " | ".join(advisories),
            })


def main():
    input_path = "npm-audit.json"
    output_path = "vulnerability-report.json"

    if len(sys.argv) >= 2:
        input_path = sys.argv[1]
    if len(sys.argv) >= 3:
        output_path = sys.argv[2]

    if not Path(input_path).exists():
        print(f"Error: input file {input_path} does not exist.")
        print(f"Usage: {Path(sys.argv[0]).name} [npm-audit.json] [vulnerability-report.json|report.csv]")
        sys.exit(1)

    audit = load_json(input_path)
    report = build_report(audit)

    if output_path.lower().endswith(".csv"):
        write_csv(report, output_path)
    else:
        with open(output_path, "w", encoding="utf-8") as handle:
            json.dump(report, handle, indent=2)

    print(f"Wrote {len(report['findings'])} package findings to {output_path}")


if __name__ == "__main__":
    main()
