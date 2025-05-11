package orchestrator_test

import (
	"bytes"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nastts/final-calculator/orchestrator"
)

func TestRegisterHandler(t *testing.T){
	req := httptest.NewRequest("POST", "/api/v1/register",bytes.NewBuffer([]byte(`{ "login":"hello" , "password":"world" }`)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.Exec(`CREATE TABLE users (user_id INTEGER PRIMARY KEY AUTOINCREMENT, login TEXT UNIQUE, password TEXT)`)
	
	handler := func(w http.ResponseWriter, r *http.Request) {
        orchestrator.RegisterHandler(w, r, db)
    }

	handler(w, req)

	expected := "регистрация прошла успешно"


	if w.Code != http.StatusOK {
        t.Errorf("expected status 200 OK, got %d", w.Code)
    }

	if strings.TrimSpace(w.Body.String()) != expected{
		t.Errorf("want %s, get %s",expected, w.Body.String())
	}

	
}


func TestRegisterHandlerBad(t *testing.T){
	req := httptest.NewRequest("POST", "/api/v1/register",bytes.NewBuffer([]byte(`{ "login":"hello"  "password":"world" }`)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.Exec(`CREATE TABLE users (user_id INTEGER PRIMARY KEY AUTOINCREMENT, login TEXT UNIQUE, password TEXT)`)
	
	handler := func(w http.ResponseWriter, r *http.Request) {
        orchestrator.RegisterHandler(w, r, db)
    }

	handler(w, req)

	expected := "you couldn't register"


	if w.Code != http.StatusUnauthorized {
        t.Errorf("expected status 401, got %d", w.Code)
    }

	if strings.TrimSpace(w.Body.String()) != expected{
		t.Errorf("want %s, get %s",expected, w.Body.String())
	}
}



