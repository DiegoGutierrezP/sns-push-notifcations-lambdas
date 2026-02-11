package httpres

import (
	"encoding/json"
	"lmbd-digital-push-notifications/internal/application/dtos"
	appErrors "lmbd-digital-push-notifications/internal/application/errors"
	"maps"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

var defaultHeaders = map[string]string{
	"Content-Type":                 "application/json",
	"Access-Control-Allow-Origin":  "*",
	"Access-Control-Allow-Methods": "OPTIONS,GET,POST,PUT,DELETE",
}

// JSON builds an API Gateway proxy response with a JSON body and merged headers.
// It applies default JSON/CORS headers and overrides them with any provided headers.
// If marshaling fails, it returns a 500 response with a standard error payload.
func JSON(status int, payload any, headers map[string]string) events.APIGatewayProxyResponse {
	raw, err := json.Marshal(payload)
	if err != nil {
		raw, _ = json.Marshal(dtos.ApiResponse[any]{
			Success: false,
			Message: "Error serializando respuesta",
			Error:   err.Error(),
		})
		status = http.StatusInternalServerError
	}

	h := make(map[string]string, len(defaultHeaders)+(len(headers)))

	maps.Copy(h, defaultHeaders)
	maps.Copy(h, headers)

	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Headers:    h,
		Body:       string(raw),
	}
}

// Success returns a standardized successful JSON response using ApiResponse[T].
func Success[T any](status int, msg string, data *T) events.APIGatewayProxyResponse {
	return JSON(status, dtos.ApiResponse[T]{
		Success: true,
		Message: msg,
		Data:    data,
	}, nil)
}

// Fail returns a standardized error JSON response using ApiResponse.
func Fail(status int, msg string, err string) events.APIGatewayProxyResponse {
	return JSON(status, dtos.ApiResponse[any]{
		Success: false,
		Message: msg,
		Error:   err,
	}, nil)
}

func AppFail(err error) events.APIGatewayProxyResponse {
	if appErr, ok := appErrors.AsApplicationError(err); ok {
		return Fail(appErr.StatusCode, appErr.Message, appErr.Error())
	}
	return Fail(http.StatusInternalServerError, "internal server error", err.Error())
}
