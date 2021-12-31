/*
  Caddy v2 module to filter requests based on source IP geographic location. This was a feature provided by the V1 ipfilter middleware.
  Complete documentation and usage examples are available at https://github.com/porech/caddy-maxmind-geolocation
*/
package caddy_maxmind_geolocation

import (
	"fmt"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/oschwald/maxminddb-golang"
	"go.uber.org/zap"
	"net"
	"net/http"
	"strconv"
)

// Interface guards
var (
	_ caddy.Module             = (*MaxmindGeolocation)(nil)
	_ caddyhttp.RequestMatcher = (*MaxmindGeolocation)(nil)
	_ caddy.Provisioner        = (*MaxmindGeolocation)(nil)
	_ caddy.CleanerUpper       = (*MaxmindGeolocation)(nil)
	_ caddyfile.Unmarshaler    = (*MaxmindGeolocation)(nil)
)

func init() {
	caddy.RegisterModule(MaxmindGeolocation{})
}

// Allows to filter requests based on source IP country.
type MaxmindGeolocation struct {

	// The path of the MaxMind GeoLite2-Country.mmdb file.
	DbPath string `json:"db_path"`

	// A list of countries that the filter will allow.
	// If you specify this, you should not specify DenyCountries.
	// If both are specified, DenyCountries will take precedence.
	// All countries that are not in this list will be denied.
	// You can specify the special value "UNK" to match unrecognized countries.
	AllowCountries []string `json:"allow_countries"`

	// A list of countries that the filter will deny.
	// If you specify this, you should not specify AllowCountries.
	// If both are specified, DenyCountries will take precedence.
	// All countries that are not in this list will be allowed.
	// You can specify the special value "UNK" to match unrecognized countries.
	DenyCountries []string `json:"deny_countries"`

	// A list of subdivisions that the filter will allow.
	// If you specify this, you should not specify DenySubdivisions.
	// If both are specified, DenySubdivisions will take precedence.
	// All subdivisions that are not in this list will be denied.
	// You can specify the special value "UNK" to match unrecognized subdivisions.
	AllowSubdivisions []string `json:"allow_subdivisions"`

	// A list of subdivisions that the filter will deny.
	// If you specify this, you should not specify AllowSubdivisions.
	// If both are specified, DenySubdivisions will take precedence.
	// All subdivisions that are not in this list will be allowed.
	// You can specify the special value "UNK" to match unrecognized subdivisions.
	DenySubdivisions []string `json:"deny_subdivisions"`

	// A list of metro codes that the filter will allow.
	// If you specify this, you should not specify DenyMetroCodes.
	// If both are specified, DenyMetroCodes will take precedence.
	// All metro codes that are not in this list will be denied.
	// You can specify the special value "UNK" to match unrecognized metro codes.
	AllowMetroCodes []string `json:"allow_metro_codes"`

	// A list of METRO CODES that the filter will deny.
	// If you specify this, you should not specify AllowMetroCodes.
	// If both are specified, DenyMetroCodes will take precedence.
	// All metro codes that are not in this list will be allowed.
	// You can specify the special value "UNK" to match unrecognized metro codes.
	DenyMetroCodes []string `json:"deny_metro_codes"`

	dbInst *maxminddb.Reader
	logger *zap.Logger
}

/*
	The matcher configuration will have a single block with the following parameters:

	- `db_path`: required, is the path to the GeoLite2-Country.mmdb file

	- `allow_countries`: a space-separated list of allowed countries

	- `deny_countries`: a space-separated list of denied countries.

	You will want specify just one of `allow_countries` or `deny_countries`. If you
	specify both of them, denied countries will take precedence over allowed ones.
	If you specify none of them, all requests will be denied.

	Examples are available at https://github.com/porech/caddy-maxmind-geolocation/
*/
func (m *MaxmindGeolocation) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	current := 0
	for d.Next() {
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "db_path":
				current = 1
			case "allow_countries":
				current = 2
			case "deny_countries":
				current = 3
			case "allow_subdivisions":
				current = 4
			case "deny_subdivisions":
				current = 5
			case "allow_metro_codes":
				current = 6
			case "deny_metro_codes":
				current = 7
			default:
				switch current {
				case 1:
					m.DbPath = d.Val()
					current = 0
				case 2:
					m.AllowCountries = append(m.AllowCountries, d.Val())
				case 3:
					m.DenyCountries = append(m.DenyCountries, d.Val())
				case 4:
					m.AllowSubdivisions = append(m.AllowSubdivisions, d.Val())
				case 5:
					m.DenySubdivisions = append(m.DenySubdivisions, d.Val())
				case 6:
					m.AllowMetroCodes = append(m.AllowMetroCodes, d.Val())
				case 7:
					m.DenyMetroCodes = append(m.DenyMetroCodes, d.Val())
				default:
					return fmt.Errorf("unexpected config parameter %s", d.Val())
				}
			}
		}
	}
	return nil
}

func (MaxmindGeolocation) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.matchers.maxmind_geolocation",
		New: func() caddy.Module { return new(MaxmindGeolocation) },
	}
}

func (m *MaxmindGeolocation) Provision(ctx caddy.Context) error {
	var err error
	m.logger = ctx.Logger(m)
	m.dbInst, err = maxminddb.Open(m.DbPath)
	if err != nil {
		return fmt.Errorf("cannot open database file %s: %v", m.DbPath, err)
	}
	return nil
}

func (m *MaxmindGeolocation) Cleanup() error {
	if m.dbInst != nil {
		return m.dbInst.Close()
	}
	return nil
}

func (m *MaxmindGeolocation) checkAllowed(item string, allowedList []string, deniedList []string) bool {
	if item == "" || item == "0" {
		item = "UNK"
	}
	if len(deniedList) > 0 {
		for _, i := range deniedList {
			if i == item {
				return false
			}
		}
		return true
	}
	if len(allowedList) > 0 {
		for _, i := range allowedList {
			if i == item {
				return true
			}
		}
		return false
	}
	return true
}

func (m *MaxmindGeolocation) Match(r *http.Request) bool {

	// If both the allow and deny fields are empty, let the request pass
	if len(m.AllowCountries) < 1 && len(m.DenyCountries) < 1 {
		return false
	}

	remoteIp, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		m.logger.Warn("cannot split IP address", zap.String("address", r.RemoteAddr), zap.Error(err))
	}

	// Get the record from the database
	addr := net.ParseIP(remoteIp)
	if addr == nil {
		m.logger.Warn("cannot parse IP address", zap.String("address", r.RemoteAddr))
		return false
	}
	var record Record
	err = m.dbInst.Lookup(addr, &record)
	if err != nil {
		m.logger.Warn("cannot lookup IP address", zap.String("address", r.RemoteAddr), zap.Error(err))
		return false
	}

	m.logger.Debug(
		"Detected MaxMind data",
		zap.String("ip", r.RemoteAddr),
		zap.String("country", record.Country.ISOCode),
		zap.String("subdivisions", record.Subdivisions.CommaSeparatedISOCodes()),
		zap.Int("metro_code", record.Location.MetroCode),
	)

	if !m.checkAllowed(record.Country.ISOCode, m.AllowCountries, m.DenyCountries) {
		m.logger.Debug("Country not allowed", zap.String("country", record.Country.ISOCode))
		return false
	}

	if len(record.Subdivisions) > 0 {
		for _, subdivision := range record.Subdivisions {
			if !m.checkAllowed(subdivision.ISOCode, m.AllowSubdivisions, m.DenySubdivisions) {
				m.logger.Debug("Subdivision not allowed", zap.String("subdivision", subdivision.ISOCode))
				return false
			}
		}
	} else {
		if !m.checkAllowed("", m.AllowSubdivisions, m.DenySubdivisions) {
			m.logger.Debug("Subdivision not allowed", zap.String("subdivision", ""))
			return false
		}
	}

	if !m.checkAllowed(strconv.Itoa(record.Location.MetroCode), m.AllowMetroCodes, m.DenyMetroCodes) {
		m.logger.Debug("Metro code not allowed", zap.Int("metro_code", record.Location.MetroCode))
		return false
	}

	return true
}
