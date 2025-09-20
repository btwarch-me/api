package utils

import (
	"btwarch/config"
	"fmt"
	"strings"
)

// ValidateSubdomainName validates if a subdomain name is valid for claiming
func ValidateSubdomainName(subdomainName string) error {
	if subdomainName == "" {
		return fmt.Errorf("subdomain name cannot be empty")
	}

	// Check if it contains only valid characters (alphanumeric and hyphens)
	for _, char := range subdomainName {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') || char == '-') {
			return fmt.Errorf("subdomain name can only contain letters, numbers, and hyphens")
		}
	}

	// Check if it starts or ends with hyphen
	if strings.HasPrefix(subdomainName, "-") || strings.HasSuffix(subdomainName, "-") {
		return fmt.Errorf("subdomain name cannot start or end with hyphen")
	}

	// Check length
	if len(subdomainName) < 1 || len(subdomainName) > 63 {
		return fmt.Errorf("subdomain name must be between 1 and 63 characters")
	}

	return nil
}

// ValidateRecordName validates if a record name is valid for the given record type and user's subdomain
func ValidateRecordName(recordName, recordType, userSubdomain string) error {
	cfg := config.LoadConfig()
	parentDomain := cfg.ParentDomain

	// Build the full subdomain name
	fullSubdomain := userSubdomain + "." + parentDomain

	// For A, AAAA, CNAME records - only allow the root subdomain
	if recordType == "A" || recordType == "AAAA" || recordType == "CNAME" {
		if recordName != fullSubdomain && recordName != userSubdomain {
			return fmt.Errorf("A, AAAA, and CNAME records can only be created for the root subdomain (%s)", fullSubdomain)
		}
	}

	// For TXT records - allow root subdomain and sub-labels
	if recordType == "TXT" {
		// Allow root subdomain
		if recordName == fullSubdomain || recordName == userSubdomain {
			return nil
		}

		// Allow sub-labels under the user's subdomain
		if strings.HasSuffix(recordName, "."+fullSubdomain) || strings.HasSuffix(recordName, "."+userSubdomain) {
			return nil
		}

		return fmt.Errorf("TXT records can only be created for the root subdomain or sub-labels under %s", fullSubdomain)
	}

	// Block NS and MX records
	if recordType == "NS" || recordType == "MX" {
		return fmt.Errorf("NS and MX records are not allowed")
	}

	return nil
}

// GetFullSubdomainName returns the full subdomain name with parent domain
func GetFullSubdomainName(subdomainName string) string {
	cfg := config.LoadConfig()
	return subdomainName + "." + cfg.ParentDomain
}

// ExtractSubdomainFromRecordName extracts the subdomain name from a full record name
func ExtractSubdomainFromRecordName(recordName string) string {
	cfg := config.LoadConfig()
	parentDomain := cfg.ParentDomain

	// If record name ends with parent domain, extract the subdomain part
	if strings.HasSuffix(recordName, "."+parentDomain) {
		subdomain := strings.TrimSuffix(recordName, "."+parentDomain)
		// Extract the first part (subdomain) from subdomain.something.xyz.com
		parts := strings.Split(subdomain, ".")
		if len(parts) > 0 {
			return parts[0]
		}
	}

	// If record name doesn't have parent domain, it might be just the subdomain
	// or a sub-label under the subdomain
	return ""
}
