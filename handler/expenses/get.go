package expenses

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *handler) GetExpensesHandlerByID(c echo.Context) error {
	var m NewsExpenses
	id := c.Param("id")

	row := h.DB.QueryRow("SELECT id, title, amount, note, tags FROM expenses where id=$1", id)

	err := row.Scan(&m.ID, &m.Title, &m.Amount, &m.Note, &m.Tags)

	switch err {
	case sql.ErrNoRows:
		return c.JSON(http.StatusNotFound, Err{Message: "expenses not found"})
	case nil:
		return c.JSON(http.StatusOK, m)
	default:
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't scan expenses:" + err.Error()})
	}
}

func (h *handler) ListExpensesHandler(c echo.Context) error {

	rows, err := h.DB.Query("SELECT id, title, amount, note, tags FROM expenses")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't query all expenses:" + err.Error()})
	}
	defer rows.Close()

	var expenses = []NewsExpenses{}
	var m = NewsExpenses{}

	for rows.Next() {
		err := rows.Scan(&m.ID, &m.Title, &m.Amount, &m.Note, &m.Tags)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Err{Message: "can't scan expenses:" + err.Error()})

		}
		expenses = append(expenses, m)
	}
	return c.JSON(http.StatusOK, expenses)
}
