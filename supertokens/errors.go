package supertokens

// GeneralError used for non specific exceptions
type GeneralError struct {
	msg         string
	actualError *error
}

func (err GeneralError) Error() string {
	return err.msg
}

// TryRefreshTokenError used for when the refresh API needs to be called
type TryRefreshTokenError struct {
	msg string
}

func (err TryRefreshTokenError) Error() string {
	return err.msg
}

// TokenTheftDetectedError used for when token theft has happened for a session
type TokenTheftDetectedError struct {
	msg           string
	SessionHandle string
	UserID        string
}

func (err TokenTheftDetectedError) Error() string {
	return err.msg
}

// UnauthorisedError used for when the user has been logged out
type UnauthorisedError struct {
	msg string
}

func (err UnauthorisedError) Error() string {
	return err.msg
}
