package product

import (
	"bytes"
	"catalog/internal/models"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/google/uuid"
)

type ProductRepository struct {
	db *sql.DB
	es *elasticsearch.Client
}

func NewRepository(db *sql.DB, es *elasticsearch.Client) *ProductRepository {
	return &ProductRepository{db: db, es: es}
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

func (rep *ProductRepository) SearchProduct(ctx context.Context, req string) ([]*models.ProductView, []string, error) {
	if req == "" {
		return nil, nil, nil
	}

	if rep.es == nil {
		return nil, nil, fmt.Errorf("elasticsearch client is not initialized")
	}

	esQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":     req,
				"fields":    []string{"title^2", "description"},
				"type":      "best_fields",
				"fuzziness": "AUTO",
			},
		},
		"size": 50,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(esQuery); err != nil {
		return nil, nil, fmt.Errorf("error encoding query: %w", err)
	}

	fmt.Printf("Elasticsearch query: %s\n", buf.String())

	res, err := rep.es.Search(
		rep.es.Search.WithContext(ctx),
		rep.es.Search.WithIndex("products"),
		rep.es.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error searching in elasticsearch: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, nil, fmt.Errorf("error parsing the elasticsearch response body: %w", err)
		}
		return nil, nil, fmt.Errorf("elasticsearch error: %v", e)
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, nil, fmt.Errorf("error parsing the elasticsearch response body: %w", err)
	}

	totalHits, _ := r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)
	fmt.Printf("Elasticsearch search results: total hits = %.0f\n", totalHits)

	var products []*models.ProductView
	var ids []string

	hits, _ := r["hits"].(map[string]interface{})
	if hits == nil {
		return nil, nil, nil
	}

	hitsList, _ := hits["hits"].([]interface{})
	if hitsList == nil {
		return nil, nil, nil
	}

	for _, hit := range hitsList {
		h, _ := hit.(map[string]interface{})
		if h == nil {
			continue
		}

		source, _ := h["_source"].(map[string]interface{})
		if source == nil {
			continue
		}

		productID, _ := source["product_id"].(string)

		fmt.Printf("Found product in Elasticsearch: ID=%s, source=%v\n", productID, source)

		product, err := rep.ViewProduct(ctx, productID)
		if err != nil {
			fmt.Printf("Error getting details for product %s: %v\n", productID, err)
			continue
		}

		products = append(products, product)
		ids = append(ids, productID)

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
