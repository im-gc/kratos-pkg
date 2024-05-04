package codec

type Response struct {
	Code    int64       `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func (e *Response) Error() string {
	return e.Message
}
