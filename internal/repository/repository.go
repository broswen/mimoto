package repository

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

type Repository struct {
	DB *gorm.DB
}

func New() (Repository, error) {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	puser := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASS")
	dbname := os.Getenv("POSTGRES_DB")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s  sslmode=disable", host, puser, pass, dbname, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return Repository{}, err
	}

	db.AutoMigrate(&User{})

	return Repository{
		DB: db,
	}, nil
}
