package game

import (
	"errors"

	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/user"

	"golang.org/x/net/context"
)

var (
	errUserExists    = errors.New("Email already has a profile")
	errUsernameTaken = errors.New("Username already taken")
)

type userData struct {
	Username string
	Games    []gameID
}

func addGameToProfile(ctx context.Context, g *gameID) error {
	u := user.Current(ctx)
	k := datastore.NewKey(ctx, "User", u.Email, 0, nil)
	var uData userData
	err := datastore.Get(ctx, k, &uData)
	if err != nil {
		return err
	}
	uData.Games = append(uData.Games, *g)
	_, err = datastore.Put(ctx, k, &uData)
	return err
}

func getUser(ctx context.Context) *userData {
	u := user.Current(ctx)
	if u == nil {
		return nil
	}
	k := datastore.NewKey(ctx, "User", u.Email, 0, nil)
	var ud userData
	err := datastore.Get(ctx, k, &ud)
	if err != nil {
		return nil
	}
	return &ud
}

func getUserByUsername(ctx context.Context, username string) (*userData, error) {
	q := datastore.NewQuery("User").Limit(1).Filter("Username =", username)
	var users []userData
	_, err := q.GetAll(ctx, &users)
	if err != nil {
		return nil, err
	}
	if len(users) < 1 {
		return nil, datastore.ErrNoSuchEntity
	}
	return &users[0], nil
}

func getUserByEmail(ctx context.Context, email string) (*userData, error) {
	k := datastore.NewKey(ctx, "User", email, 0, nil)
	var u userData
	err := datastore.Get(ctx, k, &u)
	return &u, err
}

func createUser(ctx context.Context, username string, email string) error {
	_, err := getUserByEmail(ctx, email)
	if err == nil {
		return errUserExists
	}
	// TODO Check if valid username
	_, err = getUserByUsername(ctx, username)
	if err != datastore.ErrNoSuchEntity {
		return errUsernameTaken
	}
	u := &userData{
		Username: username,
		Games:    []gameID{},
	}
	k := datastore.NewKey(ctx, "User", email, 0, nil)
	_, err = datastore.Put(ctx, k, u)
	return err
}
