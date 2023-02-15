package common

// import (
// 	"errors"
// 	"fmt"
// )

// type RequestError struct {
// 	Code int

// 	Err error
// }

// func (r *RequestError) Error() string {
// 	return fmt.Sprintf("%d: %v", r.Code, r.Err)
// }

// func EmailTakenError() error {
// 	return &RequestError{
// 		Code: 600,
// 		Err:  errors.New("email already taken"),
// 	}
// }
