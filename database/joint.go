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

func (d *Database) GetMarketDevelopers(id int64) (*pgx.Rows, error) {
	SQLStatement := `SELECT DISTINCT developers.*, COUNT(product.id), SUM(product.count)
						FROM product
						inner join developers on developers.id = product.developerId
						WHERE product.marketId = $1
						GROUP BY developers.id
	`
	conn, err := (*d).pool.Acquire(d.ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	rows, err := conn.Query(d.ctx, SQLStatement, id)
	if err != nil {
		return nil, err
	}
	return &rows, nil
}

func (d *Database) GetMarketProducts(id int64) (*pgx.Rows, error) {
	SQLStatement := `SELECT * FROM product  WHERE marketId = $1 order by marketId`
	conn, err := (*d).pool.Acquire(d.ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	rows, err := conn.Query(d.ctx, SQLStatement, id)
	if err != nil {
		return nil, err
	}
	return &rows, nil
}

func (d *Database) GetDeveloperProducts(id int64) (*pgx.Rows, error) {
	SQLStatement := `SELECT * from product where developerId = $1`
	conn, err := (*d).pool.Acquire(d.ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	rows, err := conn.Query(d.ctx, SQLStatement, id)

	if err != nil {
		return nil, err
	}
	return &rows, nil
}
