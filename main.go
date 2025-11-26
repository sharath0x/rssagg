package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sharath0x/rssagg/internal/database"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {

	godotenv.Load(".env")
	PORT := os.Getenv("PORT")
	dbURL := os.Getenv("DB_URL")

	if PORT == "" {
		log.Fatalln("PORT not found in environment")
	}
	if dbURL == "" {
		log.Fatalln("DB_URL not found in environment")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalln("failed to establish DB connection  ", err)
	}

	db := database.New(conn)
	apiCfg := apiConfig{
		DB: db,
	}

	go StartScraping(db, 10, time.Minute)

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"links"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handler_readiness)
	v1Router.Get("/err", handler_error)

	v1Router.Post("/users", apiCfg.handleCreateUser)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handleGetUser))

	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handleCreateFeed))
	v1Router.Get("/feeds", apiCfg.handlerGetFeeds)

	v1Router.Get("/posts", apiCfg.middlewareAuth(apiCfg.handleGetPostsFromUser))

	v1Router.Post("/feedfollows", apiCfg.middlewareAuth(apiCfg.handleCreateFeedFollows))
	v1Router.Get("/feedfollows", apiCfg.middlewareAuth(apiCfg.handleGetFeedFollows))
	v1Router.Delete("/feedfollows/{feedFollowsId}", apiCfg.middlewareAuth(apiCfg.handleDeleteFeedFollows))

	router.Mount("/v1", v1Router)

	srv := http.Server{
		Handler: router,
		Addr:    ":" + PORT,
	}

	log.Println("Server is up and running on PORT: ", PORT)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
