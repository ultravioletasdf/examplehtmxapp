package main

import (
	"embed"
	"log"

	"examplehtmxapp/routes"
	"examplehtmxapp/utils"

	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"
)

//go:embed assets/*
var assets embed.FS

func main() {
	cfg := utils.GetConfig()

	db := utils.ConnectDatabase(&cfg)
	defer db.Close()

	utils.CreateEmailClient(&cfg)

	app := fiber.New(fiber.Config{Prefork: cfg.Prefork, EnableTrustedProxyCheck: true})
	routes.Add(app, assets, db.Executor, &cfg)
	log.Fatal(app.Listen(cfg.ListenAddress))
}
