package services

import (
	"Alarm/internal/web/models"
	"errors"
)

type Account struct {
	db    *models.Database
	cache *models.Cache
	cfg   map[string]interface{}
}

func NewAccount(cfg map[string]interface{}) *Account {
	return &Account{
		db:    cfg["db"].(*models.Database),
		cache: cfg["cache"].(*models.Cache),
		cfg:   cfg,
	}
}

func (svc *Account) CreateUser(user *models.User) error {
	_, err := svc.db.Engine.Insert(user)
	return err
}

func (svc *Account) DeleteUser(user *models.User) error {
	has, err := svc.getUser(user)
	if err != nil {
		return err
	}
	if !has {
		return errors.New("无法找到user数据")
	}
	_, err = svc.db.Engine.ID(user.ID).Unscoped().Delete(user)
	return err
}
func (svc *Account) DeletedUser(user *models.User) error {
	has, err := svc.getUser(user)
	if err != nil {
		return err
	}
	if !has {
		return errors.New("无法找到user数据")
	}
	_, err = svc.db.Engine.ID(user.ID).Delete(user)
	return err
}

func (svc *Account) UpdateUserByID(id int, user *models.User) error {
	_, err := svc.db.Engine.ID(id).Update(user)
	return err
}

func (svc *Account) getUser(user *models.User) (bool, error) {
	if user.ID != 0 {
		has, err := svc.db.Engine.ID(user.ID).Get(user)
		if err != nil {
			return false, err
		}
		if !has {
			return false, nil
		}
		return true, nil
	}
	return svc.db.Engine.Cols("id", "username", "email", "telephone").Get(user)
}
func (svc *Account) GetUserByID(id int) (*models.UserInfo, error) {
	user := new(models.User)
	user.ID = id
	has, err := svc.getUser(user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return user.GetInfo(), nil
}
func (svc *Account) AllUser() ([]models.UserInfo, error) {
	var users []models.User
	err := svc.db.Engine.Find(&users)
	if err != nil {
		return nil, err
	}
	var usersInfo []models.UserInfo
	for _, user := range users {
		usersInfo = append(usersInfo, *user.GetInfo())
	}
	return usersInfo, nil
}

func (svc *Account) IsUserExist(user *models.User) (bool, string, error) {
	has, err := svc.db.Engine.Where("username = ?", user.Username).Exist(&models.User{})
	if err != nil {
		return true, "", err
	}
	if has {
		return true, "该用户名已被使用", nil
	}
	if user.Email != "" {
		has, err = svc.db.Engine.Where("email = ?", user.Email).Exist(&models.User{})
		if err != nil {
			return true, "", err
		}
		if has {
			return true, "该邮箱已被使用", nil
		}
	}
	if user.Phone != "" {
		has, err = svc.db.Engine.Where("phone = ?", user.Phone).Exist(&models.User{})
		if err != nil {
			return true, "", err
		}
		if has {
			return true, "该电话号码已被使用", nil
		}
	}
	return false, "", nil
}
