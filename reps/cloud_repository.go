package reps

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"mngr/models"
	"strconv"
)

type CloudRepository struct {
	Connection *redis.Client
}

func getTelegramUserKey(id string) string {
	return "cloud:telegram:users:" + id
}

func getTelegramBotKey() string {
	return "cloud:telegram:bot"
}

func getTelegramEnabledKey() string {
	return "cloud:telegram:enabled"
}

func (c *CloudRepository) IsTelegramIntegrationEnabled() bool {
	conn := c.Connection
	ctx := context.Background()
	result, err := conn.Get(ctx, getTelegramEnabledKey()).Result()
	if err != nil {
		log.Println(err.Error())
		_, err = conn.Set(ctx, getTelegramEnabledKey(), "0", 0).Result()
		if err != nil {
			log.Println(err.Error())
		}
		return false
	}

	return result == "1"
}

func (c *CloudRepository) SetTelegramIntegrationEnabled(value bool) (string, error) {
	dbValue := "0"
	if value {
		dbValue = "1"
	}
	return c.Connection.Set(context.Background(), getTelegramEnabledKey(), dbValue, 0).Result()
}

func (c *CloudRepository) GetTelegramUsers() ([]*models.TelegramUser, error) {
	ctx := context.Background()
	ret := make([]*models.TelegramUser, 0)
	conn := c.Connection
	keys, err := conn.Keys(ctx, getTelegramUserKey("*")).Result()
	if err != nil {
		return ret, err
	}
	for _, key := range keys {
		var tbUser models.TelegramUser
		err := conn.HGetAll(ctx, key).Scan(&tbUser)
		if err != nil {
			return nil, err
		}
		ret = append(ret, &tbUser)
	}

	return ret, err
}

func (c *CloudRepository) RemoveTelegramUserById(telegramUserId int) error {
	conn := c.Connection
	_, err := conn.Del(context.Background(), getTelegramUserKey(strconv.Itoa(telegramUserId))).Result()
	if err != nil {
		return err
	}
	return nil
}

func (c *CloudRepository) GetTelegramBot() (*models.TelegramBot, error) {
	conn := c.Connection
	bot := &models.TelegramBot{}
	err := conn.HGetAll(context.Background(), getTelegramBotKey()).Scan(bot)
	return bot, err
}

func (c *CloudRepository) SaveTelegramBot(telegramBot *models.TelegramBot) (int64, error) {
	if telegramBot == nil {
		return 0, nil
	}
	conn := c.Connection
	result, err := conn.HSet(context.Background(), getTelegramBotKey(), Map(telegramBot)).Result()
	return result, err
}

// ********************* GDrive *********************

func getGdriveEnabledKey() string {
	return "cloud:gdrive:enabled"
}

func getGdriveTokenKey() string {
	return "cloud:gdrive:token"
}

func getGdriveCredentialsKey() string {
	return "cloud:gdrive:credentials"
}

func getGdriveUrlKey() string {
	return "cloud:gdrive:url"
}

func getGdriveAuthCodeKey() string {
	return "cloud:gdrive:authcode"
}

func (c *CloudRepository) getValue(key string) (string, error) {
	result, err := c.Connection.Get(context.Background(), key).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			_, err := c.Connection.Set(context.Background(), key, "", 0).Result()
			return "", err
		} else {
			log.Println(err.Error())
			return "", err
		}
	}

	return result, nil
}
func (c *CloudRepository) setValue(key string, value string) (string, error) {
	return c.Connection.Set(context.Background(), key, value, 0).Result()
}

func (c *CloudRepository) IsGdriveIntegrationEnabled() bool {
	conn := c.Connection
	ctx := context.Background()
	result, err := conn.Get(ctx, getGdriveEnabledKey()).Result()
	if err != nil {
		log.Println(err.Error())
		_, err = conn.Set(ctx, getGdriveEnabledKey(), "0", 0).Result()
		if err != nil {
			log.Println(err.Error())
		}
		return false
	}

	return result == "1"
}

func (c *CloudRepository) SaveGdriveIntegrationEnabled(value bool) error {
	val := "0"
	if value {
		val = "1"
	}
	_, err := c.setValue(getGdriveEnabledKey(), val)
	return err
}

func (c *CloudRepository) GetGdriveToken() (string, error) {
	return c.getValue(getGdriveTokenKey())
}

func (c *CloudRepository) SaveGdriveToken(tokenJson string) error {
	_, err := c.setValue(getGdriveTokenKey(), tokenJson)
	return err
}

func (c *CloudRepository) GetGdriveCredentials() (string, error) {
	return c.getValue(getGdriveCredentialsKey())
}

func (c *CloudRepository) SaveGdriveCredentials(credentialsJson string) error {
	_, err := c.setValue(getGdriveCredentialsKey(), credentialsJson)
	return err
}

func (c *CloudRepository) GetGdriveUrl() (string, error) {
	return c.getValue(getGdriveUrlKey())
}

func (c *CloudRepository) SaveGdriveUrl(url string) error {
	_, err := c.setValue(getGdriveUrlKey(), url)
	return err
}

func (c *CloudRepository) GetGdriveAuthCode() (string, error) {
	return c.getValue(getGdriveAuthCodeKey())
}

func (c *CloudRepository) SaveGdriveAuthCode(authCode string) error {
	_, err := c.setValue(getGdriveAuthCodeKey(), authCode)
	return err
}
