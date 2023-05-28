package router

import (
	"net/http"
	"time"

	chi "github.com/go-chi/chi/v5"
	middleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
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

func NewRouter(opts ...Option) http.Handler {
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
	r.Use(httprate.LimitAll(routerOpts.RateLimitRequests, routerOpts.RateLimitDuration))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	return r
}
