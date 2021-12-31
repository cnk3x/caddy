package westcn

import (
	"context"
	"sync"

	"github.com/libdns/libdns"
)

// Provider implements the libdns interfaces for west.cn
type Provider struct {
	Endpoint string
	Username string
	Password string

	client *Client
	once   sync.Once
}

func (p *Provider) getClient() *Client {
	p.once.Do(func() {
		if p.Endpoint == "" {
			p.Endpoint = "https://api.west.cn/api/v2"
		}
		p.client = &Client{
			Endpoint: p.Endpoint,
			Username: p.Username,
			Password: p.Password,
		}
	})
	return p.client
}

// GetRecords lists all the records in the zone.
func (p *Provider) GetRecords(ctx context.Context, zone string) (records []libdns.Record, err error) {
	return p.getClient().GetRecords(ctx, zone)
}

// AppendRecords adds records to the zone. It returns the records that were added.
func (p *Provider) AppendRecords(ctx context.Context, zone string, records []libdns.Record) (appendedRecords []libdns.Record, err error) {
	client := p.getClient()
	for _, record := range records {
		if record, err = client.AddRecord(ctx, zone, record); err != nil {
			return
		}
		appendedRecords = append(appendedRecords, record)
	}
	return
}

// DeleteRecords deletes the records from the zone.
func (p *Provider) DeleteRecords(ctx context.Context, zone string, records []libdns.Record) (deletedRecords []libdns.Record, err error) {
	client := p.getClient()
	for _, record := range records {
		if record, err = client.DeleteRecord(ctx, zone, record); err != nil {
			return
		}
		deletedRecords = append(deletedRecords, record)
	}
	return
}

// SetRecords sets the records in the zone, either by updating existing records
// or creating new ones. It returns the updated records.
func (p *Provider) SetRecords(ctx context.Context, zone string, records []libdns.Record) (setRecords []libdns.Record, err error) {
	client := p.getClient()
	for _, record := range records {
		if record.ID != "" {
			record, err = client.AddRecord(ctx, zone, record)
		} else {
			record, err = client.UpdateRecord(ctx, zone, record)
		}
		if err != nil {
			return
		}
		setRecords = append(setRecords, record)
	}
	return
}

// Interface guards
var (
	_ libdns.RecordGetter   = (*Provider)(nil)
	_ libdns.RecordAppender = (*Provider)(nil)
	_ libdns.RecordSetter   = (*Provider)(nil)
	_ libdns.RecordDeleter  = (*Provider)(nil)
)
