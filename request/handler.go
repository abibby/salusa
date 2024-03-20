package request

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/abibby/salusa/di"
)

var defaultDP = di.NewDependencyProvider()

func init() {
	RegisterDI(defaultDP)
}

type RequestHandler[TRequest, TResponse any] struct {
	handler func(r *TRequest) (TResponse, error)
	dp      *di.DependencyProvider
}

func (h *RequestHandler[TRequest, TResponse]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// var req TRequest
	var resp TResponse
	// err := Run(r, &req)
	// if validationErr, ok := err.(ValidationError); ok {
	// 	h.respond(w, r, errorResponse(validationErr, http.StatusUnprocessableEntity, r))
	// 	return
	// } else if err != nil {
	// 	h.respond(w, r, errorResponse(err, http.StatusInternalServerError, r))
	// 	return
	// }

	ctx := r.Context()
	ctx = context.WithValue(ctx, requestKey, r)
	ctx = context.WithValue(ctx, responseKey, w)

	err := di.Use(ctx, h.dp, func(req *TRequest) error {
		err := Run(r, req)
		if validationErr, ok := err.(ValidationError); ok {
			return NewHTTPError(validationErr, http.StatusUnprocessableEntity)
		} else if err != nil {
			return NewHTTPError(err, http.StatusInternalServerError)
		}
		resp, err = h.handler(req)
		return err
	})

	// err = h.dp.Fill(ctx, &req,
	// 	di.AutoResolve[context.Context](),
	// 	di.AutoResolve[*http.Request](),
	// 	di.AutoResolve[http.ResponseWriter](),
	// )
	if err != nil {
		if responder, ok := getResponder(err); ok {
			h.respond(w, r, responder)
		} else {
			h.respond(w, r, errorResponse(err, http.StatusInternalServerError, r))
		}
		return
	}

	// resp, err := h.handler(&req)
	// if err != nil {
	// 	if responder, ok := err.(Responder); ok {
	// 		h.respond(w, r, responder)
	// 	} else {
	// 		h.respond(w, r, errorResponse(err, http.StatusInternalServerError, r))
	// 	}
	// 	return
	// }

	var anyResp any = resp
	if responder, ok := anyResp.(Responder); ok {
		h.respond(w, r, responder)
		return
	}
	h.respond(w, r, NewJSONResponse(resp))
}
func (h *RequestHandler[TRequest, TResponse]) WithDependencyProvider(dp *di.DependencyProvider) {
	h.dp = dp
}
func (h *RequestHandler[TRequest, TResponse]) Run(r *TRequest) (TResponse, error) {
	return h.handler(r)
}

// Handler is a helper to create http handlers with built in input validation
// and error handling.
func Handler[TRequest, TResponse any](callback func(r *TRequest) (TResponse, error)) *RequestHandler[TRequest, TResponse] {
	return &RequestHandler[TRequest, TResponse]{
		handler: callback,
		dp:      defaultDP,
	}
}

func (h *RequestHandler[TRequest, TResponse]) respond(w http.ResponseWriter, req *http.Request, r Responder) {
	err := r.Respond(w, req)
	if err != nil {
		logger, err := di.Resolve[*slog.Logger](req.Context(), h.dp)
		if err != nil {
			logger = slog.Default()
		}
		logger.Error("request failed", "error", err)
	}
}
