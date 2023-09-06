package request

import (
	"log"
	"net/http"
	"reflect"
)

type RequestHandler[TRequest, TResponse any] func(r *TRequest) (TResponse, error)

func (h RequestHandler[TRequest, TResponse]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req TRequest
	err := Run(r, &req)
	if validationErr, ok := err.(ValidationError); ok {
		respond(w, r, errorResponse(validationErr, http.StatusUnprocessableEntity, r))
		return
	} else if err != nil {
		respond(w, r, errorResponse(err, http.StatusInternalServerError, r))
		return
	}

	injectRequest(&req, r)
	injectResponseWriter(&req, w)
	err = di(&req, r)
	if err != nil {
		if responder, ok := err.(Responder); ok {
			respond(w, r, responder)
		} else {
			respond(w, r, errorResponse(err, http.StatusInternalServerError, r))
		}
		return
	}

	resp, err := h(&req)
	if err != nil {
		if responder, ok := err.(Responder); ok {
			respond(w, r, responder)
		} else {
			respond(w, r, errorResponse(err, http.StatusInternalServerError, r))
		}
		return
	}

	var anyResp any = resp
	if responder, ok := anyResp.(Responder); ok {
		respond(w, r, responder)
		return
	}
	respond(w, r, NewJSONResponse(resp))
}

func Handler[TRequest, TResponse any](callback func(r *TRequest) (TResponse, error)) RequestHandler[TRequest, TResponse] {
	return RequestHandler[TRequest, TResponse](callback)
}

func respond(w http.ResponseWriter, req *http.Request, r Responder) {
	err := r.Respond(w, req)
	if err != nil {
		log.Print(err)
	}
}

func injectRequest[TRequest any](req TRequest, httpR *http.Request) {
	v := reflect.ValueOf(req).Elem()

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if (f.Type() != reflect.TypeOf(&http.Request{})) {
			continue
		}
		f.Set(reflect.ValueOf(httpR))
	}

}

func injectResponseWriter[TRequest any](req TRequest, httpRW http.ResponseWriter) {
	v := reflect.ValueOf(req).Elem()

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		var rw http.ResponseWriter
		if !f.Type().Implements(reflect.TypeOf(&rw).Elem()) {
			continue
		}
		f.Set(reflect.ValueOf(httpRW))
	}

}
