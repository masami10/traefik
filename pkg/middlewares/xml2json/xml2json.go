package xml2json

import (
	"bytes"
	"context"
	xj "github.com/basgys/goxml2json"
	"github.com/containous/traefik/v2/pkg/config/dynamic"
	"github.com/containous/traefik/v2/pkg/log"
	"github.com/containous/traefik/v2/pkg/middlewares"
	"io/ioutil"
	"net/http"
)

const (
	typeName = "Xml2Json"
)

// AddPrefix is a middleware used to add prefix to an URL request.
type xml2json struct {
	next   http.Handler
	header string
	name   string
}

// New creates a new handler.
func New(ctx context.Context, next http.Handler, config dynamic.Xml2Json, name string) (http.Handler, error) {
	log.FromContext(middlewares.GetLoggerCtx(ctx, name, typeName)).Debug("Creating middleware")

	result := &xml2json{
		header: config.Header,
		next:   next,
		name:   name,
	}

	return result, nil
}

func (a *xml2json) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	logger := log.FromContext(middlewares.GetLoggerCtx(req.Context(), a.name, typeName))

	body := req.Body

	data, err := xj.Convert(body)
	if err != nil {
		logger.Errorf("Xml2Json Request Body Read Error", req.URL.Path, err)
		http.Error(rw, "Xml2Json Request Body Read Error", http.StatusInternalServerError)
		return
	}
	//fixme: 强制关闭request body, 新建一个body
	if err := body.Close(); err != nil {
		logger.Errorf("Xml2Json Body Close Error", req.URL.Path, err)
	}
	req.Header.Set("Content-Type", "application/json")

	newBodyContent := data.Bytes()
	req.Body = ioutil.NopCloser(bytes.NewReader(newBodyContent))

	req.ContentLength = int64(len(newBodyContent))

	a.next.ServeHTTP(rw, req)
}
