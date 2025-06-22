package controller

import (
	"github.com/gin-gonic/gin"

	"log/slog"
	"medrussia_news_bot/internal/config"
	"net/http"
)

type BotController interface {
	BotWebhookHandler(c *gin.Context)
}

type RestController struct {
	router           *gin.Engine
	cfg              *config.Config
	logger           *slog.Logger
	botApiController BotController
}

func NewRestController(
	cfg *config.Config,
	logger *slog.Logger,
	botApiController BotController,
) *RestController {
	router := gin.New()
	router.Use(gin.Recovery())
	router.POST("/"+cfg.Bot.Token+"/", botApiController.BotWebhookHandler)

	rest := &RestController{
		router:           router,
		cfg:              cfg,
		logger:           logger,
		botApiController: botApiController,
	}

	return rest
}

func (r RestController) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}
