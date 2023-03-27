package users

import "time"

type (
	Token struct {
		ID        uint64    `json:"id"`
		Token     string    `json:"token"`
		CreatedAt time.Time `json:"created_at"`
	}
	Tokens []Token
)

type User struct {
	ID     uint64
	Login  string
	Passwd []byte
	Groups map[uint64]struct{}
}
