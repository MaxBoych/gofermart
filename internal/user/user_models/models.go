package user_models

import "time"

type UserStorageData struct {
	UserId    int64      `db:"user_id"`
	Login     string     `db:"login"`
	Password  string     `db:"password"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

type UserRegisterRequest struct {
	UserLoginRequest
}

type UserLoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func ToUserStorageData(req UserRegisterRequest) UserStorageData {
	return UserStorageData{
		Login: req.Login,
	}
}

type CookieData struct {
	Name    string
	Value   string
	Expires time.Time
	Domain  string
}
