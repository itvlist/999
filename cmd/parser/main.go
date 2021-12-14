package main

import (
	"log"
	"os"

	"github.com/jgsqware/iptv-parser/handlers"
	"github.com/jgsqware/iptv-parser/models"
	"github.com/labstack/echo"
)

var tplt = `#EXTM3U
{{range .}}#EXTINF:-1 tvg-id="{{.TVGID}}" tvg-name="{{.TVGName}}" tvg-logo="{{.TVGLogo}}" group-title="{{.GroupTitle}}",{{.Name}}
{{.URL}}
{{end}}
`

func main() {

	channels, err := models.Parse(os.Args[1])

	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	e.File("/", "public/index.html")
	e.GET("/groups", handlers.GetGroups(channels))
	e.GET("/groups/:id/channels", handlers.GetChannels(channels))
	//e.PUT("/tasks", func(c echo.Context) error { return c.JSON(200, "PUT Tasks") })
	//e.DELETE("/tasks/:id", func(c echo.Context) error { return c.JSON(200, "DELETE Task "+c.Param("id")) })
	e.Logger.Fatal(e.Start(":1323"))
}