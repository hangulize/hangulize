package webapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {
	router := gin.New()

	v2 := router.Group("/v2")
	v2Init(v2)

	http.Handle("/", router)
}
