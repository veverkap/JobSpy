# JobSpy Examples

This document provides examples of how to use JobSpy in different scenarios.

## Basic Usage

### Simple Job Search

```bash
# Search for software engineer jobs in San Francisco
jobspy -search "software engineer" -location "San Francisco, CA"

# Search for remote data scientist positions
jobspy -search "data scientist" -remote -results 25

# Search multiple job boards
jobspy -sites "indeed,linkedin" -search "product manager" -location "New York, NY"
```

### Advanced Filtering

```bash
# Full-time positions only
jobspy -search "backend engineer" -jobtype "fulltime" -location "Seattle, WA"

# Within 10 miles of location
jobspy -search "frontend developer" -location "Austin, TX" -distance 10

# In specific country
jobspy -search "golang developer" -country "canada" -location "Toronto"
```

### Output Options

```bash
# JSON output to file
jobspy -search "devops engineer" -format json -output jobs.json

# CSV with custom description format
jobspy -search "qa engineer" -descformat "plain" -format csv -output jobs.csv

# Verbose logging
jobspy -search "mobile developer" -verbose 2
```

## Go Library Usage

### Basic Search

```go
package main

import (
    "fmt"
    "log"

    "github.com/veverkap/jobspy"
)

func main() {
    input := &jobspy.ScrapeJobsInput{
        SiteName:      []string{"indeed"},
        SearchTerm:    stringPtr("golang developer"),
        Location:      stringPtr("San Francisco, CA"),
        ResultsWanted: intPtr(20),
        CountryIndeed: "usa",
        Verbose:       1,
    }

    result, err := jobspy.ScrapeJobs(input)
    if err != nil {
        log.Fatalf("Error: %v", err)
    }

    fmt.Printf("Found %d jobs\n", len(result.Jobs))
    for _, job := range result.Jobs {
        fmt.Printf("• %s at %s\n", job.Title, stringValue(job.CompanyName))
    }
}
```

### Remote Jobs Only

```go
input := &jobspy.ScrapeJobsInput{
    SiteName:      []string{"indeed", "linkedin"},
    SearchTerm:    stringPtr("remote software engineer"),
    IsRemote:      true,
    ResultsWanted: intPtr(50),
    CountryIndeed: "usa",
}
```

### With Salary Information

```go
result, err := jobspy.ScrapeJobs(input)
if err != nil {
    log.Fatal(err)
}

for _, job := range result.Jobs {
    fmt.Printf("Job: %s\n", job.Title)
    if job.Compensation != nil {
        if job.Compensation.MinAmount != nil && job.Compensation.MaxAmount != nil {
            fmt.Printf("Salary: $%.0f - $%.0f %s\n", 
                *job.Compensation.MinAmount, 
                *job.Compensation.MaxAmount,
                stringValue(job.Compensation.Currency))
        }
    }
    fmt.Println()
}
```

### Custom User Agent and Proxies

```go
input := &jobspy.ScrapeJobsInput{
    SiteName:   []string{"indeed"},
    SearchTerm: stringPtr("python developer"),
    Location:   stringPtr("Boston, MA"),
    UserAgent:  stringPtr("MyJobBot/1.0"),
    Proxies: []string{
        "http://proxy1.example.com:8080",
        "http://proxy2.example.com:8080",
    },
}
```

### Processing Job Descriptions

```go
for _, job := range result.Jobs {
    if job.Description != nil {
        // Extract emails from description
        if len(job.Emails) > 0 {
            fmt.Printf("Contact emails: %v\n", job.Emails)
        }
        
        // Check for specific technologies
        desc := strings.ToLower(*job.Description)
        if strings.Contains(desc, "kubernetes") {
            fmt.Printf("Kubernetes job: %s\n", job.Title)
        }
    }
}
```

### Export to CSV

```go
package main

import (
    "encoding/csv"
    "os"
    "strconv"

    "github.com/veverkap/jobspy"
)

func exportToCSV(jobs []jobspy.JobPost, filename string) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    // Header
    header := []string{"Title", "Company", "Location", "URL", "Salary Min", "Salary Max"}
    writer.Write(header)

    // Data
    for _, job := range jobs {
        record := []string{
            job.Title,
            stringValue(job.CompanyName),
            job.Location.String(),
            job.JobURL,
            floatToString(job.Compensation, true),  // min
            floatToString(job.Compensation, false), // max
        }
        writer.Write(record)
    }

    return nil
}

func floatToString(comp *jobspy.Compensation, isMin bool) string {
    if comp == nil {
        return ""
    }
    var amount *float64
    if isMin {
        amount = comp.MinAmount
    } else {
        amount = comp.MaxAmount
    }
    if amount == nil {
        return ""
    }
    return strconv.FormatFloat(*amount, 'f', 0, 64)
}
```

### Job Type Filtering

```go
// Filter by job type after scraping
var fullTimeJobs []jobspy.JobPost
for _, job := range result.Jobs {
    for _, jt := range job.JobType {
        if jt == jobspy.JobTypeFullTime {
            fullTimeJobs = append(fullTimeJobs, job)
            break
        }
    }
}
```

### Location-based Analysis

```go
// Count jobs by state
stateCount := make(map[string]int)
for _, job := range result.Jobs {
    if job.Location != nil && job.Location.State != nil {
        stateCount[*job.Location.State]++
    }
}

for state, count := range stateCount {
    fmt.Printf("%s: %d jobs\n", state, count)
}
```

## Error Handling

```go
result, err := jobspy.ScrapeJobs(input)
if err != nil {
    switch {
    case strings.Contains(err.Error(), "unknown site"):
        log.Fatal("Invalid job board specified")
    case strings.Contains(err.Error(), "context deadline exceeded"):
        log.Fatal("Request timed out - try reducing results or using proxies")
    default:
        log.Fatalf("Scraping failed: %v", err)
    }
}

if len(result.Jobs) == 0 {
    log.Println("No jobs found with the specified criteria")
}
```

## Utility Functions

```go
// Helper functions for working with pointers
func stringPtr(s string) *string {
    if s == "" {
        return nil
    }
    return &s
}

func intPtr(i int) *int {
    return &i
}

func boolPtr(b bool) *bool {
    return &b
}

func stringValue(s *string) string {
    if s == nil {
        return ""
    }
    return *s
}

func intValue(i *int) int {
    if i == nil {
        return 0
    }
    return *i
}

func boolValue(b *bool) bool {
    if b == nil {
        return false
    }
    return *b
}
```