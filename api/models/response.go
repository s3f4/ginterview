package models

// Response holds response data
type Response struct {
	Code    int               `json:"code"`
	Msg     string            `json:"msg"`
	Records []*ResponseRecord `json:"records,omitempty"`
}
