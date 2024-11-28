package routes

import (
	"github.com/ssss-tantalum/todos-api/api/handler"
	"github.com/ssss-tantalum/todos-api/internal/todos"
)

func InitRoutes(app *todos.App) {
	// r := app.Router()
	g := app.APIRouter()

	todoHandler := handler.NewTodoHandler(app)
	g.POST("/todos", todoHandler.Create)
	g.GET("/todos", todoHandler.List)

	g.GET("/todo/:id", todoHandler.Show)
	g.PUT("/todo/:id", todoHandler.Update)
	g.DELETE("/todo/:id", todoHandler.Delete)
}
