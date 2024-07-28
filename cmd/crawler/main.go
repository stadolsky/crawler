package main

import (
	"crawler/internal/pkg/crawler"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/url"
	"os"
)

type cliArgs struct {
	sourceURL    string
	outputFormat string
	limit        int
	logLevel     string
}

type format string

const (
	JSON format = "json"
	CSV  format = "csv"
)

func parseFormat(s string) (format, error) {
	switch s {
	case "json":
		return JSON, nil
	case "csv":
		return CSV, nil
	}

	return "", fmt.Errorf("unknown format: %s", s)
}

func main() {
	args := cliArgs{}

	flag.StringVar(&args.sourceURL, "url", "", "Initial sourceURL URL to parse")
	flag.StringVar(&args.outputFormat, "format", "json", "Output format for results")
	flag.IntVar(&args.limit, "limit", 100, "Limit the number of links to parse (0 - no limit)")
	flag.StringVar(&args.logLevel, "loglevel", zerolog.DebugLevel.String(), "log level (debug, info, warn, error, ...)")

	flag.Parse()

	logLevel := zerolog.DebugLevel
	if level, err := zerolog.ParseLevel(args.logLevel); err != nil {
		logLevel = level
	}

	zerolog.SetGlobalLevel(logLevel)

	frmt, err := parseFormat(args.outputFormat)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse output format")
	}

	urlObj, err := url.Parse(args.sourceURL)
	if err != nil {
		log.Fatal().Err(err).Str("source url", args.sourceURL).Msg("invalid source URL")
	}

	crwl := crawler.NewCrawler(urlObj, args.limit)

	links := crwl.ScrapeAllDomainLinks()

	writeResults(links, frmt)
}

func writeResults(links []string, frmt format) {
	switch frmt {
	case JSON:
		b, err := json.Marshal(links)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to marshal links")
		}
		fmt.Println(string(b))
	case CSV:
		wr := csv.NewWriter(os.Stdout)
		err := wr.Write(links)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to write links")
		}
		wr.Flush()
	}
}
