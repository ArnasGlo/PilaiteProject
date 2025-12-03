package server

import (
	"PilaiteProject/internal/db_config"
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	httpServer     *http.Server
	sessionManager *scs.SessionManager
}

type ServerConfig struct {
	Host string
	Port string
}

func NewServer(serverConfiq ServerConfig, conn *db_config.Connection) *Server {

	sessionManager := initSessionManager()

	router := chi.NewRouter()

	applyGlobalMiddleware(router, sessionManager)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		http.ServeFile(w, r, "./frontend/index.html")
	})

	router.Handle("/static/*", http.StripPrefix("/static/", cacheControlFileServer(http.Dir("./frontend/static"))))

	setupRoutes(router, conn, sessionManager)

	// SPA fallback: if no other route matched (and not an API/static route), serve index.html.
	// This lets client-side routes like /dashboard work on refresh.
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		// Don't return index for requests that look like they are requesting a real file
		// e.g. /something.png or /file.ext — optional safety:
		if looksLikeFile(r.URL.Path) {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		http.ServeFile(w, r, "./frontend/index.html")
	})

	addr := fmt.Sprintf("%s:%s", serverConfiq.Host, serverConfiq.Port)
	httpServer := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	return &Server{
		httpServer:     httpServer,
		sessionManager: sessionManager,
	}
}

func initSessionManager() *scs.SessionManager {
	sessionManager := scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Name = "session_id"
	sessionManager.Cookie.HttpOnly = true                 // Prevent JavaScript access
	sessionManager.Cookie.Secure = false                  // Set to true in production with HTTPS
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode // CSRF protection
	return sessionManager
}

func (s *Server) Start() error {
	fmt.Printf("Server starting on %s\n", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	fmt.Println("Server shutting down...")
	return s.httpServer.Shutdown(ctx)
}

// cacheControlFileServer wraps http.FileServer and sets cache headers for assets.
// Files (CSS/JS/images) will get a long max-age; override or remove if you prefer different cache policy.
func cacheControlFileServer(fs http.FileSystem) http.Handler {
	fileServer := http.FileServer(fs)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only set long cache for GET/HEAD requests
		if r.Method == http.MethodGet || r.Method == http.MethodHead {
			// Example: cache for 7 days
			w.Header().Set("Cache-Control", "public, max-age=10, immutable")
		}
		fileServer.ServeHTTP(w, r)
	})
}

// looksLikeFile tries to detect if the path is a request for a static file (has an extension).
// If you want all unknown paths to map to index.html (even /foo.png if missing), then remove this check.
func looksLikeFile(path string) bool {
	ext := filepath.Ext(path)
	return ext != ""
}
