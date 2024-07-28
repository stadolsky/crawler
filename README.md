# Crawler

Crawler is a simple web scraping tool written in Go. It fetches and parses web pages by pointed URL to extract links which belongs to the host URL.

## Quick start

Requirements

- Go 1.22.3 or higher

Clone the repository and run app:

```bash
git clone https://github.com/stadolsky/crawler.git
cd crawler
make run url="YOUR_URL"
```

Note: By default, the application will crawl pointed website and output the results in CSV format, up to a limit of 10 links. The log level is set to 0 (debug).

## Usage
You can customize the behavior using the following command-line arguments:  
- url: The initial URL to parse.
- format: The output format for results. Can be 'json' or 'csv'. Default is 'json'.
- limit: The maximum number of links to parse. Set to 0 for no limit. Default is 100.
- loglevel: The log level. Can be 'debug', 'info', 'warn', 'error', etc. Default is 'debug'.