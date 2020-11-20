package models

import (
	"brackets/internal/db"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const BcryptCost = 12

type User struct {
	gorm.Model

	Name     string `gorm:"unique"`
	Email    string
	Password []byte
}

func (r *User) JSON() map[string]interface{} {
	return map[string]interface{}{
		"id":    r.ID,
		"name":  r.Name,
		"email": r.Email,
	}
}

func NewUser(name, email, password string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), BcryptCost)
	if err != nil {
		return nil, err
	}

	u := User{
		Name:     name,
		Email:    email,
		Password: hashedPassword,
	}

	result := db.DB.Create(&u)
	if result.Error != nil {
		return nil, result.Error
	}
	return &u, nil
}

func CheckPassword(name, password string) (bool, error) {
	u := User{}
	result := db.DB.Where("name = ?", name).First(&u)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, result.Error
	}

	err := bcrypt.CompareHashAndPassword(u.Password, []byte(password))
	if err != nil {
		return false, nil
	}
	return true, nil
}
