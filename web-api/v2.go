package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/sublee/hangulize2/hangulize"
)

func v2Init(r *gin.RouterGroup) {
	r.GET("/version", v2Version)
	r.GET("/hangulized/:lang/:word", v2Hangulized)
	r.GET("/specs", v2Specs)
	r.GET("/specs/:path", v2SpecHGL)
}

// v2Version handles: GET /version
//
// Accept: application/json, text/plain by default
//
// It returns the version of the "hangulize" package.
//
func v2Version(c *gin.Context) {
	switch c.NegotiateFormat(gin.MIMEJSON) {

	case gin.MIMEJSON:
		c.JSON(http.StatusOK, gin.H{
			"version": hangulize.Version,
		})

	default:
		c.String(http.StatusOK, hangulize.Version)
	}
}

// v2Hangulized handles: GET /hangulized/:lang/:word[?trace]
//
// Accept: application/json, text/plain by default
//
// It hangulizes a word.
//
func v2Hangulized(c *gin.Context) {
	lang := c.Param("lang")
	word := c.Param("word")

	transcribed := hangulize.Hangulize(lang, word)

	switch c.NegotiateFormat(gin.MIMEJSON) {

	case gin.MIMEJSON:
		c.JSON(http.StatusOK, gin.H{
			"lang":        lang,
			"word":        word,
			"transcribed": transcribed,
		})

	default:
		c.String(http.StatusOK, transcribed)
	}
}

// v2Specs handles: GET /specs
//
// Accept: application/json, text/plain by default
//
// It returns the list of available specs.
//
func v2Specs(c *gin.Context) {
	switch c.NegotiateFormat(gin.MIMEJSON) {

	case gin.MIMEJSON:
		specs := make(gin.H)

		for _, lang := range hangulize.ListLangs() {
			spec, _ := hangulize.LoadSpec(lang)
			specs[lang] = v2PackSpec(spec)
		}

		c.JSON(http.StatusOK, gin.H{
			"specs": specs,
		})

	default:
		c.String(http.StatusOK, strings.Join(hangulize.ListLangs(), "\n"))
	}
}

func v2PackSpec(s *hangulize.Spec) *gin.H {
	test := make([]gin.H, len(s.Test))

	for i, example := range s.Test {
		test[i] = gin.H{
			"word":        example[0],
			"transcribed": example[1],
		}
	}

	return &gin.H{
		"lang": gin.H{
			"id":      s.Lang.ID,
			"codes":   s.Lang.Codes,
			"english": s.Lang.English,
			"korean":  s.Lang.Korean,
			"script":  s.Lang.Script,
		},

		"config": gin.H{
			"authors": s.Config.Authors,
			"stage":   s.Config.Stage,
		},

		"test": test,
	}
}

// v2SpecHGL handles: GET /specs/:lang.hgl
//
// Accept: text/vnd.hgl
//
// It serves the HGL source of the spec.
//
func v2SpecHGL(c *gin.Context) {
	// Should look like "ita.hgl".
	path := c.Param("path")

	if !strings.HasSuffix(path, ".hgl") {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	lang := strings.TrimSuffix(path, ".hgl")
	spec, ok := hangulize.LoadSpec(lang)

	if !ok {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.Header("Content-Type", "text/vnd.hgl")
	c.String(http.StatusOK, spec.Source)
}
