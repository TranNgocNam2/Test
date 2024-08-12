package web

type BaseResponse struct {
	ResultCode    string      `json:"resultCode"`
	ResultMessage string      `json:"resultMessage"`
	Data          interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	ResultCode    string `json:"resultCode"`
	ResultMessage string `json:"resultMessage"`
}
