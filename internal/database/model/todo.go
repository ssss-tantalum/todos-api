package model

import "github.com/uptrace/bun"

type Todo struct {
	bun.BaseModel

	ID          uint64 `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
