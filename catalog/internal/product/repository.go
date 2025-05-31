package product

import (
	"catalog/internal/models"
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
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

func (rep *ProductRepository) ViewProduct(ctx context.Context, id string) (*models.ProductView, error) {
	const query = `
		SELECT 
			p.id,
			p.name,
			p.description,
			p.price,
			p.seller_name,
			c.name,
			p.created_at
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.id = $1
	`

	var product models.ProductView

	err := rep.db.QueryRowContext(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.SellerName,
		&product.CategoryName,
		&product.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (rep *ProductRepository) AddProduct(ctx context.Context, sellerName, name, description, categoryName string, price int) (uuid.UUID, error) {
	const op = "product.repository.AddProduct"

	var categoryID uuid.UUID
	queryCategory := `SELECT id FROM categories WHERE name = $1`
	err := rep.db.QueryRowContext(ctx, queryCategory, categoryName).Scan(&categoryID)
	if err == sql.ErrNoRows {
		return uuid.Nil, fmt.Errorf("%s: category not found", op)
	}
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: failed to query category: %w", op, err)
	}

	productID := uuid.New()

	const queryInsert = `
		INSERT INTO products (id, name, description, price, seller_name, category_id)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err = rep.db.ExecContext(ctx, queryInsert,
		productID, name, description, price, sellerName, categoryID,
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: failed to insert product: %w", op, err)
	}

	return productID, nil
}
