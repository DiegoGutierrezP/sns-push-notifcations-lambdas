package gateways

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"lmbd-digital-push-notifications/internal/application/contracts/services"
	"lmbd-digital-push-notifications/internal/shared/config"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type OptimoveGateway struct {
	client  *http.Client
	baseUrl string
	apiKey  string
}

func NewOptimoveGateway(cfg *config.Config) *OptimoveGateway {
	return &OptimoveGateway{
		client:  &http.Client{Timeout: 30 * time.Second},
		baseUrl: cfg.Optimove.ApiUrl,
		apiKey:  cfg.Optimove.ApiKey,
	}
}

func (s *OptimoveGateway) GetExecutedCampaignChannelDetails(ctx context.Context, campaignId int, channelId int) ([]services.ExecutedCampaignDetailsResponse, error) {

	endpoint, err := url.JoinPath(s.baseUrl, "/Actions/GetExecutedCampaignChannelDetails")
	if err != nil {
		return nil, fmt.Errorf("invalid base url: %w", err)
	}

	q := url.Values{}
	q.Set("campaignId", fmt.Sprint(campaignId))
	q.Set("channelId", fmt.Sprint(channelId))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint+"?"+q.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("x-api-key", s.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("error en request: %w", err)
	}
	defer resp.Body.Close()

	var result []services.ExecutedCampaignDetailsResponse

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return result, nil
}

func (s *OptimoveGateway) GetCampaignDetails(ctx context.Context, campaignId int) (*services.OptimoveCampaignDetails, error) {
	endpoint, err := url.JoinPath(s.baseUrl, "/Actions/GetCampaignDetails")
	if err != nil {
		return nil, fmt.Errorf("invalid base url: %w", err)
	}

	q := url.Values{}
	q.Set("campaignId", fmt.Sprint(campaignId))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint+"?"+q.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("x-api-key", s.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("error en request: %w", err)
	}
	defer resp.Body.Close()

	var result services.OptimoveCampaignDetails

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &result, nil
}

func (s *OptimoveGateway) GetCustomerExecutionDetailsByCampaign(ctx context.Context, campaignId, channelId int, top, skip *int, customerParams []string) (any, error) {
	endpoint, err := url.JoinPath(s.baseUrl, "/Actions/GetCustomerExecutionDetailsByCampaign")
	if err != nil {
		return nil, fmt.Errorf("invalid base url: %w", err)
	}

	defaultTop := 5000
	defaultSkip := 0

	if top == nil {
		top = &defaultTop
	}
	if skip == nil {
		skip = &defaultSkip
	}

	q := url.Values{}
	q.Set("campaignId", fmt.Sprint(campaignId))
	q.Set("channelId", fmt.Sprint(channelId))
	q.Set("top", fmt.Sprint(top))
	q.Set("skip", fmt.Sprint(skip))
	q.Set("customerAttributes", strings.Join(customerParams, ";"))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint+"?"+q.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("x-api-key", s.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("error en request: %w", err)
	}
	defer resp.Body.Close()

	var result []services.OptimoveCustomerDetailsResponse

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return result, nil
}

func (s *OptimoveGateway) UpdateCustomerPromotionStatus(ctx context.Context, data services.OptimoveUpdateCustomerDto) error {
	endpoint, err := url.JoinPath(s.baseUrl, "/Integrations/UpdateCustomerPromotionStatus")
	if err != nil {
		return fmt.Errorf("invalid base url: %w", err)
	}

	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("x-api-key", s.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)

	if err != nil {
		return fmt.Errorf("optimove request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return fmt.Errorf("optimove error: status=%d body=%s", resp.StatusCode, string(b))
	}

	return nil
}

func (s *OptimoveGateway) GetCustomerAttributes(ctx context.Context, customerId string) ([]services.OptimoveCustomerAttributeResponse, error) {

	endpoint, err := url.JoinPath(s.baseUrl, fmt.Sprintf("/Customers/%s/Attributes", customerId))
	if err != nil {
		return nil, fmt.Errorf("invalid base url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("x-api-key", s.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("error en request: %w", err)
	}
	defer resp.Body.Close()

	var result []services.OptimoveCustomerAttributeResponse

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return result, nil
}

func (s *OptimoveGateway) UpdateCustomerAttributes(ctx context.Context, data services.OptimoveUpdateCustomerAttributesDto) error {

	endpoint, err := url.JoinPath(s.baseUrl, "/Customers/UpdateCustomerAttributes")
	if err != nil {
		return fmt.Errorf("invalid base url: %w", err)
	}

	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("x-api-key", s.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)

	if err != nil {
		return fmt.Errorf("error en request: %w", err)
	}
	defer resp.Body.Close()

	return nil
}
