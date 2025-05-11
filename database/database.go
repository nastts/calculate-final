package database

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nastts/final-calculator/calculate"
	"github.com/nastts/final-calculator/structs"

	"golang.org/x/crypto/bcrypt"
)
const HmacSampleSecret = "token for expressions"


func CreateTables(ctx context.Context, db *sql.DB) error {
	const (
		usersTable = `
	CREATE TABLE IF NOT EXISTS users(
		user_id INTEGER PRIMARY KEY AUTOINCREMENT, 
		login TEXT UNIQUE,
		password TEXT

	);`

		expressionsTable = `
	CREATE TABLE IF NOT EXISTS expressions(
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		login TEXT NOT NULL,
		expression TEXT NOT NULL,
		result TEXT NOT NULL,
	
		FOREIGN KEY (login)  REFERENCES users (login)
	);`

		TasksTable = `
	CREATE TABLE IF NOT EXISTS tasks(
		id INTEGER,
		first_argument INTEGER,
		second_argument INTEGER,
		operation TEXT,

		FOREIGN KEY (id)  REFERENCES expressions (id)
	);`
	)

	if _, err := db.ExecContext(ctx, usersTable); err != nil {
		return err
	}

	if _, err := db.ExecContext(ctx, expressionsTable); err != nil {
		return err
	}

	if _, err := db.ExecContext(ctx, TasksTable); err != nil {
		return err
	}

	return nil
}


func InsertUser(ctx context.Context, db *sql.DB, user *structs.User) (int64, error) {
	var q = `
	INSERT INTO users (login, password) values ($1, $2)
	`
	result, err := db.ExecContext(ctx, q, user.Login, user.Password)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func SelectUser(ctx context.Context, db *sql.DB, login string) (structs.User, error) {
	var (
		user structs.User
		err  error
	)

	var q = "SELECT user_id, login, password FROM users WHERE login=$1"
	err = db.QueryRowContext(ctx, q, login).Scan(&user.User_ID, &user.Login, &user.Password)
	return user, err
}

func Generate(s string) (string, error) {
	saltedBytes := []byte(s)
	hashedBytes, err := bcrypt.GenerateFromPassword(saltedBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	hash := string(hashedBytes[:])
	return hash, nil
}

func Compare(hash string, s string) error {
	incoming := []byte(s)
	existing := []byte(hash)
	return bcrypt.CompareHashAndPassword(existing, incoming)
}


func LoginUser(ctx context.Context, db *sql.DB, login, password string)(string, error){
	var hash string
	
	err := db.QueryRow("SELECT password FROM users WHERE login = ?", login).Scan(&hash)
	if err == sql.ErrNoRows{
		return "", errors.New("пользователь не найден")
	}else if err != nil{
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil{
		return "", errors.New("неверный пароль")
	}
	
	
	
	
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": login,
		"nbf":  now.Unix(),
		"exp":  now.Add(24*time.Hour).Unix(),
		"iat":  now.Unix(),
	})

	mainToken, err := token.SignedString([]byte(HmacSampleSecret))
	if err != nil {
		return "", errors.New("не удалось создать токен")
	}

	return mainToken, nil
	
}

func InsertExpression(ctx context.Context, db *sql.DB, expression *structs.Expression, login string) (int64, error) {
	var q = `
	INSERT INTO expressions (expression, login, result) VALUES (?, ?, ?) RETURNING id
	`
	var id int64
	result, err := calculate.Calc(expression.Expression)
	if err != nil{
		return 0, err
	}
	err = db.QueryRowContext(ctx, q, expression.Expression, login, result).Scan(&id)
	if err != nil {
		
		return 0, err
	}

	return id, nil
}



func SelectExpressions(ctx context.Context, db *sql.DB, login string) ([]structs.Expression, error) {
	var expressions []structs.Expression
	

	rows, err := db.Query("SELECT id, expression, result FROM expressions WHERE login = ?", login)
	if err != nil {
		log.Println("1:",err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		e := structs.Expression{}
		err := rows.Scan(&e.ID, &e.Expression, &e.Result)
		if err != nil {
			log.Println("2:",err)
			return nil, err
		}
		expressions = append(expressions, e)
	}

	return expressions, nil
}

func SelectExpression(ctx context.Context, db *sql.DB, login string, id int64) (structs.Expression, error) {
	var (
		ex structs.Expression
		err  error
	)

	var q = "SELECT id, expression, result FROM expressions WHERE id=$1"
	err = db.QueryRowContext(ctx, q, id).Scan(&ex.ID, &ex.Expression, &ex.Result)
	return ex, err
}


func InsertTask(ctx context.Context, db *sql.DB, arg1,arg2 float64, operation string)(int64, error){
	
	var q = `
	INSERT INTO expressions (first_argument, second_argument, operation) VALUES (?, ?, ?)
	`
	var id int64
	err := db.QueryRowContext(ctx, q, arg1, arg2, operation).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}