package database

import (
	"github.com/jackc/pgx/v4"
)

func (d *Database) GetProductsCount(id int64) (*pgx.Row, error) {
	SQLStatement := `SELECT DISTINCT markets.id as ID, SUM(product.count) AS cnt
										FROM markets join product 
										ON markets.id = product.marketId
										WHERE markets.id = $1
										GROUP BY markets.id`
	conn, err := (*d).pool.Acquire(d.ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	rows := conn.QueryRow(d.ctx, SQLStatement, id)

	return &rows, nil
}
