package dto

import (
	"time"
)

type TokenRsp struct {
	Access   string        `json:"access_token,omitempty"`
	Refresh  string        `json:"refresh_token,omitempty"`
	Duration time.Duration `json:"expires_in,omitempty"`
}
