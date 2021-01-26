package main

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"text/template"
)

var progress_template string = `<?xml version="1.0" encoding="UTF-8"?>
<svg width="118" height="20" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" preserveAspectRatio="xMidYMid">
  <linearGradient id="a" x2="0" y2="100%">
    <stop offset="0" stop-color="#bbb" stop-opacity=".1"/>
    <stop offset="1" stop-opacity=".1"/>
  </linearGradient>

  <rect rx="4" x="0" width="118" height="20" fill="#428bca"/>
  <rect rx="4" x="58" width="60" height="20" fill="#555" />
  <rect rx="4" x="58" width="{{.Progress_width}}" height="20" fill="{{.Progress_color}}" />

    <path fill="{{.Progress_color}}" d="M58 0h4v20h-4z" />

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

type template_params struct {
	Progress       int
	Progress_width int
	Progress_color string
}

func main() {
	r := gin.Default()
	r.GET("/:progress/", render_progress)
	r.Run()
}

func render_progress(context *gin.Context) {
	progress_str := context.Param("progress")

	progress, err := strconv.ParseInt(progress_str, 10, 32)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	text, err := render_svg(int(progress))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	context.Header("Content-Type", "image/svg+xml")
	context.String(http.StatusOK, text)
}

func progress_color(progress int) string {
	if progress < 30 {
		return "#d9534f"
	}
	if progress < 70 {
		return "#f0ad4e"
	}
	return "#5cb85c"
}

func render_svg(progress int) (string, error) {

	params := template_params{
		Progress:       progress,
		Progress_width: int(60 * progress / 100),
		Progress_color: progress_color(progress),
	}

	prog_template, err := template.New("progress").Parse(progress_template)

	if err != nil {
		return "", err
	}

	buffer := new(bytes.Buffer)

	if err := prog_template.Execute(buffer, params); err != nil {
		return "", err
	}

	return buffer.String(), nil
}
