package controller

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"urlshortner/internal/model"
	"urlshortner/internal/service"
)

type ShortenController struct {
	shortenService *service.ShortnerService
	logger         *zap.Logger
}

func NewShortenController(svc *service.ShortnerService, logger *zap.Logger) *ShortenController {
	return &ShortenController{
		shortenService: svc,
		logger:         logger,
	}
}

func (s *ShortenController) CreateShortLink(ctx *fiber.Ctx) error {
	var req model.Request
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	if req.URL == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "URL must not be nil"})
	}

	resp, err := s.shortenService.CreateShortLink(ctx.Context(), req.URL)
	if err != nil {
		s.logger.Error("some server CreateShortLink error:", zap.Error(err))
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}

func (s *ShortenController) GetOriginalByShort(ctx *fiber.Ctx) error {
	shortenURL := ctx.Params("shortenURL")

	if shortenURL == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ShortenURL must not be nil"})
	}

	resp, err := s.shortenService.GetOriginalByShort(ctx.Context(), shortenURL)
	if err != nil {
		s.logger.Error("some server GetOriginalByShort error:", zap.Error(err))
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}
