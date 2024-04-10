package generateclient_test

import (
	"bytes"
	"testing"

	"github.com/abibby/salusa/request"
	"github.com/abibby/salusa/router"
	"github.com/abibby/salusa/router/generateclient"
	"github.com/stretchr/testify/assert"
)

func TestGenerateClient(t *testing.T) {
	testCases := []struct {
		Name     string
		Router   func(r *router.Router)
		Expected string
	}{
		// 		{
		// 			Name: "handlerFunc",
		// 			Router: func(r *router.Router) {
		// 				r.Get("/path", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		// 			},
		// 			Expected: `export async function path(): Promise<unknown> {
		// 	const response = await fetcher("/path")
		// 	if (response.status < 200 || response.status >= 300) {
		// 		throw new Error("invalid status")
		// 	}
		// 	return await response.json()
		// }`,
		// 		},
		// 		{
		// 			Name: "handlerFunc multi part",
		// 			Router: func(r *router.Router) {
		// 				r.Get("/path/section", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		// 			},
		// 			Expected: `export async function pathSection(): Promise<unknown> {
		// 	const response = await fetcher("/path/section")
		// 	if (response.status < 200 || response.status >= 300) {
		// 		throw new Error("invalid status")
		// 	}
		// 	return await response.json()
		// }`,
		// 		},
		// 		{
		// 			Name: "path param",
		// 			Router: func(r *router.Router) {
		// 				r.Get("/path/{param1}/{param2}/more_path", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		// 			},
		// 			Expected: `export async function pathMorePath(param1: string|number, param2: string|number): Promise<unknown> {
		// 	const response = await fetcher("/path/" + encodeURIComponent(param1) + "/" + encodeURIComponent(param2) + "/more_path")
		// 	if (response.status < 200 || response.status >= 300) {
		// 		throw new Error("invalid status")
		// 	}
		// 	return await response.json()
		// }`,
		// 		},
		{
			Name: "request.Handler",
			Router: func(r *router.Router) {
				type Req struct {
					Foo string `json:"foo"`
				}
				type Resp struct {
					Bar string `json:"bar"`
				}
				r.Get("/path/section", request.Handler(func(r *Req) (*Resp, error) {
					return nil, nil
				}))
			},
			Expected: `export type PathSectionRequest = {foo:string}
export type PathSectionResponse = {bar:string}
export async function pathSection(req: PathSectionRequest): Promise<PathSectionResponse> {
	const response = await fetcher("/path/section")
	if (response.status < 200 || response.status >= 300) {
		throw new Error("invalid status")
	}
	return await response.json()
}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			r := router.New()
			tc.Router(r)
			b := &bytes.Buffer{}
			err := generateclient.GenerateClient(b, r)
			assert.NoError(t, err)
			assert.Equal(t, tc.Expected, b.String())
		})
	}
}
