package db

type DBActivityLogInterface interface {
	CreateActivityLog(serviceCode, requestNo string, requestBody, responseBody []byte) error
}
