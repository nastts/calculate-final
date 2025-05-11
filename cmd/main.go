package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/nastts/final-calculator/agent"
	"github.com/nastts/final-calculator/database"
	"github.com/nastts/final-calculator/orchestrator"
)

func main() {
	ctx := context.Background()

	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.PingContext(ctx)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 2; i++{
		go agent.Worker()
	}

	if err = database.CreateTables(ctx, db); err != nil {
		panic(err)
	}
	http.HandleFunc("/api/v1/register", func(w http.ResponseWriter, r *http.Request) {
		orchestrator.RegisterHandler(w, r, db)
	})
	
	
	http.HandleFunc("/api/v1/login", func(w http.ResponseWriter, r *http.Request) {
		orchestrator.LoginHandler(w, r, db)
	})
	http.Handle("/api/v1/calculate", orchestrator.LoginMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		orchestrator.IDHandler(w,r,db)
	})))
	http.Handle("/api/v1/expressions", orchestrator.LoginMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		orchestrator.GetExpressionsHandler(w,r,db)
	})))
	http.Handle("/api/v1/expressions/{id}", orchestrator.LoginMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		orchestrator.GetExpressionByIDHandler(w,r,db)
	})))
	log.Printf("сервер запущен")
	if err := http.ListenAndServe(":8080", nil); err != nil{
		log.Fatal("ошибка при запуске сервера", err)
	}
	
}