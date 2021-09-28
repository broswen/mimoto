package repository

import (
	"errors"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type User struct {
	Email            string `json:"email" gorm:"primaryKey"`
	Password         string `json:"-" gorm:"-"`
	HashedPassword   string `json:"-"`
	Name             string `json:"name"`
	RefreshToken     string `json:"refreshToken"`
	ConfirmationCode string `json:"confirmationCode"`
	Confirmed        bool   `json:"confirmed"`
	ResetCode        string `json:"ResetCode"`
}

type UserRepository interface {
	FindByEmail(email string) (User, error)
	Create(*User) error
	Save(*User) error
}

type MapRepository struct {
	M map[string]User
}

func NewMap() (MapRepository, error) {
	return MapRepository{
		M: make(map[string]User),
	}, nil
}

func (mr MapRepository) FindByEmail(email string) (User, error) {
	user, ok := mr.M[email]
	if !ok {
		return User{}, ErrUserNotFound
	}
	return user, nil
}

func (mr MapRepository) Create(user *User) error {
	_, ok := mr.M[user.Email]
	if ok {
		return ErrUserAlreadyExists
	}
	mr.M[user.Email] = *user
	return nil
}

func (mr MapRepository) Save(user *User) error {
	mr.M[user.Email] = *user
	return nil
}

type PostgresRepository struct {
	DB *gorm.DB
}

func NewPostgres() (PostgresRepository, error) {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	puser := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASS")
	dbname := os.Getenv("POSTGRES_DB")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s  sslmode=disable", host, puser, pass, dbname, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return PostgresRepository{}, err
	}

	db.AutoMigrate(&User{})

	return PostgresRepository{
		DB: db,
	}, nil
}

func (r PostgresRepository) FindByEmail(email string) (User, error) {
	user := User{Email: email}
	tx := r.DB.First(&user)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return User{}, ErrUserNotFound
		}
		return User{}, tx.Error
	}
	return user, nil
}

func (r PostgresRepository) Create(user *User) error {
	tx := r.DB.Create(user)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (r PostgresRepository) Save(user *User) error {
	tx := r.DB.Save(user)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
