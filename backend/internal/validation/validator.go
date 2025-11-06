package validation

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	// Valid Docker image format
	imageRegex = regexp.MustCompile(`^[a-z0-9]+(?:[._-][a-z0-9]+)*(?:/[a-z0-9]+(?:[._-][a-z0-9]+)*)*$`)
	// Valid service name (alphanumeric, hyphens, underscores)
	nameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,64}$`)
	// Valid version format (semver-like)
	versionRegex = regexp.MustCompile(`^[a-zA-Z0-9._-]{1,32}$`)
)

var validRegions = map[string]bool{
	"us-east-1":    true,
	"us-west-1":    true,
	"us-west-2":    true,
	"eu-west-1":    true,
	"eu-central-1": true,
	"ap-southeast-1": true,
	"ap-northeast-1": true,
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

type ValidationErrors []ValidationError

func (e ValidationErrors) Error() string {
	var msgs []string
	for _, err := range e {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

func ValidateServiceName(name string) error {
	if name == "" {
		return ValidationError{Field: "name", Message: "name is required"}
	}
	if !nameRegex.MatchString(name) {
		return ValidationError{Field: "name", Message: "name must be 3-64 alphanumeric characters, hyphens, or underscores"}
	}
	return nil
}

func ValidateRegion(region string) error {
	if region == "" {
		return ValidationError{Field: "region", Message: "region is required"}
	}
	if !validRegions[region] {
		return ValidationError{Field: "region", Message: fmt.Sprintf("invalid region, must be one of: %v", getRegionList())}
	}
	return nil
}

func ValidateImage(image string) error {
	if image == "" {
		return ValidationError{Field: "image", Message: "image is required"}
	}
	if len(image) > 255 {
		return ValidationError{Field: "image", Message: "image name too long (max 255 characters)"}
	}
	// Allow image:tag format
	parts := strings.Split(image, ":")
	if len(parts) > 2 {
		return ValidationError{Field: "image", Message: "invalid image format"}
	}
	if !imageRegex.MatchString(parts[0]) {
		return ValidationError{Field: "image", Message: "invalid image name format"}
	}
	return nil
}

func ValidateVersion(version string) error {
	if version == "" {
		return ValidationError{Field: "version", Message: "version is required"}
	}
	if !versionRegex.MatchString(version) {
		return ValidationError{Field: "version", Message: "invalid version format"}
	}
	return nil
}

func getRegionList() []string {
	regions := make([]string, 0, len(validRegions))
	for region := range validRegions {
		regions = append(regions, region)
	}
	return regions
}
