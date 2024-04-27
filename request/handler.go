package request

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/abibby/salusa/di"
)

func init() {
	err := Register(context.Background())
	if err != nil {
		panic(err)
	}
}

type RequestHandler[TRequest, TResponse any] struct {
	handler func(r *TRequest) (TResponse, error)
}

func (h *RequestHandler[TRequest, TResponse]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err, status := h.serveHTTP(w, r)
	if err == nil {
		return
	}
	if responder, ok := getResponder(err); ok {
		h.respond(w, r, responder)
	} else if handler, ok := err.(http.Handler); ok {
		handler.ServeHTTP(w, r)
	} else {
		h.respond(w, r, errorResponse(err, status, r))
	}
	addError(r, err)
}
func (h *RequestHandler[TRequest, TResponse]) serveHTTP(w http.ResponseWriter, r *http.Request) (error, int) {
	var req TRequest
	err := Run(r, &req)
	if validationErr, ok := err.(ValidationError); ok {
		return validationErr, http.StatusUnprocessableEntity
	} else if err != nil {
		return err, http.StatusInternalServerError
	}

	ctx := r.Context()
	ctx = context.WithValue(ctx, requestKey, r)
	ctx = context.WithValue(ctx, responseKey, w)

	err = di.Fill(ctx, &req,
		di.AutoResolve[context.Context](),
		di.AutoResolve[*http.Request](),
		di.AutoResolve[http.ResponseWriter](),
	)
	if err != nil {
		return err, http.StatusInternalServerError
	}

	resp, err := h.handler(&req)
	if err != nil {
		return err, http.StatusInternalServerError
	}

	var anyResp any = resp
	switch resp := anyResp.(type) {
	case Responder:
		h.respond(w, r, resp)
	case http.Handler:
		resp.ServeHTTP(w, r)
	default:
		h.respond(w, r, NewJSONResponse(resp))
	}

	return nil, http.StatusOK
}
func (h *RequestHandler[TRequest, TResponse]) Run(r *TRequest) (TResponse, error) {
	return h.handler(r)
}

// Handler is a helper to create http handlers with built in input validation
// and error handling.
//
//	type Request struct {
//		Foo string `path:"foo" validate:"required"`
//		Bar string `query:"bar"`
//		Baz string `json:"baz"`
//	}
//	type Response struct{}
//	request.Handler(func(r *Request) (*Response, error) {
//		return nil, nil
//	})
func Handler[TRequest, TResponse any](callback func(r *TRequest) (TResponse, error)) *RequestHandler[TRequest, TResponse] {
	return &RequestHandler[TRequest, TResponse]{
		handler: callback,
	}
}

func (h *RequestHandler[TRequest, TResponse]) respond(w http.ResponseWriter, req *http.Request, r Responder) {
	err := r.Respond(w, req)
	if err != nil {
		logger, resolveErr := di.Resolve[*slog.Logger](req.Context())
		if resolveErr != nil {
			logger = slog.Default()
		}
		logger.Error("request failed", "error", err)
	}
}
func (r *RequestHandler[TRequest, TResponse]) Validate(ctx context.Context) error {
	var req TRequest
	return di.Validate(ctx, &req,
		di.AutoResolve[context.Context](),
		di.AutoResolve[*http.Request](),
		di.AutoResolve[http.ResponseWriter](),
	)
}
