package jobspy

import (
	"strings"
	"time"
)

// JobType represents the type of employment
type JobType string

const (
	JobTypeFullTime   JobType = "fulltime"
	JobTypePartTime   JobType = "parttime"
	JobTypeContract   JobType = "contract"
	JobTypeTemporary  JobType = "temporary"
	JobTypeInternship JobType = "internship"
	JobTypePerDiem    JobType = "perdiem"
	JobTypeNights     JobType = "nights"
	JobTypeOther      JobType = "other"
	JobTypeSummer     JobType = "summer"
	JobTypeVolunteer  JobType = "volunteer"
)

// JobTypeFromString attempts to match a string to a JobType
func JobTypeFromString(s string) *JobType {
	s = strings.ToLower(s)
	switch {
	case strings.Contains(s, "fulltime") || strings.Contains(s, "full-time") || strings.Contains(s, "períodointegral"):
		jt := JobTypeFullTime
		return &jt
	case strings.Contains(s, "parttime") || strings.Contains(s, "part-time") || strings.Contains(s, "teilzeit"):
		jt := JobTypePartTime
		return &jt
	case strings.Contains(s, "contract") || strings.Contains(s, "contractor"):
		jt := JobTypeContract
		return &jt
	case strings.Contains(s, "temporary"):
		jt := JobTypeTemporary
		return &jt
	case strings.Contains(s, "internship") || strings.Contains(s, "prácticas"):
		jt := JobTypeInternship
		return &jt
	}
	return nil
}

// Site represents the job board site
type Site string

const (
	SiteLinkedIn     Site = "linkedin"
	SiteIndeed       Site = "indeed"
	SiteZipRecruiter Site = "zip_recruiter"
	SiteGlassdoor    Site = "glassdoor"
	SiteGoogle       Site = "google"
	SiteBayt         Site = "bayt"
	SiteNaukri       Site = "naukri"
	SiteBDJobs       Site = "bdjobs"
)

// Country represents the country for job searches
type Country string

const (
	CountryUSA     Country = "usa"
	CountryCanada  Country = "canada"
	CountryUK      Country = "uk"
	CountryGermany Country = "germany"
	CountryFrance  Country = "france"
	// Add more countries as needed
)

// CompensationInterval represents how often compensation is paid
type CompensationInterval string

const (
	IntervalYearly  CompensationInterval = "yearly"
	IntervalMonthly CompensationInterval = "monthly"
	IntervalWeekly  CompensationInterval = "weekly"
	IntervalDaily   CompensationInterval = "daily"
	IntervalHourly  CompensationInterval = "hourly"
)

// Location represents a geographic location
type Location struct {
	City    *string  `json:"city,omitempty"`
	State   *string  `json:"state,omitempty"`
	Country *Country `json:"country,omitempty"`
}

// String returns a formatted location string
func (l *Location) String() string {
	var parts []string
	if l.City != nil {
		parts = append(parts, *l.City)
	}
	if l.State != nil {
		parts = append(parts, *l.State)
	}
	if l.Country != nil {
		countryStr := string(*l.Country)
		if countryStr == "usa" || countryStr == "uk" {
			countryStr = strings.ToUpper(countryStr)
		} else {
			countryStr = strings.Title(countryStr)
		}
		parts = append(parts, countryStr)
	}
	return strings.Join(parts, ", ")
}

// Compensation represents salary/wage information
type Compensation struct {
	Interval  *CompensationInterval `json:"interval,omitempty"`
	MinAmount *float64              `json:"min_amount,omitempty"`
	MaxAmount *float64              `json:"max_amount,omitempty"`
	Currency  *string               `json:"currency,omitempty"`
}

// DescriptionFormat represents the format of job descriptions
type DescriptionFormat string

const (
	FormatMarkdown DescriptionFormat = "markdown"
	FormatHTML     DescriptionFormat = "html"
	FormatPlain    DescriptionFormat = "plain"
)

// JobPost represents a single job posting
type JobPost struct {
	ID           *string       `json:"id,omitempty"`
	Title        string        `json:"title"`
	CompanyName  *string       `json:"company_name,omitempty"`
	JobURL       string        `json:"job_url"`
	JobURLDirect *string       `json:"job_url_direct,omitempty"`
	Location     *Location     `json:"location,omitempty"`
	Description  *string       `json:"description,omitempty"`
	CompanyURL   *string       `json:"company_url,omitempty"`
	JobType      []JobType     `json:"job_type,omitempty"`
	Compensation *Compensation `json:"compensation,omitempty"`
	DatePosted   *time.Time    `json:"date_posted,omitempty"`
	Emails       []string      `json:"emails,omitempty"`
	IsRemote     *bool         `json:"is_remote,omitempty"`
	ListingType  *string       `json:"listing_type,omitempty"`

	// Site-specific fields
	JobLevel            *string  `json:"job_level,omitempty"`
	CompanyIndustry     *string  `json:"company_industry,omitempty"`
	CompanyAddresses    *string  `json:"company_addresses,omitempty"`
	CompanyNumEmployees *string  `json:"company_num_employees,omitempty"`
	CompanyRevenue      *string  `json:"company_revenue,omitempty"`
	CompanyDescription  *string  `json:"company_description,omitempty"`
	CompanyLogo         *string  `json:"company_logo,omitempty"`
	BannerPhotoURL      *string  `json:"banner_photo_url,omitempty"`
	JobFunction         *string  `json:"job_function,omitempty"`
	Skills              []string `json:"skills,omitempty"`
	ExperienceRange     *string  `json:"experience_range,omitempty"`
	CompanyRating       *float64 `json:"company_rating,omitempty"`
	CompanyReviewsCount *int     `json:"company_reviews_count,omitempty"`
	VacancyCount        *int     `json:"vacancy_count,omitempty"`
	WorkFromHomeType    *string  `json:"work_from_home_type,omitempty"`
}

// JobResponse contains the results of a job scraping operation
type JobResponse struct {
	Jobs []JobPost `json:"jobs"`
}

// ScraperInput contains all the parameters for a job scraping operation
type ScraperInput struct {
	SiteType                  []Site             `json:"site_type"`
	SearchTerm                *string            `json:"search_term,omitempty"`
	GoogleSearchTerm          *string            `json:"google_search_term,omitempty"`
	Location                  *string            `json:"location,omitempty"`
	Country                   *Country           `json:"country,omitempty"`
	Distance                  *int               `json:"distance,omitempty"`
	IsRemote                  bool               `json:"is_remote"`
	JobType                   *JobType           `json:"job_type,omitempty"`
	EasyApply                 *bool              `json:"easy_apply,omitempty"`
	Offset                    int                `json:"offset"`
	LinkedInFetchDescription  bool               `json:"linkedin_fetch_description"`
	LinkedInCompanyIDs        []int              `json:"linkedin_company_ids,omitempty"`
	DescriptionFormat         *DescriptionFormat `json:"description_format,omitempty"`
	RequestTimeout            int                `json:"request_timeout"`
	ResultsWanted             int                `json:"results_wanted"`
	HoursOld                  *int               `json:"hours_old,omitempty"`
}

// NewScraperInput creates a new ScraperInput with default values
func NewScraperInput() *ScraperInput {
	defaultFormat := FormatMarkdown
	defaultCountry := CountryUSA
	return &ScraperInput{
		SiteType:          []Site{},
		Country:           &defaultCountry,
		IsRemote:          false,
		Offset:            0,
		LinkedInFetchDescription: false,
		DescriptionFormat: &defaultFormat,
		RequestTimeout:    60,
		ResultsWanted:     15,
	}
}