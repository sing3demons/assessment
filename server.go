package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sing3demons/assessment/handler/expenses"
)

func initDB() *sql.DB {
	connStr := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	createTb := `
	CREATE TABLE IF NOT EXISTS expenses (
		id SERIAL PRIMARY KEY,
		title TEXT,
		amount FLOAT,
		note TEXT,
		tags TEXT[]
	);
	`
	_, err = db.Exec(createTb)
	if err != nil {
		log.Fatal("can't create table", err)
	}

	return db
}

func main() {
	db := initDB()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	h := expenses.NewApplication(db)
	e.POST("/expenses", h.CreateExpensesHandler)
	e.GET("/expenses/:id", h.GetExpensesHandlerByID)
	e.PUT("/expenses/:id", h.UpdateExpensesHandler)
	e.DELETE("/expenses/:id", h.DeleteExpenseHandlerByID)
	// e.GET("/expenses", h.ListExpensesHandler)

	e.Use(echo.MiddlewareFunc(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			headers := c.Request().Header.Get("Authorization")
			if headers != "admin" {
				return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Unauthorized"})
			}
			return next(c)
		}
	}))
	e.GET("/expenses", h.ListExpensesHandler)

	fmt.Println("start at port:", os.Getenv("PORT"))
	e.Start(":" + os.Getenv("PORT"))
}
