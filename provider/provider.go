package provider

import "fmt"

const (
	RRTypeA     RRType = "A"
	RRTypeAAAA  RRType = "AAAA"
	RRTypeCNAME RRType = "CNAME"
	RRTypeTXT   RRType = "TXT"
	RRTypeMX    RRType = "MX"
	RRTypeNS    RRType = "NS"
	RRTypeSRV   RRType = "SRV"
	RRTypeCAA   RRType = "CAA"
)

var (
	gProviderFactories map[string]ProviderFactory
)

type RRType string
type Record struct {
	Id       string `json:"id"`
	Host     string `json:"host" binding:"required"`
	Domain   string `json:"domain" binding:"required"`
	Type     RRType `json:"type" binding:"required"`
	Value    string `json:"value" binding:"required"`
	Priority int    `json:"priority,omitempty"`
	TTL      int    `json:"ttl" binding:"required"`
	Remark   string `json:"remark"`
}

type Page struct {
	Number int `form:"pageNumber"`
	Size   int `form:"pageSize"`
}

type Provider interface {
	AddRecord(record *Record) error
	DeleteRecord(id string) error
	ModifyRecord(record *Record) error
	GetRecord(id string) (record *Record, err error)
	ListRecords(page Page) (list []*Record, total int, err error)
}

type ProviderConfig interface {
	Unmarshal(v any) error
}

type ProviderFactory interface {
	NewProvider(domain string, c ProviderConfig) (provider Provider, err error)
}

type ProviderFactoryFunc func(domain string, c ProviderConfig) (provider Provider, err error)

func (f ProviderFactoryFunc) NewProvider(domain string, c ProviderConfig) (provider Provider, err error) {
	return f(domain, c)
}

func init() {
	gProviderFactories = make(map[string]ProviderFactory)
}

func RegisterProvider(name string, factory ProviderFactory) {
	if _, ok := gProviderFactories[name]; ok {
		panic(fmt.Sprintf("duplicate provider name %v", name))
	}
	gProviderFactories[name] = factory
}

func GetProviderFactory(name string) (factory ProviderFactory, ok bool) {
	factory, ok = gProviderFactories[name]
	return
}
