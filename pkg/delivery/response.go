package delivery

type ResCode int

const (
	RES_OK ResCode = 0
)

type ResError struct {
	Title string `json:"title"`
	Desc  string `json:"description"`
}

type ResBody struct {
	ResCode ResCode     `json:"resp_code"`
	Error   *ResError   `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
