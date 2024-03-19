package api

import (
	"errors"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	recoverMiddleware "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/kbgod/coinbot/config"
	"github.com/kbgod/coinbot/internal/observer"
	"github.com/kbgod/coinbot/internal/service"
	"net/http"
)

func Run(svc *service.Service) {
	router := newRouter(svc.Observer, svc.CFG)

	h := newHandler(svc)

	router.Post("/authorize", h.authorize)

	private := router.Group("/", h.userMiddleware)
	//private.Post("/makeEnergy", h.makeEnergy)
	//private.Post("/mine", h.mine)
	private.Get("/getMe", h.getMe)
	private.Post("/mine", h.mine)
	private.Get("/boosts", h.boosts)
	private.Post("/boosts", h.purchaseBoost)
	private.Post("/dailyBooster", h.openDailyBooster)
	private.Get("/leaderboard/daily", h.dailyLeaderboard)
	private.Get("/leaderboard/monthly", h.monthlyLeaderboard)
	private.Get("/leaderboard", h.leaderboard)
	private.Get("/channels", h.channels)
	if svc.CFG.API.TLS {
		if err := router.ListenTLS(svc.CFG.API.Host+":"+svc.CFG.API.Port, "cert.pem", "private.key"); err != nil {
			svc.Observer.Logger.Fatal().Err(err).Msg("failed to start api server")
		}
	}
	if err := router.Listen(svc.CFG.API.Host + ":" + svc.CFG.API.Port); err != nil {
		svc.Observer.Logger.Fatal().Err(err).Msg("failed to start api server")
	}
}

func newRouter(observer *observer.Observer, cfg *config.Config) *fiber.App {
	router := fiber.New(fiber.Config{
		AppName:               "API",
		ErrorHandler:          errorHandler(observer),
		DisableStartupMessage: !cfg.Debug,
	})
	router.Use(recoverMiddleware.New(recoverMiddleware.Config{
		EnableStackTrace: true,
	}))
	router.Use(cors.New())

	router.Use(requestid.New())

	router.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: observer.Logger,
		Fields: append(fiberzerolog.ConfigDefault.Fields, fiberzerolog.FieldRequestID),
	}))

	router.Use(func(ctx *fiber.Ctx) error {
		err := ctx.Next()
		ctx.Set("content-type", "application/json; charset=utf-8")
		return err
	})

	return router
}

func errorHandler(observer *observer.Observer) fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		{
			var e *fiber.Error
			if errors.As(err, &e) {
				ctx.Status(e.Code)
				return newMessageResponse(ctx, e.Code, e.Message)
			}
		}
		{
			var e *ResponseError
			if errors.As(err, &e) {
				ctx.Status(e.StatusCode)
				return newMessageResponse(ctx, e.Code, e.Message)
			}
		}
		observer.Logger.Error().
			Err(err).
			Ctx(ctx.UserContext()).
			Str("request-id", ctx.Locals(requestid.ConfigDefault.ContextKey).(string)).
			Str("url", ctx.Path()).
			Msg("unhandled error")

		ctx.Context().SetStatusCode(http.StatusInternalServerError)
		return nil
	}
}
