package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"os"

	logger "go-rengan/pkg/logger"
	todo_http "go-rengan/todo/delivery/http"
	response "go-rengan/utils/response"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/riandyrn/otelchi"
	"github.com/sirupsen/logrus"
)

type HTTPServer interface {
	PrintAllRoutes()
	Run() error
	GracefulStop(ctx context.Context) error
	GetRouter() *chi.Mux
}

type HTTPServerImpl struct {
	router *chi.Mux
	svr    *http.Server
	logger logger.Logger
}

func New(
	logger logger.Logger,
	todoHandler todo_http.HTTPHandler,
) HTTPServer {
	router := chi.NewRouter()
	router.Use(otelchi.Middleware(os.Getenv("APP_NAME"), otelchi.WithChiRoutes(router)))
	router.Use(
		render.SetContentType(render.ContentTypeJSON), // Set content-Type headers as application/json
		middleware.Logger,          // Log API request calls
		middleware.Compress(5),     // Compress results, mostly gzipping assets and json
		middleware.RedirectSlashes, // Redirect slashes to no slash URL versions
		middleware.Recoverer,       // Recover from panics without crashing server
	)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, response.H{
			"success": "true",
			"code":    200,
			"message": fmt.Sprintf("Services %s run properly", os.Getenv("APP_NAME")),
		})
	})

	// Register TodoHTTPHandler routes
	todoHandler.RegisterRoutes(router)

	s := &HTTPServerImpl{
		router: router,
		logger: logger,
	}

	s.PrintAllRoutes()

	return s
}

// PrintAllRoutes - Walk and print out all routes
func (hs *HTTPServerImpl) PrintAllRoutes() {
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		logrus.Printf("%s %s\n", method, route)
		return nil
	}
	router := hs.GetRouter()
	if err := chi.Walk(router, walkFunc); err != nil {
		hs.logger.Error(err)
	}
}

// Run - running server
func (hs *HTTPServerImpl) Run() error {
	addr := fmt.Sprintf("%s%s", ":", os.Getenv("PORT"))
	logrus.Infoln("HTTP server listening on", addr)

	router := hs.GetRouter()
	hs.svr = &http.Server{
		Addr:    addr,
		Handler: router,
	}

	err := hs.svr.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

// GracefulStop the server
func (hs *HTTPServerImpl) GracefulStop(ctx context.Context) error {
	return hs.svr.Shutdown(ctx)
}

func (hs *HTTPServerImpl) GetRouter() *chi.Mux {
	return hs.router
}
