package dao

import (
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type DBType uint8

const (
	_ DBType = iota
	GORMDB
	MONGO
)

type Driver interface {
	Connect() (*gorm.DB, *mongo.Client, error)
	Type() uint8
}
