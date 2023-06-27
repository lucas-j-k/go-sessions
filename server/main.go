package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/lucas-j-k/go-sessions/user"
	"github.com/redis/go-redis/v9"
	"github.com/rs/cors"
	"github.com/spf13/viper"
)

func main() {

	// initialize env vars and SQL connection
	viper.SetConfigFile("../.env")
	viper.ReadInConfig()

	port := viper.Get("PORT")
	sqlPass := viper.Get("MYSQL_PASS")
	sqlHost := viper.Get("MYSQL_HOST")
	sqlUser := viper.Get("MYSQL_USER")
	redisPass := viper.Get("REDIS_PASSWORD")

	// set up CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true, // accept cookies
		Debug:            true, // should be development env only
	})

	// initialize Database connection and database services
	dataSourceName := fmt.Sprintf("%v:%v@tcp(%v:3306)/default_db?parseTime=True", sqlUser, sqlPass, sqlHost)
	connection := sqlx.MustConnect("mysql", dataSourceName)

	userService := user.UserService{
		Db: connection,
	}

	// Initialize Redis client and services
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",             // TODO - env vars and concatenate
		Password: fmt.Sprintf("%v", redisPass), // no password set
		DB:       0,                            // use default DB
	})

	redisSessionManager := user.RedisSessionManager{
		Client: redisClient,
	}

	// setup Chi router and global middlewares
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(c.Handler)

	// initialise custom middleware. This is in an interface so we can inject our redis session manager
	cacheMiddleware := user.CacheMiddleware{SessionManager: &redisSessionManager}

	// protected test - delete this
	r.Route("/protected", func(r chi.Router) {
		r.Use(cacheMiddleware.SessionGuard)
		r.Get("/", user.Protected(&userService))
	})

	r.Post("/users/signup", user.Signup(&userService))
	r.Post("/users/login", user.Login(&userService, redisSessionManager))
	r.Post("/users/logout", user.Logout(&userService, redisSessionManager))

	////
	// HEALTHCHECK
	////
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Pong"))
	})

	fmt.Printf("Server running on port [%v]\n\n", port)
	http.ListenAndServe(fmt.Sprintf(":%v", port), r)
}
