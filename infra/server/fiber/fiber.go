package fiber

import (
	"fmt"
	"log"
	"nft/config"

	"github.com/gofiber/fiber/v2"
)

type Server struct {
	App *fiber.App
}

func (s *Server) ListenAndServe() error {
	go func() {
		if err := s.App.Listen(fmt.Sprintf(":%s", config.C().App.Http.Port)); err != nil {
			log.Println(err)
		}
		log.Println("http server started")
	}()
	return nil
}

func (s *Server) Shutdown() error {
	return s.App.Shutdown()
}
