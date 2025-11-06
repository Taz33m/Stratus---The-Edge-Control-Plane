package validation

import (
	"testing"
)

func TestValidateServiceName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid name", "my-service", false},
		{"valid with underscore", "my_service", false},
		{"valid alphanumeric", "service123", false},
		{"empty name", "", true},
		{"too short", "ab", true},
		{"too long", "a123456789012345678901234567890123456789012345678901234567890123456", true},
		{"invalid characters", "my service!", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateServiceName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateServiceName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateRegion(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid region", "us-east-1", false},
		{"another valid region", "eu-west-1", false},
		{"invalid region", "invalid-region", true},
		{"empty region", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRegion(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRegion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateImage(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid image", "nginx", false},
		{"valid with tag", "nginx:latest", false},
		{"valid with registry", "myregistry/nginx", false},
		{"empty image", "", true},
		{"too long", string(make([]byte, 300)), true},
		{"multiple colons", "nginx:1.0:latest", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateImage(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateImage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateVersion(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid semver", "1.0.0", false},
		{"valid with v prefix", "v1.2.3", false},
		{"valid simple", "latest", false},
		{"empty version", "", true},
		{"invalid characters", "1.0.0@beta", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateVersion(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateVersion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
