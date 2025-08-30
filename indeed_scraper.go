package jobspy

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// IndeedScraper implements the Scraper interface for Indeed
type IndeedScraper struct {
	*BaseScraper
	apiURL         string
	apiCountryCode string
	jobsPerPage    int
	seenURLs       map[string]bool
}

// NewIndeedScraper creates a new Indeed scraper
func NewIndeedScraper(proxies []string, userAgent string) *IndeedScraper {
	baseScraper := NewBaseScraper(SiteIndeed, "https://www.indeed.com", proxies, userAgent)
	
	return &IndeedScraper{
		BaseScraper:    baseScraper,
		apiURL:         "https://apis.indeed.com/graphql",
		jobsPerPage:    100,
		seenURLs:       make(map[string]bool),
	}
}

// Scrape implements the Scraper interface
func (is *IndeedScraper) Scrape(ctx context.Context, input *ScraperInput) (*JobResponse, error) {
	// Set up country-specific settings
	domain, countryCode := is.getCountryDomain(input.Country)
	is.apiCountryCode = countryCode
	baseURL := fmt.Sprintf("https://%s.indeed.com", domain)
	
		Logger.Info("Starting Indeed scrape", 
		"search_term", stringValue(input.SearchTerm),
		"location", stringValue(input.Location),
		"results_wanted", input.ResultsWanted)
	
	var allJobs []JobPost
	page := 1
	cursor := ""
	
	for len(allJobs) < input.ResultsWanted+input.Offset {
		Logger.Info("Scraping page", "page", page)
		
		jobs, nextCursor, err := is.scrapePage(ctx, input, cursor, baseURL)
		if err != nil {
			return nil, fmt.Errorf("error scraping page %d: %w", page, err)
		}
		
		if len(jobs) == 0 {
			Logger.Info("No more jobs found", "page", page)
			break
		}
		
		allJobs = append(allJobs, jobs...)
		cursor = nextCursor
		page++
		
		if cursor == "" {
			break
		}
	}
	
	// Apply offset and limit
	start := input.Offset
	end := input.Offset + input.ResultsWanted
	if start > len(allJobs) {
		start = len(allJobs)
	}
	if end > len(allJobs) {
		end = len(allJobs)
	}
	
	return &JobResponse{
		Jobs: allJobs[start:end],
	}, nil
}

func (is *IndeedScraper) scrapePage(ctx context.Context, input *ScraperInput, cursor string, baseURL string) ([]JobPost, string, error) {
	// Build search URL
	params := url.Values{}
	if input.SearchTerm != nil {
		params.Set("q", *input.SearchTerm)
	}
	if input.Location != nil {
		params.Set("l", *input.Location)
	}
	if input.Distance != nil {
		params.Set("radius", strconv.Itoa(*input.Distance))
	}
	if input.IsRemote {
		params.Set("remotejob", "1")
	}
	if cursor != "" {
		params.Set("start", cursor)
	}
	
	searchURL := fmt.Sprintf("%s/jobs?%s", baseURL, params.Encode())
	
	// Make request
	headers := map[string]string{
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
		"Accept-Language": "en-US,en;q=0.5",
	}
	
	resp, err := is.BaseScraper.httpClient.Get(ctx, searchURL, headers)
	if err != nil {
		return nil, "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		return nil, "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	
	// Parse HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("error parsing HTML: %w", err)
	}
	
	var jobs []JobPost
	
	// Extract job cards
	doc.Find("[data-jk]").Each(func(i int, s *goquery.Selection) {
		job := is.parseJobCard(s, baseURL, input)
		if job != nil {
			jobURL := job.JobURL
			if !is.seenURLs[jobURL] {
				is.seenURLs[jobURL] = true
				jobs = append(jobs, *job)
			}
		}
	})
	
	// Extract next page cursor
	nextCursor := ""
	nextLink := doc.Find("a[aria-label='Next Page']")
	if nextLink.Length() > 0 {
		href, exists := nextLink.Attr("href")
		if exists {
			u, err := url.Parse(href)
			if err == nil {
				nextCursor = u.Query().Get("start")
			}
		}
	}
	
	return jobs, nextCursor, nil
}

func (is *IndeedScraper) parseJobCard(s *goquery.Selection, baseURL string, input *ScraperInput) *JobPost {
	// Extract job key
	jobKey, exists := s.Attr("data-jk")
	if !exists {
		return nil
	}
	
	// Extract title
	titleEl := s.Find("h2.jobTitle a span[title]")
	if titleEl.Length() == 0 {
		titleEl = s.Find("h2.jobTitle a")
	}
	title := strings.TrimSpace(titleEl.Text())
	if title == "" {
		return nil
	}
	
	// Extract company
	companyEl := s.Find(".companyName")
	companyName := strings.TrimSpace(companyEl.Text())
	
	// Extract location
	locationEl := s.Find("[data-testid='job-location']")
	if locationEl.Length() == 0 {
		locationEl = s.Find(".companyLocation")
	}
	locationText := strings.TrimSpace(locationEl.Text())
	location := is.parseLocation(locationText, input.Country)
	
	// Build job URL
	jobURL := fmt.Sprintf("%s/viewjob?jk=%s", baseURL, jobKey)
	
	// Extract salary info
	salaryEl := s.Find(".salary-snippet")
	var compensation *Compensation
	if salaryEl.Length() > 0 {
		compensation = is.parseSalary(salaryEl.Text())
	}
	
	// Extract description snippet
	descriptionEl := s.Find(".summary")
	description := strings.TrimSpace(descriptionEl.Text())
	
	// Convert description format if needed
	if input.DescriptionFormat != nil {
		switch *input.DescriptionFormat {
		case FormatMarkdown:
			description = MarkdownConverter(description)
		case FormatPlain:
			description = PlainConverter(description)
		}
	}
	
	// Check if remote
	isRemote := is.isJobRemote(description, locationText)
	
	// Extract emails from description
	emails := ExtractEmailsFromText(description)
	
	// Extract job type
	jobTypes := is.extractJobType(description)
	
	job := &JobPost{
		ID:           &jobKey,
		Title:        title,
		CompanyName:  &companyName,
		JobURL:       jobURL,
		Location:     location,
		Description:  &description,
		Compensation: compensation,
		IsRemote:     &isRemote,
		Emails:       emails,
		JobType:      jobTypes,
	}
	
	return job
}

func (is *IndeedScraper) parseLocation(locationText string, country *Country) *Location {
	if locationText == "" {
		return nil
	}
	
	parts := strings.Split(locationText, ",")
	location := &Location{}
	
	if len(parts) >= 1 {
		city := strings.TrimSpace(parts[0])
		location.City = &city
	}
	
	if len(parts) >= 2 {
		state := strings.TrimSpace(parts[1])
		// Remove any postal codes
		state = regexp.MustCompile(`\s+\d+.*`).ReplaceAllString(state, "")
		location.State = &state
	}
	
	if country != nil {
		location.Country = country
	}
	
	return location
}

func (is *IndeedScraper) parseSalary(salaryText string) *Compensation {
	if salaryText == "" {
		return nil
	}
	
	// Extract numbers from salary text
	re := regexp.MustCompile(`\$?([\d,]+(?:\.\d{2})?)`)
	matches := re.FindAllStringSubmatch(salaryText, -1)
	
	if len(matches) == 0 {
		return nil
	}
	
	compensation := &Compensation{
		Currency: stringPtr("USD"),
	}
	
	// Parse amounts
	var amounts []float64
	for _, match := range matches {
		if len(match) > 1 {
			amountStr := strings.Replace(match[1], ",", "", -1)
			if amount, err := strconv.ParseFloat(amountStr, 64); err == nil {
				amounts = append(amounts, amount)
			}
		}
	}
	
	if len(amounts) == 1 {
		compensation.MinAmount = &amounts[0]
		compensation.MaxAmount = &amounts[0]
	} else if len(amounts) >= 2 {
		compensation.MinAmount = &amounts[0]
		compensation.MaxAmount = &amounts[1]
	}
	
	// Determine interval
	salaryLower := strings.ToLower(salaryText)
	if strings.Contains(salaryLower, "year") || strings.Contains(salaryLower, "annual") {
		interval := IntervalYearly
		compensation.Interval = &interval
	} else if strings.Contains(salaryLower, "hour") {
		interval := IntervalHourly
		compensation.Interval = &interval
	} else if strings.Contains(salaryLower, "month") {
		interval := IntervalMonthly
		compensation.Interval = &interval
	}
	
	return compensation
}

func (is *IndeedScraper) isJobRemote(description, location string) bool {
	combined := strings.ToLower(description + " " + location)
	remoteKeywords := []string{"remote", "work from home", "wfh", "telecommute", "virtual"}
	
	for _, keyword := range remoteKeywords {
		if strings.Contains(combined, keyword) {
			return true
		}
	}
	return false
}

func (is *IndeedScraper) extractJobType(description string) []JobType {
	description = strings.ToLower(description)
	var jobTypes []JobType
	
	if strings.Contains(description, "full-time") || strings.Contains(description, "fulltime") {
		jobTypes = append(jobTypes, JobTypeFullTime)
	}
	if strings.Contains(description, "part-time") || strings.Contains(description, "parttime") {
		jobTypes = append(jobTypes, JobTypePartTime)
	}
	if strings.Contains(description, "contract") {
		jobTypes = append(jobTypes, JobTypeContract)
	}
	if strings.Contains(description, "intern") {
		jobTypes = append(jobTypes, JobTypeInternship)
	}
	if strings.Contains(description, "temporary") {
		jobTypes = append(jobTypes, JobTypeTemporary)
	}
	
	return jobTypes
}

func (is *IndeedScraper) getCountryDomain(country *Country) (string, string) {
	if country == nil {
		return "www", "us"
	}
	
	switch *country {
	case CountryUSA:
		return "www", "us"
	case CountryCanada:
		return "ca", "ca"
	case CountryUK:
		return "uk", "gb"
	case CountryGermany:
		return "de", "de"
	case CountryFrance:
		return "fr", "fr"
	default:
		return "www", "us"
	}
}

// Helper functions
func stringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func stringPtr(s string) *string {
	return &s
}