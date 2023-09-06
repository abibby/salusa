package request

import "bytes"

type HTMLResponse struct {
	Response
}

func NewHTMLResponse(data []byte) *HTMLResponse {
	return &HTMLResponse{
		Response: *NewResponse(bytes.NewBuffer(data)).AddHeader("Content-Type", "text/html"),
	}
}
