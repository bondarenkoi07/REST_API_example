package database

import (
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"strings"
)

type IterableModel interface {
	GetFields() (string, []interface{})
}

//структура БД
type Database struct {
	//пул подклчений
	pool *pgxpool.Pool
	ctx  context.Context
}

func New(user, password, connection, databaseName string) (*Database, error) {
	var ctx context.Context
	var pool *pgxpool.Pool

	ctx = context.Background()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:5432/%s", user, password, connection, databaseName)
	var err error
	pool, err = pgxpool.Connect(ctx, dsn)

	if err != nil {
		return nil, err
	} else {
		conn := &Database{pool: pool, ctx: ctx}
		return conn, nil
	}
}

func (d *Database) Create(tableName string, model IterableModel) error {

	cols, values := model.GetFields()
	valuesPlaceholders := ""
	count := len(values) - 1
	index := 0

	for ; index < count; index++ {
		valuesPlaceholders += fmt.Sprintf("$%d", index+1)
		valuesPlaceholders += ","
	}
	valuesPlaceholders += fmt.Sprintf("$%d", index+1)

	SQLStatement := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, cols, valuesPlaceholders)

	conn, err := (*d).pool.Acquire(d.ctx)

	if err != nil {
		return errors.New(err.Error() + "\n" + SQLStatement)
	}

	defer conn.Release()

	tx, err := conn.Begin(d.ctx)

	if err != nil {
		return err
	}

	ct, err := tx.Exec(d.ctx, SQLStatement, values...)

	if err != nil {
		var superErr = tx.Rollback(d.ctx)
		if superErr != nil {
			superErr = errors.New(superErr.Error() + " thrown while handling ZeroRowsAffectedError" + "\n" + SQLStatement)
			return errors.New(superErr.Error() + "\n" + err.Error() + "\n" + SQLStatement)
		} else {
			return err
		}

	} else if ct.RowsAffected() == 0 {
		var superErr = tx.Rollback(d.ctx)
		if superErr != nil {
			superErr = errors.New(superErr.Error() + " thrown while handling ZeroRowsAffectedError" + "\n" + SQLStatement)
		}
		return superErr
	} else {
		err := tx.Commit(d.ctx)
		if err != nil {
			return err
		} else {
			return nil
		}
	}
}

func (d *Database) ReadAll(tableName string) (*pgx.Rows, error) {
	SQLStatement := fmt.Sprintf("SELECT * FROM %s ORDER BY id", tableName)
	conn, err := (*d).pool.Acquire(d.ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	rows, err := conn.Query(d.ctx, SQLStatement)

	if err != nil {
		return nil, err
	} else {
		return &rows, nil
	}
}

func (d *Database) ReadOne(tableName string, id int64) (*pgx.Row, error) {
	SQLStatement := fmt.Sprintf("SELECT * FROM %s where id = $1", tableName)
	conn, err := (*d).pool.Acquire(d.ctx)

	if err != nil {
		return nil, err
	}

	defer conn.Release()

	row := conn.QueryRow(d.ctx, SQLStatement, id)

	return &row, nil
}

func (d *Database) Update(tableName string, id int64, model IterableModel) error {
	cols, values := model.GetFields()
	colSlice := strings.Split(cols, ",")
	SQLStatement := fmt.Sprintf("UPDATE %s SET ", tableName)
	placeholders := fmt.Sprintf("$%d", len(colSlice)+1)
	for index, colName := range colSlice {
		SQLStatement += fmt.Sprintf("%s = $%d,", colName, index+1)
	}
	l := len(SQLStatement)
	SQLStatement = SQLStatement[:l-1]
	conn, err := (*d).pool.Acquire(d.ctx)

	if err != nil {
		return err
	} else if conn == nil {
		return errors.New("Acquired nil pointer.")
	}

	defer conn.Release()

	tx, err := conn.Begin(d.ctx)

	if err != nil {
		return err
	} else if tx == nil {
		return errors.New("Cannot start transaction.")
	}

	values = append(values, id)

	ct, err := tx.Exec(d.ctx, SQLStatement+" where id ="+placeholders, values...)

	if err != nil {
		var superErr = tx.Rollback(d.ctx)
		if superErr != nil {
			superErr = errors.New(superErr.Error() + " thrown while handling ZeroRowsAffectedError" + "\n" + SQLStatement)
			return errors.New(superErr.Error() + "\n" + err.Error() + "\n" + SQLStatement)
		} else {
			return err
		}

	} else if ct.RowsAffected() == 0 {
		var superErr = tx.Rollback(d.ctx)
		if superErr != nil {
			superErr = errors.New(superErr.Error() + " thrown while handling ZeroRowsAffectedError" + "\n" + SQLStatement)
		}
		return superErr
	} else {
		err := tx.Commit(d.ctx)
		if err != nil {
			return err
		} else {
			return nil
		}
	}
}

func (d *Database) DeleteAll(tableName string) error {
	SQLStatement := fmt.Sprintf("DELETE FROM %s", tableName)
	conn, err := (*d).pool.Acquire(d.ctx)
	if err != nil {
		return err
	}
	defer conn.Release()
	_, err = conn.Exec(d.ctx, SQLStatement)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (d *Database) DeleteOne(tableName string, id int64) error {
	SQLStatement := fmt.Sprintf("DELETE FROM %s WHERE id = $1", tableName)

	conn, err := (*d).pool.Acquire(d.ctx)

	if err != nil {
		return err
	}

	defer conn.Release()

	ct, err := conn.Exec(d.ctx, SQLStatement, id)

	if err != nil {
		return err
	} else if ct.RowsAffected() == 0 {
		return errors.New(fmt.Sprintf("Seems where are any row with id = %d", id))
	} else {
		return nil
	}
}
