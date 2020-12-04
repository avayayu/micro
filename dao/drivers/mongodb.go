package drivers

import (
	"context"
	"fmt"
	"time"

	"github.com/avayayu/micro/dao"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"gorm.io/gorm"
)

type MongoDrivers struct{}

func (d *MongoDrivers) Connect(config *dao.DBConfigs) (*gorm.DB, *mongo.Client, error) {

	sqlFullConnection := config.String()
	client, err := newMongoClient(sqlFullConnection, config.UserName, config.Password, config.MongoIsReplicated)
	if err != nil {
		panic(err)
	}

	return nil, client, nil
}

//NewMongoClient 根据config中的mongodb信息初始化连接
func newMongoClient(mongodbFullURL, userName, password string, IsReplicated bool) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(mongodbFullURL)
	var moptions *options.ClientOptions = clientOptions.SetMinPoolSize(5)
	if userName != "" && password != "" {
		moptions = moptions.SetAuth(options.Credential{
			AuthSource: "admin", Username: userName, Password: password,
		})
	}
	if IsReplicated {
		moptions = moptions.SetReadConcern(readconcern.Available())
		//写策略 所有数据写入副本集所有节点 才完成 性能慢 但数据一致性高
		moptions = moptions.SetWriteConcern(writeconcern.New(writeconcern.W(3)))
		// 全局读策略 从最近节点读取数据
		moptions = moptions.SetReadPreference(readpref.Nearest())
	}

	client, err := mongo.NewClient(moptions)

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		fmt.Println(err)
	}

	// client.connect()
	return client, nil
}
