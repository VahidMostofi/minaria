package repositories

import (
	"context"
	"crypto/md5"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/vahidmostofi/minaria/domain"
)

type inMemoryUserRepository struct {
	cache []*domain.User
}

func newInMemoryUserRepository(ima *InMemoryArgs) *inMemoryUserRepository {
	im := &inMemoryUserRepository{}
	if ima != nil {
		im.cache = ima.Data
	} else {
		im.cache = []*domain.User{{
			ID:       "54215f2a-b752-11eb-8529-0242ac130003",
			Username: "jack",
			Email:    "jack@gmail.com",
			Password: "8bb0cf6eb9b17d0f7d22b456f121257dc1254e1f01665370476383ea776df414",
		}, {
			ID:       "5a823a9c-b752-11eb-8529-0242ac130003",
			Username: "john",
			Email:    "john@gmail.com",
			Password: "8bb0cf6eb9b17d0f7d22b456f121257dc1254e1f01665370476383ea776df414",
		}, {
			ID:       "601427c2-b752-11eb-8529-0242ac130003",
			Username: "jill",
			Email:    "jill@gmail.com",
			Password: "8bb0cf6eb9b17d0f7d22b456f121257dc1254e1f01665370476383ea776df414",
		}}
	}
	return im
}

func (im *inMemoryUserRepository) GetByID(ctx context.Context, ID string) (*domain.User, error) {
	for _, u := range im.cache {
		if u.ID == ID {
			return u, nil
		}
	}
	return nil, ErrNoUserFound
}

func (im *inMemoryUserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	for _, u := range im.cache {
		if u.Username == username {
			return u, nil
		}
	}
	return nil, ErrNoUserFound
}
func (im *inMemoryUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	for _, u := range im.cache {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, ErrNoUserFound
}
func (im *inMemoryUserRepository) Store(ctx context.Context, u *domain.User) (*domain.User, error) {
	if u, errU := im.GetByUsername(ctx, u.Username); u != nil && errU != ErrNoUserFound {
		return nil, ErrUsernameNotUnique
	}

	if u, errU := im.GetByEmail(ctx, u.Email); u != nil && errU != ErrNoUserFound {
		return nil, ErrEmailNotUnique
	}

	if _, err := uuid.Parse(u.ID); err == nil {
		return nil, fmt.Errorf("can't store the the object already has an ID.")
	}

	u.ID = uuid.New().String()

	im.cache = append(im.cache, u)

	return u, nil
}

func (im *inMemoryUserRepository) Update(ctx context.Context, u *domain.User) (*domain.User, error) {
	if _, errU := im.GetByUsername(ctx, u.Username); errU != nil {
		return nil, ErrUsernameNotUnique
	}

	if _, errU := im.GetByEmail(ctx, u.Email); errU != nil {
		return nil, ErrEmailNotUnique
	}

	for i, candidateUser := range im.cache {
		if candidateUser.ID == u.ID {
			im.cache[i].Username = u.Username
			im.cache[i].Email = u.Email
			if len(im.cache[i].Password) > 0 {
				im.cache[i].Password = string(md5.New().Sum([]byte(im.cache[i].Password)))
			}
			im.cache[i].UpdatedAt = time.Now()
			return im.cache[i], nil
		}
	}

	return nil, ErrNoUserFound
}
