package rest

import (
	"io"
	"net/http"
	"time"

	"github.com/ContainerSolutions/bookinfo/bookStockAPI/application"

	middleware "github.com/ContainerSolutions/bookinfo/bookStockAPI/adapters/comm/rest/middleware"
	"github.com/gorilla/mux"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"

	openapimw "github.com/go-openapi/runtime/middleware"

	"github.com/uber/jaeger-client-go"
	jprom "github.com/uber/jaeger-lib/metrics/prometheus"

	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
)

// APIContext handler for getting and updating Ratings
type APIContext struct {
	validation *middleware.Validation
	//dbContext  DBContext
	healthRepo    application.HealthRepository
	bookInfoRepo  application.BookStockRepository
	configuration map[string]string
	delay         bool
}

// NewAPIContext returns a new APIContext handler with the given logger
// func NewAPIContext(dc DBContext, bindAddress *string, ur application.UserRepository) *http.Server {
func NewAPIContext(bindAddress *string, hr application.HealthRepository, pr application.BookStockRepository, dly bool) (*http.Server, io.Closer) {
	apiContext := &APIContext{
		healthRepo:   hr,
		bookInfoRepo: pr,
		delay:        dly,
	}
	s, c := apiContext.prepareContext(bindAddress)
	return s, c
}

func (apiContext *APIContext) prepareContext(bindAddress *string) (*http.Server, io.Closer) {
	// Example logger and metrics factory. Use github.com/uber/jaeger-client-go/log
	// and github.com/uber/jaeger-lib/metrics respectively to bind to real logging and metrics
	// frameworks.
	jLogger := jaegerlog.StdLogger
	jMetricsFactory := jprom.New()

	// Sample configuration for testing. Use constant sampling to sample every trace
	// and enable LogSpan to log every span via configured Logger.
	cfg, err := jaegercfg.FromEnv()
	if err != nil || cfg.ServiceName == "" {
		log.Error().Err(err).Msg("Cannot load tracer config from env")
		cfg = &jaegercfg.Configuration{
			ServiceName: "bookInfoAPI",
			Sampler:     &jaegercfg.SamplerConfig{},
			Reporter:    &jaegercfg.ReporterConfig{},
		}
	}
	cfg.Sampler.Type = jaeger.SamplerTypeConst
	cfg.Sampler.Param = 1
	cfg.Reporter.LogSpans = true

	// Initialize tracer with a logger and a metrics factory
	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
	)
	if err != nil {
		log.Error().Err(err).Msg("Cannot create tracer")
	}
	// Set the singleton opentracing.Tracer with the Jaeger tracer.
	opentracing.SetGlobalTracer(tracer)

	apiContext.validation = middleware.NewValidation()

	// create a new serve mux and register the handlers
	sm := mux.NewRouter()
	sm.Use(middleware.MetricsMiddleware)

	// handlers for API
	getR := sm.Methods(http.MethodGet).Subrouter()
	// Generic handlers
	getR.HandleFunc("/", apiContext.Index)
	getR.HandleFunc("/version", apiContext.Version)
	getR.HandleFunc("/health/live", apiContext.Live)
	getR.HandleFunc("/health/ready", apiContext.Ready)
	// BookStock handlers
	getR.HandleFunc("/book/{id}", apiContext.GetBookStock)
	// Documentation handler
	opts := openapimw.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := openapimw.Redoc(opts, nil)
	getR.Handle("/docs", sh)
	getR.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	// create a new server
	s := &http.Server{
		Addr:         *bindAddress,      // configure the bind address
		Handler:      sm,                // set the default handler
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	sm.PathPrefix("/metrics").Handler(promhttp.Handler())
	prometheus.MustRegister(middleware.RequestCounterVec)
	prometheus.MustRegister(middleware.RequestDurationGauge)

	return s, closer
}

// createSpan creates a new openTracing.Span with the given name and returns it
func createSpan(spanName string, r *http.Request) (span opentracing.Span) {
	tracer := opentracing.GlobalTracer()

	wireContext, err := opentracing.GlobalTracer().Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(r.Header))
	if err != nil {
		// The method is called without a span context in the http header.
		//
		span = tracer.StartSpan(spanName)
	} else {
		// Create the span referring to the RPC client if available.
		// If wireContext == nil, a root span will be created.
		span = opentracing.StartSpan(spanName, ext.RPCServerOption(wireContext))
	}
	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, r.URL.RequestURI())
	ext.HTTPMethod.Set(span, r.Method)
	return span
}
