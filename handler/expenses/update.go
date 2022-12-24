package expenses

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *handler) UpdateExpensesHandler(c echo.Context) error {
	var m NewsExpenses

	id, _ := strconv.Atoi(c.Param("id"))

	if err := c.Bind(&m); err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}
	m.ID = id
	stmt, err := h.DB.Prepare("UPDATE expenses SET title=$2, amount=$3, note=$4, tags=$5 WHERE id=$1")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	_, err = stmt.Exec(m.ID, m.Title, m.Amount, m.Note, m.Tags)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, m)
}
