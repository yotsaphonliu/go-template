package postgresql

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"go-template/src/core/model"
)

func (pgdb *PostgresqlDB) InsertApiKeys(apiKeysList []model.ApiKey, isRoot bool) error {

	ctx := context.Background()

	tx, err := pgdb.DB.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "Unable to make a transaction")
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
			if err != nil {
				err = errors.Wrap(err, "Unable to commit a transaction")
			}
		}
	}()

	for _, apiKey := range apiKeysList {
		_, err = tx.Exec(ctx, `
				INSERT INTO api_keys(
					key,
					azure_user_id,
					user_id,
					email_address,
					user_role_name,
					expire_time,
					user_profile_pic
				)
				VALUES ($1, $2, $3, $4, $5, $6, $7)
			`,
			apiKey.Key,
			apiKey.AzureUserID,
			apiKey.UserID,
			apiKey.EmailAddress,
			apiKey.UserRoleName,
			apiKey.ExpireTime,
			apiKey.UserProfilePic,
		)
		if err != nil {
			return err
		}

	}

	return nil
}

func (pgdb *PostgresqlDB) VerifyApiKey(key string, newExpireTime time.Time) ([]*model.ApiKey, error) {
	result := make([]*model.ApiKey, 0)
	_, err := pgdb.DB.Exec(context.Background(), `
		UPDATE api_keys
		SET expire_time = $1
		WHERE key = $2 AND expire_time >= NOW()
	`,
		newExpireTime,
		key,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	err = pgdb.DB.QueryRow(context.Background(),
		`
			WITH cte AS (
				SELECT 
					*
				FROM 
					api_keys
				WHERE 
					key = $1
			)
			SELECT
				(
					SELECT
						COALESCE(jsonb_agg(d.*), '[]')
					FROM
						(
							SELECT
								*
							FROM
								cte
						) as d
				) as rows;
		`,
		key,
	).Scan(
		&result,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Can not select api key list from database")
	}

	return result, nil
}

func (pgdb *PostgresqlDB) DeleteApiKey(key string) error {
	_, err := pgdb.DB.Exec(context.Background(), `
		DELETE FROM api_keys WHERE key = $1
	`,
		key,
	)
	if err != nil {
		return err
	}

	return nil
}

func (pgdb *PostgresqlDB) DeleteExpireApiKey() error {
	result, err := pgdb.DB.Exec(context.Background(), `
		DELETE FROM api_keys WHERE expire_time IS NOT NULL AND expire_time < NOW()
	`,
	)
	if err != nil {
		return err
	}

	if result.Delete() && result.RowsAffected() > 0 {
		pgdb.logger.Infof("Deleted %v api key expire", result.RowsAffected())
	}

	return nil
}
