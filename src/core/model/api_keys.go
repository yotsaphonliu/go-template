package model

import "time"

type ApiKey struct {
	Key            string    `json:"key"`
	AzureUserID    string    `json:"azure_user_id"`
	UserID         int64     `json:"user_id"`
	EmailAddress   string    `json:"email_address"`
	UserRoleName   string    `json:"user_role_name"`
	ExpireTime     time.Time `json:"expire_time"`
	CreatedTime    time.Time `json:"created_time"`
	UserProfilePic string    `json:"user_profile_pic"`
}
