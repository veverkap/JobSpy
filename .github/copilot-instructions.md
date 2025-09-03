# JobSpy - Job Scraping Library

JobSpy is a Python library for scraping job postings from multiple job boards (LinkedIn, Indeed, Glassdoor, Google, ZipRecruiter, Bayt, Naukri, BDJobs) concurrently and aggregating results in a pandas DataFrame.

Always reference these instructions first and fallback to search or bash commands only when you encounter unexpected information that does not match the info here.

## Working Effectively

### Environment Setup and Dependencies
- Install Poetry package manager: `pip install poetry`
- Install all dependencies: `poetry install` -- NEVER CANCEL: takes 5 minutes to complete. Set timeout to 10+ minutes.
- Activate virtual environment: `poetry shell` (optional, `poetry run` works too)

### Development Workflow
- Format code: `poetry run black jobspy/` -- takes 5-10 seconds
- Check formatting: `poetry run black --check jobspy/` -- takes 2-3 seconds  
- Build package: `poetry build` -- takes 1 second, creates wheel and sdist in `dist/`
- Run basic functionality test: `poetry run python -c "from jobspy import scrape_jobs; print('Import successful')"` -- takes 1 second

### Testing and Validation
- Test core imports and configuration (recommended after changes): `poetry run python /tmp/test_jobspy.py` -- takes ~30 seconds if using full test suite
- Actual job scraping will fail in sandboxed environments due to network restrictions -- this is EXPECTED behavior
- Network requests typically fail quickly (~1 second) with DNS resolution errors -- this is NORMAL in restricted environments

### Package Structure Navigation
- `/jobspy/` -- Main package directory containing all scrapers
- `/jobspy/__init__.py` -- Main scrape_jobs() function and public API
- `/jobspy/model.py` -- Core data models (JobPost, Site enum, JobType enum, etc.)
- `/jobspy/util.py` -- Shared utilities for logging, session management, etc.
- `/jobspy/{site}/` -- Individual scraper implementations (indeed, linkedin, google, etc.)
- `pyproject.toml` -- Poetry configuration with dependencies and metadata
- `.pre-commit-config.yaml` -- Code formatting hooks (may fail in restricted environments)

## Validation Scenarios

### Basic Functionality Validation
Always run these after making changes to verify the package works:

1. **Import Test**: `poetry run python -c "from jobspy import scrape_jobs, Indeed, LinkedIn, JobType, Site; print('All imports successful')"`
2. **Configuration Test**: Test parameter validation with a minimal scrape_jobs() call (expect network failure)
3. **Enum Test**: Verify JobType and Site enums are accessible and contain expected values

### Example Working Code Validation
```python
from jobspy import scrape_jobs

# This will validate configuration but fail at network level (expected)
jobs = scrape_jobs(
    site_name=['indeed'],
    search_term='python developer', 
    location='San Francisco, CA',
    results_wanted=5,
    hours_old=24
)
```

Expected error: "Max retries exceeded" or "Failed to resolve hostname" - this is NORMAL in restricted environments.

### Complete Manual Test
Use this comprehensive test script to validate all functionality:

```python
# Test all major components
from jobspy import scrape_jobs, JobType, Site, Indeed, LinkedIn, ZipRecruiter

# Test enums
print("JobType values:", [jt.value[0] for jt in [JobType.FULL_TIME, JobType.PART_TIME, JobType.CONTRACT, JobType.INTERNSHIP]])
print("Site values:", [s.value for s in [Site.INDEED, Site.LINKEDIN, Site.ZIP_RECRUITER, Site.GOOGLE]])

# Test configuration (will fail at network, this is expected)
try:
    jobs = scrape_jobs(site_name=['indeed'], search_term='test', results_wanted=1)
except Exception as e:
    if "address associated with hostname" in str(e):
        print("✓ Configuration valid, network failed as expected")
    else:
        print(f"✗ Unexpected error: {e}")
```

## Timing Expectations and Warnings

### CRITICAL - NEVER CANCEL These Operations:
- `poetry install` -- **NEVER CANCEL**: Takes 5 minutes. Set timeout to 10+ minutes.
- Network operations (when functional) -- **NEVER CANCEL**: Job scraping can take 1-5 minutes per site
- Large job scraping requests -- **NEVER CANCEL**: Can take 10+ minutes for hundreds of jobs

### Quick Operations (under 10 seconds):
- `poetry run black jobspy/` -- 5-10 seconds
- `poetry build` -- 1 second  
- Basic imports and enum tests -- 1 second

### Expected Network Timeouts:
- In restricted environments, network requests fail quickly (~1 second) with DNS errors
- This is EXPECTED and NORMAL behavior - do not attempt to fix network issues
- Focus on validating code structure, imports, and configuration instead

## Common Tasks

### Making Code Changes
1. Always format code after changes: `poetry run black jobspy/`
2. Test imports: `poetry run python -c "from jobspy import scrape_jobs"`
3. Run validation script to test core functionality
4. Build package to ensure no syntax errors: `poetry build`

### Working with Different Job Sites
- Each job site has its own scraper in `/jobspy/{site}/` directory
- Main scrapers: `indeed`, `linkedin`, `ziprecruiter`, `google`, `glassdoor`, `bayt`, `naukri`, `bdjobs`
- All scrapers inherit from base `Scraper` class in `model.py`
- Site-specific parameters documented in main `scrape_jobs()` function

### Key Parameters for scrape_jobs()
- `site_name`: List of job sites (default: all available)
- `search_term`: Job search query
- `location`: Geographic location for search
- `results_wanted`: Number of jobs to retrieve (default: 15)
- `hours_old`: Filter by posting age
- `job_type`: Filter by employment type (fulltime, parttime, contract, internship)
- `is_remote`: Filter for remote jobs
- `proxies`: List of proxy servers for rate limiting avoidance

### Understanding Limitations
- Network access required for actual job scraping
- Most job sites have rate limiting (LinkedIn is most restrictive)
- Job sites limit results to ~1000 jobs per search
- Some parameters are mutually exclusive (documented in README.md)

## Important Development Notes

### Code Quality
- Always run `poetry run black jobspy/` before committing
- Code follows Black formatting with 88 character line length
- No existing test infrastructure - rely on manual validation scripts

### Package Management
- Uses Poetry for dependency management
- Python 3.10+ required
- Main dependencies: requests, beautifulsoup4, pandas, pydantic, tls-client

### Network Dependencies
- All functionality requires external network access to job sites
- In restricted environments, focus on code structure and import validation
- Network timeouts are expected and should not be treated as code errors

### Troubleshooting Common Issues
- "Max retries exceeded" errors: Normal in restricted environments
- "Failed to resolve hostname" errors: Expected when network access blocked  
- Import errors: Check `poetry install` completed successfully
- Formatting errors: Run `poetry run black jobspy/` to fix
- Build errors: Usually syntax issues, check Python syntax in changed files

## Repository Structure Overview

```
/jobspy/                 # Main package
├── __init__.py         # Public API and scrape_jobs() function  
├── model.py            # Data models and base classes
├── util.py             # Shared utilities
├── exception.py        # Custom exceptions
├── indeed/             # Indeed scraper
├── linkedin/           # LinkedIn scraper  
├── ziprecruiter/       # ZipRecruiter scraper
├── google/             # Google Jobs scraper
├── glassdoor/          # Glassdoor scraper
├── bayt/               # Bayt scraper
├── naukri/             # Naukri scraper
└── bdjobs/             # BDJobs scraper

pyproject.toml          # Poetry configuration
.pre-commit-config.yaml # Code formatting hooks  
.github/workflows/      # CI/CD for PyPI publishing
README.md               # User documentation
```

This structure allows for easy navigation and understanding of the codebase. Each scraper is self-contained but follows common patterns defined in the base model.