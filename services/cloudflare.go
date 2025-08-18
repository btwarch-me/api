package services

import (
	"btwarch/config"
	"context"
	"fmt"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/dns"
	"github.com/cloudflare/cloudflare-go/v4/option"
	"github.com/cloudflare/cloudflare-go/v4/zones"
)

type CloudflareService struct {
	client *cloudflare.Client
}

func NewCloudflareService(apiKey string) (*CloudflareService, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("cloudflare api key is required")
	}

	service := &CloudflareService{
		client: cloudflare.NewClient(option.WithAPIToken(apiKey)),
	}

	return service, nil
}

func (s *CloudflareService) AddTXTRecord(name string, content string) (*dns.RecordResponse, error) {
	cfg := config.LoadConfig()
	if cfg.CloudFlareZoneId == "" {
		return nil, fmt.Errorf("cloudflare zone id is required")
	}

	ctx := context.Background()

	_, err := s.client.Zones.Get(ctx, zones.ZoneGetParams{
		ZoneID: cloudflare.F(cfg.CloudFlareZoneId),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get cloudflare zone: %w", err)
	}

	record, err := s.client.DNS.Records.New(ctx, dns.RecordNewParams{
		ZoneID: cloudflare.F(cfg.CloudFlareZoneId),
		Body: dns.RecordNewParamsBody{
			Type:    cloudflare.F(dns.RecordNewParamsBodyTypeTXT),
			Name:    cloudflare.F(name),
			Content: cloudflare.F(content),
			TTL:     cloudflare.F(dns.TTL(1)),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create TXT record: %w", err)
	}

	return record, nil
}

func (s *CloudflareService) AddARecord(name string, content string) (*dns.RecordResponse, error) {
	cfg := config.LoadConfig()
	if cfg.CloudFlareZoneId == "" {
		return nil, fmt.Errorf("cloudflare zone id is required")
	}

	ctx := context.Background()

	_, err := s.client.Zones.Get(ctx, zones.ZoneGetParams{
		ZoneID: cloudflare.F(cfg.CloudFlareZoneId),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get cloudflare zone: %w", err)
	}

	record, err := s.client.DNS.Records.New(ctx, dns.RecordNewParams{
		ZoneID: cloudflare.F(cfg.CloudFlareZoneId),
		Body: dns.RecordNewParamsBody{
			Type:    cloudflare.F(dns.RecordNewParamsBodyTypeA),
			Name:    cloudflare.F(name),
			Content: cloudflare.F(content),
			TTL:     cloudflare.F(dns.TTL(1)),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create A record: %w", err)
	}

	return record, nil
}

func (s *CloudflareService) AddAAAARecord(name string, content string) (*dns.RecordResponse, error) {
	cfg := config.LoadConfig()
	if cfg.CloudFlareZoneId == "" {
		return nil, fmt.Errorf("cloudflare zone id is required")
	}

	ctx := context.Background()

	_, err := s.client.Zones.Get(ctx, zones.ZoneGetParams{
		ZoneID: cloudflare.F(cfg.CloudFlareZoneId),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get cloudflare zone: %w", err)
	}

	record, err := s.client.DNS.Records.New(ctx, dns.RecordNewParams{
		ZoneID: cloudflare.F(cfg.CloudFlareZoneId),
		Body: dns.RecordNewParamsBody{
			Type:    cloudflare.F(dns.RecordNewParamsBodyTypeAAAA),
			Name:    cloudflare.F(name),
			Content: cloudflare.F(content),
			TTL:     cloudflare.F(dns.TTL(1)),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AAAA record: %w", err)
	}

	return record, nil
}
