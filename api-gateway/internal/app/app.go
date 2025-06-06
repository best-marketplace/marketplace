package app

import (
	"api-gateway/internal/config"
	"api-gateway/internal/lib/auth"
	geturl "api-gateway/internal/lib/handlers/getUrl"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type App struct {
	httpServer *http.Server
	log        *slog.Logger
	routes     map[string]string
}

func NewApp(log *slog.Logger, cfg *config.Config) *App {
	routes := cfg.GetRoutes()

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.HandleFunc("/proxy/*", func(w http.ResponseWriter, r *http.Request) {
		requestPath := r.URL.Path
		targetPath := strings.TrimPrefix(requestPath, "/proxy")

		r.URL.Path = targetPath

		serviceURL, matchedPrefix, err := geturl.GetServiceURL(r, routes)
		if err != nil {
			log.Error("Servce not found",
				slog.String("path", targetPath),
				slog.Any("error", err))
			http.Error(w, "Servce not found", http.StatusBadGateway)
			return
		}

		reverseProxy := httputil.NewSingleHostReverseProxy(serviceURL)

		reverseProxy.Director = func(req *http.Request) {
			req.URL.Scheme = serviceURL.Scheme
			req.URL.Host = serviceURL.Host
			req.URL.Path = targetPath
			req.URL.RawQuery = r.URL.RawQuery
			req.Header = r.Header.Clone()

			if strings.HasPrefix(matchedPrefix, "/product") ||
				strings.HasPrefix(matchedPrefix, "/products") ||
				strings.HasPrefix(matchedPrefix, "/search") ||
				strings.HasPrefix(matchedPrefix, "/comment") ||
				strings.HasPrefix(matchedPrefix, "/comments") {

				authHeader := r.Header.Get("Authorization")
				if authHeader != "" {
					tokenInfo, verifyErr := auth.VerifyToken(routes["/auth"], authHeader)
					if verifyErr == nil && tokenInfo != nil {
						log.Info("Adding user info to request",
							slog.String("user_id", tokenInfo.UserID),
							slog.String("username", tokenInfo.Username))
						req.Header.Set("X-User-ID", tokenInfo.UserID)
						req.Header.Set("X-Username", tokenInfo.Username)
					} else {
						log.Warn("Failed to verify token",
							slog.String("error", verifyErr.Error()))
					}
				}
			}

			log.Info("request to service",
				slog.String("host", serviceURL.Host),
				slog.String("path", targetPath))
		}

		reverseProxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			log.Error("error ", slog.Any("error", err))
			http.Error(w, fmt.Sprintf("error: %v", err), http.StatusBadGateway)
		}

		reverseProxy.ServeHTTP(w, r)
	})

	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%s", cfg.ProxyPort),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return &App{
		httpServer: srv,
		log:        log,
		routes:     routes,
	}
}

func (a *App) Run() error {
	a.log.Info("Start api-gateway", slog.String("address", a.httpServer.Addr))
	return a.httpServer.ListenAndServe()
}

func (a *App) Shutdown(ctx context.Context) error {
	return a.httpServer.Shutdown(ctx)
}
