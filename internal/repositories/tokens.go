package repositories

import (
	"fmt"
	"github.com/WildEgor/gAuth/internal/db"
	"github.com/WildEgor/gAuth/internal/models"
	"github.com/pkg/errors"
	"time"
)

type TokensRepository struct {
	db *db.RedisConnection
}

const rtPrefix = "refresh_token:"
const atPrefix = "access_token:"

func NewTokensRepository(
	db *db.RedisConnection,
) *TokensRepository {
	return &TokensRepository{
		db,
	}
}

func (tr *TokensRepository) SetRT(token *models.TokenDetails) error {
	err := tr.db.Client.Set(fmt.Sprintf("%s%s", rtPrefix, token.TokenUuid), token.UserID, time.Unix(token.ExpiresIn, 0).Sub(time.Now())).Err()
	if err != nil {
		return errors.Wrap(err, "set refresh token")
	}

	return nil
}

func (tr *TokensRepository) GetRT(tokenUuid string) (string, error) {
	userId, err := tr.db.Client.Get(fmt.Sprintf("%s%s", rtPrefix, tokenUuid)).Result()
	if err != nil {
		return "", errors.Wrap(err, "get refresh token")
	}

	return userId, nil
}

func (tr *TokensRepository) SetAT(token *models.TokenDetails) error {
	err := tr.db.Client.Set(fmt.Sprintf("%s%s", atPrefix, token.TokenUuid), token.UserID, time.Unix(token.ExpiresIn, 0).Sub(time.Now())).Err()
	if err != nil {
		return errors.Wrap(err, "set access token")
	}

	return nil
}

func (tr *TokensRepository) GetAT(tokenUuid string) (string, error) {
	userId, err := tr.db.Client.Get(fmt.Sprintf("%s%s", atPrefix, tokenUuid)).Result()
	if err != nil {
		return "", errors.Wrap(err, "get access token")
	}

	return userId, nil
}

func (tr *TokensRepository) DeleteTokens(tokens ...string) error {
	if len(tokens) != 2 {
		return errors.New("invalid tokens count")
	}

	err := tr.db.Client.Del(fmt.Sprintf("%s%s", atPrefix, tokens[0]), fmt.Sprintf("%s%s", rtPrefix, tokens[1])).Err()
	if err != nil {
		return errors.Wrap(err, "delete refresh token")
	}

	return nil
}
