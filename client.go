package alidns

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/libdns/libdns"
)

//Client ...
type Client struct {
	client *alidns.Client
	mutex  sync.Mutex
}

func (p *Provider) newSession() error {
	var err error
	if p.client == nil {
		if len(p.RegionID) == 0 {
			p.RegionID = "cn-hangzhou"
		}
		p.client, err = alidns.NewClientWithAccessKey(p.RegionID, p.AccessKeyID, p.AccessKeySecret)
	}
	return err
}

func (p *Provider) addDNSRecord(ctx context.Context, zone string, record libdns.Record) (libdns.Record, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.newSession()

	req := alidns.CreateAddDomainRecordRequest()
	req.Scheme = "https"

	req.DomainName = zone
	req.RR = strings.TrimRight(strings.TrimSuffix(record.Name, zone), ".")
	req.Type = record.Type
	req.Value = record.Value

	response, err := p.client.AddDomainRecord(req)
	if err != nil {
		return record, fmt.Errorf("Add record error. Zone:%s, Name:%s, Error:%s, %v", zone, req.RR, err.Error(), record)
	}
	record.ID = response.RecordId
	return record, nil
}

func (p *Provider) delDNSRecord(ctx context.Context, record libdns.Record) (libdns.Record, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.newSession()

	req := alidns.CreateDeleteDomainRecordRequest()
	req.Scheme = "https"

	req.RecordId = record.ID
	_, err := p.client.DeleteDomainRecord(req)
	if err != nil {
		return record, fmt.Errorf("Delete record error. Error: %s, %v", err.Error(), record)
	}
	return record, nil
}

func (p *Provider) updateDNSRecord(ctx context.Context, zone string, record libdns.Record) (libdns.Record, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.newSession()

	req := alidns.CreateUpdateDomainRecordRequest()
	req.RecordId = record.ID
	req.RR = strings.TrimRight(strings.TrimSuffix(record.Name, zone), ".")
	req.Type = record.Type
	req.Value = record.Value

	_, err := p.client.UpdateDomainRecord(req)
	if err != nil {
		return record, fmt.Errorf("Update record error. Zone:%s, Error:%s, %v", zone, err.Error(), record)
	}
	return record, nil
}

func (p *Provider) getDNSRecords(ctx context.Context, zone string) ([]libdns.Record, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.newSession()
	var aliRecords []alidns.Record

	req := alidns.CreateDescribeDomainRecordsRequest()
	req.Scheme = "https"

	req.DomainName = zone
	pageNumber := 1

	for {
		req.PageNumber = requests.NewInteger(pageNumber)
		response, err := p.client.DescribeDomainRecords(req)
		if err != nil {
			return nil, fmt.Errorf("Get records error. Zone:%s, Error:%s", zone, err.Error())
		}
		aliRecords = append(aliRecords, response.DomainRecords.Record...)
		if response.PageNumber*response.PageSize >= response.TotalCount {
			break
		}
		pageNumber++
	}
	var records []libdns.Record
	for _, r := range aliRecords {
		record := libdns.Record{
			ID:    r.RecordId,
			Type:  r.Type,
			Name:  r.RR + "." + r.DomainName,
			Value: r.Value,
			TTL:   time.Duration(r.TTL) * time.Second,
		}
		records = append(records, record)
	}
	return records, nil
}
