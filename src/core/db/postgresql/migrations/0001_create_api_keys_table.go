package migrations

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

var createApiKeysTableMigration = &Migration{
	Number: 1,
	Name:   "Create api_keys table",
	Forwards: func(db *gorm.DB) error {
		const sql = `
			CREATE TABLE api_keys(
			    key TEXT NOT NULL PRIMARY KEY,
			    azure_user_id TEXT NOT NULL,
				user_id BIGINT,
				email_address TEXT,
				user_role_name TEXT,
				expire_time timestamptz NOT NULL,
				created_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
				user_profile_pic TEXT
			);
			
			create index if not exists ak_key_idx on api_keys (key);
			create index if not exists ak_user_id_idx on api_keys (user_id);
			create index if not exists ak_expire_time_idx on api_keys (expire_time);
		`

		return errors.Wrap(db.Exec(sql).Error, "unable to create api_keys table")
	},
}

func init() {
	Migrations = append(Migrations, createApiKeysTableMigration)
}
