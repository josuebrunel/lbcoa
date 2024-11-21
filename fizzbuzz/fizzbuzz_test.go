package fizzbuzz

import (
	"context"
	"database/sql"
	"encoding/json"
	"fizzbuzz/pkg/assert"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// MockStorer is a mock implementation of the Storer interface
type MockStorer struct {
	DB   *sql.DB
	mock sqlmock.Sqlmock
}

// NewMockStorer creates a new MockStorer
func NewMockStorer() (*MockStorer, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, err
	}
	return &MockStorer{DB: db, mock: mock}, nil
}

// Exec mocks the Exec method of Storer
func (m *MockStorer) Exec(ctx context.Context, query string, args ...any) (int64, error) {
	result := m.mock.ExpectExec(query)
	if len(args) > 0 {
		result.WithArgs(args)
	}
	result.WillReturnResult(sqlmock.NewResult(1, 1)) // Mock affected rows
	return 1, nil
}

// SelectOne mocks the SelectOne method of Storer
func (m *MockStorer) SelectOne(ctx context.Context, query string, args ...any) (*sql.Row, error) {
	rows := sqlmock.NewRows([]string{"qs", "hits"}).AddRow("int1=3&int2=5&limit=20&str1=fizz&str2=buzz", 18)
	m.mock.ExpectQuery(query).WillReturnRows(rows)

	row := m.DB.QueryRowContext(ctx, query, args...)
	return row, row.Err()
}

// Close closes the mock database
func (m *MockStorer) Close() error {
	if err := m.DB.Close(); err != nil {
		return err
	}
	return nil
}

var testCases = []struct {
	name string
	path string
	code int
	err  string
	data string
}{
	{name: "WithMissingParams", path: "/", code: 400, err: ErrMissingParam.Error(), data: ""},
	{name: "WithMissingStrParam", path: "/?int1=3&int2=5&limit=20&str1=fizz", code: 400, err: ErrMissingParam.Error(), data: ""},
	{name: "WithInvalidIntParam", path: "/?int1=a&int2=5&limit=20&str1=fizz&str2=buzz", code: 400, err: ErrValMustBeInt.Error(), data: ""},
	{name: "WithInvalidLimitParam", path: "/?int1=3&int2=5&limit=-20&str1=fizz&str2=buzz", code: 400, err: ErrValMustBeGTZero.Error(), data: ""},
	{name: "WithValidParam", path: "/?int1=3&int2=5&limit=20&str1=fizz&str2=buzz", code: 200, err: "",
		data: "1,2,fizz,4,buzz,fizz,7,8,fizz,buzz,11,fizz,13,14,fizzbuzz,16,17,fizz,19,buzz"},
}

func TestFizzBuzz(t *testing.T) {
	store, _ := NewMockStorer()
	defer store.Close()
	// test handler
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tc.path, nil)
			w := httptest.NewRecorder()
			t.Logf("Request: GET %v", tc.path)
			Handler(context.Background(), store)(w, req)
			t.Logf("Response: %d %v", w.Code, w.Body.String())
			assert.Eq(t, w.Code, tc.code)
			resp := make(map[string]string)
			_ = json.Unmarshal(w.Body.Bytes(), &resp)
			assert.Eq(t, resp["error"], tc.err)
			assert.Eq(t, resp["data"], tc.data)
		})
	}
	// test stat handler
	store, _ = NewMockStorer()
	defer store.Close()
	t.Run("Stat", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/stat", nil)
		w := httptest.NewRecorder()
		t.Logf("Request: GET %v", "/stat")
		StatHandler(context.Background(), store)(w, req)
		t.Logf("Response: %d %v", w.Code, w.Body.String())
		assert.Eq(t, w.Code, 200)
		resp := make(map[string]any)
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		t.Logf("resp: %v", resp)
		assert.Eq(t, resp["error"].(string), "")
		data := resp["data"].(map[string]any)
		assert.Eq(t, data["hits"].(float64), 18)
	})
}
