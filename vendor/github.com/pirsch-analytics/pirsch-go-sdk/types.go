package pirsch

import (
	"github.com/emvi/null"
	"time"
)

const (
	// ScaleDay groups results by day.
	ScaleDay = "day"

	// ScaleWeek groups results by week.
	ScaleWeek = "week"

	// ScaleMonth groups results by month.
	ScaleMonth = "month"

	// ScaleYear groups results by year.
	ScaleYear = "year"
)

// Scale is used to group results in the Filter.
// Use one of the constants ScaleDay, ScaleWeek, ScaleMonth, ScaleYear.
type Scale string

// Hit are the parameters to send a page hit to Pirsch.
type Hit struct {
	Hostname       string
	URL            string `json:"url"`
	IP             string `json:"ip"`
	CFConnectingIP string `json:"cf_connecting_ip"`
	XForwardedFor  string `json:"x_forwarded_for"`
	Forwarded      string `json:"forwarded"`
	XRealIP        string `json:"x_real_ip"`
	UserAgent      string `json:"user_agent"`
	AcceptLanguage string `json:"accept_language"`
	Title          string `json:"title"`
	Referrer       string `json:"referrer"`
	ScreenWidth    int    `json:"screen_width"`
	ScreenHeight   int    `json:"screen_height"`
}

// Event represents a single data point for custom events.
// It's basically the same as Hit, but with some additional fields (event name, time, and meta fields).
type Event struct {
	Hit
	Name            string            `json:"event_name"`
	DurationSeconds int               `json:"event_duration"`
	Metadata        map[string]string `json:"event_meta"`
}

// Filter is used to filter statistics.
// DomainID, From, and To are required dates (the time is ignored).
type Filter struct {
	DomainID             string    `json:"id"`
	From                 time.Time `json:"from"`
	To                   time.Time `json:"to"`
	Start                int       `json:"start,omitempty"`
	Scale                Scale     `json:"scale,omitempty"`
	Path                 string    `json:"path,omitempty"`
	Pattern              string    `json:"pattern,omitempty"`
	EntryPath            string    `json:"entry_path,omitempty"`
	ExitPath             string    `json:"exit_path,omitempty"`
	Event                string    `json:"event,omitempty"`
	EventMetaKey         string    `json:"event_meta_key,omitempty"`
	Language             string    `json:"language,omitempty"`
	Country              string    `json:"country,omitempty"`
	City                 string    `json:"city,omitempty"`
	Referrer             string    `json:"referrer,omitempty"`
	ReferrerName         string    `json:"referrer_name,omitempty"`
	OS                   string    `json:"os,omitempty"`
	Browser              string    `json:"browser,omitempty"`
	Platform             string    `json:"platform,omitempty"`
	ScreenClass          string    `json:"screen_class,omitempty"`
	ScreenWidth          string    `json:"screen_width,omitempty"`
	ScreenHeight         string    `json:"screen_height,omitempty"`
	UTMSource            string    `json:"utm_source,omitempty"`
	UTMMedium            string    `json:"utm_medium,omitempty"`
	UTMCampaign          string    `json:"utm_campaign,omitempty"`
	UTMContent           string    `json:"utm_content,omitempty"`
	UTMTerm              string    `json:"utm_term,omitempty"`
	Limit                int       `json:"limit,omitempty"`
	IncludeAvgTimeOnPage bool      `json:"include_avg_time_on_page,omitempty"`
}

// BaseEntity contains the base data for all entities.
type BaseEntity struct {
	ID      string    `json:"id"`
	DefTime time.Time `json:"def_time"`
	ModTime time.Time `json:"mod_time"`
}

// Domain is a domain on the dashboard.
type Domain struct {
	BaseEntity

	UserID             string      `json:"user_id"`
	Hostname           string      `json:"hostname"`
	Subdomain          string      `json:"subdomain"`
	IdentificationCode string      `json:"identification_code"`
	Public             bool        `json:"public"`
	GoogleUserID       null.String `json:"google_user_id"`
	GoogleUserEmail    null.String `json:"google_user_email"`
	GSCDomain          null.String `json:"gsc_domain"`
	NewOwner           null.Int64  `json:"new_owner"`
	Timezone           null.String `json:"timezone"`
}

// TimeSpentStats is the time spent on the website or specific pages.
type TimeSpentStats struct {
	Day                     null.Time `json:"day"`
	Week                    null.Time `json:"week"`
	Month                   null.Time `json:"month"`
	Year                    null.Time `json:"year"`
	Path                    string    `json:"path"`
	Title                   string    `json:"title"`
	AverageTimeSpentSeconds int       `json:"average_time_spent_seconds"`
}

// MetaStats is the base for meta result types (languages, countries, ...).
type MetaStats struct {
	Visitors         int     `json:"visitors"`
	RelativeVisitors float64 `json:"relative_visitors"`
}

// UTMSourceStats is the result type for utm source statistics.
type UTMSourceStats struct {
	MetaStats
	UTMSource string `json:"utm_source"`
}

// UTMMediumStats is the result type for utm medium statistics.
type UTMMediumStats struct {
	MetaStats
	UTMMedium string `json:"utm_medium"`
}

// UTMCampaignStats is the result type for utm campaign statistics.
type UTMCampaignStats struct {
	MetaStats
	UTMCampaign string `json:"utm_campaign"`
}

// UTMContentStats is the result type for utm content statistics.
type UTMContentStats struct {
	MetaStats
	UTMContent string `json:"utm_content"`
}

// UTMTermStats is the result type for utm term statistics.
type UTMTermStats struct {
	MetaStats
	UTMTerm string `json:"utm_term"`
}

// TotalVisitorStats is the result type for total visitor statistics.
type TotalVisitorStats struct {
	Visitors   int     `json:"visitors"`
	Views      int     `json:"views"`
	Sessions   int     `json:"sessions"`
	Bounces    int     `json:"bounces"`
	BounceRate float64 `json:"bounce_rate"`
}

// VisitorStats is the result type for visitor statistics.
type VisitorStats struct {
	Day        null.Time `json:"day"`
	Week       null.Time `json:"week"`
	Month      null.Time `json:"month"`
	Year       null.Time `json:"year"`
	Visitors   int       `json:"visitors"`
	Views      int       `json:"views"`
	Sessions   int       `json:"sessions"`
	Bounces    int       `json:"bounces"`
	BounceRate float64   `json:"bounce_rate"`
}

// PageStats is the result type for page statistics.
type PageStats struct {
	Path                    string  `json:"path"`
	Visitors                int     `json:"visitors"`
	Views                   int     `json:"views"`
	Sessions                int     `json:"sessions"`
	Bounces                 int     `json:"bounces"`
	RelativeVisitors        float64 `json:"relative_visitors"`
	RelativeViews           float64 `json:"relative_views"`
	BounceRate              float64 `json:"bounce_rate"`
	AverageTimeSpentSeconds int     `json:"average_time_spent_seconds"`
}

// EntryStats is the result type for entry page statistics.
type EntryStats struct {
	Path                    string  `json:"path"`
	Title                   string  `json:"title"`
	Visitors                int     `json:"visitors"`
	Sessions                int     `json:"sessions"`
	Entries                 int     `json:"entries"`
	EntryRate               float64 `json:"entry_rate"`
	AverageTimeSpentSeconds int     `json:"average_time_spent_seconds"`
}

// ExitStats is the result type for exit page statistics.
type ExitStats struct {
	Path     string  `json:"path"`
	Title    string  `json:"title"`
	Visitors int     `json:"visitors"`
	Sessions int     `json:"sessions"`
	Exits    int     `json:"exits"`
	ExitRate float64 `json:"exit_rate"`
}

// ConversionGoal is a conversion goal as configured on the dashboard.
type ConversionGoal struct {
	BaseEntity

	PageGoal struct {
		DomainID      string       `json:"domain_id"`
		Name          string       `json:"name"`
		PathPattern   string       `json:"path_pattern"`
		Pattern       string       `json:"pattern"`
		VisitorGoal   null.Int64   `json:"visitor_goal"`
		CRGoal        null.Float64 `json:"cr_goal"`
		DeleteReached bool         `json:"delete_reached"`
		EmailReached  bool         `json:"email_reached"`
	} `json:"page_goal"`
	Stats struct {
		Visitors int     `json:"visitors"`
		Views    int     `json:"views"`
		CR       float64 `json:"cr"`
	} `json:"stats"`
}

// EventStats is the result type for custom events.
type EventStats struct {
	Name                   string   `json:"name"`
	Visitors               int      `json:"visitors"`
	Views                  int      `json:"views"`
	CR                     float64  `json:"cr"`
	AverageDurationSeconds int      `json:"average_duration_seconds"`
	MetaKeys               []string `json:"meta_keys"`
	MetaValue              string   `json:"meta_value"`
}

// EventListStats is the result type for a custom event list.
type EventListStats struct {
	Name     string            `json:"name"`
	Meta     map[string]string `json:"meta"`
	Visitors int               `json:"visitors"`
	Count    int               `json:"count"`
}

// PageConversionsStats is the result type for page conversions.
type PageConversionsStats struct {
	Visitors int     `json:"visitors"`
	Views    int     `json:"views"`
	CR       float64 `json:"cr"`
}

// ConversionGoalStats are the statistics for a conversion goal.
type ConversionGoalStats struct {
	ConversionGoal *ConversionGoal       `json:"page_goal"` // page_goal is returned by the API, but we name it differently here
	Stats          *PageConversionsStats `json:"stats"`
}

// Growth represents the visitors, views, sessions, bounces, and average session duration growth between two time periods.
type Growth struct {
	VisitorsGrowth  float64 `json:"visitors_growth"`
	ViewsGrowth     float64 `json:"views_growth"`
	SessionsGrowth  float64 `json:"sessions_growth"`
	BouncesGrowth   float64 `json:"bounces_growth"`
	TimeSpentGrowth float64 `json:"time_spent_growth"`
}

// ActiveVisitorStats is the result type for active visitor statistics.
type ActiveVisitorStats struct {
	Path     string `json:"path"`
	Title    string `json:"title"`
	Visitors int    `json:"visitors"`
}

// ActiveVisitorsData contains the active visitors data.
type ActiveVisitorsData struct {
	Stats    []ActiveVisitorStats `json:"stats"`
	Visitors int                  `json:"visitors"`
}

// VisitorHourStats is the result type for visitor statistics grouped by time of day.
type VisitorHourStats struct {
	Hour       int     `json:"hour"`
	Visitors   int     `json:"visitors"`
	Views      int     `json:"views"`
	Sessions   int     `json:"sessions"`
	Bounces    int     `json:"bounces"`
	BounceRate float64 `json:"bounce_rate"`
}

// LanguageStats is the result type for language statistics.
type LanguageStats struct {
	MetaStats
	Language string `json:"language"`
}

// CountryStats is the result type for country statistics.
type CountryStats struct {
	MetaStats
	CountryCode string `json:"country_code"`
}

// CityStats is the result type for city statistics.
type CityStats struct {
	MetaStats
	City string `json:"city"`
}

// BrowserStats is the result type for browser statistics.
type BrowserStats struct {
	MetaStats
	Browser string `json:"browser"`
}

// BrowserVersionStats is the result type for browser version statistics.
type BrowserVersionStats struct {
	MetaStats
	Browser        string `json:"browser"`
	BrowserVersion string `db:"browser_version" json:"browser_version"`
}

// OSStats is the result type for operating system statistics.
type OSStats struct {
	MetaStats
	OS string `json:"os"`
}

// OSVersionStats is the result type for operating system version statistics.
type OSVersionStats struct {
	MetaStats
	OS        string `json:"os"`
	OSVersion string `db:"os_version" json:"os_version"`
}

// ReferrerStats is the result type for referrer statistics.
type ReferrerStats struct {
	Referrer         string  `json:"referrer"`
	ReferrerName     string  `json:"referrer_name"`
	ReferrerIcon     string  `json:"referrer_icon"`
	Visitors         int     `json:"visitors"`
	Sessions         int     `json:"sessions"`
	RelativeVisitors float64 `json:"relative_visitors"`
	Bounces          int     `json:"bounces"`
	BounceRate       float64 `json:"bounce_rate"`
}

// PlatformStats is the result type for platform statistics.
type PlatformStats struct {
	PlatformDesktop         int     `json:"platform_desktop"`
	PlatformMobile          int     `json:"platform_mobile"`
	PlatformUnknown         int     `json:"platform_unknown"`
	RelativePlatformDesktop float64 `json:"relative_platform_desktop"`
	RelativePlatformMobile  float64 `json:"relative_platform_mobile"`
	RelativePlatformUnknown float64 `json:"relative_platform_unknown"`
}

// ScreenClassStats is the result type for screen class statistics.
type ScreenClassStats struct {
	MetaStats
	ScreenClass string `json:"screen_class"`
}

// Keyword is the result type for keyword statistics.
type Keyword struct {
	Keys        []string `json:"keys"`
	Clicks      int      `json:"clicks"`
	Impressions int      `json:"impressions"`
	CTR         float64  `json:"ctr"`
	Position    float64  `json:"position"`
}
