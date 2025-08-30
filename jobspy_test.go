package jobspy

import (
	"testing"
)

func TestJobTypeFromString(t *testing.T) {
	tests := []struct {
		input    string
		expected *JobType
	}{
		{"full-time", &[]JobType{JobTypeFullTime}[0]},
		{"FULLTIME", &[]JobType{JobTypeFullTime}[0]},
		{"part-time", &[]JobType{JobTypePartTime}[0]},
		{"contract", &[]JobType{JobTypeContract}[0]},
		{"internship", &[]JobType{JobTypeInternship}[0]},
		{"temporary", &[]JobType{JobTypeTemporary}[0]},
		{"unknown", nil},
		{"", nil},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := JobTypeFromString(tt.input)
			if tt.expected == nil {
				if result != nil {
					t.Errorf("Expected nil, got %v", result)
				}
			} else {
				if result == nil {
					t.Errorf("Expected %v, got nil", *tt.expected)
				} else if *result != *tt.expected {
					t.Errorf("Expected %v, got %v", *tt.expected, *result)
				}
			}
		})
	}
}

func TestMapStringToSite(t *testing.T) {
	tests := []struct {
		input     string
		expected  Site
		shouldErr bool
	}{
		{"indeed", SiteIndeed, false},
		{"linkedin", SiteLinkedIn, false},
		{"glassdoor", SiteGlassdoor, false},
		{"google", SiteGoogle, false},
		{"ziprecruiter", SiteZipRecruiter, false},
		{"zip_recruiter", SiteZipRecruiter, false},
		{"bayt", SiteBayt, false},
		{"naukri", SiteNaukri, false},
		{"bdjobs", SiteBDJobs, false},
		{"unknown", "", true},
		{"", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := MapStringToSite(tt.input)
			if tt.shouldErr {
				if err == nil {
					t.Errorf("Expected error for input %q", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for input %q: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("Expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}

func TestCountryFromString(t *testing.T) {
	tests := []struct {
		input    string
		expected *Country
	}{
		{"usa", &[]Country{CountryUSA}[0]},
		{"US", &[]Country{CountryUSA}[0]},
		{"united states", &[]Country{CountryUSA}[0]},
		{"canada", &[]Country{CountryCanada}[0]},
		{"ca", &[]Country{CountryCanada}[0]},
		{"uk", &[]Country{CountryUK}[0]},
		{"united kingdom", &[]Country{CountryUK}[0]},
		{"germany", &[]Country{CountryGermany}[0]},
		{"france", &[]Country{CountryFrance}[0]},
		{"unknown", nil},
		{"", nil},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := CountryFromString(tt.input)
			if tt.expected == nil {
				if result != nil {
					t.Errorf("Expected nil, got %v", result)
				}
			} else {
				if result == nil {
					t.Errorf("Expected %v, got nil", *tt.expected)
				} else if *result != *tt.expected {
					t.Errorf("Expected %v, got %v", *tt.expected, *result)
				}
			}
		})
	}
}

func TestExtractEmailsFromText(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"Contact us at jobs@example.com", []string{"jobs@example.com"}},
		{"Email: test@company.com or hr@company.com", []string{"test@company.com", "hr@company.com"}},
		{"No emails here", nil},
		{"", nil},
		{"Invalid email: not-an-email", nil},
		{"Reach out to john.doe+hiring@example-company.com", []string{"john.doe+hiring@example-company.com"}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ExtractEmailsFromText(tt.input)
			if len(tt.expected) == 0 && len(result) == 0 {
				return // Both nil/empty - success
			}
			if len(tt.expected) != len(result) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
				return
			}
			for i, expected := range tt.expected {
				if result[i] != expected {
					t.Errorf("Expected %v, got %v", tt.expected, result)
					return
				}
			}
		})
	}
}

func TestLocationString(t *testing.T) {
	tests := []struct {
		name     string
		location *Location
		expected string
	}{
		{
			name: "full location",
			location: &Location{
				City:    testStringPtr("San Francisco"),
				State:   testStringPtr("CA"),
				Country: &[]Country{CountryUSA}[0],
			},
			expected: "San Francisco, CA, USA",
		},
		{
			name: "city and state only",
			location: &Location{
				City:  testStringPtr("New York"),
				State: testStringPtr("NY"),
			},
			expected: "New York, NY",
		},
		{
			name: "city only",
			location: &Location{
				City: testStringPtr("Boston"),
			},
			expected: "Boston",
		},
		{
			name:     "nil location",
			location: nil,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result string
			if tt.location != nil {
				result = tt.location.String()
			}
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestNewScraperInput(t *testing.T) {
	input := NewScraperInput()
	
	if input == nil {
		t.Fatal("NewScraperInput returned nil")
	}
	
	if input.Country == nil || *input.Country != CountryUSA {
		t.Errorf("Expected default country to be USA, got %v", input.Country)
	}
	
	if input.DescriptionFormat == nil || *input.DescriptionFormat != FormatMarkdown {
		t.Errorf("Expected default format to be markdown, got %v", input.DescriptionFormat)
	}
	
	if input.ResultsWanted != 15 {
		t.Errorf("Expected default results to be 15, got %d", input.ResultsWanted)
	}
	
	if input.RequestTimeout != 60 {
		t.Errorf("Expected default timeout to be 60, got %d", input.RequestTimeout)
	}
}

// Helper function for tests
func testStringPtr(s string) *string {
	return &s
}