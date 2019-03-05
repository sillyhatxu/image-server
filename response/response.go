package response

type ResponseEntity struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"message"`
}

const SUCCESS int = 1

const ERROR int = 0 //Server error

const PARAMS_VALIDATE_ERROR int = -1 //system.validate.error

func Success(data interface{}) *ResponseEntity {
	return &ResponseEntity{Code: SUCCESS, Data: data, Msg: "Success"}
}

func Error(data interface{}, msg string) *ResponseEntity {
	return &ResponseEntity{Code: ERROR, Data: data, Msg: msg}
}

func ErrorParamsValidate(data interface{}, msg string) *ResponseEntity {
	return &ResponseEntity{Code: PARAMS_VALIDATE_ERROR, Data: data, Msg: msg}
}
