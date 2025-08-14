package postgresql

import (
	"context"
)

func (pgdb *PostgresqlDB) CreateActivityLog(serviceCode, requestNo string, requestBody, responseBody []byte) error {
	_, err := pgdb.DB.Exec(
		context.Background(),
		`
		INSERT INTO activity_log(request_no, service_code, request_body, response_body) 
		VALUES ($1, $2, $3, $4)
		`,
		requestNo,
		serviceCode,
		requestBody,
		responseBody,
	)
	if err != nil {
		return err
	}

	return nil
}
