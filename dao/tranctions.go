package dao

import (
	"context"
)

type transactions struct {
	query           *QueryOptions
	Ctx             context.Context
	subTransactions []SubTransactions
}

type SubTransactions func(query Query) error

func (tran *transactions) Commit() error {

	tran.query.session = tran.query.session.Begin()

	if tran.Ctx != nil {
		tran.query.session = tran.query.session.WithContext(tran.Ctx)
	}

	for _, transaction := range tran.subTransactions {
		if err := transaction(tran.query); err != nil {
			tran.query.session.Rollback()
			return err
		}
	}

	if err := tran.query.session.Commit().Error; err != nil {
		tran.query.session.Rollback()
		return err
	}
	return nil
}

func (tran *transactions) Execute(sub func(query Query) error) {
	tran.subTransactions = append(tran.subTransactions, sub)
}
