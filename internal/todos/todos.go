package todos

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/ssss-tantalum/todos-api/internal/config"
	"github.com/uptrace/bun"
)

type appCtxKey struct{}

func ContextWithApp(ctx context.Context, app *App) context.Context {
	ctx = context.WithValue(ctx, appCtxKey{}, app)

	return ctx
}

type App struct {
	ctx context.Context
	cfg *config.Config

	router    *echo.Echo
	apiRouter *echo.Group

	db *bun.DB
}

func New(ctx context.Context, cfg *config.Config, db *bun.DB) *App {
	app := &App{
		cfg: cfg,
		db:  db,
	}
	app.ctx = ContextWithApp(ctx, app)
	app.initRouter()

	return app
}

func Start(ctx context.Context, cfg *config.Config, db *bun.DB) (context.Context, *App) {
	app := New(ctx, cfg, db)

	return app.Context(), app
}

func (app *App) Context() context.Context {
	return app.ctx
}

func (app *App) Config() *config.Config {
	return app.cfg
}

func (app *App) Router() *echo.Echo {
	return app.router
}

func (app *App) APIRouter() *echo.Group {
	return app.apiRouter
}

func (app *App) DB() *bun.DB {
	return app.db
}

func (app *App) WaitExitSignal() os.Signal {
	ch := make(chan os.Signal, 3)
	signal.Notify(
		ch,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)
	return <-ch
}
