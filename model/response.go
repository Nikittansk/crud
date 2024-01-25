package model

type Response struct {
	StatusCode int         `json:"status_code"`
	Data       interface{} `json:"data"`
}