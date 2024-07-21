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
	err := h.serveHTTP(w, r)
	if err == nil {
		return
	}
	if responder, ok := getResponder(err); ok {
		h.respond(w, r, responder)
	} else if handler, ok := err.(http.Handler); ok {
		handler.ServeHTTP(w, r)
	} else {
		h.respond(w, r, NewHTTPError(err, http.StatusInternalServerError))
	}
	addError(r, err)
}
func (h *RequestHandler[TRequest, TResponse]) serveHTTP(w http.ResponseWriter, r *http.Request) error {
	var req TRequest
	err := Run(r, &req)
	if validationErr, ok := err.(ValidationError); ok {
		return NewHTTPError(validationErr, http.StatusUnprocessableEntity)
	} else if err != nil {
		return err
	}

	resp, err := h.handler(&req)
	if validationErr, ok := err.(ValidationError); ok {
		return NewHTTPError(validationErr, http.StatusUnprocessableEntity)
	} else if err != nil {
		return err
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

	return nil
}

func (h *RequestHandler[TRequest, TResponse]) Run(r *TRequest) (TResponse, error) {
	return h.handler(r)
}

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
