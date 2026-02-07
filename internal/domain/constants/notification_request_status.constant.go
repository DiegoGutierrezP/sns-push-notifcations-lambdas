package constants

type NotificationRequestStatus int

const (
	NotificationRequestStatusPending NotificationRequestStatus = 0
	NotificationRequestStatusSent    NotificationRequestStatus = 1
	NotificationRequestStatusFailed  NotificationRequestStatus = 2
)
