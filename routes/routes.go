package routes

import (
	"context"
	"embed"
	"examplehtmxapp/frontend"
	sqlc "examplehtmxapp/sql"
	"examplehtmxapp/utils"
	"log"
	"net/http"
	"time"

	"github.com/a-h/templ"
	_snowflake "github.com/bwmarrin/snowflake"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var ctx = context.Background()
var executor *sqlc.Queries
var snowflake *_snowflake.Node
var cfg *utils.Config

func Add(app *fiber.App, assets embed.FS, _executor *sqlc.Queries, _cfg *utils.Config) {
	executor = _executor
	cfg = _cfg
	var err error
	snowflake, err = _snowflake.NewNode(1)
	if err != nil {
		log.Fatalln(err)
	}

	app.Use(limiter.New(limiter.Config{LimiterMiddleware: limiter.SlidingWindow{}, Max: 20, Expiration: 3 * time.Second}), logger.New(), helmet.New(), recover.New())

	app.Get("/", landing)
	app.Get("/sign/out", signOut)
	app.Get("/sign/in", signIn)
	app.Post("/sign/in", postSignIn)
	app.Get("/sign/up", signUp)
	app.Post("/sign/up", postSignUp)
	app.Get("/onboarding", onboarding)
	app.Put("/onboarding", putOnboarding)

	app.Use(etag.New(), compress.New(), adaptor.HTTPHandler(http.FileServer(http.FS(assets))))
	app.Use(monitor.New())
}
func landing(c *fiber.Ctx) error {
	if c.Cookies("session") != "" {
		usr, _ := executor.GetUserFromSession(ctx, c.Cookies("session"))
		return Render(c, frontend.Landing(usr))
	}
	return Render(c, frontend.Landing(sqlc.User{}))
}

func Render(c *fiber.Ctx, component templ.Component) error {
	c.Set("Content-Type", "text/html")
	return component.Render(c.Context(), c.Response().BodyWriter())
}
