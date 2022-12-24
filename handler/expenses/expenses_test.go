package expenses

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestListExpensesHandler(t *testing.T) {
	// Arrange
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/expenses", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	newsMockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).AddRow("1", "strawberry smoothie", "79", "night market promotion discount 10 bath", pq.Array([]string{"food", "beverage"}))

	db, mock, err := sqlmock.New()
	mock.ExpectQuery("SELECT (.+) FROM expenses").WillReturnRows(newsMockRows)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	h := handler{db}
	c := e.NewContext(req, rec)

	expected := "[{\"id\":1,\"title\":\"strawberry smoothie\",\"amount\":79,\"note\":\"night market promotion discount 10 bath\",\"tags\":[\"food\",\"beverage\"]}]"

	// Act
	err = h.ListExpensesHandler(c)

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}
}

func TestCreateExpensesHandler(t *testing.T) {
	// Arrange
	data := `{
		"title": "strawberry smoothie",
		"amount": 79,
		"note": "night market promotion discount 10 bath", 
		"tags": ["food", "beverage"]
	}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(data))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	newsMockRows := sqlmock.NewRows([]string{"id"}).AddRow("1")
	db, mock, err := sqlmock.New()
	mock.ExpectQuery(
		"INSERT INTO expenses \\(title, amount, note, tags\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\) RETURNING id").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnRows(newsMockRows)

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	h := handler{db}
	ctx := e.NewContext(req, rec)

	// Act
	h.CreateExpensesHandler(ctx)

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}
