package orchestrator

import (
	"database/sql"
	"encoding/json"
	"log"
	"strconv"

	"net/http"
	"strings"

	"github.com/nastts/final-calculator/database"
	"github.com/nastts/final-calculator/structs"
)




func RegisterHandler(w http.ResponseWriter, r *http.Request, db *sql.DB){
	var reg structs.Register
	ctx := r.Context()
	defer r.Body.Close()
	
	err := json.NewDecoder(r.Body).Decode(&reg)
	if err != nil{
		http.Error(w, "you couldn't register", http.StatusUnauthorized)
		return 
	}
	if _, err := database.SelectUser(ctx, db, reg.Login); err == nil {
		http.Error(w, "User already exists", http.StatusNotAcceptable)
		return
	} else if err != sql.ErrNoRows {
		http.Error(w, "err bd", http.StatusBadRequest)
		return
	}

	hashPass, err := database.Generate(reg.Password)
	if err != nil{
		http.Error(w, "coudn't hash", http.StatusBadRequest)
		return
	}
	user := &structs.User{
		Login: reg.Login,
		Password: hashPass,
		OriginPassword: reg.Password,
	}

	_, err = database.InsertUser(ctx, db, user)
	if err != nil{
		http.Error(w, "Failed to save user", http.StatusBadRequest)
		return
	}


	w.WriteHeader(http.StatusOK)
	w.Write([]byte("регистрация прошла успешно"))
}


func LoginMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		authHeader := r.Header.Get("Authorization")
		
        if authHeader == "" {
            http.Error(w, "токена нет", http.StatusUnauthorized)
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        token, err := database.ParseToken(tokenString)
        if err != nil || !token.Valid {
            http.Error(w, "неверный токен", http.StatusUnauthorized)
            return
        }
		next.ServeHTTP(w, r)
	})
}




func LoginHandler(w http.ResponseWriter, r *http.Request, db *sql.DB){
	var l structs.Login
	ctx := r.Context()
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&l)
	if err != nil{
		http.Error(w, "you couldn't login", http.StatusUnauthorized)
		return 
	}
	token, err := database.LoginUser(ctx, db, l.Login,l.Password)
	if err != nil{
		http.Error(w, "не получается зарегаться", http.StatusUnauthorized)
		return 
	}
	answer := &structs.Token{
		Token: token,
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(answer)
}


func IDHandler(w http.ResponseWriter, r *http.Request, db *sql.DB){
	var ex structs.Expression
	ctx := r.Context()
	defer r.Body.Close()
	tokenString := r.Header.Get("Authorization")
		
        if tokenString == "" {
            http.Error(w, "токена нет", http.StatusUnauthorized)
            return
        }
	token := strings.TrimPrefix(tokenString, "Bearer ")
	err := json.NewDecoder(r.Body).Decode(&ex)
	if err != nil{
		http.Error(w, "проблемы с выражением", http.StatusUnauthorized)
		return 
	}
	login, err := database.ParseTokenForLogin(token)
	if err != nil{
		http.Error(w, "логин не парситься", http.StatusUnauthorized)
            return
	}
	id, err := database.InsertExpression(ctx,db, &ex, login)
	if err != nil{
		http.Error(w, "ну получается получить id", http.StatusUnauthorized)
        return
	}
	
	
	

	
	answer := &structs.IDDD{
		ID: id,
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(answer)
}
	

func GetExpressionsHandler(w http.ResponseWriter, r *http.Request, db *sql.DB){
	
	ctx := r.Context()
	defer r.Body.Close()
	tokenString := r.Header.Get("Authorization")
		
        if tokenString == "" {
            http.Error(w, "токена нет", http.StatusUnauthorized)
            return
        }
	token := strings.TrimPrefix(tokenString, "Bearer ")
	login, err := database.ParseTokenForLogin(token)
	if err != nil{
		http.Error(w, "логин не парситься", http.StatusUnauthorized)
            return
	}

	expressions, err := database.SelectExpressions(ctx, db, login)
	if err != nil{
		log.Println(err)
		http.Error(w, "не получилось получить выражения", http.StatusUnauthorized)
        return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expressions)
}


func GetExpressionByIDHandler(w http.ResponseWriter, r *http.Request, db *sql.DB){
	ctx := r.Context()
	defer r.Body.Close()
	tokenString := r.Header.Get("Authorization")
		
        if tokenString == "" {
            http.Error(w, "токена нет", http.StatusUnauthorized)
            return
        }
	token := strings.TrimPrefix(tokenString, "Bearer ")
	login, err := database.ParseTokenForLogin(token)
	if err != nil{
		http.Error(w, "логин не парситься", http.StatusUnauthorized)
            return
	}
	id := r.PathValue("id")
	idParsed, err := strconv.ParseInt(id, 10, 64)
	if err != nil{
		http.Error(w, "не получилось поменять id на int64", http.StatusUnauthorized)
        return
	}
	result, err := database.SelectExpression(ctx, db, login, idParsed)
	if err != nil{
		http.Error(w, "не получилось получить одно выражение", http.StatusUnauthorized)
        return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}



func GetTaskHandler(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()
	if len(structs.TasksQueue) == 0{
		http.Error(w, "tasks is not found", http.StatusNotFound)
		return
	}
	task := structs.TasksQueue[0]
	structs.TasksQueue = structs.TasksQueue[1:]
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]*structs.Task{"task":task})
}



