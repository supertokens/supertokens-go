package errors

import "reflect"

// GeneralError used for non specific exceptions
type GeneralError struct {
	Msg         string
	ActualError error
}

func (err GeneralError) Error() string {
	return err.Msg
}

// TryRefreshTokenError used for when the refresh API needs to be called
type TryRefreshTokenError struct {
	Msg string
}

func (err TryRefreshTokenError) Error() string {
	return err.Msg
}

// TokenTheftDetectedError used for when token theft has happened for a session
type TokenTheftDetectedError struct {
	Msg           string
	SessionHandle string
	UserID        string
}

func (err TokenTheftDetectedError) Error() string {
	return err.Msg
}

// UnauthorisedError used for when the user has been logged out
type UnauthorisedError struct {
	Msg string
}

func (err UnauthorisedError) Error() string {
	return err.Msg
}

// IsTokenTheftDetectedError returns true if error is a TokenTheftDetectedError
func IsTokenTheftDetectedError(err error) bool {
	return reflect.TypeOf(err) == reflect.TypeOf(TokenTheftDetectedError{})
}

// IsUnauthorisedError returns true if error is a UnauthorisedError
func IsUnauthorisedError(err error) bool {
	return reflect.TypeOf(err) == reflect.TypeOf(UnauthorisedError{})
}

// IsTryRefreshTokenError returns true if error is a TryRefreshTokenError
func IsTryRefreshTokenError(err error) bool {
	return reflect.TypeOf(err) == reflect.TypeOf(TryRefreshTokenError{})
}
