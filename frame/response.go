package frame

var (
	DefaultArray  = []int{}
	DefaultObject = make(map[int]struct{})
)

type DefaultResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

func NewDefaultResponse() DefaultResponse {
	return DefaultResponse{
		Data: DefaultObject,
	}
}
