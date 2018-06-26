package main

import (
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/sublee/hangulize2/hangulize"
)

func v1Init(r *gin.RouterGroup) {
	r.GET("/", v1Hangulize)
	r.GET("/example", v1Example)
	r.GET("/langs", v1Langs)
}

// v1Hangulize handles: GET /?lang=:lang&word=:word
func v1Hangulize(c *gin.Context) {
	lang := c.Query("lang")
	word := c.Query("word")

	transcribed := hangulize.Hangulize(lang, word)
	spec, _ := hangulize.LoadSpec(lang)

	c.JSON(http.StatusOK, gin.H{
		"result":  transcribed,
		"word":    word,
		"lang":    v1PackSpec(spec),
		"success": true,
		"reason":  nil,
	})
}

func v1PackSpec(s *hangulize.Spec) *gin.H {
	return &gin.H{
		"code":     s.Lang.ID,
		"name":     s.Lang.English,
		"label":    s.Lang.English,
		"iso639-1": s.Lang.Codes[0],
		"iso639-2": nil,
		"iso639-3": s.Lang.Codes[1],
	}
}

// v1Example handles: GET /example?lang=:lang
func v1Example(c *gin.Context) {
	lang := c.Query("lang")
	spec, ok := hangulize.LoadSpec(lang)

	if !ok {
		c.AbortWithStatus(http.StatusNotFound)
	}

	i := rand.Intn(len(spec.Test))
	example := spec.Test[i]

	c.JSON(http.StatusOK, gin.H{
		"result":  example[1],
		"word":    example[0],
		"lang":    v1PackSpec(spec),
		"success": true,
		"reason":  nil,
	})
}

// v1Langs handles: GET /langs
func v1Langs(c *gin.Context) {
	langs := hangulize.ListLangs()

	packedSpecs := make([]gin.H, len(langs))
	for i, lang := range langs {
		spec, _ := hangulize.LoadSpec(lang)
		packedSpecs[i] = *v1PackSpec(spec)
	}

	c.JSON(http.StatusOK, gin.H{
		"length":  len(langs),
		"langs":   packedSpecs,
		"success": true,
		"reason":  nil,
	})
}
