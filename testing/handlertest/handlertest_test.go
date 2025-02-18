package handlertest_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/abibby/salusa/testing/handlertest"
)

func TestHandlerTest(t *testing.T) {
	responseBody := `{"foo":{"bar":1}}`
	ctx := context.Background()
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, responseBody)
	})

	handlertest.New(ctx, t, h).
		Get("/test").
		AssertJSONString(responseBody).
		AssertJSONContains("foo.bar", 1.0)
}
