package database

import (
	"errors"
	"fmt"
)

func (d *Database) Save(tableName string, model IterableModel) (int64, error) {

	cols, values := model.GetFields()
	valuesPlaceholders := ""
	count := len(values) - 1
	index := 0

	for ; index < count; index++ {
		valuesPlaceholders += fmt.Sprintf("$%d", index+1)
		valuesPlaceholders += ","
	}
	valuesPlaceholders += fmt.Sprintf("$%d", index+1)

	SQLStatement := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING id", tableName, cols, valuesPlaceholders)

	conn, err := (*d).pool.Acquire(d.ctx)

	if err != nil {
		return -1, errors.New(err.Error() + "\n" + SQLStatement)
	}

	defer conn.Release()

	tx, err := conn.Begin(d.ctx)

	if err != nil {
		return -1, err
	}

	row := tx.QueryRow(d.ctx, SQLStatement, values...)

	var Id int64

	err = row.Scan(&Id)

	if err != nil {
		newErr := tx.Rollback(d.ctx)
		if newErr != nil {
			err = errors.New(newErr.Error() + "caugth while handling err:" + err.Error())
		}
		return -1, err
	} else if Id <= 0 {
		err := tx.Rollback(d.ctx)
		if err != nil {
			err = errors.New(err.Error() + "caugth while handling err: wrong id")
		} else {
			err = errors.New("wrong id")
		}
		return -1, err
	} else {
		err = tx.Commit(d.ctx)
		if err != nil {
			return -1, err
		} else {
			return Id, nil
		}
	}
}
