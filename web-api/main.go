package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"
)

func init() {
	router := gin.New()

	v2 := router.Group("/v2")
	v2Init(v2)

	v1 := router.Group("/v1")
	v1Init(v1)

	http.Handle("/", router)
}

func main() {
	appengine.Main()
}
