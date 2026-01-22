package usecases

import (
	"context"
	"errors"
	"lmbd-digital-push-notifications/internal/application/contracts/repositories"
	"lmbd-digital-push-notifications/internal/application/contracts/services"
	"lmbd-digital-push-notifications/internal/application/dtos"
	"lmbd-digital-push-notifications/internal/domain/entities"
	"lmbd-digital-push-notifications/internal/shared/config"
	"strconv"
)

type RegisterDeviceUseCase struct {
	snsService       services.ISnsService
	deviceRepository repositories.IDeviceRepository
	platformAppArn   string
}

func NewRegisterDeviceUseCase(
	config config.Config,
	snsService services.ISnsService,
	deviceRepository repositories.IDeviceRepository,
) *RegisterDeviceUseCase {
	return &RegisterDeviceUseCase{
		snsService:       snsService,
		deviceRepository: deviceRepository,
		platformAppArn:   config.Sns.PlatformAppArn,
	}
}

func (uc *RegisterDeviceUseCase) Execute(ctx context.Context, request dtos.RegisterDeviceRequest) (*dtos.RegisterDeviceResponse, error) {

	exists, err := uc.deviceRepository.ExistsByToken(ctx, request.DeviceToken)

	if err != nil {
		return nil, errors.New("An Error ocurred ")
	}

	if exists {
		return nil, errors.New("token already registered")
	}

	// create endpoint arn for new devices
	endpointArn, err := uc.snsService.CreateEndpoint(ctx, request.DeviceToken, uc.platformAppArn)

	if err != nil {
		return nil, errors.New("Ocurrio un error al crear el endpoint ARN")
	}

	// create new device instance
	deviceEntity := entities.NewDevice(
		request.DeviceToken,
		endpointArn,
		request.DeviceName,
		request.ApplicationVersion,
		strconv.Itoa(int(request.CalimacoId)),
		"",
		request.OperatingSystem,
		request.SystemVersion,
	)

	if err := uc.deviceRepository.Save(ctx, deviceEntity); err != nil {
		return nil, errors.New("Ocurrio un error al registrar el dispositivo")
	}

	return &dtos.RegisterDeviceResponse{
		DeviceId:    deviceEntity.ID.String(),
		DeviceToken: deviceEntity.DeviceToken,
		EndpointArn: deviceEntity.EndpointArn,
	}, nil
}
