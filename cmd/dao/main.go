package main

import "github.com/avayayu/micro/dao"

func main() {
	dbConfig := dao.DBOptions{
		Oracle: true,
	}

	db := dao.NewDatabase(&dbConfig)

	db.SetOracleConfig(&dao.OracleConfig{
		Base: dao.Base{
			URL:      "172.27.232.73",
			Port:     "1521",
			UserName: "system",
			Password: "123456",
			DBName:   "ORCLCDB",
		},
	})

	db.Connect()

}
