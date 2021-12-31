package vultr

import (
	"context"
	"strconv"
	"sync"
	"time"

	"golang.org/x/oauth2"

	"github.com/libdns/libdns"
	"github.com/vultr/govultr/v2"
)

type Client struct {
	vultr *govultr.Client
	mutex sync.Mutex
}

func (p *Provider) getClient() error {
	if p.client.vultr == nil {
		oauth_cfg := &oauth2.Config{}
		oauth_token_source := oauth_cfg.TokenSource(nil, &oauth2.Token{AccessToken: p.APIToken})

		p.client.vultr = govultr.NewClient(oauth2.NewClient(nil, oauth_token_source))
	}

	return nil
}

func (p *Provider) getDNSEntries(ctx context.Context, domain string) ([]libdns.Record, error) {
	p.client.mutex.Lock()
	defer p.client.mutex.Unlock()

	p.getClient()

	listOptions := &govultr.ListOptions{}

	var records []libdns.Record
	for {
		dns_entries, meta, err := p.client.vultr.DomainRecord.List(ctx, domain, listOptions)
		if err != nil {
			return records, err
		}

		for _, entry := range dns_entries {
			record := libdns.Record{
				Name:  entry.Name,
				Value: entry.Data,
				Type:  entry.Type,
				TTL:   time.Duration(entry.TTL) * time.Second,
				ID:    entry.ID,
			}
			records = append(records, record)
		}

		if meta.Links.Next == "" {
			break
		}

		listOptions.Cursor = meta.Links.Next
	}

	return records, nil
}

func (p *Provider) addDNSRecord(ctx context.Context, domain string, record libdns.Record) (libdns.Record, error) {
	p.client.mutex.Lock()
	defer p.client.mutex.Unlock()

	p.getClient()

	domainRecordReq := &govultr.DomainRecordReq{
		Name: record.Name,
		Type: record.Type,
		Data: strconv.Quote(record.Value),
		TTL:  int(record.TTL.Seconds()),
	}

	rec, err := p.client.vultr.DomainRecord.Create(ctx, domain, domainRecordReq)
	if err != nil {
		return record, err
	}

	record.ID = rec.ID

	return record, nil
}

func (p *Provider) removeDNSRecord(ctx context.Context, domain string, record libdns.Record) (libdns.Record, error) {
	p.client.mutex.Lock()
	defer p.client.mutex.Unlock()

	p.getClient()

	err := p.client.vultr.DomainRecord.Delete(ctx, domain, record.ID)
	if err != nil {
		return record, err
	}

	return record, nil
}

func (p *Provider) updateDNSRecord(ctx context.Context, domain string, record libdns.Record) (libdns.Record, error) {
	p.client.mutex.Lock()
	defer p.client.mutex.Unlock()

	p.getClient()

	domainRecordReq := &govultr.DomainRecordReq{
		Name: record.Name,
		Type: record.Type,
		Data: strconv.Quote(record.Value),
		TTL:  int(record.TTL.Seconds()),
	}

	err := p.client.vultr.DomainRecord.Update(ctx, domain, record.ID, domainRecordReq)
	if err != nil {
		return record, err
	}

	return record, nil
}
