package request

import "bytes"

type HTMLResponse struct {
	ResponseBuilder
}

func NewHTMLResponse(data []byte) *HTMLResponse {
	return &HTMLResponse{
		ResponseBuilder: *NewResponse(bytes.NewBuffer(data)).AddHeader("Content-Type", "text/html"),
	}
}
