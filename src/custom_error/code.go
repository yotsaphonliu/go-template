package custom_error

const errorCodeBase = 10000

const (
	UnknownError int = errorCodeBase + iota + 1
	InvalidJSONString
	InputValidationError
	Unauthorized
	InvalidAuthData
	ReadCsvError
	InvalidCsvHeader
	DBError
	ConvertStringToInt64Error
	InvalidUsernameOrPassword
	NotInAllowGroup
	DuplicateRole
	TransactionNotFound
	ExternalServiceError
	InvalidParameter
	InternalServerError
	MissingRequiredField
)
