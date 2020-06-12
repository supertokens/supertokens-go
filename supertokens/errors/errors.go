/*
 * Copyright (c) 2020, VRAI Labs and/or its affiliates. All rights reserved.
 *
 * This software is licensed under the Apache License, Version 2.0 (the
 * "License") as published by the Apache Software Foundation.
 *
 * You may not use this file except in compliance with the License. You may
 * obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations
 * under the License.
 */

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

// UnauthorizedError used for when the user has been logged out
type UnauthorizedError struct {
	Msg string
}

func (err UnauthorizedError) Error() string {
	return err.Msg
}

// IsTokenTheftDetectedError returns true if error is a TokenTheftDetectedError
func IsTokenTheftDetectedError(err error) bool {
	return reflect.TypeOf(err) == reflect.TypeOf(TokenTheftDetectedError{})
}

// IsUnauthorizedError returns true if error is a UnauthorizedError
func IsUnauthorizedError(err error) bool {
	return reflect.TypeOf(err) == reflect.TypeOf(UnauthorizedError{})
}

// IsTryRefreshTokenError returns true if error is a TryRefreshTokenError
func IsTryRefreshTokenError(err error) bool {
	return reflect.TypeOf(err) == reflect.TypeOf(TryRefreshTokenError{})
}
