package main

import (
	"github.com/VictorH97/devfullcycle/goexpert/Lab-Temp-CEP/internal/infra/web"
	"github.com/VictorH97/devfullcycle/goexpert/Lab-Temp-CEP/internal/infra/web/webserver"
)

func main() {
	findWeatherHandler := web.NewFindWeatherHandler("2083c7dd8e734e46971234222242102")

	webserver := webserver.NewWebServer("8080")
	webserver.AddHandler("/", findWeatherHandler.FindWeather)

	webserver.Start()
}
