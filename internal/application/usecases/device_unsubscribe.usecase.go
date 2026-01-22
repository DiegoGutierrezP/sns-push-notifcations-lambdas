package usecases

import (
	"context"
	"fmt"
	"lmbd-digital-push-notifications/internal/application/contracts/services"
	"lmbd-digital-push-notifications/internal/application/dtos"
)

type DeviceUnsubscribeUseCase struct {
	snsService services.ISnsService
}

func NewDeviceUnsubscribeUseCase(snsService services.ISnsService) *DeviceUnsubscribeUseCase {
	return &DeviceUnsubscribeUseCase{
		snsService: snsService,
	}
}

func (uc *DeviceUnsubscribeUseCase) Execute(ctx context.Context, data dtos.DeviceUnsubscribeRequest) error {
	// Lógica para enviar la notificación
	fmt.Print(data)
	return nil
}
