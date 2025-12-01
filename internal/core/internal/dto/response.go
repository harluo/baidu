package dto

type Response struct {
	Log     uint64 `json:"log_id,omitempty"`
	Message string `json:"error_msg,omitempty"`
	Code    uint32 `json:"error_code,omitempty"`

	Data any `json:",inline"`
}

func (r *Response) IsError() bool {
	return 0 != r.Code
}
