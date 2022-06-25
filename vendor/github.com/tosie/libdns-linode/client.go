package linode

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"time"

	"fmt"
	"net/http"

	"github.com/linode/linodego"
	"golang.org/x/oauth2"

	"github.com/libdns/libdns"
)

type Client struct {
	client *linodego.Client
	mutex  sync.Mutex
}

func (p *Provider) getClient() error {
	if p.client == nil {
		tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: p.APIToken})

		oauth2Client := &http.Client{
			Transport: &oauth2.Transport{
				Source: tokenSource,
			},
		}

		client := linodego.NewClient(oauth2Client)
		p.client = &client
	}

	return nil
}

func (p *Provider) getDomainId(ctx context.Context, zone string) (int, error) {
	opt := &linodego.ListOptions{}

	domains, err := p.client.ListDomains(ctx, opt)
	if err != nil {
		return 0, err
	}

	var domain *linodego.Domain
	for _, d := range domains {
		if d.Domain == zone {
			domain = &d
			break
		}
	}

	if domain == nil {
		return 0, fmt.Errorf("Did not find a zone with the name '%s'", zone)
	}

	return domain.ID, nil
}

func (p *Provider) getDNSEntries(ctx context.Context, zone string) ([]libdns.Record, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.getClient()

	domainId, err := p.getDomainId(ctx, zone)
	if err != nil {
		return nil, err
	}

	// opt := &linodego.ListOptions{}
	entries, err := p.client.ListDomainRecords(ctx, domainId, nil)
	if err != nil {
		return nil, err
	}

	var records []libdns.Record
	for _, entry := range entries {
		record := libdns.Record{
			Name:  entry.Name, // + "." + strings.Trim(zone, ".") + ".",
			Value: entry.Target,
			Type:  string(entry.Type),
			TTL:   time.Duration(entry.TTLSec) * time.Second,
			ID:    strconv.Itoa(entry.ID),
		}
		records = append(records, record)
	}

	return records, nil
}

func (p *Provider) addDNSEntry(ctx context.Context, zone string, record libdns.Record) (libdns.Record, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.getClient()

	domainId, err := p.getDomainId(ctx, zone)
	if err != nil {
		return libdns.Record{}, err
	}

	entry := linodego.DomainRecordCreateOptions{
		Name:   strings.Trim(strings.ReplaceAll(record.Name, zone, ""), "."),
		Target: record.Value,
		Type:   linodego.DomainRecordType(record.Type),
		TTLSec: int(record.TTL.Seconds()),
	}

	rec, err := p.client.CreateDomainRecord(ctx, domainId, entry)
	if err != nil {
		return record, err
	}

	record = p.mergeRecordWithLinodeDomainRecord(record, rec)

	return record, nil
}

func (p *Provider) removeDNSEntry(ctx context.Context, zone string, record libdns.Record) (libdns.Record, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.getClient()

	domainId, err := p.getDomainId(ctx, zone)
	if err != nil {
		return record, err
	}

	id, err := strconv.Atoi(record.ID)
	if err != nil {
		return record, err
	}

	err = p.client.DeleteDomainRecord(ctx, domainId, id)
	if err != nil {
		return record, err
	}

	return record, nil
}

func (p *Provider) updateDNSEntry(ctx context.Context, zone string, record libdns.Record) (libdns.Record, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.getClient()

	domainId, err := p.getDomainId(ctx, zone)
	if err != nil {
		return record, err
	}

	id, err := strconv.Atoi(record.ID)
	if err != nil {
		// If the ID field is empty, just add a new record instead (do not call p.addDNSEntry here, as it will
		// deadlock because of the mutex).
		if record.ID == "" {
			entry := linodego.DomainRecordCreateOptions{
				Name:   strings.Trim(strings.ReplaceAll(record.Name, zone, ""), "."),
				Target: record.Value,
				Type:   linodego.DomainRecordType(record.Type),
				TTLSec: int(record.TTL.Seconds()),
			}

			rec, err := p.client.CreateDomainRecord(ctx, domainId, entry)
			if err != nil {
				return record, err
			}

			record = p.mergeRecordWithLinodeDomainRecord(record, rec)

			return record, nil
		} else {
			return record, err
		}
	}

	entry := linodego.DomainRecordUpdateOptions{
		Name:   strings.Trim(strings.ReplaceAll(record.Name, zone, ""), "."),
		Target: record.Value,
		Type:   linodego.DomainRecordType(record.Type),
		TTLSec: int(record.TTL.Seconds()),
	}

	rec, err := p.client.UpdateDomainRecord(ctx, domainId, id, entry)
	if err != nil {
		return record, err
	}

	record = p.mergeRecordWithLinodeDomainRecord(record, rec)

	return record, nil
}

func (p *Provider) mergeRecordWithLinodeDomainRecord(record libdns.Record, rec *linodego.DomainRecord) libdns.Record {
	record.ID = strconv.Itoa(rec.ID)
	record.Name = rec.Name
	record.Value = rec.Target
	record.Type = string(rec.Type)
	record.TTL = time.Duration(rec.TTLSec) * time.Second

	return record
}
