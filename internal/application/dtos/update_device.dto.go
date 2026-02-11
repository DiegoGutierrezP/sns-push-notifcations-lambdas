package dtos

type UpdateDeviceRequest struct {
	DeviceId           string  `json:"deviceId" validate:"required"`
	DeviceName         *string `json:"deviceName"`
	DeviceToken        *string `json:"deviceToken"`
	ApplicationVersion *string `json:"applicationVersion" `
	CalimacoId         *int    `json:"calimacoId" `
	OperatingSystem    *string `json:"operatingSystem" `
	SystemVersion      *string `json:"systemVersion" `
}

type UpdateDeviceResponse struct {
	DeviceId string `json:"deviceId"`
}
