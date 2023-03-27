package npm

import "time"

type (
	LoginRequest struct {
		ID       string    `json:"_id"`
		Name     string    `json:"name"`
		Password string    `json:"password"`
		Type     string    `json:"type"`
		Roles    []string  `json:"roles"`
		Date     time.Time `json:"date"`
	}
	LoginResponse struct {
		Token string `json:"token"`
		Ok    bool   `json:"ok"`
		ID    string `json:"id"`
		Rev   string `json:"rev"`
	}

	ErrorResponse struct {
		Message string `json:"message"`
		Error   string `json:"error"`
		Ok      bool   `json:"ok"`
	}
)
