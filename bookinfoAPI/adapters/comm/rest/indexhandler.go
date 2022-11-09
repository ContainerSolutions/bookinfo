package rest

import (
	"fmt"
	"io/ioutil"
	"net/http"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

// swagger:route GET / Index index
// Returns OK if there's no problem
// responses:
//	200: OK

// Index returns OK handles GET requests
func (p *APIContext) Index(rw http.ResponseWriter, r *http.Request) {
	tracer := opentracing.GlobalTracer()
	spanname := "BookInfoAPI.Index"
	var span opentracing.Span

	wireContext, err := opentracing.GlobalTracer().Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(r.Header))
	if err != nil {
		// The method is called without a span context in the http header.
		span = tracer.StartSpan(spanname)
	} else {
		// Create the span referring to the RPC client if available.
		// If wireContext == nil, a root span will be created.
		span = opentracing.StartSpan(spanname, ext.RPCServerOption(wireContext))
	}
	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, r.URL.RequestURI())
	ext.HTTPMethod.Set(span, r.Method)
	defer span.Finish()

	rw.WriteHeader(200)
}

// swagger:route GET /version Index version
// Returns version information
// responses:
//	200: OK

// Version returns the version info for the service by reading from a static file
func (p *APIContext) Version(rw http.ResponseWriter, r *http.Request) {
	tracer := opentracing.GlobalTracer()
	spanname := "BookInfoAPI.Version"
	var span opentracing.Span

	wireContext, err := opentracing.GlobalTracer().Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(r.Header))
	if err != nil {
		// The method is called without a span context in the http header.
		span = tracer.StartSpan(spanname)
	} else {
		// Create the span referring to the RPC client if available.
		// If wireContext == nil, a root span will be created.
		span = opentracing.StartSpan(spanname, ext.RPCServerOption(wireContext))
	}
	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, r.URL.RequestURI())
	ext.HTTPMethod.Set(span, r.Method)
	defer span.Finish()
	dat, err := ioutil.ReadFile("./static/version.txt")
	if err != nil {
		dat = append(dat, '0')
	}
	fmt.Fprintf(rw, "Welcome to BookInfo API! Version:%s", dat)
}
