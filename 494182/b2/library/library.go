package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func main() {
	e := echo.New()
	e.GET("/", handler)
	e.Logger.Fatal(e.Start(":8080"))
}

func handler(c echo.Context) error {
    // Simulate an error
    return echo.NewHTTPError(http.StatusBadRequest, "Invalid input data")
}