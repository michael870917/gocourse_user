package bootstrap

import (
	"fmt"
	"log"
	"os"

	"github.com/gorilla/mux"
	"github.com/michael870917/gocourse_domain/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func DBConnector() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	if os.Getenv("DB_DEBUG") == "true" {
		db = db.Debug()
	}

	if os.Getenv("DB_MIGRATE") == "true" {
		if err := db.AutoMigrate(&domain.User{}); err != nil {
			log.Fatal(err)
		}
	}

	return db, nil
}

func InitLogger() *log.Logger {
	return log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
}

func LoadRouter() *mux.Router {
	router := mux.NewRouter()
	return router
}
