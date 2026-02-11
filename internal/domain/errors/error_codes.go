package errors

type ErrorCode string

const (

	// GEN
	GENInternal        ErrorCode = "GEN_0001"
	GENNotImplemented  ErrorCode = "GEN_0002"
	GENFeatureDisabled ErrorCode = "GEN_0003"

	// API
	APIInvalidRequestFormat ErrorCode = "API_2001"
	APIMissingParams        ErrorCode = "API_2002"
	APINotFound             ErrorCode = "API_2003"
	APIMethodNotAllowed     ErrorCode = "API_2004"
	APIUnsupportedMediaType ErrorCode = "API_2005"

	// CFG
	CFGEnvMissing   ErrorCode = "CFG_3001"
	CFGEnvInvalid   ErrorCode = "CFG_3002"
	CFGInconsistent ErrorCode = "CFG_3003"

	// RATE
	RATELimitExceeded    ErrorCode = "RATE_4001"
	RATEBurstExceeded    ErrorCode = "RATE_4002"
	RATETokenBucketEmpty ErrorCode = "RATE_4003"

	// DDB
	DDBConditionalFailed ErrorCode = "DDB_5001"
	DDBThrottled         ErrorCode = "DDB_5002"
	DDBItemNotFound      ErrorCode = "DDB_5003"
	DDBMarshalError      ErrorCode = "DDB_5004"
	DDBQueryFailed       ErrorCode = "DDB_5005"
	DDBInsertFailed      ErrorCode = "DDB_5006"
	DDBUpdateFailed      ErrorCode = "DDB_5007"

	// SNS
	SNSPublishFailed        ErrorCode = "SNS_6001"
	SNSInvalidArn           ErrorCode = "SNS_6002"
	SNSAccessDenied         ErrorCode = "SNS_6003"
	SNSCreateEndpointFailed ErrorCode = "SNS_6004"
	SNSUpdateEndpointFailed ErrorCode = "SNS_6005"

	// OPTI
	OPTIInvalidApiKey ErrorCode = "OPTI_7001"
	OPTIUnavailable   ErrorCode = "OPTI_7002"
	OPTIBadPayload    ErrorCode = "OPTI_7003"

	// SUBS
	SUBSAlreadySubscribed ErrorCode = "SUBS_8001"
	SUBSNotSubscribed     ErrorCode = "SUBS_8002"
	SUBSInvalidState      ErrorCode = "SUBS_8003"

	// DEV
	DEVNotFound           ErrorCode = "DEV_9001"
	DEVUserWithoutDevices ErrorCode = "DEV_9002"
	DEVTokenAlreadyExists ErrorCode = "DEV_9003"
)
