package utils

var (
	success = "success"
	fail    = "fail"
	error   = "error"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func NewSuccessResponse(message string, data any) Response {
	return Response{
		Status:  success,
		Message: message,
		Data:    data,
	}
}

func NewFailResponse(message string, data any) Response {
	return Response{
		Status:  fail,
		Data:    data,
		Message: message,
	}
}

func NewErrorResponse(message string, data any) Response {
	return Response{
		Status:  error,
		Message: message,
		Data:    data,
	}
}
