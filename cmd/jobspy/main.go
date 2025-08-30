package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/veverkap/jobspy"
)

func main() {
	// Define command line flags
	var (
		sites         = flag.String("sites", "indeed", "Comma-separated list of sites to scrape (indeed,linkedin,glassdoor,google,ziprecruiter,bayt,naukri,bdjobs)")
		searchTerm    = flag.String("search", "", "Search term for jobs")
		location      = flag.String("location", "", "Location to search for jobs")
		distance      = flag.Int("distance", 50, "Distance in miles from location")
		isRemote      = flag.Bool("remote", false, "Search for remote jobs only")
		jobType       = flag.String("jobtype", "", "Job type (fulltime, parttime, contract, etc.)")
		resultsWanted = flag.Int("results", 15, "Number of results to return")
		country       = flag.String("country", "usa", "Country for the search")
		format        = flag.String("format", "csv", "Output format (csv, json)")
		output        = flag.String("output", "", "Output file (defaults to stdout)")
		verbose       = flag.Int("verbose", 1, "Verbosity level (0=errors, 1=warnings, 2=info)")
		userAgent     = flag.String("useragent", "", "Custom user agent")
		descFormat    = flag.String("descformat", "markdown", "Description format (markdown, html, plain)")
	)
	flag.Parse()

	// Parse sites
	siteList := strings.Split(*sites, ",")
	for i, site := range siteList {
		siteList[i] = strings.TrimSpace(site)
	}

	// Create scrape input
	input := &jobspy.ScrapeJobsInput{
		SiteName:          siteList,
		SearchTerm:        stringPtr(*searchTerm),
		Location:          stringPtr(*location),
		Distance:          intPtr(*distance),
		IsRemote:          *isRemote,
		JobType:           stringPtr(*jobType),
		ResultsWanted:     intPtr(*resultsWanted),
		CountryIndeed:     *country,
		Verbose:           *verbose,
		UserAgent:        stringPtr(*userAgent),
		DescriptionFormat: stringPtr(*descFormat),
	}

	// Perform scraping
	fmt.Printf("Scraping jobs from: %s\n", strings.Join(siteList, ", "))
	if *searchTerm != "" {
		fmt.Printf("Search term: %s\n", *searchTerm)
	}
	if *location != "" {
		fmt.Printf("Location: %s\n", *location)
	}

	result, err := jobspy.ScrapeJobs(input)
	if err != nil {
		log.Fatalf("Error scraping jobs: %v", err)
	}

	fmt.Printf("Found %d jobs\n", len(result.Jobs))

	if len(result.Jobs) == 0 {
		fmt.Println("No jobs found with the specified criteria")
		return
	}

	// Output results
	var outputFile *os.File
	if *output != "" {
		outputFile, err = os.Create(*output)
		if err != nil {
			log.Fatalf("Error creating output file: %v", err)
		}
		defer outputFile.Close()
	} else {
		outputFile = os.Stdout
	}

	switch *format {
	case "csv":
		err = outputCSV(result.Jobs, outputFile)
	case "json":
		err = outputJSON(result.Jobs, outputFile)
	default:
		err = fmt.Errorf("unsupported format: %s", *format)
	}

	if err != nil {
		log.Fatalf("Error outputting results: %v", err)
	}

	if *output != "" {
		fmt.Printf("Results saved to %s\n", *output)
	}
}

func outputCSV(jobs []jobspy.JobPost, file *os.File) error {
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{
		"site", "title", "company", "location", "job_type", "compensation_interval",
		"min_amount", "max_amount", "currency", "is_remote", "job_url", "description",
	}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write data
	for _, job := range jobs {
		record := []string{
			"indeed", // For now, we only support Indeed
			job.Title,
			stringValue(job.CompanyName),
			locationString(job.Location),
			jobTypesString(job.JobType),
			compensationInterval(job.Compensation),
			compensationMinAmount(job.Compensation),
			compensationMaxAmount(job.Compensation),
			compensationCurrency(job.Compensation),
			boolString(job.IsRemote),
			job.JobURL,
			stringValue(job.Description),
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

func outputJSON(jobs []jobspy.JobPost, file *os.File) error {
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(jobs)
}

// Helper functions
func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func intPtr(i int) *int {
	return &i
}

func stringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func locationString(loc *jobspy.Location) string {
	if loc == nil {
		return ""
	}
	return loc.String()
}

func jobTypesString(jobTypes []jobspy.JobType) string {
	if len(jobTypes) == 0 {
		return ""
	}
	var strs []string
	for _, jt := range jobTypes {
		strs = append(strs, string(jt))
	}
	return strings.Join(strs, "; ")
}

func compensationInterval(comp *jobspy.Compensation) string {
	if comp == nil || comp.Interval == nil {
		return ""
	}
	return string(*comp.Interval)
}

func compensationMinAmount(comp *jobspy.Compensation) string {
	if comp == nil || comp.MinAmount == nil {
		return ""
	}
	return strconv.FormatFloat(*comp.MinAmount, 'f', 2, 64)
}

func compensationMaxAmount(comp *jobspy.Compensation) string {
	if comp == nil || comp.MaxAmount == nil {
		return ""
	}
	return strconv.FormatFloat(*comp.MaxAmount, 'f', 2, 64)
}

func compensationCurrency(comp *jobspy.Compensation) string {
	if comp == nil || comp.Currency == nil {
		return ""
	}
	return *comp.Currency
}

func boolString(b *bool) string {
	if b == nil {
		return ""
	}
	if *b {
		return "true"
	}
	return "false"
}