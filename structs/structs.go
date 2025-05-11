package structs

import (
	

	
)





type User struct {
	User_ID      int64
	Login    string
	Password string
	OriginPassword string
}

type Expression struct {
	ID         int64 `json:"id,omitempty"`
	Expression string `json:"expression,omitempty"`
	Login     int64 `json:"login,omitempty"`
	Result float64 `json:"result,omitempty"`
	
	
}


type ExpressionsList struct{
	Expressions []Expression `json:"expressions"`
}
type IDDD struct{
	ID int64 `json:"id"`
}


type Request struct{
	Expression string `json:"expression"`
}
type Error struct{
	Error string `json:"error"`
}

type Task struct{
	ID string `json:"id"`
	Arg1 float64 `json:"arg1"`
	Arg2 float64 `json:"arg2"`
	Operation string `json:"operation"`
	OperationTime int `json:"operationTime"`
}

type Result struct{
	ID string  `json:"id"`
	Result float64 `json:"result"`
}


type Times struct {
	TimeAdditionMs int `json:"time_addition_ms"`
	TimeSubtractionMs int `json:"time_subtraction_ms"`
	TimeMultiplicationsMs int `json:"time_multiplications_ms"`
	TimeDivisionsMs int `json:"time_divisions_ms"`
}


type Register struct{
	Login string `json:"login"`
	Password string `json:"password"`
}
type Login struct{
	Login string `json:"login"`
	Password string `json:"password"`
	
}
type Token struct{
	Token string `json:"token"`
}


var(
	TasksQueue  = make([]*Task, 0)
	Expressions = make(map[string][]*Expression)
)

