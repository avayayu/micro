package dao

import (
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type Driver interface {
	Connect() (*gorm.DB, *mongo.Client, error)
}
