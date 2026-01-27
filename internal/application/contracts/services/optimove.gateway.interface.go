package services

import (
	"context"
)

type ExecutedCampaignDetailsResponse struct {
	SendID        string
	TemplateID    string
	ScheduledTime string
}

type OptimoveCampaignDetails struct {
	TestGroupSize    int
	ControlGroupSize int
	Channels         []map[string]any
	Tags             []string
	TargetGroupID    int
	CampaignType     string
	Duration         int
	LeadTime         int
	Notes            string
	IsMultiChannel   bool
	IsRecurrence     bool
	Error            string
}

type OptimoveCustomerDetailsResponse struct {
	CustomerI          string
	ActionI            int
	ChannelI           int
	TemplateI          string
	TemplateNam        string
	ScheduleTime       any
	PromoCod           string
	CustomerAttributes []string
}

type OptimoveUpdateCustomerDto struct {
	TimeStamp     *string
	CampaignID    int
	PromosUpdated []OptimovePromoUpdateDto
}

type OptimovePromoUpdateDto struct {
	CustomerId string
	Status     bool
}

type OptimoveCustomerAttributeResponse struct {
	RealFieldName string
	Alias         string
	FieldType     string
	Value         any
}

type OptimoveCustomerAttributeValue struct {
	RealFieldName string
	Value         any
}

type OptimoveCustomerNewAttributesValues struct {
	CustomerID string
	Attributes []OptimoveCustomerAttributeValue
}

type OptimoveUpdateCustomerAttributesDto struct {
	CustomerNewAttributesValuesList []OptimoveCustomerNewAttributesValues
	CallbackURL                     *string
}

type IOptimoveGateway interface {
	GetExecutedCampaignChannelDetails(ctx context.Context, campaignId int, channelId int) ([]ExecutedCampaignDetailsResponse, error)
}
