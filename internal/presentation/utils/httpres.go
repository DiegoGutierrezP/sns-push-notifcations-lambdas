package httpres

import (
	"encoding/json"
	"lmbd-digital-push-notifications/internal/application/dtos"
	"maps"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

var defaultHeaders = map[string]string{
	"Content-Type":                 "application/json",
	"Access-Control-Allow-Origin":  "*",
	"Access-Control-Allow-Methods": "OPTIONS,GET,POST,PUT,DELETE",
}

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
	// for k, v := range defaultHeaders {
	// 	h[k] = v
	// }
	// for k, v := range headers {
	// 	h[k] = v
	// }
	maps.Copy(h, defaultHeaders)
	maps.Copy(h, headers)

	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Headers:    h,
		Body:       string(raw),
	}
}

func Success[T any](status int, msg string, data *T) events.APIGatewayProxyResponse {
	return JSON(status, dtos.ApiResponse[T]{
		Success: true,
		Message: msg,
		Data:    data,
	}, nil)
}

func Fail(status int, msg string, err string) events.APIGatewayProxyResponse {
	return JSON(status, dtos.ApiResponse[any]{
		Success: false,
		Message: msg,
		Error:   err,
	}, nil)
}
