package providers

import "testing"

func TestAWSOpenSourceDetectionValues(t *testing.T) {
    if protocolValue(nil) != "*" {
        t.Fatalf("expected wildcard protocol")
    }
    p := "tcp"
    if protocolValue(&p) != "tcp" {
        t.Fatalf("expected tcp protocol")
    }
    if awsPort(nil, nil) != "*" {
        t.Fatalf("expected wildcard port")
    }
    if awsPort(int32ptr(443), int32ptr(443)) != 443 {
        t.Fatalf("expected single port integer 443")
    }
    if awsPort(int32ptr(80), int32ptr(85)) != "80-85" {
        t.Fatalf("expected port range 80-85")
    }
}

func TestGCPPortsFormatting(t *testing.T) {
    if gcpPorts(nil) != "all" {
        t.Fatalf("expected 'all' for nil ports")
    }
    if gcpPorts([]string{}) != "all" {
        t.Fatalf("expected 'all' for empty ports")
    }
    if gcpPorts([]string{"80", "443"}) != "80,443" {
        t.Fatalf("expected comma-joined ports")
    }
}

func TestMissingEnvironmentError(t *testing.T) {
    err := &missingEnvironmentError{"GOOGLE_CLOUD_PROJECT"}
    expected := "required environment variable is not set: GOOGLE_CLOUD_PROJECT"
    if err.Error() != expected {
        t.Fatalf("expected %q, got %q", expected, err.Error())
    }
}

func TestAzureOpenSourceDetection(t *testing.T) {
    cases := []struct {
        source   string
        expected bool
    }{
        {"0.0.0.0/0", true},
        {"*", true},
        {"Internet", true},
        {"internet", true},
        {"10.0.0.0/16", false},
        {"VirtualNetwork", false},
        {"", false},
    }

    for _, c := range cases {
        if isAzureOpenSource(c.source) != c.expected {
            t.Fatalf("isAzureOpenSource(%q): expected %v", c.source, c.expected)
        }
    }
}

func TestStringValueHelper(t *testing.T) {
    if stringValue((*string)(nil)) != "" {
        t.Fatalf("expected empty string for nil pointer")
    }
    val := "test"
    if stringValue(&val) != "test" {
        t.Fatalf("expected 'test'")
    }
}

func int32ptr(v int32) *int32 {
    return &v
}
