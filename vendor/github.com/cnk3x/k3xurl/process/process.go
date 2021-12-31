package process

import (
	"encoding/xml"
	"io"
	"net/http"

	"github.com/goccy/go-json"
)

// Process 响应处理器
type Process = func(resp *http.Response) error

// JSON 解析JSON
func JSON(out any) Process {
	return func(resp *http.Response) error {
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return json.Unmarshal(data, out)
	}
}

// XML 解析XML
func XML(out any) Process {
	return func(resp *http.Response) error {
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return xml.Unmarshal(data, out)
	}
}
