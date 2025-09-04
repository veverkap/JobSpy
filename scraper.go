package jobspy

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Scraper is the interface that all job site scrapers must implement
type Scraper interface {
	Scrape(ctx context.Context, input *ScraperInput) (*JobResponse, error)
	GetSite() Site
}

// BaseScraper provides common functionality for all scrapers
type BaseScraper struct {
	site       Site
	httpClient *HTTPClient
	baseURL    string
}

// NewBaseScraper creates a new base scraper
func NewBaseScraper(site Site, baseURL string, proxies []string, userAgent string) *BaseScraper {
	return &BaseScraper{
		site:       site,
		httpClient: NewHTTPClient(proxies, 60*time.Second, userAgent),
		baseURL:    baseURL,
	}
}

// GetSite returns the site this scraper handles
func (bs *BaseScraper) GetSite() Site {
	return bs.site
}

// ScrapeJobs is the main function to scrape jobs from multiple sites concurrently
func ScrapeJobs(input *ScrapeJobsInput) (*JobResponse, error) {
	ctx := context.Background()
	
	// Set up logging
	SetLogLevel(input.Verbose)
	
	// Convert site names to Site enums
	sites, err := parseSites(input.SiteName)
	if err != nil {
		return nil, fmt.Errorf("error parsing sites: %w", err)
	}
	
	// Set up defaults
	if input.ResultsWanted == nil {
		defaultResults := 15
		input.ResultsWanted = &defaultResults
	}
	if input.Offset == nil {
		defaultOffset := 0
		input.Offset = &defaultOffset
	}
	if input.LinkedInFetchDescription == nil {
		defaultLinkedIn := false
		input.LinkedInFetchDescription = &defaultLinkedIn
	}

	// Create scraper input
	scraperInput := &ScraperInput{
		SiteType:                 sites,
		SearchTerm:               input.SearchTerm,
		GoogleSearchTerm:         input.GoogleSearchTerm,
		Location:                 input.Location,
		Country:                  CountryFromString(input.CountryIndeed),
		Distance:                 input.Distance,
		IsRemote:                 input.IsRemote,
		JobType:                  JobTypeFromString(getValue(input.JobType)),
		EasyApply:                input.EasyApply,
		Offset:                   *input.Offset,
		LinkedInFetchDescription: *input.LinkedInFetchDescription,
		LinkedInCompanyIDs:       input.LinkedInCompanyIDs,
		RequestTimeout:           60,
		ResultsWanted:            *input.ResultsWanted,
		HoursOld:                 input.HoursOld,
	}
	
	// Set description format
	if input.DescriptionFormat != nil {
		switch *input.DescriptionFormat {
		case "markdown":
			format := FormatMarkdown
			scraperInput.DescriptionFormat = &format
		case "html":
			format := FormatHTML
			scraperInput.DescriptionFormat = &format
		case "plain":
			format := FormatPlain
			scraperInput.DescriptionFormat = &format
		}
	}
	
	// Create scrapers map
	scrapers := map[Site]Scraper{
		SiteIndeed: NewIndeedScraper(input.Proxies, getValue(input.UserAgent)),
		// TODO: Add other scrapers
	}
	
	// Scrape from each site concurrently
	var wg sync.WaitGroup
	resultsChan := make(chan *JobResponse, len(sites))
	errorsChan := make(chan error, len(sites))
	
	for _, site := range sites {
		scraper, ok := scrapers[site]
		if !ok {
			Logger.Warn("Scraper not implemented for site", "site", site)
			continue
		}
		
		wg.Add(1)
		go func(s Scraper) {
			defer wg.Done()
			
			Logger.Info("Starting scrape", "site", s.GetSite())
			result, err := s.Scrape(ctx, scraperInput)
			if err != nil {
				Logger.Error("Scrape failed", "site", s.GetSite(), "error", err)
				errorsChan <- err
				return
			}
			
			Logger.Info("Scrape completed", "site", s.GetSite(), "jobs_found", len(result.Jobs))
			resultsChan <- result
		}(scraper)
	}
	
	wg.Wait()
	close(resultsChan)
	close(errorsChan)
	
	// Collect results
	var allJobs []JobPost
	for result := range resultsChan {
		allJobs = append(allJobs, result.Jobs...)
	}
	
	// Log any errors (but don't fail the entire operation)
	for err := range errorsChan {
		Logger.Error("Scraping error", "error", err)
	}
	
	return &JobResponse{Jobs: allJobs}, nil
}

// ScrapeJobsInput contains all the input parameters for the ScrapeJobs function
type ScrapeJobsInput struct {
	SiteName                  interface{} // string, []string, Site, or []Site
	SearchTerm                *string
	GoogleSearchTerm          *string
	Location                  *string
	Distance                  *int
	IsRemote                  bool
	JobType                   *string
	EasyApply                 *bool
	ResultsWanted             *int
	CountryIndeed             string
	Proxies                   []string
	CaCert                    *string
	DescriptionFormat         *string
	LinkedInFetchDescription  *bool
	LinkedInCompanyIDs        []int
	Offset                    *int
	HoursOld                  *int
	EnforceAnnualSalary       bool
	Verbose                   int
	UserAgent                 *string
}

// Helper functions
func parseSites(siteNames interface{}) ([]Site, error) {
	switch v := siteNames.(type) {
	case nil:
		return []Site{SiteIndeed, SiteLinkedIn, SiteGlassdoor, SiteGoogle}, nil
	case string:
		site, err := MapStringToSite(v)
		if err != nil {
			return nil, err
		}
		return []Site{site}, nil
	case []string:
		var sites []Site
		for _, s := range v {
			site, err := MapStringToSite(s)
			if err != nil {
				return nil, err
			}
			sites = append(sites, site)
		}
		return sites, nil
	case Site:
		return []Site{v}, nil
	case []Site:
		return v, nil
	default:
		return nil, fmt.Errorf("unsupported site name type: %T", v)
	}
}

func getValue[T any](ptr *T) T {
	var zero T
	if ptr == nil {
		return zero
	}
	return *ptr
}