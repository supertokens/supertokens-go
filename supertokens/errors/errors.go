package errors

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
