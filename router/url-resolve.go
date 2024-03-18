package router

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/abibby/salusa/database/model"
	"github.com/abibby/salusa/di"
	"github.com/abibby/salusa/internal/helpers"
)

type URLResolver interface {
	Resolve(name string, params ...any) string
	ResolveHandler(h http.Handler, params ...any) string
}

func (r *Router) Register(dp *di.DependencyProvider) {
	di.Register(dp, func(ctx context.Context, tag string) (URLResolver, error) {
		origin := ""
		req, err := di.Resolve[*http.Request](ctx, dp)
		if err != nil {
			origin = "/"
		} else {
			origin = req.URL.Scheme + "://" + req.URL.Host + "/"
		}
		return &SalusaResolver{
			origin: origin,
			r:      r,
		}, nil
	})
}

type Attr struct {
	Key   string
	Value string
}

func ToAttrs(params []any) ([]*Attr, error) {
	var active *Attr
	attrs := []*Attr{}
	for _, p := range params {
		if active == nil {
			switch p := p.(type) {
			case string:
				active = &Attr{
					Key: p,
				}
			case *Attr:
				attrs = append(attrs, p)
			default:
				return nil, fmt.Errorf("invalid attrs: expected key string or attr received %s", reflect.TypeOf(p))
			}
		} else {
			switch p := p.(type) {
			case string, fmt.Stringer,
				int, int8, int16, int32, int64,
				uint, uint8, uint16, uint32, uint64,
				float32, float64:

				active.Value = fmt.Sprint(p)
				attrs = append(attrs, active)
				active = nil
			case model.Model:
				pKeyValues := helpers.PrimaryKeyValue(p)
				if len(pKeyValues) != 1 {
					return nil, fmt.Errorf("invalid attrs: model must have only 1 primary key found %d", len(pKeyValues))
				}
				active.Value = fmt.Sprint(pKeyValues[0])
				attrs = append(attrs, active)
				active = nil
			default:
				return nil, fmt.Errorf("invalid attrs: expected value string or number received %s", reflect.TypeOf(p))
			}
		}
	}
	return attrs, nil
}

type SalusaResolver struct {
	origin string
	r      *Router
}

func NewResolver(origin string, r *Router) *SalusaResolver {
	return &SalusaResolver{
		origin: origin,
		r:      r,
	}
}

func (r *SalusaResolver) Resolve(name string, params ...any) string {
	attrs, err := ToAttrs(params)
	if err != nil {
		panic(err)
	}
	v := url.Values{}

	path := name
	for _, a := range attrs {
		k := "{" + a.Key + "}"
		if strings.Contains(path, k) {
			path = strings.ReplaceAll(path, k, a.Value)
		} else {
			v.Add(a.Key, a.Value)
		}
	}

	u := strings.TrimSuffix(r.origin, "/") + "/" + strings.TrimPrefix(path, "/")
	if len(v) > 0 {
		u += "?" + v.Encode()
	}
	return u
}
func (r *SalusaResolver) ResolveHandler(h http.Handler, params ...any) string {
	s, ok := r.HandlerPath(h)
	if !ok {
		s = "/"
	}
	return r.Resolve(s, params...)
}

func (r *SalusaResolver) HandlerPath(h http.Handler) (string, bool) {
	p, ok := r.r.handlers[h]
	return p, ok
}
