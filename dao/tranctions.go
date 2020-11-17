package dao

import (
	"context"

	"gorm.io/gorm"
)

type Transactions struct {
	session         *gorm.DB
	Ctx             context.Context
	subTransactions []SubTransactions
}

type SubTransactions func(session *gorm.DB) error

func (tran *Transactions) Run() error {

	session := tran.session.Begin()

	if tran.Ctx != nil {
		session = session.WithContext(tran.Ctx)
	}

	for _, transaction := range tran.subTransactions {
		if err := transaction(session); err != nil {
			session.Rollback()
			return err
		}
	}

	if err := tran.session.Commit().Error; err != nil {
		tran.session.Rollback()
		return err
	}
	return nil
}

func (tran *Transactions) SubExecute(sub func(session *gorm.DB) error) {
	tran.subTransactions = append(tran.subTransactions, sub)
}
