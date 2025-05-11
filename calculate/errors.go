package calculate

import "errors"

var (
	ErrExpressionIsNotValid = errors.New("expression is not valid")//422
	ErrInternalServerError  = errors.New("internal server error")//500
)