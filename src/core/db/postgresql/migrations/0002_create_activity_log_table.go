package migrations

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

var createActivityLogTableMigration = &Migration{
	Number: 2,
	Name:   "Create activity_log table",
	Forwards: func(db *gorm.DB) error {
		const sql = `
			CREATE TABLE activity_log (
				request_no uuid PRIMARY KEY NOT NULL,
				service_code text NOT NULL,
				service_name text NOT NULL,
				service_endpoint text NOT NULL,
				http_method http_type NOT NULL,
				request_header text,
				request_body text,
				http_status_code int NOT NULL,
				response_code int NOT NULL,
				response_message text,
				response_body text,
				active_status varchar(1) NOT NULL DEFAULT 'Y',
				create_date timestamptz NOT NULL DEFAULT now(),
				create_by varchar(255) NOT NULL,
				update_date timestamptz NOT NULL DEFAULT now(),
				update_by varchar(255) NOT NULL
			);
			
			CREATE INDEX service_code_idx ON activity_log (service_code);
		`

		return errors.Wrap(db.Exec(sql).Error, "unable to create activity_log table")
	},
}

func init() {
	Migrations = append(Migrations, createActivityLogTableMigration)
}
