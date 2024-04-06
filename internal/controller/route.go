package controller

import "github.com/gofiber/fiber/v2"

func (s *ShortenController) Register(router fiber.Router) {
	router.Post("/", s.CreateShortLink)
	router.Get("/:shortenURL", s.GetOriginalByShort)

}

func (s *ShortenController) Name() string {
	return ""
}
