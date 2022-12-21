package expenses

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *handler) GetExpensesHandlerByID(c echo.Context) error {
	var m NewsExpenses
	id := c.Param("id")
	stmt, err := h.DB.Prepare("SELECT id, title, amount, note, tags FROM expenses where id=$1")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't prepare query expenses statement:" + err.Error()})
	}
	row := stmt.QueryRow(id)

	err = row.Scan(&m.ID, &m.Title, &m.Amount, &m.Note, &m.Tags)

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
	stmt, err := h.DB.Prepare("SELECT id, title, amount, note, tags FROM expenses")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't prepare query all expenses statement:" + err.Error()})
	}
	rows, err := stmt.Query()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't query all expenses:" + err.Error()})
	}

	expenses := []NewsExpenses{}

	for rows.Next() {
		m := NewsExpenses{}
		err = rows.Scan(&m.ID, &m.Title, &m.Amount, &m.Note, &m.Tags)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Err{Message: "can't scan expenses:" + err.Error()})
		}
		expenses = append(expenses, m)
	}

	return c.JSON(http.StatusOK, expenses)
}
