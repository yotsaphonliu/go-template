package model

import "time"

// ActivityLog represents the activity_log table
type ActivityLog struct {
	ID            int64     `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	UserID        *int64    `json:"user_id" gorm:"column:user_id"`
	EmailAddress  *string   `json:"email_address" gorm:"column:email_address"`
	UserRoleName  *string   `json:"user_role_name" gorm:"column:user_role_name"`
	HTTPMethod    string    `json:"http_method" gorm:"column:http_method"`
	RequestURI    string    `json:"request_uri" gorm:"column:request_uri"`
	RequestBody   *string   `json:"request_body" gorm:"column:request_body"`
	ResponseCode  int       `json:"response_code" gorm:"column:response_code"`
	ResponseBody  *string   `json:"response_body" gorm:"column:response_body"`
	IPAddress     string    `json:"ip_address" gorm:"column:ip_address"`
	UserAgent     *string   `json:"user_agent" gorm:"column:user_agent"`
	ExecutionTime int64     `json:"execution_time" gorm:"column:execution_time"`
	CreatedTime   time.Time `json:"created_time" gorm:"column:created_time;default:now()"`
}
