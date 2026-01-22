package dtos

type DeviceUnsubscribeRequest struct {
	DeviceToken     *string
	SubscriptionArn string
}
