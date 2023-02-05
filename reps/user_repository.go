package reps

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"log"
	"mngr/models"
	"mngr/utils"
)

type UserRepository struct {
	Connection *redis.Client
}

func getUserKey(id string) string {
	return "users:" + id
}

func (u *UserRepository) addToRedis(user *models.User) (int64, error) {
	return u.Connection.HSet(context.Background(), getUserKey(user.Id), Map(user)).Result()
}

func (u *UserRepository) Login(lu *models.LoginUserViewModel) (*models.User, error) {
	validate := validator.New()
	err := validate.Struct(lu)
	if err != nil {
		return nil, err
	}

	users, err := u.GetUsers()
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		if user.Username == lu.Username && utils.VerifyPassword(user.Password, lu.Password) {
			user.LastLoginAt = utils.DatetimeNow()
			user.VisitCount += 1
			u.addToRedis(user)
			return user, nil
		}
	}

	return nil, nil
}

func (u *UserRepository) LoginByToken(token string) (*models.User, error) {
	users, err := u.GetUsers()
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		if user.Token == token {
			user.LastLoginAt = utils.DatetimeNow()
			user.VisitCount += 1
			u.addToRedis(user)
			return user, nil
		}
	}

	return nil, nil
}

func (u *UserRepository) Register(uv *models.RegisterUserViewModel) (*models.User, error) {
	validate := validator.New()
	err := validate.Struct(uv)
	if err != nil {
		return nil, err
	}
	if uv.Password != uv.RePassword {
		return nil, errors.New("password and re-password does not match")
	}
	user := &models.User{}
	user.Id = utils.NewId()
	user.Username = uv.Username
	user.Password, _ = utils.EncryptPassword(uv.Password)
	user.Email = uv.Email
	user.Token = utils.GenerateSecureToken(4)
	user.LastLoginAt = utils.DatetimeNow()
	user.Ip = uv.Ip
	user.Uag = uv.Uag
	user.Location = uv.Location
	user.DataCenterLocation = uv.DataCenterLocation

	_, err = u.addToRedis(user)
	if err != nil {
		return nil, err
	}

	return user, err
}

func (u *UserRepository) GetUser(userId string) (*models.User, error) {
	conn := u.Connection
	var us models.User
	err := conn.HGetAll(context.Background(), getUserKey(userId)).Scan(&us)
	if err != nil {
		log.Println("Error getting stream from redis: ", err)
		return nil, err
	}
	return &us, nil
}

func (u *UserRepository) GetUsers() ([]*models.User, error) {
	conn := u.Connection
	keys, err := conn.Keys(context.Background(), getUserKey("*")).Result()
	list := make([]*models.User, 0)
	if err != nil {
		if err.Error() == "redis: nil" {
			conn.Set(context.Background(), getUserKey(""), list, 0)
			return list, nil
		} else {
			log.Println("Error getting all stream from redis: ", err)
			return nil, err
		}
	}

	for _, key := range keys {
		var us models.User
		err := conn.HGetAll(context.Background(), key).Scan(&us)
		if err != nil {
			log.Println("Error getting stream from redis: ", err)
			return nil, err
		}
		list = append(list, &us)
	}
	return list, nil
}

func (u *UserRepository) RemoveById(userId string) (int64, error) {
	result, err := u.Connection.Del(context.Background(), getUserKey(userId)).Result()
	if err != nil {
		log.Println("Error while deleting source: ", err)
	}
	return result, err
}
