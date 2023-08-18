package model

const (
    SuccessCode = 0
    ErrorCode   = 1
)

type Response struct {
    StatusCode int32  `json:"status_code"`
    StatusMsg  string `json:"status_msg,omitempty"`
}

func SuccessResponse() Response {
    return Response{
        StatusCode: SuccessCode,
        StatusMsg:  "success",
    }
}

func ErrorResponse(statusCode int32, statusMsg string) Response {
    return Response{
        StatusCode: statusCode,
        StatusMsg:  statusMsg,
    }
}
