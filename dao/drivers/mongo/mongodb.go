package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"gorm.io/gorm"
)

type MongoDrivers struct {
	Configs *MongoConfigs
}

type MongoConfigs struct {
	URL                  string `json:"URL" yaml:"URL"`
	Port                 string `json:"Port" yarml:"Port"`
	UserName             string `json:"userName" yaml:"userName"`
	Password             string `json:"password" yaml:"password"`
	DBName               string `json:"dbName" yaml:"DBName"`
	MongoIsReplicated    bool   `json:"isReplicated"`
	MongoReplicatedName  string `json:"replicatedName"`
	FullConnectionString string `json:"fullConnectionString"`
}

func (c *MongoConfigs) String() string {
	if c.MongoIsReplicated {
		c.FullConnectionString = fmt.Sprintf("mongodb://%s/?replicaSet=%s", c.URL, c.MongoReplicatedName)
	} else {
		c.FullConnectionString = fmt.Sprintf("mongodb://%s:%s", c.URL, c.Port)
	}
	return c.FullConnectionString
}

func (d *MongoDrivers) Connect() (*gorm.DB, *mongo.Client, error) {

	sqlFullConnection := d.Configs.String()
	client, err := newMongoClient(sqlFullConnection, d.Configs.UserName, d.Configs.Password, d.Configs.MongoIsReplicated)
	if err != nil {
		panic(err)
	}
	return nil, client, nil
}

func (d *MongoDrivers) Type() uint8 {
	return 2
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
