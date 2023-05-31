package router

import (
	"net/http"
	"time"

	chi "github.com/go-chi/chi/v5"
	middleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/mariadb-operator/agent/pkg/handler"
)

type Options struct {
	CompressLevel     int
	RateLimitRequests int
	RateLimitDuration time.Duration
}

type Option func(*Options)

func WithCompressLevel(level int) Option {
	return func(o *Options) {
		o.CompressLevel = level
	}
}

func WithRateLimit(requests int, duration time.Duration) Option {
	return func(o *Options) {
		o.RateLimitRequests = requests
		o.RateLimitDuration = duration
	}
}

func NewRouter(handler *handler.Handler, opts ...Option) http.Handler {
	routerOpts := Options{
		CompressLevel:     5,
		RateLimitRequests: 100,
		RateLimitDuration: 1 * time.Minute,
	}
	for _, setOpt := range opts {
		setOpt(&routerOpts)
	}
	r := chi.NewRouter()
	r.Use(middleware.Compress(routerOpts.CompressLevel))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	r.Mount("/api", apiRouter(handler, &routerOpts))

	return r
}

func apiRouter(h *handler.Handler, opts *Options) http.Handler {
	r := chi.NewRouter()
	r.Use(httprate.LimitAll(opts.RateLimitRequests, opts.RateLimitDuration))

	r.Route("/bootstrap", func(r chi.Router) {
		r.Put("/", h.Bootstrap.Put)
		r.Delete("/", h.Bootstrap.Delete)
	})
	r.Get("/galerastate", h.GaleraState.Get)
	r.Route("/recovery", func(r chi.Router) {
		r.Put("/", h.Recovery.Put)
		r.Delete("/", h.Recovery.Delete)
	})

	return r
}
