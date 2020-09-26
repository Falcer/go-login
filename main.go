package main

import (
	"log"

	"github.com/Falcer/go-login/routes"
)

func main() {

	userRoute := routes.NewUserRouter()
	log.Fatal(userRoute.Listen(":8080"))
}
