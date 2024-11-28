package todos

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (app *App) initRouter() {
	app.router = echo.New()
	app.router.Use(middleware.Logger())

	app.apiRouter = app.router.Group("/api")
	app.apiRouter.Use(errorHandler)
}

func errorHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err == nil {
			return nil
		}

		httpErr := from(err)

		return echo.NewHTTPError(httpErr.Code, httpErr.Message)
	}
}

func from(err error) *echo.HTTPError {
	switch err {
	case io.EOF:
		return echo.NewHTTPError(http.StatusBadRequest, "EOF")
	case sql.ErrNoRows:
		return echo.NewHTTPError(http.StatusNotFound, "not found")
	}

	switch err := err.(type) {
	case *echo.HTTPError:
		return err
	case *json.SyntaxError:
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
}
