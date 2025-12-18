package server

import (
	"PilaiteProject/internal/dbConfig"
	"PilaiteProject/internal/handler"
	"PilaiteProject/internal/service"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
)

func setupRoutes(router *chi.Mux, conn *dbConfig.Connection, sessionManager *scs.SessionManager) {

	userService := service.NewUserService(conn.Queries)

	spotService := service.NewSpotService(conn.Queries)

	spotHandler := handler.NewSpotHandler(spotService)

	authHandler := handler.NewAuthHandler(userService, sessionManager)

	authMiddleware := NewAuthMiddleware(sessionManager)

	setupSpotRoutes(router, spotHandler, authMiddleware)

	setupPublicRoutes(router)
	setupAuthRoutes(router, authHandler, authMiddleware)

}

func setupPublicRoutes(router *chi.Mux) {
	router.Route("/public", func(router chi.Router) {
		router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"healthy"}`))
		})
	})
}

func setupAuthRoutes(router *chi.Mux, authHandler *handler.AuthHandler, authMiddleware *AuthMiddleware) {
	//No authentication required
	router.Group(func(router chi.Router) {
		router.Use(authMiddleware.RequireGuest)
		router.Post("/register", authHandler.Register)
		router.Post("/login", authHandler.Login)
	})

	router.Group(func(router chi.Router) {
		router.Use(authMiddleware.RequireAuth)
		router.Get("/me", authHandler.GetCurrentUser)
		router.Get("/logout", authHandler.Logout)
	})
}

func setupSpotRoutes(router *chi.Mux, spotHandler *handler.SpotHandler, authMiddleware *AuthMiddleware) {
	router.Route("/spots", func(r chi.Router) {
		//r.Group(func(router chi.Router) {
		//	//router.Use(authMiddleware.RequireAdmin)
		//	r.Post("/", spotHandler.CreateSpot)
		//})

		//public routes
		r.Get("/", spotHandler.GetPublicSpotsWithDetails) // no auth, but will filter secret spots
		r.Get("/public/category/{category}", spotHandler.GetPublicSpotsByCategoryWithDetails)
		r.Get("/{id}", spotHandler.GetSpotById)
		//r.Get("/{id}/location", spotHandler.GetSpotWithLocation)

		// Secret spots category - requires auth
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.RequireAuth)
			r.Get("/all", spotHandler.GetSpotsWithDetails)
			r.Get("/category/{category}", spotHandler.GetSpotsByCategoryWithDetails)
			//r.Get("/category/secret", spotHandler.GetSecretSpotsByCategory)
		})

	})
}

//func setupUserRoutes(router *chi.Mux, userHandler *handler.UserHandler, authMiddleware *authorization.AuthMiddleware) {
//	router.Group(func(r chi.Router) {
//		r.Use(authMiddleware.RequireAuth)
//		r.Route("/users", func(r chi.Router) {
//			r.Get("/", userHandler.GetAllUsers)
//			r.Get("/{id}", userHandler.GetUserById)
//		})
//	})
//}
// Admin routes (admin only)
//func setupAdminRoutes(router *chi.Mux, adminHandler *handler.AdminHandler, authMiddleware *authorization.AuthMiddleware) {
//	router.Group(func(r chi.Router) {
//		r.Use(authMiddleware.RequireAdmin)
//		r.Route("/admin", func(r chi.Router) {
//			r.Get("/dashboard", adminHandler.Dashboard)
//			r.Get("/users", adminHandler.GetAllUsers)
//			r.Delete("/users/{id}", adminHandler.DeleteUser)
//		})
//	})
//}
