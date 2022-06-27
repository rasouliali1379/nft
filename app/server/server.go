package server

import (
	"fmt"
	"log"
	"maskan/config"
	auth "maskan/src/auth/contract"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

type Server struct {
	app *fiber.App
}

type ControllerContainer struct {
	fx.In
	AuthController auth.IAuthController
}

func New(cc ControllerContainer) IServer {

	app := fiber.New()
	router := app.Group(config.C().App.BaseURL)
	
	authRouter := router.Group("/auth")
	authRouter.Post("/signup", cc.AuthController.SignUp)
	authRouter.Post("/login", cc.AuthController.Login)

	return &Server{
		app: app,
	}
}

func (s Server) ListenAndServe() error {
	go func() {
		if err := s.app.Listen(fmt.Sprintf(":%s", config.C().App.Http.Port)); err != nil {
			log.Println(err)
		}
	}()
	return nil
}

func (s Server) Shutdown() error {
	return s.app.Shutdown()
}

func corsHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", config.C().App.Http.Cors)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS, HEAD")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, TE, User-Agent, Cache-Control, Sec-Fetch-Dest, Sec-Fetch-Mode, Sec-Fetch-Site, Referer, Content-Type, Pragma, Connection, Content-Length, Accept-Language, Accept-Encoding, Authorization, ResponseType")

		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}
