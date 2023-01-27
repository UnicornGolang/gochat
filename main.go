package main

import (
	"gochat/router"
	"gochat/utils"
)

func main() {
  utils.InitConfig()
  utils.InitMySQL()
  utils.InitRedis()

	r := router.Router() // router.Router()
	r.Run(":8081")       // listen and server on 0.0.0.0:8080 （for windows "localhost:8080")
}
