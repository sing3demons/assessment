package expenses

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *handler) UpdateExpensesHandler(c echo.Context) error {
	var m NewsExpenses

	id := c.Param("id")

	if err := c.Bind(&m); err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	stmt, err := h.DB.Prepare("UPDATE expenses SET title=$2, amount=$3, note=$4, tags=$5 WHERE id=$1")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	_, err = stmt.Exec(id, m.Title, m.Amount, m.Note, m.Tags)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"id":     m.ID,
		"title":  m.Title,
		"amount": m.Amount,
		"note":   m.Note,
		"tags":   m.Tags,
	})
}
