package duckdns

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/libdns/libdns"
	"github.com/miekg/dns"
)

func (p *Provider) getDomain(ctx context.Context, zone string) ([]libdns.Record, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// we trim the dot at the end of the zone name to get the fqdn
	// and then use it to fetch domain information through the api
	fqdn := strings.TrimRight(zone, ".")

	var libRecords []libdns.Record

	// fetch the details for the domain
	result, err := p.doRequest(ctx, fqdn, map[string]string{"verbose": "true"})
	if err != nil {
		return libRecords, err
	}

	// append the A and AAAA records which we should have (may be blank)
	if result[1] != "" {
		libRecords = append(libRecords, libdns.Record{
			Type:  "A",
			Name:  fqdn,
			Value: result[1],
		})
	}
	if result[2] != "" {
		libRecords = append(libRecords, libdns.Record{
			Type:  "AAAA",
			Name:  fqdn,
			Value: result[2],
		})
	}

	return libRecords, nil
}

func (p *Provider) setRecord(ctx context.Context, zone string, record libdns.Record, clear bool) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// sanitize the domain, combines the zone and record names
	// the record name should typically be relative to the zone
	domain := libdns.AbsoluteName(record.Name, zone)

	params := map[string]string{"verbose": "true"}

	switch record.Type {
	case "TXT":
		params["txt"] = record.Value
	case "A":
		params["ip"] = record.Value
	case "AAAA":
		params["ipv6"] = record.Value
	default:
		return fmt.Errorf("unsupported record type: %s", record.Type)
	}

	if clear {
		params["clear"] = "true"
	}

	// make the request to duckdns to set the records according to the params
	_, err := p.doRequest(ctx, domain, params)
	if err != nil {
		return err
	}
	return nil
}

func (p *Provider) doRequest(ctx context.Context, domain string, params map[string]string) ([]string, error) {
	u, _ := url.Parse("https://www.duckdns.org/update")

	// extract the main domain
	var mainDomain string
	if p.OverrideDomain != "" {
		mainDomain = p.OverrideDomain
	} else {
		mainDomain = getMainDomain(domain)
	}

	if len(mainDomain) == 0 {
		return nil, fmt.Errorf("unable to find the main domain for: %s", domain)
	}

	// set up the query with the params we always set
	query := u.Query()
	query.Set("domains", mainDomain)
	query.Set("token", p.APIToken)

	// add the remaining ones for this request
	for key, val := range params {
		query.Set(key, val)
	}

	// set the query back on the URL
	u.RawQuery = query.Encode()

	// make the request
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	body := string(bodyBytes)
	bodyParts := strings.Split(body, "\n")
	if bodyParts[0] != "OK" {
		return nil, fmt.Errorf("DuckDNS request failed, expected (OK) but got (%s), url: [%s], body: %s", bodyParts[0], u, body)
	}

	return bodyParts, nil
}

// DuckDNS only lets you write to your subdomain.
// It must be in format subdomain.duckdns.org,
// not in format subsubdomain.subdomain.duckdns.org.
// So strip off everything that is not top 3 levels.
func getMainDomain(domain string) string {
	domain = strings.TrimSuffix(domain, ".")
	split := dns.Split(domain)
	if strings.HasSuffix(strings.ToLower(domain), "duckdns.org") {
		if len(split) < 3 {
			return ""
		}

		firstSubDomainIndex := split[len(split)-3]
		return domain[firstSubDomainIndex:]
	}

	return domain[split[len(split)-1]:]
}
