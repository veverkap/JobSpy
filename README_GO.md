# JobSpy - Go

JobSpy is a Go library and CLI tool for scraping job postings from major job boards including Indeed, LinkedIn, Glassdoor, Google Jobs, ZipRecruiter, and more.

## Features

- **Multiple Job Boards**: Support for Indeed, LinkedIn, Glassdoor, Google Jobs, ZipRecruiter, Bayt, Naukri, and BDJobs
- **Concurrent Scraping**: Scrape multiple job boards simultaneously using goroutines
- **Flexible Search**: Search by keywords, location, job type, and more
- **Multiple Output Formats**: Export results as CSV or JSON
- **Proxy Support**: Built-in proxy rotation for avoiding rate limits
- **Type Safety**: Full Go type safety with structured data models
- **CLI Interface**: Easy-to-use command-line tool

## Installation

### Go Module

```bash
go get github.com/veverkap/jobspy
```

### CLI Tool

```bash
go install github.com/veverkap/jobspy/cmd/jobspy@latest
```

Or build from source:

```bash
git clone https://github.com/veverkap/jobspy.git
cd jobspy
go build -o jobspy ./cmd/jobspy
```

## Usage

### Command Line Interface

```bash
# Basic usage
jobspy -search "software engineer" -location "San Francisco, CA" -results 20

# Specify job board
jobspy -sites "indeed,linkedin" -search "data scientist" -location "New York, NY"

# Remote jobs only
jobspy -search "golang developer" -remote -results 50

# Export to file
jobspy -search "product manager" -location "Seattle, WA" -format json -output jobs.json

# Full-time jobs only
jobspy -search "backend engineer" -jobtype "fulltime" -location "Austin, TX"
```

#### CLI Options

```
  -country string
        Country for the search (default "usa")
  -descformat string
        Description format (markdown, html, plain) (default "markdown")
  -distance int
        Distance in miles from location (default 50)
  -format string
        Output format (csv, json) (default "csv")
  -jobtype string
        Job type (fulltime, parttime, contract, etc.)
  -location string
        Location to search for jobs
  -output string
        Output file (defaults to stdout)
  -remote
        Search for remote jobs only
  -results int
        Number of results to return (default 15)
  -search string
        Search term for jobs
  -sites string
        Comma-separated list of sites to scrape (default "indeed")
  -useragent string
        Custom user agent
  -verbose int
        Verbosity level (0=errors, 1=warnings, 2=info) (default 1)
```

### Go Library

```go
package main

import (
    "fmt"
    "log"

    "github.com/veverkap/jobspy"
)

func main() {
    // Create search input
    input := &jobspy.ScrapeJobsInput{
        SiteName:      []string{"indeed", "linkedin"},
        SearchTerm:    stringPtr("software engineer"),
        Location:      stringPtr("San Francisco, CA"),
        Distance:      intPtr(25),
        IsRemote:      false,
        ResultsWanted: intPtr(50),
        CountryIndeed: "usa",
        Verbose:       1,
    }

    // Scrape jobs
    result, err := jobspy.ScrapeJobs(input)
    if err != nil {
        log.Fatalf("Error scraping jobs: %v", err)
    }

    // Process results
    fmt.Printf("Found %d jobs\n", len(result.Jobs))
    for _, job := range result.Jobs {
        fmt.Printf("Title: %s\nCompany: %s\nLocation: %s\nURL: %s\n\n",
            job.Title,
            stringValue(job.CompanyName),
            job.Location.String(),
            job.JobURL)
    }
}

func stringPtr(s string) *string {
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
```

## Data Models

### JobPost

```go
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
    // ... additional fields
}
```

### Location

```go
type Location struct {
    City    *string  `json:"city,omitempty"`
    State   *string  `json:"state,omitempty"`
    Country *Country `json:"country,omitempty"`
}
```

### Compensation

```go
type Compensation struct {
    Interval  *CompensationInterval `json:"interval,omitempty"`
    MinAmount *float64              `json:"min_amount,omitempty"`
    MaxAmount *float64              `json:"max_amount,omitempty"`
    Currency  *string               `json:"currency,omitempty"`
}
```

## Supported Job Boards

- **Indeed** ✅ (Fully implemented)
- **LinkedIn** 🚧 (In development)
- **Glassdoor** 🚧 (In development) 
- **Google Jobs** 🚧 (In development)
- **ZipRecruiter** 🚧 (In development)
- **Bayt** 🚧 (In development)
- **Naukri** 🚧 (In development)
- **BDJobs** 🚧 (In development)

## Proxy Support

JobSpy supports proxy rotation to avoid rate limiting:

```go
input := &jobspy.ScrapeJobsInput{
    // ... other fields
    Proxies: []string{
        "http://proxy1.example.com:8080",
        "http://proxy2.example.com:8080",
        "socks5://proxy3.example.com:1080",
    },
}
```

## Country Support

Currently supports job searches in:
- United States (usa)
- Canada (canada) 
- United Kingdom (uk)
- Germany (germany)
- France (france)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Migration from Python

This is a Go port of the original Python JobSpy library. The Go version provides:

- Better performance through compiled code and concurrency
- Type safety and better error handling
- Smaller memory footprint
- Easy deployment as a single binary
- Modern CLI interface

### Key Differences from Python Version

- Uses Go structs instead of Pydantic models
- Native Go HTTP client instead of requests/tls-client
- Structured logging with slog instead of Python logging
- CLI built with Go's flag package
- Direct CSV/JSON output instead of pandas DataFrames