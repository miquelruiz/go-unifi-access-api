package schema

const (
	ResponseSuccess = "SUCCESS"
)

type Response[T any] struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}
