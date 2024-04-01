package main

import (
	"clubhouse/internal/database"
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"

	firebase "firebase.google.com/go/v4"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB          *database.Queries
	JWTSecret   string
	firebaseApp *firebase.App
}

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("PORT environment variable is not set")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("could not connect to database")
	}
	dbQueries := database.New(db)

	fbApp, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Fatalf("error initializing firebase: %v\n", err)
	}

	apiConf := apiConfig{DB: dbQueries, JWTSecret: jwtSecret, firebaseApp: fbApp}
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	apiRouter := chi.NewRouter()

	apiRouter.Get("/status", apiConf.handlerStatus)

	// Auth routes
	apiRouter.Post("/login", apiConf.handlerLogin)
	apiRouter.Post("/refresh", apiConf.handlerRefresh)
	apiRouter.Post("/revoke", apiConf.handlerRevoke)
	apiRouter.Post("/register", apiConf.middlewareAuth(apiConf.handlerNotificationsRegisterDevice))

	// User routes
	apiRouter.Post("/users", apiConf.handlerUsersCreate)
	apiRouter.Get("/users", apiConf.middlewareAuth(apiConf.handlerUsersGet))
	apiRouter.Get("/users/self", apiConf.middlewareAuth(apiConf.handlerUsersGetSelf))

	// Messages routes
	apiRouter.Get("/threads/{id}/messages", apiConf.middlewareAuth(apiConf.handlerMessagesGet))
	apiRouter.Post("/messages", apiConf.middlewareAuth(apiConf.handlerMessagesCreate))

	// Threads routes
	apiRouter.Post("/threads", apiConf.middlewareAuth(apiConf.handlerThreadsCreate))
	apiRouter.Get("/threads", apiConf.middlewareAuth(apiConf.handlerThreadsGet))
	apiRouter.Delete("/threads/{id}", apiConf.middlewareAuth(apiConf.handlerThreadsDelete))
	apiRouter.Get("/threads/{id}/subscribe", apiConf.middlewareAuth(apiConf.handlerUnsubscribedUsersGet))
	apiRouter.Post("/threads/{id}/subscribe", apiConf.middlewareAuth(apiConf.handlerThreadsAddUsers))
	apiRouter.Delete("/threads/{id}/subscribe", apiConf.middlewareAuth(apiConf.handlerUnsubscribeUsers))
	apiRouter.Get("/threads/{id}/members", apiConf.middlewareAuth(apiConf.handlerThreadsGetMembers))

	// Events routes
	apiRouter.Get("/events", apiConf.middlewareAuth(apiConf.handlerEventsGet))
	apiRouter.Post("/events", apiConf.middlewareAuth(apiConf.handlerEventsCreate))

	// Images routes
	apiRouter.Get("/images/{image}", apiConf.handlerImagesGet)
	apiRouter.Post("/upload", apiConf.middlewareAuth(apiConf.handlerImagesCreate))

	router.Mount("/v1", apiRouter)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	log.Printf("Listing on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
