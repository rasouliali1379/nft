package server

import (
	"net/http"
	"nft/config"
	"nft/contract"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"go.uber.org/fx"
	fiberapp "nft/client/server/fiber"
)

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

type ControllerContainer struct {
	fx.In
	JwtMiddleware      contract.IJwtMiddleware
	AuthController     contract.IAuthController
	UserController     contract.IUserController
	CategoryController contract.ICategoryController
}

func New(cc ControllerContainer) contract.IServer {

	app := fiber.New()
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))
	app.Get("/swagger/*", swagger.HandlerDefault) // default

	router := app.Group(config.C().App.BaseURL)

	authRouter := router.Group("/auth")
	authRouter.Post("/signup", cc.AuthController.SignUp)
	authRouter.Post("/login", cc.AuthController.Login)
	authRouter.Post("/refresh", cc.AuthController.Refresh)
	authRouter.Post("/verify-email", cc.AuthController.VerifyEmail)
	authRouter.Post("/resend-email", cc.AuthController.ResendEmail)

	userRouter := router.Group("/user")
	userRouter.Get("/", cc.UserController.GetAllUsers)
	userRouter.Get("/:id", cc.UserController.GetUser)
	userRouter.Post("/", cc.UserController.AddUser)
	userRouter.Patch("/:id", cc.UserController.UpdateUser)
	userRouter.Delete("/:id", cc.UserController.DeleteUser)

	categoryRouter := router.Group("/category")
	categoryRouter.Use(cc.JwtMiddleware.Handle)
	categoryRouter.Get("/", cc.CategoryController.GetAllCategories)
	categoryRouter.Get("/:id", cc.CategoryController.GetCategory)
	categoryRouter.Post("/", cc.CategoryController.AddCategory)
	categoryRouter.Patch("/:id", cc.CategoryController.UpdateCategory)
	categoryRouter.Delete("/:id", cc.CategoryController.DeleteCategory)

	return &fiberapp.Server{
		App: app,
	}
}
