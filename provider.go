package alidns

import (
	"context"

	"github.com/libdns/libdns"
)

// Provider implements the libdns interfaces for AliDNS.
type Provider struct {
	AccessKeyID     string `json:"accesskey_id"`
	AccessKeySecret string `json:"accesskey_secret"`
	RegionID        string `json:"region_id,omitempty"`
	Client
}

// AppendRecords adds records to the zone. It returns the records that were added.
func (p *Provider) AppendRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	var createdRecords []libdns.Record
	for _, record := range records {
		r, err := p.addDNSRecord(ctx, zone, record)
		if err != nil {
			return nil, err
		}
		createdRecords = append(createdRecords, r)
	}
	return createdRecords, nil
}

// DeleteRecords deletes the records from the zone.
func (p *Provider) DeleteRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	var deletedRecords []libdns.Record
	for _, record := range records {
		r, err := p.delDNSRecord(ctx, record)
		if err != nil {
			return nil, err
		}
		deletedRecords = append(deletedRecords, r)
	}
	return deletedRecords, nil

}

// GetRecords lists all the records in the zone.
func (p *Provider) GetRecords(ctx context.Context, zone string) ([]libdns.Record, error) {
	records, err := p.getDNSRecords(ctx, zone)
	if err != nil {
		return nil, err
	}
	return records, nil
}

// SetRecords sets the records in the zone, either by updating existing records
// or creating new ones. It returns the updated records.
func (p *Provider) SetRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	var setRecords []libdns.Record
	for _, record := range records {
		r, err := p.updateDNSRecord(ctx, zone, record)
		if err != nil {
			return nil, err
		}
		setRecords = append(setRecords, r)
	}
	return setRecords, nil
}

// Interface guards
var (
	_ libdns.RecordGetter   = (*Provider)(nil)
	_ libdns.RecordAppender = (*Provider)(nil)
	_ libdns.RecordSetter   = (*Provider)(nil)
	_ libdns.RecordDeleter  = (*Provider)(nil)
)
