package router

import (
	"net/http"
	"time"

	chi "github.com/go-chi/chi/v5"
	middleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/go-logr/logr"
	"github.com/mariadb-operator/agent/pkg/handler"
	"github.com/mariadb-operator/agent/pkg/kubernetesauth"
	"k8s.io/client-go/kubernetes"
)

type Options struct {
	CompressLevel     int
	RateLimitRequests *int
	RateLimitDuration *time.Duration
	KubernetesAuth    bool
	KubernetesTrusted *kubernetesauth.Trusted
}

type Option func(*Options)

func WithCompressLevel(level int) Option {
	return func(o *Options) {
		o.CompressLevel = level
	}
}

func WithRateLimit(requests int, duration time.Duration) Option {
	return func(o *Options) {
		if requests != 0 && duration != 0 {
			o.RateLimitRequests = &requests
			o.RateLimitDuration = &duration
		}
	}
}

func WithKubernetesAuth(auth bool, trusted *kubernetesauth.Trusted) Option {
	return func(o *Options) {
		o.KubernetesAuth = auth
		o.KubernetesTrusted = trusted
	}
}

func NewRouter(handler *handler.Handler, clientset *kubernetes.Clientset, logger logr.Logger, opts ...Option) http.Handler {
	routerOpts := Options{
		CompressLevel:     5,
		KubernetesAuth:    false,
		KubernetesTrusted: nil,
	}
	for _, setOpt := range opts {
		setOpt(&routerOpts)
	}
	r := chi.NewRouter()
	r.Use(middleware.Compress(routerOpts.CompressLevel))
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	r.Mount("/api", apiRouter(handler, clientset, logger, &routerOpts))

	return r
}

func apiRouter(h *handler.Handler, clientset *kubernetes.Clientset, logger logr.Logger, opts *Options) http.Handler {
	r := chi.NewRouter()
	if opts.RateLimitRequests != nil && opts.RateLimitDuration != nil {
		r.Use(httprate.LimitAll(*opts.RateLimitRequests, *opts.RateLimitDuration))
	}
	r.Use(middleware.Logger)
	if opts.KubernetesAuth && opts.KubernetesTrusted != nil {
		kauth := kubernetesauth.NewKubernetesAuth(clientset, opts.KubernetesTrusted, logger)
		r.Use(kauth.Handler)
	}

	r.Route("/bootstrap", func(r chi.Router) {
		r.Put("/", h.Bootstrap.Put)
		r.Delete("/", h.Bootstrap.Delete)
	})
	r.Get("/galerastate", h.GaleraState.Get)
	r.Route("/recovery", func(r chi.Router) {
		r.Put("/", h.Recovery.Put)
		r.Post("/", h.Recovery.Post)
		r.Delete("/", h.Recovery.Delete)
	})

	return r
}
