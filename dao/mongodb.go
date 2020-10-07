package dao

// readPreference 主要控制客户端 Driver 从复制集的哪个节点读取数据，这个特性可方便的实现读写分离、就近读取等策略。

// primary 只从 primary 节点读数据，这个是默认设置
// primaryPreferred 优先从 primary 读取，primary 不可服务，从 secondary 读
// secondary 只从 scondary 节点读数据
// secondaryPreferred 优先从 secondary 读取，没有 secondary 成员时，从 primary 读取
// nearest 根据网络距离就近读取
// readConcern 决定到某个读取数据时，能读到什么样的数据。

// local 能读取任意数据，这个是默认设置。
// available 对于没有做shard的数据库，它跟local一样！但如果是shared情况下，在chunk迁移的过程中，mongod实例异常宕机，导致迁移过程失败或者部分完成，会产生孤儿文件，而available可能会返回孤儿文件查询，而local根本不会去请求primary，哪怕是config server更新后。
// majority 只能读取到成功写入到大多数节点的数据；使用 majority 的方式会有诸多的限制，必须使用 WiredTiger 存储引擎，而且必须将选举协议设置为 1。
// linearizable 和 majority 类似，读取成功写入到大多数节点的数据；但是它修复了majority的一些bug,当然也要比majority在性能上损耗更多。关于它与majority的区别

//代码示例
// var client *mongo.Client // assume client is configured with write concern majority and read preference primary

// Specify the DefaultReadConcern option so any transactions started through the session will have read concern
// majority.
// The DefaultReadPreference and DefaultWriteConcern options aren't specified so they will be inheritied from client
// and be set to primary and majority, respectively.
// opts := options.Session().SetDefaultReadConcern(readconcern.Majority())
// sess, err := client.StartSession(opts)
// if err != nil {
// 	log.Fatal(err)
// }
// defer sess.EndSession(context.TODO())

// // Specify the ReadPreference option to set the read preference to primary preferred for this transaction.
// txnOpts := options.Transaction().SetReadPreference(readpref.PrimaryPreferred())
// result, err := sess.WithTransaction(context.TODO(), func(sessCtx mongo.SessionContext) (interface{}, error) {
// 	// Use sessCtx as the Context parameter for InsertOne and FindOne so both operations are run in a
// 	// transaction.

// 	coll := client.Database("db").Collection("coll")
// 	res, err := coll.InsertOne(sessCtx, bson.D{{"x", 1}})
// 	if err != nil {
// 		return nil, err
// 	}

// 	var result bson.M
// 	if err = coll.FindOne(sessCtx, bson.D{{"_id", res.InsertedID}}).Decode(result); err != nil {
// 		return nil, err
// 	}
// 	return result, err
// }, txnOpts)
// if err != nil {
// 	log.Fatal(err)
// }
// fmt.Printf("result: %v\n", result)

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.uber.org/zap"
)

//NewMongoClient 根据config中的mongodb信息初始化连接
func (db *DB) NewMongoClient(mongodbFullURL, userName, password string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(mongodbFullURL)
	var moptions *options.ClientOptions
	if userName != "" && password != "" {
		moptions = clientOptions.SetMinPoolSize(5).SetAuth(options.Credential{
			AuthSource: "admin", Username: userName, Password: password,
		})
	}

	moptions = moptions.SetReadConcern(readconcern.Available())
	//写策略 所有数据写入副本集所有节点 才完成 性能慢 但数据一致性高
	moptions = moptions.SetWriteConcern(writeconcern.New(writeconcern.W(3)))
	// 全局读策略 从最近节点读取数据
	moptions = moptions.SetReadPreference(readpref.Nearest())

	client, err := mongo.NewClient(moptions)

	if err != nil {
		db.logger.Error("fali to connect to mongodb")
		return nil, err
	}

	db.logger.Info("connected to nosql database:", zap.String("url", mongodbFullURL))

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
