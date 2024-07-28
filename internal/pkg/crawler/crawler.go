package crawler

import (
	"github.com/rs/zerolog/log"
	"golang.org/x/net/html"
	"net/http"
	"net/url"
	"sync"
)

type Crawler struct {
	links map[string]struct{}

	sourceURL *url.URL
	limit     int

	mu sync.RWMutex

	wg sync.WaitGroup
}

func NewCrawler(sourceURL *url.URL, limit int) *Crawler {
	return &Crawler{
		links:     make(map[string]struct{}),
		sourceURL: sourceURL,
		limit:     limit,
	}
}

func (c *Crawler) ScrapeAllDomainLinks() []string {
	c.storeLink(c.sourceURL.String())

	// Initial start
	c.wg.Add(1)
	c.scrapePage(c.sourceURL.String(), c.sourceURL.Hostname(), c.limit)

	c.wg.Wait()

	return c.getLinks()
}

func (c *Crawler) storeLink(link string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.links[link] = struct{}{}
}

func (c *Crawler) getLinks() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	links := make([]string, 0, len(c.links))
	for link := range c.links {
		links = append(links, link)
	}

	return links
}

func (c *Crawler) isLinkExist(link string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, ok := c.links[link]
	return ok
}

func (c *Crawler) isLimitExceeded() bool {
	if c.limit == 0 {
		return false
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.links) >= c.limit
}

func (c *Crawler) scrapePage(pageURL, domain string, limit int) {
	defer c.wg.Done()

	if c.isLimitExceeded() {
		return
	}

	log.Debug().Str("url", pageURL).Msg("downloading page")

	resp, err := http.Get(pageURL)
	if err != nil {
		log.Error().Err(err).Msg("failed to download page")
		return
	}
	defer resp.Body.Close()

	z := html.NewTokenizer(resp.Body)
	for {
		tt := z.Next()

		switch tt {
		case html.ErrorToken:
			return
		case html.StartTagToken, html.EndTagToken:
			token := z.Token()
			if "a" == token.Data {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						if !c.isLinkExist(attr.Val) {
							urlObj, err := url.Parse(attr.Val)
							if err != nil {
								log.Error().Err(err).Str("raw_url", attr.Val).Msg("failed to parse URL")
								continue
							}
							if urlObj.Hostname() == domain {
								c.storeLink(attr.Val)

								c.wg.Add(1)
								go c.scrapePage(attr.Val, domain, limit)
							}
						}

					}
				}
			}
		default:
			// do nothing
		}
	}
}
