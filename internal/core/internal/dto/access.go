package dto

import (
	"time"
)

type AccessRsp struct {
	Token   string    `json:"token,omitempty"`
	Expired time.Time `json:"expired_in,omitempty"`
}
