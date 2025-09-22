package repository

import "errors"

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrPlaceNotFound        = errors.New("place not found")
	ErrPostNotFound         = errors.New("post not found")
	ErrCommentNotFound      = errors.New("comment not found")
	ErrRatingNotFound       = errors.New("rating not found")
	ErrComparisonNotFound   = errors.New("comparison not found")
	ErrOAuthAccountNotFound = errors.New("OAuth account not found")
	ErrOAuthStateNotFound   = errors.New("OAuth state not found")
	ErrSessionNotFound      = errors.New("session not found")
)
