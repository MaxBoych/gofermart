package token_models

import "time"

type SecretKeyStorageData struct {
	Value     string    `db:"key_value"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type TokenStorageData struct {
	TokenId   int64     `db:"token_id"`
	UserId    int64     `db:"user_id"`
	Value     string    `db:"token_value"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
