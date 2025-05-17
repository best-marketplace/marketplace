package product

import (
	"catalog/internal/models"
	"context"
	"database/sql"
)

type ProductRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (rep *ProductRepository) ViewListProducts(ctx context.Context, offset int, limit int) ([]*models.ProductListView, []string, error) {
	const query = `
		SELECT id, name, price
		FROM products
		OFFSET $1 LIMIT $2
	`

	rows, err := rep.db.QueryContext(ctx, query, offset, limit)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var (
		products []*models.ProductListView
		ids      []string
	)

	for rows.Next() {
		var p models.ProductListView
		err := rows.Scan(&p.ID, &p.Name, &p.Price)
		if err != nil {
			return nil, nil, err
		}
		products = append(products, &p)
		ids = append(ids, p.ID)
	}

	if err = rows.Err(); err != nil {
		return nil, nil, err
	}

	return products, ids, nil
}
