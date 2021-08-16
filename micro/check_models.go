package micro

import (
	"gogs.buffalo-robot.com/zouhy/micro/dao"
	"gogs.buffalo-robot.com/zouhy/micro/models"
)

func CheckEnityExists(db dao.DAO, model models.MicroModels, id models.Int64Str) (bool, error) {
	if flag, err := db.NewQuery().WhereQuery("id = ?", id).First(model, model); err != nil {
		return false, err
	} else {
		if !flag {
			return false, nil
		}
	}

	return true, nil
}
