package aliyun

import (
	"strings"

	"github.com/abxuz/dns-manager/provider"

	alidns "github.com/alibabacloud-go/alidns-20150109/v4/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
)

const ProviderName = "aliyun"

type ProviderConfig struct {
	AccessKeyId     string `yaml:"access_key_id"`
	AccessKeySecret string `yaml:"access_key_secret"`
}

type Provider struct {
	domain string
	cfg    *ProviderConfig
	api    *alidns.Client
}

func init() {
	factory := provider.ProviderFactoryFunc(NewProvider)
	provider.RegisterProvider(ProviderName, factory)
}

func NewProvider(domain string, cfg provider.ProviderConfig) (provider.Provider, error) {
	p := &Provider{
		domain: domain,
		cfg:    &ProviderConfig{},
	}
	err := cfg.Unmarshal(p.cfg)
	if err != nil {
		return nil, err
	}

	p.api, err = alidns.NewClient(&openapi.Config{
		AccessKeyId:     tea.String(p.cfg.AccessKeyId),
		AccessKeySecret: tea.String(p.cfg.AccessKeySecret),
		Endpoint:        tea.String("alidns.aliyuncs.com"),
	})
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Provider) AddRecord(record *provider.Record) error {
	request := &alidns.AddDomainRecordRequest{
		DomainName: tea.String(record.Domain),
		RR:         tea.String(record.Host),
		Type:       tea.String(string(record.Type)),
		Value:      tea.String(record.Value),
	}

	if record.TTL > 0 {
		request.TTL = tea.Int64(int64(record.TTL))
	}
	if record.Type == provider.RRTypeMX {
		request.Priority = tea.Int64(int64(record.Priority))
	}

	result, err := p.api.AddDomainRecord(request)
	if err != nil {
		return err
	}

	if record.Remark != "" {
		p.SetRemark(tea.StringValue(result.Body.RecordId), record.Remark)
	}
	return nil
}

func (p *Provider) DeleteRecord(id string) error {
	request := &alidns.DeleteDomainRecordRequest{
		RecordId: tea.String(id),
	}
	_, err := p.api.DeleteDomainRecord(request)
	return err
}

func (p *Provider) ModifyRecord(record *provider.Record) error {
	request := &alidns.UpdateDomainRecordRequest{
		RecordId: tea.String(record.Id),
		RR:       tea.String(record.Host),
		Type:     tea.String(string(record.Type)),
		Value:    tea.String(record.Value),
	}

	if record.TTL > 0 {
		request.TTL = tea.Int64(int64(record.TTL))
	}
	if record.Type == provider.RRTypeMX {
		request.Priority = tea.Int64(int64(record.Priority))
	}

	_, err := p.api.UpdateDomainRecord(request)
	if err != nil && !strings.Contains(err.Error(), "DomainRecordDuplicate") {
		return err
	}

	p.SetRemark(record.Id, record.Remark)
	return nil
}

func (p *Provider) GetRecord(id string) (record *provider.Record, err error) {
	request := &alidns.DescribeDomainRecordInfoRequest{
		RecordId: tea.String(id),
	}
	result, err := p.api.DescribeDomainRecordInfo(request)
	if err != nil {
		return nil, err
	}
	record = ConvertToProviderRecord(result.Body)
	return
}

func (p *Provider) ListRecords(page provider.Page) (list []*provider.Record, total int, err error) {
	request := &alidns.DescribeDomainRecordsRequest{
		DomainName: tea.String(p.domain),
	}
	if page.Number > 0 {
		request.PageNumber = tea.Int64(int64(page.Number))
	}
	if page.Size > 0 {
		request.PageSize = tea.Int64(int64(page.Size))
	}

	result, err := p.api.DescribeDomainRecords(request)
	if err != nil {
		return
	}
	total = int(tea.Int64Value(result.Body.TotalCount))
	list = make([]*provider.Record, 0)
	for _, r := range result.Body.DomainRecords.Record {
		list = append(list, ConvertToProviderRecord(r))
	}
	return
}

func (p *Provider) SetRemark(id string, remark string) error {
	request := &alidns.UpdateDomainRecordRemarkRequest{
		RecordId: tea.String(id),
		Remark:   tea.String(remark),
	}
	_, err := p.api.UpdateDomainRecordRemark(request)
	return err
}

func ConvertToProviderRecord(r any) *provider.Record {
	switch r := r.(type) {
	case *alidns.DescribeDomainRecordsResponseBodyDomainRecordsRecord:
		return &provider.Record{
			Id:       tea.StringValue(r.RecordId),
			Host:     tea.StringValue(r.RR),
			Domain:   tea.StringValue(r.DomainName),
			Type:     provider.RRType(tea.StringValue(r.Type)),
			Value:    tea.StringValue(r.Value),
			Priority: int(tea.Int64Value(r.Priority)),
			TTL:      int(tea.Int64Value(r.TTL)),
			Remark:   tea.StringValue(r.Remark),
		}
	case *alidns.DescribeDomainRecordInfoResponseBody:
		return &provider.Record{
			Id:       tea.StringValue(r.RecordId),
			Host:     tea.StringValue(r.RR),
			Domain:   tea.StringValue(r.DomainName),
			Type:     provider.RRType(tea.StringValue(r.Type)),
			Value:    tea.StringValue(r.Value),
			Priority: int(tea.Int64Value(r.Priority)),
			TTL:      int(tea.Int64Value(r.TTL)),
			Remark:   tea.StringValue(r.Remark),
		}
	}
	return nil
}
