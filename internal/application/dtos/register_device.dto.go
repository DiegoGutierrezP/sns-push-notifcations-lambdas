package dtos

type RegisterDeviceRequest struct {
	DeviceToken        string `json:"deviceToken" validate:"required,min=10"`
	DeviceName         string `json:"deviceName" validate:"required,min=2"`
	ApplicationVersion string `json:"applicationVersion" validate:"required"`
	CalimacoId         int64  `json:"calimacoId" validate:"required"`
	OperatingSystem    string `json:"operatingSystem" validate:"required"`
	SystemVersion      string `json:"systemVersion" validate:"required"`
}

type RegisterDeviceResponse struct {
	DeviceId    string
	EndpointArn string
}
