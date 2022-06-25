package pirsch

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

const (
	defaultBaseURL          = "https://api.pirsch.io"
	authenticationEndpoint  = "/api/v1/token"
	hitEndpoint             = "/api/v1/hit"
	eventEndpoint           = "/api/v1/event"
	sessionEndpoint         = "/api/v1/session"
	domainEndpoint          = "/api/v1/domain"
	sessionDurationEndpoint = "/api/v1/statistics/duration/session"
	timeOnPageEndpoint      = "/api/v1/statistics/duration/page"
	utmSourceEndpoint       = "/api/v1/statistics/utm/source"
	utmMediumEndpoint       = "/api/v1/statistics/utm/medium"
	utmCampaignEndpoint     = "/api/v1/statistics/utm/campaign"
	utmContentEndpoint      = "/api/v1/statistics/utm/content"
	utmTermEndpoint         = "/api/v1/statistics/utm/term"
	totalVisitorsEndpoint   = "/api/v1/statistics/total"
	visitorsEndpoint        = "/api/v1/statistics/visitor"
	pagesEndpoint           = "/api/v1/statistics/page"
	entryPagesEndpoint      = "/api/v1/statistics/page/entry"
	exitPagesEndpoint       = "/api/v1/statistics/page/exit"
	conversionGoalsEndpoint = "/api/v1/statistics/goals"
	eventsEndpoint          = "/api/v1/statistics/events"
	eventMetadataEndpoint   = "/api/v1/statistics/event/meta"
	listEventsEndpoint      = "/api/v1/statistics/event/list"
	growthRateEndpoint      = "/api/v1/statistics/growth"
	activeVisitorsEndpoint  = "/api/v1/statistics/active"
	timeOfDayEndpoint       = "/api/v1/statistics/hours"
	languageEndpoint        = "/api/v1/statistics/language"
	referrerEndpoint        = "/api/v1/statistics/referrer"
	osEndpoint              = "/api/v1/statistics/os"
	osVersionEndpoint       = "/api/v1/statistics/os/version"
	browserEndpoint         = "/api/v1/statistics/browser"
	browserVersionEndpoint  = "/api/v1/statistics/browser/version"
	countryEndpoint         = "/api/v1/statistics/country"
	cityEndpoint            = "/api/v1/statistics/city"
	platformEndpoint        = "/api/v1/statistics/platform"
	screenEndpoint          = "/api/v1/statistics/screen"
	keywordsEndpoint        = "/api/v1/statistics/keywords"
	requestRetries          = 5
)

var referrerQueryParams = []string{
	"ref",
	"referer",
	"referrer",
	"source",
	"utm_source",
}

// Client is used to access the Pirsch API.
type Client struct {
	baseURL      string
	logger       *log.Logger
	clientID     string
	clientSecret string
	hostname     string
	accessToken  string
	expiresAt    time.Time
	m            sync.RWMutex
}

// ClientConfig is used to configure the Client.
type ClientConfig struct {
	// BaseURL is optional and can be used to configure a different host for the API.
	// This is usually left empty in production environments.
	BaseURL string

	// Logger is an optional logger for debugging.
	Logger *log.Logger
}

// HitOptions optional parameters to send with the hit request.
type HitOptions struct {
	Hostname       string
	URL            string
	IP             string
	CFConnectingIP string
	XForwardedFor  string
	Forwarded      string
	XRealIP        string
	UserAgent      string
	AcceptLanguage string
	Title          string
	Referrer       string
	ScreenWidth    int
	ScreenHeight   int
}

// NewClient creates a new client for given client ID, client secret, hostname, and optional configuration.
// A new client ID and secret can be generated on the Pirsch dashboard.
// The hostname must match the hostname you configured on the Pirsch dashboard (e.g. example.com).
// The clientID is optional when using single access tokens.
func NewClient(clientID, clientSecret, hostname string, config *ClientConfig) *Client {
	if config == nil {
		config = &ClientConfig{
			BaseURL: defaultBaseURL,
		}
	}

	if config.BaseURL == "" {
		config.BaseURL = defaultBaseURL
	}

	c := &Client{
		baseURL:      config.BaseURL,
		logger:       config.Logger,
		clientID:     clientID,
		clientSecret: clientSecret,
		hostname:     hostname,
	}

	// single access tokens do not require to query an access token using oAuth
	if clientID == "" {
		c.accessToken = clientSecret
	}

	return c
}

// Hit sends a page hit to Pirsch for given http.Request.
func (client *Client) Hit(r *http.Request) error {
	return client.HitWithOptions(r, nil)
}

// HitWithOptions sends a page hit to Pirsch for given http.Request and options.
func (client *Client) HitWithOptions(r *http.Request, options *HitOptions) error {
	if r.Header.Get("DNT") == "1" {
		return nil
	}

	if options == nil {
		options = new(HitOptions)
	}

	hit := client.getHit(r, options)
	return client.performPost(client.baseURL+hitEndpoint, &hit, requestRetries)
}

// Event sends an event to Pirsch for given http.Request.
func (client *Client) Event(name string, durationSeconds int, meta map[string]string, r *http.Request) error {
	return client.EventWithOptions(name, durationSeconds, meta, r, nil)
}

// EventWithOptions sends an event to Pirsch for given http.Request and options.
func (client *Client) EventWithOptions(name string, durationSeconds int, meta map[string]string, r *http.Request, options *HitOptions) error {
	if r.Header.Get("DNT") == "1" {
		return nil
	}

	if options == nil {
		options = new(HitOptions)
	}

	return client.performPost(client.baseURL+eventEndpoint, &Event{
		Name:            name,
		DurationSeconds: durationSeconds,
		Metadata:        meta,
		Hit:             client.getHit(r, options),
	}, requestRetries)
}

// Session keeps a session alive for the given http.Request.
func (client *Client) Session(r *http.Request) error {
	return client.HitWithOptions(r, nil)
}

// SessionWithOptions keeps a session alive for the given http.Request and options.
func (client *Client) SessionWithOptions(r *http.Request, options *HitOptions) error {
	if r.Header.Get("DNT") == "1" {
		return nil
	}

	if options == nil {
		options = new(HitOptions)
	}

	return client.performPost(client.baseURL+sessionEndpoint, &Hit{
		Hostname:       client.hostname,
		URL:            r.URL.String(),
		IP:             r.RemoteAddr,
		CFConnectingIP: r.Header.Get("CF-Connecting-IP"),
		XForwardedFor:  r.Header.Get("X-Forwarded-For"),
		Forwarded:      r.Header.Get("Forwarded"),
		XRealIP:        r.Header.Get("X-Real-IP"),
		UserAgent:      r.Header.Get("User-Agent"),
	}, requestRetries)
}

// Domain returns the domain for this client.
func (client *Client) Domain() (*Domain, error) {
	domains := make([]Domain, 0, 1)

	if err := client.performGet(client.baseURL+domainEndpoint, requestRetries, &domains); err != nil {
		return nil, err
	}

	if len(domains) != 1 {
		return nil, errors.New("domain not found")
	}

	return &domains[0], nil
}

// SessionDuration returns the session duration grouped by day.
func (client *Client) SessionDuration(filter *Filter) ([]TimeSpentStats, error) {
	stats := make([]TimeSpentStats, 0)

	if err := client.performGet(client.getStatsRequestURL(sessionDurationEndpoint, filter), requestRetries, &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

// TimeOnPage returns the time spent on pages.
func (client *Client) TimeOnPage(filter *Filter) ([]TimeSpentStats, error) {
	stats := make([]TimeSpentStats, 0)

	if err := client.performGet(client.getStatsRequestURL(timeOnPageEndpoint, filter), requestRetries, &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

// UTMSource returns the utm sources.
func (client *Client) UTMSource(filter *Filter) ([]UTMSourceStats, error) {
	stats := make([]UTMSourceStats, 0)

	if err := client.performGet(client.getStatsRequestURL(utmSourceEndpoint, filter), requestRetries, &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

// UTMMedium returns the utm medium.
func (client *Client) UTMMedium(filter *Filter) ([]UTMMediumStats, error) {
	stats := make([]UTMMediumStats, 0)

	if err := client.performGet(client.getStatsRequestURL(utmMediumEndpoint, filter), requestRetries, &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

// UTMCampaign returnst he utm campaigns.
func (client *Client) UTMCampaign(filter *Filter) ([]UTMCampaignStats, error) {
	stats := make([]UTMCampaignStats, 0)

	if err := client.performGet(client.getStatsRequestURL(utmCampaignEndpoint, filter), requestRetries, &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

// UTMContent returns the utm content.
func (client *Client) UTMContent(filter *Filter) ([]UTMContentStats, error) {
	stats := make([]UTMContentStats, 0)

	if err := client.performGet(client.getStatsRequestURL(utmContentEndpoint, filter), requestRetries, &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

// UTMTerm returns the utm term.
func (client *Client) UTMTerm(filter *Filter) ([]UTMTermStats, error) {
	stats := make([]UTMTermStats, 0)

	if err := client.performGet(client.getStatsRequestURL(utmTermEndpoint, filter), requestRetries, &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

// TotalVisitors returns the total visitor statistics.
func (client *Client) TotalVisitors(filter *Filter) (*TotalVisitorStats, error) {
	stats := new(TotalVisitorStats)

	if err := client.performGet(client.getStatsRequestURL(totalVisitorsEndpoint, filter), requestRetries, stats); err != nil {
		return nil, err
	}

	return stats, nil
}

// Visitors returns the visitor statistics grouped by day.
func (client *Client) Visitors(filter *Filter) ([]VisitorStats, error) {
	stats := make([]VisitorStats, 0)

	if err := client.performGet(client.getStatsRequestURL(visitorsEndpoint, filter), requestRetries, &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

// Pages returns the page statistics grouped by page.
func (client *Client) Pages(filter *Filter) ([]PageStats, error) {
	stats := make([]PageStats, 0)

	if err := client.performGet(client.getStatsRequestURL(pagesEndpoint, filter), requestRetries, &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

// EntryPages returns the entry page statistics grouped by page.
func (client *Client) EntryPages(filter *Filter) ([]EntryStats, error) {
	stats := make([]EntryStats, 0)

	if err := client.performGet(client.getStatsRequestURL(entryPagesEndpoint, filter), requestRetries, &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

// ExitPages returns the exit page statistics grouped by page.
func (client *Client) ExitPages(filter *Filter) ([]ExitStats, error) {
	stats := make([]ExitStats, 0)

	if err := client.performGet(client.getStatsRequestURL(exitPagesEndpoint, filter), requestRetries, &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

// ConversionGoals returns all conversion goals.
func (client *Client) ConversionGoals(filter *Filter) ([]ConversionGoal, error) {
	stats := make([]ConversionGoal, 0)

	if err := client.performGet(client.getStatsRequestURL(conversionGoalsEndpoint, filter), requestRetries, &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

// Events returns all events.
func (client *Client) Events(filter *Filter) ([]EventStats, error) {
	stats := make([]EventStats, 0)

	if err := client.performGet(client.getStatsRequestURL(eventsEndpoint, filter), requestRetries, &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

// EventMetadata returns the metadata values for an event and key.
func (client *Client) EventMetadata(filter *Filter) ([]EventStats, error) {
	stats := make([]EventStats, 0)

	if err := client.performGet(client.getStatsRequestURL(eventMetadataEndpoint, filter), requestRetries, &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

// ListEvents returns a list of all events including metadata.
func (client *Client) ListEvents(filter *Filter) ([]EventListStats, error) {
	stats := make([]EventListStats, 0)

	if err := client.performGet(client.getStatsRequestURL(listEventsEndpoint, filter), requestRetries, &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

// Growth returns the growth rates for visitors, bounces, ...
func (client *Client) Growth(filter *Filter) (*Growth, error) {
	growth := new(Growth)

	if err := client.performGet(client.getStatsRequestURL(growthRateEndpoint, filter), requestRetries, growth); err != nil {
		return nil, err
	}

	return growth, nil
}

// ActiveVisitors returns the active visitors and what pages they're on.
func (client *Client) ActiveVisitors(filter *Filter) (*ActiveVisitorsData, error) {
	active := new(ActiveVisitorsData)

	if err := client.performGet(client.getStatsRequestURL(activeVisitorsEndpoint, filter), requestRetries, active); err != nil {
		return nil, err
	}

	return active, nil
}

// TimeOfDay returns the number of unique visitors grouped by time of day.
func (client *Client) TimeOfDay(filter *Filter) ([]VisitorHourStats, error) {
	stats := make([]VisitorHourStats, 0)

	if err := client.performGet(client.getStatsRequestURL(timeOfDayEndpoint, filter), requestRetries, &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

// Languages returns language statistics.
func (client *Client) Languages(filter *Filter) ([]LanguageStats, error) {
	stats := make([]LanguageStats, 0)

	if err := client.performGet(client.getStatsRequestURL(languageEndpoint, filter), requestRetries, &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

// Referrer returns referrer statistics.
func (client *Client) Referrer(filter *Filter) ([]ReferrerStats, error) {
	stats := make([]ReferrerStats, 0)

	if err := client.performGet(client.getStatsRequestURL(referrerEndpoint, filter), requestRetries, &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

// OS returns operating system statistics.
func (client *Client) OS(filter *Filter) ([]OSStats, error) {
	stats := make([]OSStats, 0)

	if err := client.performGet(client.getStatsRequestURL(osEndpoint, filter), requestRetries, &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

// OSVersions returns operating system version statistics.
func (client *Client) OSVersions(filter *Filter) ([]OSVersionStats, error) {
	stats := make([]OSVersionStats, 0)

	if err := client.performGet(client.getStatsRequestURL(osVersionEndpoint, filter), requestRetries, &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

// Browser returns browser statistics.
func (client *Client) Browser(filter *Filter) ([]BrowserStats, error) {
	stats := make([]BrowserStats, 0)

	if err := client.performGet(client.getStatsRequestURL(browserEndpoint, filter), requestRetries, &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

// BrowserVersions returns browser version statistics.
func (client *Client) BrowserVersions(filter *Filter) ([]BrowserVersionStats, error) {
	stats := make([]BrowserVersionStats, 0)

	if err := client.performGet(client.getStatsRequestURL(browserVersionEndpoint, filter), requestRetries, &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

// Country returns country statistics.
func (client *Client) Country(filter *Filter) ([]CountryStats, error) {
	stats := make([]CountryStats, 0)

	if err := client.performGet(client.getStatsRequestURL(countryEndpoint, filter), requestRetries, &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

// City returns city statistics.
func (client *Client) City(filter *Filter) ([]CityStats, error) {
	stats := make([]CityStats, 0)

	if err := client.performGet(client.getStatsRequestURL(cityEndpoint, filter), requestRetries, &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

// Platform returns the platforms used by visitors.
func (client *Client) Platform(filter *Filter) (*PlatformStats, error) {
	platforms := new(PlatformStats)

	if err := client.performGet(client.getStatsRequestURL(platformEndpoint, filter), requestRetries, platforms); err != nil {
		return nil, err
	}

	return platforms, nil
}

// Screen returns the screen classes used by visitors.
func (client *Client) Screen(filter *Filter) ([]ScreenClassStats, error) {
	stats := make([]ScreenClassStats, 0)

	if err := client.performGet(client.getStatsRequestURL(screenEndpoint, filter), requestRetries, &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

// Keywords returns the Google keywords, rank, and CTR.
func (client *Client) Keywords(filter *Filter) ([]Keyword, error) {
	stats := make([]Keyword, 0)

	if err := client.performGet(client.getStatsRequestURL(keywordsEndpoint, filter), requestRetries, &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

func (client *Client) getHit(r *http.Request, options *HitOptions) Hit {
	return Hit{
		Hostname:       client.selectField(options.Hostname, client.hostname),
		URL:            client.selectField(options.URL, r.URL.String()),
		IP:             client.selectField(options.IP, r.RemoteAddr),
		CFConnectingIP: client.selectField(options.CFConnectingIP, r.Header.Get("CF-Connecting-IP")),
		XForwardedFor:  client.selectField(options.XForwardedFor, r.Header.Get("X-Forwarded-For")),
		Forwarded:      client.selectField(options.Forwarded, r.Header.Get("Forwarded")),
		XRealIP:        client.selectField(options.XRealIP, r.Header.Get("X-Real-IP")),
		UserAgent:      client.selectField(options.UserAgent, r.Header.Get("User-Agent")),
		AcceptLanguage: client.selectField(options.AcceptLanguage, r.Header.Get("Accept-Language")),
		Title:          options.Title,
		Referrer:       client.selectField(options.Referrer, client.getReferrerFromHeaderOrQuery(r)),
		ScreenWidth:    options.ScreenWidth,
		ScreenHeight:   options.ScreenHeight,
	}
}

func (client *Client) getReferrerFromHeaderOrQuery(r *http.Request) string {
	referrer := r.Header.Get("Referer")

	if referrer == "" {
		for _, param := range referrerQueryParams {
			referrer = r.URL.Query().Get(param)

			if referrer != "" {
				return referrer
			}
		}
	}

	return referrer
}

func (client *Client) refreshToken() error {
	client.m.Lock()
	defer client.m.Unlock()
	client.accessToken = ""
	client.expiresAt = time.Time{}
	body := struct {
		ClientId     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
	}{
		client.clientID,
		client.clientSecret,
	}
	bodyJson, err := json.Marshal(&body)

	if err != nil {
		return err
	}

	c := http.Client{}
	resp, err := c.Post(client.baseURL+authenticationEndpoint, "application/json", bytes.NewBuffer(bodyJson))

	if err != nil {
		return err
	}

	respJson := struct {
		AccessToken string    `json:"access_token"`
		ExpiresAt   time.Time `json:"expires_at"`
	}{}

	decoder := json.NewDecoder(resp.Body)

	if err := decoder.Decode(&respJson); err != nil {
		return err
	}

	client.accessToken = respJson.AccessToken
	client.expiresAt = respJson.ExpiresAt
	return nil
}

func (client *Client) performPost(url string, body interface{}, retry int) error {
	client.m.RLock()
	accessToken := client.accessToken
	client.m.RUnlock()

	if client.clientID != "" && retry > 0 && accessToken == "" {
		client.waitBeforeNextRequest(retry)

		if err := client.refreshToken(); err != nil {
			if client.logger != nil {
				client.logger.Printf("error refreshing token: %s", err)
			}

			return errors.New(fmt.Sprintf("error refreshing token (attempt %d/%d): %s", requestRetries-retry, requestRetries, err))
		}

		return client.performPost(url, body, retry-1)
	}

	reqBody, err := json.Marshal(body)

	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(reqBody))

	if err != nil {
		return err
	}

	client.m.RLock()
	req.Header.Set("Authorization", "Bearer "+client.accessToken)
	client.m.RUnlock()
	c := http.Client{}
	resp, err := c.Do(req)

	if err != nil {
		return err
	}

	// refresh access token and retry
	if client.clientID != "" && retry > 0 && resp.StatusCode != http.StatusOK {
		client.waitBeforeNextRequest(retry)

		if err := client.refreshToken(); err != nil {
			if client.logger != nil {
				client.logger.Printf("error refreshing token: %s", err)
			}

			return errors.New(fmt.Sprintf("error refreshing token (attempt %d/%d): %s", requestRetries-retry, requestRetries, err))
		}

		return client.performPost(url, body, retry-1)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return client.requestError(url, resp.StatusCode, string(body))
	}

	return nil
}

func (client *Client) performGet(url string, retry int, result interface{}) error {
	client.m.RLock()
	accessToken := client.accessToken
	client.m.RUnlock()

	if client.clientID != "" && retry > 0 && accessToken == "" {
		client.waitBeforeNextRequest(retry)

		if err := client.refreshToken(); err != nil {
			if client.logger != nil {
				client.logger.Printf("error refreshing token: %s", err)
			}

			return errors.New(fmt.Sprintf("error refreshing token (attempt %d/%d): %s", requestRetries-retry, requestRetries, err))
		}

		return client.performGet(url, retry-1, result)
	}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return err
	}

	client.m.RLock()
	req.Header.Set("Authorization", "Bearer "+client.accessToken)
	client.m.RUnlock()
	req.Header.Set("Content-Type", "application/json")
	c := http.Client{}
	resp, err := c.Do(req)

	if err != nil {
		return err
	}

	// refresh access token and retry
	if client.clientID != "" && retry > 0 && resp.StatusCode != http.StatusOK {
		client.waitBeforeNextRequest(retry)

		if err := client.refreshToken(); err != nil {
			if client.logger != nil {
				client.logger.Printf("error refreshing token: %s", err)
			}

			return errors.New(fmt.Sprintf("error refreshing token (attempt %d/%d): %s", requestRetries-retry, requestRetries, err))
		}

		return client.performGet(url, retry-1, result)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return client.requestError(url, resp.StatusCode, string(body))
	}

	decoder := json.NewDecoder(resp.Body)

	if err := decoder.Decode(result); err != nil {
		return err
	}

	return nil
}

func (client *Client) requestError(url string, statusCode int, body string) error {
	if body != "" {
		return errors.New(fmt.Sprintf("%s: received status code %d on request: %s", url, statusCode, body))
	}

	return errors.New(fmt.Sprintf("%s: received status code %d on request", url, statusCode))
}

func (client *Client) getStatsRequestURL(endpoint string, filter *Filter) string {
	u := fmt.Sprintf("%s%s", client.baseURL, endpoint)
	v := url.Values{}
	v.Set("id", filter.DomainID)
	v.Set("from", filter.From.Format("2006-01-02"))
	v.Set("to", filter.To.Format("2006-01-02"))
	v.Set("path", filter.Path)
	v.Set("entry_path", filter.EntryPath)
	v.Set("exit_path", filter.ExitPath)
	v.Set("pattern", filter.Pattern)
	v.Set("event", filter.Event)
	v.Set("event_meta_key", filter.EventMetaKey)
	v.Set("language", filter.Language)
	v.Set("country", filter.Country)
	v.Set("city", filter.City)
	v.Set("referrer", filter.Referrer)
	v.Set("referrer_name", filter.ReferrerName)
	v.Set("os", filter.OS)
	v.Set("browser", filter.Browser)
	v.Set("platform", filter.Platform)
	v.Set("screen_class", filter.ScreenClass)
	v.Set("utm_source", filter.UTMSource)
	v.Set("utm_medium", filter.UTMMedium)
	v.Set("utm_campaign", filter.UTMCampaign)
	v.Set("utm_content", filter.UTMContent)
	v.Set("utm_term", filter.UTMTerm)
	v.Set("limit", strconv.Itoa(filter.Limit))

	if filter.IncludeAvgTimeOnPage {
		v.Set("include_avg_time_on_page", "true")
	} else {
		v.Set("include_avg_time_on_page", "false")
	}

	return u + "?" + v.Encode()
}

func (client *Client) waitBeforeNextRequest(retry int) {
	time.Sleep(time.Second * time.Duration(requestRetries-retry))
}

func (client *Client) selectField(a, b string) string {
	if a != "" {
		return a
	}

	return b
}
