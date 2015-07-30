package game

import (
	"google.golang.org/appengine/user"

	"golang.org/x/net/context"
)

type userData struct {
	Username string
	Games    []gameID
}

func getUser(ctx context.Context) *userData {
	// TODO Get user from datastore
	u := user.Current(ctx)
	if u == nil {
		return nil
	}
	return &userData{"TestUser", []gameID{gameID{1, "Test"}}}
}
