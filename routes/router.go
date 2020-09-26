package routes

import (
	"log"

	"github.com/Falcer/go-login/model"
	"github.com/Falcer/go-login/repository/mongo"
	"github.com/Falcer/go-login/service"
	"github.com/gofiber/fiber/v2"
	"github.com/kelseyhightower/envconfig"
)

type (
	// UserApp interface
	UserApp interface {
		getUser(c *fiber.Ctx) error
		login(c *fiber.Ctx) error
		register(c *fiber.Ctx) error
	}
	userAppImpl struct {
		service service.UserService
	}
	userConfig struct {
		MongoURL string `envconfig:"MONGO_URL"`
	}
)

func (a *userAppImpl) getUser(c *fiber.Ctx) error {
	users, err := a.service.Users()
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": true,
			"message": "Failed to Get Users",
		})
	}
	return c.Status(200).JSON(&fiber.Map{
		"success": true,
		"message": "Failed to Get Users",
		"data":    users,
	})
}

func (a *userAppImpl) login(c *fiber.Ctx) error {
	userReq := new(model.User)
	if err := c.BodyParser(userReq); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": true,
			"message": "Failed to Decode Request",
		})
	}

	token, err := a.service.Login(*userReq)
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": true,
			"message": "Failed to Login",
		})
	}

	return c.Status(200).JSON(&fiber.Map{
		"success": true,
		"message": "Login Done 🔥",
		"data":    *token,
	})
}

func (a *userAppImpl) register(c *fiber.Ctx) error {
	userReq := new(model.User)
	if err := c.BodyParser(userReq); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": "Failed to Decode Request",
		})
	}
	if err := a.service.Register(*userReq); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": "Failed to Register",
		})
	}

	return c.Status(201).JSON(&fiber.Map{
		"success": true,
		"message": "Register Done 🔥",
	})
}

// NewUserRouter function
func NewUserRouter() *fiber.App {

	var cfg userConfig

	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	userRepo, _ := mongo.NewMongoRepository(cfg.MongoURL)
	userService := service.NewUserService(userRepo)
	app := &userAppImpl{userService}

	route := fiber.New()
	route.Get("/user", app.getUser)
	route.Post("/user/login", app.login)
	route.Post("/user/register", app.register)

	return route
}
