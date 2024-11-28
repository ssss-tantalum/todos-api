package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ssss-tantalum/todos-api/internal/database/model"
	"github.com/ssss-tantalum/todos-api/internal/todos"
)

type TodoHandler struct {
	app *todos.App
}

func NewTodoHandler(app *todos.App) TodoHandler {
	return TodoHandler{
		app: app,
	}
}

func (h TodoHandler) List(c echo.Context) error {
	ctx := c.Request().Context()

	todos := new([]model.Todo)
	err := h.app.DB().NewSelect().
		Model(todos).
		Scan(ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, todos)
}

func (h TodoHandler) Show(c echo.Context) error {
	ctx := c.Request().Context()
	todoID := c.Param("id")

	todo := new(model.Todo)
	err := h.app.DB().NewSelect().
		Model(todo).
		Where("id = ?", todoID).
		Scan(ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, todo)
}

func (h TodoHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()

	var todo *model.Todo

	if err := c.Bind(&todo); err != nil {
		return err
	}

	// TODO: validate todo

	if _, err := h.app.DB().NewInsert().
		Model(todo).
		Exec(ctx); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, todo)
}

func (h TodoHandler) Update(c echo.Context) error {
	return nil
}

func (h TodoHandler) Delete(c echo.Context) error {
	return nil
}
