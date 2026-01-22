package usecases

import (
	"context"
	"fmt"
	"lmbd-digital-push-notifications/internal/application/contracts/services"
	"lmbd-digital-push-notifications/internal/application/dtos"
)

type DeviceSubscribeUseCase struct {
	snsService services.ISnsService
}

func NewDeviceSubscribeUseCase(snsService services.ISnsService) *DeviceSubscribeUseCase {
	return &DeviceSubscribeUseCase{
		snsService: snsService,
	}
}

func (uc *DeviceSubscribeUseCase) Execute(ctx context.Context, data dtos.DeviceSubscribeRequest) error {
	// Lógica para enviar la notificación
	fmt.Print(data)
	return nil
}
