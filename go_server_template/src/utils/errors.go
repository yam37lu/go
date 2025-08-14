package utils

import (
	"fmt"

	"github.com/pkg/errors"
)

type ErrorType string

const (
	Normal                        ErrorType = "10000"
	InvaildUserKey                ErrorType = "10001"
	ServiceNotAvailable           ErrorType = "10002"
	DailyQueryOverLimit           ErrorType = "10003"
	AccessTooFrequent             ErrorType = "10004"
	InvaildUserIp                 ErrorType = "10005"
	InvaildUserDomain             ErrorType = "10006"
	InvaildUserSignature          ErrorType = "10007"
	InvaildUserScode              ErrorType = "10008"
	UserKeyPlatNoMatch            ErrorType = "10009"
	IpQueryOverLimit              ErrorType = "10010"
	NotSupportHttps               ErrorType = "10011"
	InsufficientPrivileges        ErrorType = "10012"
	UserKeyRecycled               ErrorType = "10013"
	QpsHasExceededTheLimit        ErrorType = "10014"
	GatewayTimeout                ErrorType = "10015"
	ServerIsBuy                   ErrorType = "10016"
	ResourceUnavailable           ErrorType = "10017"
	CQpsHasExceededTheLimit       ErrorType = "10019"
	CKQpsHasExceededTheLimit      ErrorType = "10020"
	CUQpsHasExceededTheLimit      ErrorType = "10021"
	InvaildRequest                ErrorType = "10026"
	AbroadDailyQueryOverLimit     ErrorType = "10029"
	UserDailyQueryOverLimit       ErrorType = "10044"
	UserAbroadDailyQueryOverLimit ErrorType = "10045"

	InvaildParams                ErrorType = "20000"
	MissingRequiredParams        ErrorType = "20001"
	IllegalRequest               ErrorType = "20002"
	UnknownError                 ErrorType = "20003"
	InsufficientAbroadPrivileges ErrorType = "20011"
	IllegalContent               ErrorType = "20012"
	OutOfService                 ErrorType = "20800"
	NoRoadsNearby                ErrorType = "20801"
	RouteFail                    ErrorType = "20802"
	OverDirectionRange           ErrorType = "20803"
	ParamsConflict               ErrorType = "20004"

	QuotaPlanRunOut         ErrorType = "40000"
	GeofenceMaxCountReached ErrorType = "40001"
	ServiceExpired          ErrorType = "40002"
	AbroadQuotaPlanRunOut   ErrorType = "40003"

	UpstreamServiceRequestError  ErrorType = "60001"
	UpstreamServiceResponseError ErrorType = "60002"
	ServiceInternalError         ErrorType = "60003"
	ServiceResponseDataEmpty     ErrorType = "60004"

	DataAlreadyExist   ErrorType = "70001"
	DataOperateError   ErrorType = "70002"
	DataQueryError     ErrorType = "70003"
	DataIsNotExist     ErrorType = "70004"
	SqlOperateError    ErrorType = "70005"
	UserIsNotLoginedIn ErrorType = "70006"
	UserIsNotExist     ErrorType = "70007"
	DataIsError        ErrorType = "70008"
	OperateUnsupported ErrorType = "70009"
	ApproveUnBinded    ErrorType = "70010"

	DistrictBinded               ErrorType = "80001"
	OperateObjectHasBindRelation ErrorType = "80002"

	OtherUnknownError ErrorType = "99999"
	DetailMessage     ErrorType = "-1"
)

var CodeDef map[ErrorType]string

func init() {
	CodeDef = make(map[ErrorType]string)

	CodeDef[Normal] = "ok"
	CodeDef[InvaildUserKey] = "INVALID_USER_KEY"
	CodeDef[ServiceNotAvailable] = "SERVICE_NOT_AVAILABLE"
	CodeDef[DailyQueryOverLimit] = "DAILY_QUERY_OVER_LIMIT"
	CodeDef[AccessTooFrequent] = "ACCESS_TOO_FREQUENT"
	CodeDef[InvaildUserIp] = "INVALID_USER_IP"
	CodeDef[InvaildUserDomain] = "INVALID_USER_DOMAIN"
	CodeDef[InvaildUserSignature] = "INVALID_USER_SIGNATURE"
	CodeDef[InvaildUserScode] = "INVALID_USER_SCODE"
	CodeDef[UserKeyPlatNoMatch] = "USERKEY_PLAT_NOMATCH"
	CodeDef[IpQueryOverLimit] = "IP_QUERY_OVER_LIMIT"
	CodeDef[NotSupportHttps] = "NOT_SUPPORT_HTTPS"
	CodeDef[InsufficientPrivileges] = "INSUFFICIENT_PRIVILEGES"
	CodeDef[UserKeyRecycled] = "USER_KEY_RECYCLED"
	CodeDef[QpsHasExceededTheLimit] = "QPS_HAS_EXCEEDED_THE_LIMIT"
	CodeDef[GatewayTimeout] = "GATEWAY_TIMEOUT"
	CodeDef[ServerIsBuy] = "SERVER_IS_BUSY"
	CodeDef[ResourceUnavailable] = "RESOURCE_UNAVAILABLE"
	CodeDef[CQpsHasExceededTheLimit] = "CQPS_HAS_EXCEEDED_THE_LIMIT"
	CodeDef[CKQpsHasExceededTheLimit] = "CKQPS_HAS_EXCEEDED_THE_LIMIT"
	CodeDef[CUQpsHasExceededTheLimit] = "CUQPS_HAS_EXCEEDED_THE_LIMIT"
	CodeDef[InvaildRequest] = "INVALID_REQUEST"
	CodeDef[AbroadDailyQueryOverLimit] = "ABROAD_DAILY_QUERY_OVER_LIMIT"
	CodeDef[UserDailyQueryOverLimit] = "USER_DAILY_QUERY_OVER_LIMIT"
	CodeDef[UserAbroadDailyQueryOverLimit] = "USER_ABROAD_DAILY_QUERY_OVER_LIMIT"
	CodeDef[InvaildParams] = "INVALID_PARAMS"
	CodeDef[MissingRequiredParams] = "MISSING_REQUIRED_PARAMS"
	CodeDef[IllegalRequest] = "ILLEGAL_REQUEST"
	CodeDef[UnknownError] = "UNKNOWN_ERROR"
	CodeDef[InsufficientAbroadPrivileges] = "INSUFFICIENT_ABROAD_PRIVILEGES"
	CodeDef[IllegalContent] = "ILLEGAL_CONTENT"
	CodeDef[OutOfService] = "OUT_OF_SERVICE"
	CodeDef[NoRoadsNearby] = "NO_ROADS_NEARBY"
	CodeDef[RouteFail] = "ROUTE_FAIL"
	CodeDef[OverDirectionRange] = "OVER_DIRECTION_RANGE"
	CodeDef[QuotaPlanRunOut] = "QUOTA_PLAN_RUN_OUT"
	CodeDef[GeofenceMaxCountReached] = "GEOFENCE_MAX_COUNT_REACHED"
	CodeDef[ServiceExpired] = "SERVICE_EXPIRED"
	CodeDef[AbroadQuotaPlanRunOut] = "ABROAD_QUOTA_PLAN_RUN_OUT"
	CodeDef[ParamsConflict] = "PARAMS_CONFLICT"

	CodeDef[UpstreamServiceRequestError] = "UPSTREAM_SERVICE_REQUEST_ERROR"
	CodeDef[UpstreamServiceResponseError] = "UPSTREAM_SERVICE_RESPONSE_ERROR"
	CodeDef[ServiceInternalError] = "SERVICE_INTERNAL_ERROR"
	CodeDef[ServiceResponseDataEmpty] = "SERVICE_RESPONSE_DATA_EMPTY"

	CodeDef[DataAlreadyExist] = "DATA_ALREADY_EXIST"
	CodeDef[DataOperateError] = "DATA_OPERATE_ERROR"
	CodeDef[DataQueryError] = "DATA_QUERY_ERROR"
	CodeDef[DataIsNotExist] = "DATA_IS_NOT_EXIST"
	CodeDef[SqlOperateError] = "SQL_OPERATE_ERROR"
	CodeDef[UserIsNotLoginedIn] = "USER_IS_NOT_LOGINED_IN"
	CodeDef[UserIsNotExist] = "USER_IS_NOT_EXIST"
	CodeDef[DataIsError] = "DATA_IS_ERROR"
	CodeDef[OperateUnsupported] = "OPERATE_UNSUPPORTED"
	CodeDef[ApproveUnBinded] = "APPROVE_UNBINDED"
	CodeDef[DistrictBinded] = "DISTRICT_BINDED"

	CodeDef[OtherUnknownError] = "OTHER_UNKNOWN_ERROR"
}

type customError struct {
	errorType     ErrorType
	originalError error
	context       errorContext
}

type errorContext struct {
	Field   string
	Message string
}

// New creates a new customError
func (errorType ErrorType) New(msg string) error {
	return customError{errorType: errorType, originalError: errors.New(msg)}
}

// New creates a new customError with formatted message
func (errorType ErrorType) Newf(msg string, args ...interface{}) error {
	return customError{errorType: errorType, originalError: fmt.Errorf(msg, args...)}
}

// Wrap creates a new wrapped error
func (errorType ErrorType) Wrap(err error, msg string) error {
	return errorType.Wrapf(err, msg)
}

// Wrap creates a new wrapped error with formatted message
func (errorType ErrorType) Wrapf(err error, msg string, args ...interface{}) error {
	return customError{errorType: errorType, originalError: errors.Wrapf(err, msg, args...)}
}

// Error returns the mssage of a customError
func (error customError) Error() string {
	return error.originalError.Error()
}

// New creates a no type error
func New(msg string) error {
	return customError{errorType: Normal, originalError: errors.New(msg)}
}

// Newf creates a no type error with formatted message
func Newf(msg string, args ...interface{}) error {
	return customError{errorType: Normal, originalError: errors.New(fmt.Sprintf(msg, args...))}
}

// Wrap an error with a string
func Wrap(err error, msg string) error {
	return Wrapf(err, msg)
}

// Cause gives the original error
func Cause(err error) error {
	return errors.Cause(err)
}

// Wrapf an error with format string
func Wrapf(err error, msg string, args ...interface{}) error {
	wrappedError := errors.Wrapf(err, msg, args...)
	if customErr, ok := err.(customError); ok {
		return customError{
			errorType:     customErr.errorType,
			originalError: wrappedError,
			context:       customErr.context,
		}
	}

	return customError{errorType: Normal, originalError: wrappedError}
}

// AddErrorContext adds a context to an error
func AddErrorContext(err error, field, message string) error {
	context := errorContext{Field: field, Message: message}
	if customErr, ok := err.(customError); ok {
		return customError{errorType: customErr.errorType, originalError: customErr.originalError, context: context}
	}

	return customError{errorType: Normal, originalError: err, context: context}
}

// GetErrorContext returns the error context
func GetErrorContext(err error) map[string]string {
	emptyContext := errorContext{}
	if customErr, ok := err.(customError); ok || customErr.context != emptyContext {

		return map[string]string{"field": customErr.context.Field, "message": customErr.context.Message}
	}

	return nil
}

// GetType returns the error type
func GetType(err error) ErrorType {
	if customErr, ok := err.(customError); ok {
		return customErr.errorType
	}

	return Normal
}
