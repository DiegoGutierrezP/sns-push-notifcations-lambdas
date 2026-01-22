package usecases

import (
	"context"
	"fmt"
	"lmbd-digital-push-notifications/internal/application/contracts/services"
	"lmbd-digital-push-notifications/internal/application/dtos"
)

type PublishNotificationUseCase struct {
	snsService services.ISnsService
}

func NewPublishNotificationUseCase(snsService services.ISnsService) *PublishNotificationUseCase {
	return &PublishNotificationUseCase{
		snsService: snsService,
	}
}

func (uc *PublishNotificationUseCase) Execute(ctx context.Context, data dtos.PublishNotificationRequest) error {
	// Lógica para enviar la notificación
	fmt.Print(data)
	return nil
}
