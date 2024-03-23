package jsonrpc

import (
	"encoding/json"
	"fmt"
)

const jsonRPCVersion = "2.0"

type Request struct {
	Version string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
	ID      string `json:"id"`
}

func NewRequest(method string, params []any, id ...string) Request {
	r := Request{
		Version: jsonRPCVersion,
		Method:  method,
		Params:  params,
	}

	if len(id) > 0 {
		r.ID = id[0]
	}

	return r
}

func (r Request) JSON() ([]byte, error) {
	bytes, err := json.Marshal(&r)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

type Response struct {
	ID      string         `json:"id"`
	Version string         `json:"jsonrpc"`
	Result  any            `json:"result,omitempty"`
	Error   *ResponseError `json:"error,omitempty"`
}

type ResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("code=%v msg=%s", e.Code, e.Message)
}
