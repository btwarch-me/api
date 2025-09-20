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

func NewCloudflareService(apiToken string) (*CloudflareService, error) {
	if apiToken == "" {
		return nil, fmt.Errorf("cloudflare api token is required")
	}

	service := &CloudflareService{
		client: cloudflare.NewClient(option.WithAPIToken(apiToken)),
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

func (s *CloudflareService) AddCNAMERecord(name string, content string) (*dns.RecordResponse, error) {
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
			Type:    cloudflare.F(dns.RecordNewParamsBodyTypeCNAME),
			Name:    cloudflare.F(name),
			Content: cloudflare.F(content),
			TTL:     cloudflare.F(dns.TTL(1)),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create CNAME record: %w", err)
	}

	return record, nil
}

func (s *CloudflareService) DeleteRecordByID(recordID string) (*dns.RecordDeleteResponse, error) {
	cfg := config.LoadConfig()
	if cfg.CloudFlareZoneId == "" {
		return nil, fmt.Errorf("cloudflare zone id is required")
	}
	ctx := context.Background()
	record, err := s.client.DNS.Records.Delete(ctx, recordID, dns.RecordDeleteParams{ZoneID: cloudflare.F(cfg.CloudFlareZoneId)})
	if err != nil {
		return nil, fmt.Errorf("failed to delete record: %w", err)
	}
	return record, nil
}

func (s *CloudflareService) UpdateARecord(recordID string, name string, content string) (*dns.RecordResponse, error) {
	cfg := config.LoadConfig()
	if cfg.CloudFlareZoneId == "" {
		return nil, fmt.Errorf("cloudflare zone id is required")
	}
	ctx := context.Background()
	record, err := s.client.DNS.Records.Update(ctx, recordID, dns.RecordUpdateParams{
		ZoneID: cloudflare.F(cfg.CloudFlareZoneId),
		Body: dns.RecordUpdateParamsBody{
			Type:    cloudflare.F(dns.RecordUpdateParamsBodyTypeA),
			Name:    cloudflare.F(name),
			Content: cloudflare.F(content),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update A record: %w", err)
	}
	return record, nil
}

func (s *CloudflareService) UpdateAAAARecord(recordID string, name string, content string) (*dns.RecordResponse, error) {
	cfg := config.LoadConfig()
	if cfg.CloudFlareZoneId == "" {
		return nil, fmt.Errorf("cloudflare zone id is required")
	}
	ctx := context.Background()
	record, err := s.client.DNS.Records.Update(ctx, recordID, dns.RecordUpdateParams{
		ZoneID: cloudflare.F(cfg.CloudFlareZoneId),
		Body: dns.RecordUpdateParamsBody{
			Type:    cloudflare.F(dns.RecordUpdateParamsBodyTypeAAAA),
			Name:    cloudflare.F(name),
			Content: cloudflare.F(content),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update AAAA record: %w", err)
	}
	return record, nil
}

func (s *CloudflareService) UpdateCNAMERecord(recordID string, name string, content string) (*dns.RecordResponse, error) {
	cfg := config.LoadConfig()
	if cfg.CloudFlareZoneId == "" {
		return nil, fmt.Errorf("cloudflare zone id is required")
	}
	ctx := context.Background()
	record, err := s.client.DNS.Records.Update(ctx, recordID, dns.RecordUpdateParams{
		ZoneID: cloudflare.F(cfg.CloudFlareZoneId),
		Body: dns.RecordUpdateParamsBody{
			Type:    cloudflare.F(dns.RecordUpdateParamsBodyTypeCNAME),
			Name:    cloudflare.F(name),
			Content: cloudflare.F(content),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update CNAME record: %w", err)
	}
	return record, nil
}

func (s *CloudflareService) UpdateTXTRecord(recordID string, name string, content string) (*dns.RecordResponse, error) {
	cfg := config.LoadConfig()
	if cfg.CloudFlareZoneId == "" {
		return nil, fmt.Errorf("cloudflare zone id is required")
	}
	ctx := context.Background()
	record, err := s.client.DNS.Records.Update(ctx, recordID, dns.RecordUpdateParams{
		ZoneID: cloudflare.F(cfg.CloudFlareZoneId),
		Body: dns.RecordUpdateParamsBody{
			Type:    cloudflare.F(dns.RecordUpdateParamsBodyTypeTXT),
			Name:    cloudflare.F(name),
			Content: cloudflare.F(content),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update TXT record: %w", err)
	}
	return record, nil
}
