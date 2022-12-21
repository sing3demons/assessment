package expenses

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *handler) DeleteExpenseHandlerByID(c echo.Context) error {
	id := c.Param("id")

	stmt, err := h.DB.Prepare("DELETE FROM expenses where id=$1")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})

	}
	_, err = stmt.Exec(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})

	}
	return c.JSON(http.StatusNoContent, "Delete")
}
