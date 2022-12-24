//go:build integration
// +build integration

package expenses

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

const serverPort = 2565
const connStr = "postgresql://sing:12345678@db/goapi?sslmode=disable"

type Response struct {
	*http.Response
	err error
}

func request(method, url string, body io.Reader) *Response {
	req, _ := http.NewRequest(method, url, body)

	req.Header.Add("Authorization", "November 10, 2009")
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	return &Response{res, err}
}

func uri(paths ...string) string {
	host := "http://localhost:2565"
	if paths == nil {
		return host
	}
	url := append([]string{host}, paths...)
	return strings.Join(url, "/")
}

func seedData(t *testing.T) NewsExpenses {
	db := initDB()
	var m NewsExpenses = NewsExpenses{
		Title: "strawberry smoothie", Amount: 79,
		Note: "night market promotion discount 10 bath",
		Tags: []string{"food", "beverage"},
	}

	row := db.QueryRow("INSERT INTO expenses (title, amount, note, tags) VALUES ($1, $2, $3, $4) RETURNING id", m.Title, m.Amount, m.Note, m.Tags)

	if err := row.Scan(&m.ID); err != nil {
		t.Fatal("can't create expenses", err)

	}
	return m
}

func (r *Response) Decode(v interface{}) error {
	if r.err != nil {
		return r.err
	}
	return json.NewDecoder(r.Body).Decode(v)
}

func initDB() *sql.DB {
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

func TestListExpensesHandler(t *testing.T) {
	// Setup server
	eh := echo.New()
	go func(e *echo.Echo) {
		db := initDB()
		h := NewApplication(db)
		e.GET("/expenses", h.ListExpensesHandler)
		e.Start(fmt.Sprintf(":%d", serverPort))
	}(eh)
	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", serverPort), 30*time.Second)
		if err != nil {
			log.Println(err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}
	seedData(t)
	var expenses []NewsExpenses

	res := request(http.MethodGet, uri("expenses"), nil)
	err := res.Decode(&expenses)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Greater(t, len(expenses), 0)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = eh.Shutdown(ctx)
	assert.NoError(t, err)

}

func TestCreateExpensesHandler(t *testing.T) {
	// Setup server
	eh := echo.New()
	go func(e *echo.Echo) {
		db := initDB()

		h := NewApplication(db)

		e.POST("/expenses", h.CreateExpensesHandler)
		e.Start(fmt.Sprintf(":%d", serverPort))
	}(eh)
	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", serverPort), 30*time.Second)
		if err != nil {
			log.Println(err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}

	var c NewsExpenses

	body := bytes.NewBufferString(`{
		"title": "strawberry smoothie",
		"amount": 79,
		"note": "night market promotion discount 10 bath", 
		"tags": ["food", "beverage"]
	}`)

	res := request(http.MethodPost, uri("expenses"), body)
	err := res.Decode(&c)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.NotEqual(t, 0, c.ID)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = eh.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestGetExpensesHandlerByID(t *testing.T) {
	// Setup server
	eh := echo.New()
	go func(e *echo.Echo) {
		db := initDB()

		h := NewApplication(db)

		e.GET("/expenses/:id", h.GetExpensesHandlerByID)
		e.Start(fmt.Sprintf(":%d", serverPort))
	}(eh)
	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", serverPort), 30*time.Second)
		if err != nil {
			log.Println(err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}

	e := seedData(t)

	var c NewsExpenses

	res := request(http.MethodGet, uri("expenses", strconv.Itoa(e.ID)), nil)
	err := res.Decode(&c)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, e.ID, c.ID)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = eh.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestUpdateExpensesHandler(t *testing.T) {
	// Setup server
	eh := echo.New()
	go func(e *echo.Echo) {
		db := initDB()

		h := NewApplication(db)

		e.PUT("/expenses/:id", h.UpdateExpensesHandler)
		e.Start(fmt.Sprintf(":%d", serverPort))
	}(eh)
	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", serverPort), 30*time.Second)
		if err != nil {
			log.Println(err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}

	e := seedData(t)

	update := bytes.NewBufferString(`{
    	"title": "apple smoothie",
    	"amount": 89,
    	"note": "no discount",
    	"tags": ["beverage"]
	}`)

	var a NewsExpenses

	res := request(http.MethodPut, uri("expenses", strconv.Itoa(e.ID)), update)
	err := res.Decode(&a)

	expected := NewsExpenses{
		Title:  "apple smoothie",
		Amount: 89,
		Note:   "no discount",
		Tags:   []string{"beverage"},
	}

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusOK, res.StatusCode)

	assert.Equal(t, expected.Title, a.Title)
	assert.Equal(t, expected.Amount, a.Amount)
	assert.Equal(t, expected.Note, a.Note)
	assert.Equal(t, expected.Tags, a.Tags)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = eh.Shutdown(ctx)
	assert.NoError(t, err)
}
