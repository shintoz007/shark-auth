package apis

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"shark-auth/internal"
	"shark-auth/internal/accesstoken"
	"shark-auth/internal/refreshtoken"
	"shark-auth/internal/user"

	"github.com/go-redis/redis/v7"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"shark-auth/apis/handlers"
	"shark-auth/apis/middlewares"
)

func API(db *sqlx.DB, redisClient *redis.Client) http.Handler {

	userRepo := user.NewUserRepository(db)
	refreshTokenRepo := refreshtoken.NewRefreshTokenStore(db)
	accessTokenBlacklistStore := accesstoken.NewBlacklistStore(redisClient)

	tokenService := internal.TokenService{
		UserRepo:          userRepo,
		RefreshTokenStore: refreshTokenRepo,
		BlacklistStore:    accessTokenBlacklistStore,
	}
	tokenServer := handlers.NewTokenServer(tokenService)

	logrus.Info("starting server")
	r := mux.NewRouter()
	r.Use(middlewares.PanicHandlerMiddleware, middlewares.LoggingMiddleware)
	r.HandleFunc("/user", handlers.HandleUserSignup(userRepo)).Methods(http.MethodPost)

	r.HandleFunc("/token", tokenServer.HandleTokenCreate()).Methods(http.MethodPost)
	r.HandleFunc("/token", tokenServer.HandleTokenRefresh()).Methods(http.MethodPatch)
	r.HandleFunc("/token", tokenServer.HandleTokenDelete()).Methods(http.MethodDelete)

	r.HandleFunc("/welcome", handlers.HandleWelcome(accessTokenBlacklistStore)).Methods(http.MethodGet)
	r.Path("/metrics").Handler(promhttp.Handler())

	return r
}
