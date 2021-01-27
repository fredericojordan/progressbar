package main

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"text/template"
)

var progressTemplate string = `<?xml version="1.0" encoding="UTF-8"?>
<svg width="118" height="20" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" preserveAspectRatio="xMidYMid">
  <linearGradient id="a" x2="0" y2="100%">
    <stop offset="0" stop-color="#bbb" stop-opacity=".1"/>
    <stop offset="1" stop-opacity=".1"/>
  </linearGradient>

  <rect rx="4" x="0" width="118" height="20" fill="#428bca"/>
  <rect rx="4" x="58" width="60" height="20" fill="#555" />
  <rect rx="4" x="58" width="{{.ProgressWidth}}" height="20" fill="{{.ProgressColor}}" />

    <path fill="{{.ProgressColor}}" d="M58 0h4v20h-4z" />

  <rect rx="4" width="118" height="20" fill="url(#a)" />

    <g fill="#fff" text-anchor="left" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="11">
      <text x="4" y="15" fill="#010101" fill-opacity=".3">
        progress
      </text>
      <text x="4" y="14">
        progress
      </text>
    </g>

  <g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="11">
    <text x="88" y="15" fill="#010101" fill-opacity=".3">
      {{.Progress}}%
    </text>
    <text x="88" y="14">
      {{.Progress}}%
    </text>
  </g>
</svg>`

type templateParams struct {
	Progress      int
	ProgressWidth int
	ProgressColor string
}

func main() {
	setupServer().Run()
}

func setupServer() *gin.Engine {
	r := gin.Default()
	r.GET("/:progress/", renderProgress)
	r.GET("/", redirectToGithub)
	return r
}

func redirectToGithub(context *gin.Context) {
	context.Redirect(http.StatusFound, "https://github.com/fredericojordan/progressbar")
}

func renderProgress(context *gin.Context) {
	progressStr := context.Param("progress")

	progress, err := strconv.ParseInt(progressStr, 10, 32)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	text, err := renderSVG(int(progress))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	context.Header("Content-Type", "image/svg+xml")
	context.String(http.StatusOK, text)
}

func progressColor(progress int) string {
	if progress < 30 {
		return "#d9534f"
	}
	if progress < 70 {
		return "#f0ad4e"
	}
	return "#5cb85c"
}

func renderSVG(progress int) (string, error) {

	params := templateParams{
		Progress:      progress,
		ProgressWidth: int(60 * progress / 100),
		ProgressColor: progressColor(progress),
	}

	progTemplate, err := template.New("progress").Parse(progressTemplate)

	if err != nil {
		return "", err
	}

	buffer := new(bytes.Buffer)

	if err := progTemplate.Execute(buffer, params); err != nil {
		return "", err
	}

	return buffer.String(), nil
}
