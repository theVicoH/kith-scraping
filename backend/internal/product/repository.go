package product

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type Repository interface {
	GetAllReferences(ctx context.Context) ([]string, error)
	UpsertProducts(ctx context.Context, products []Product) ([]Product, error)
	MarkOutOfStock(ctx context.Context, missingRefs []string) error
	WasOutOfStock(ctx context.Context, reference string) (bool, error)
	GetAllProducts(ctx context.Context) ([]Product, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetAllReferences(ctx context.Context) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT reference FROM products`)
	if err != nil {
		return nil, fmt.Errorf("GetAllReferences: %w", err)
	}
	defer rows.Close()

	var refs []string
	for rows.Next() {
		var ref string
		if err := rows.Scan(&ref); err != nil {
			return nil, fmt.Errorf("GetAllReferences.Scan: %w", err)
		}
		refs = append(refs, ref)
	}
	return refs, nil
}

func (r *PostgresRepository) UpsertProducts(ctx context.Context, products []Product) ([]Product, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("UpsertProducts.BeginTx: %w", err)
	}

	stmt := `
    INSERT INTO products (reference, title, price, image_url, product_url, category, event_type, event_date, in_stock)
    VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
    ON CONFLICT (reference) DO UPDATE SET
        title = EXCLUDED.title,
        price = EXCLUDED.price,
        image_url = EXCLUDED.image_url,
        product_url = EXCLUDED.product_url,
        category = EXCLUDED.category,
        event_type = EXCLUDED.event_type,
        event_date = EXCLUDED.event_date,
        in_stock = EXCLUDED.in_stock
    RETURNING id;
    `

	result := make([]Product, len(products))
	copy(result, products)

	for i, p := range products {
		var id int
		if err := tx.QueryRowContext(ctx, stmt,
			p.Reference, p.Title, p.Price, p.ImageURL,
			p.ProductURL, p.Category, p.EventType, p.EventDate, p.InStock,
		).Scan(&id); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("UpsertProducts.QueryRow: %w", err)
		}
		result[i].ID = id
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("UpsertProducts.Commit: %w", err)
	}

	return result, nil
}

func (r *PostgresRepository) MarkOutOfStock(ctx context.Context, missingRefs []string) error {
	if len(missingRefs) == 0 {
		return nil
	}
	placeholders := make([]string, len(missingRefs))
	args := make([]interface{}, len(missingRefs))
	for i, ref := range missingRefs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = ref
	}
	query := fmt.Sprintf(
		"UPDATE products SET in_stock = FALSE WHERE reference IN (%s)",
		strings.Join(placeholders, ","),
	)
	if _, err := r.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("MarkOutOfStock: %w", err)
	}
	return nil
}

func (r *PostgresRepository) WasOutOfStock(ctx context.Context, reference string) (bool, error) {
	var inStock bool
	err := r.db.QueryRowContext(ctx,
		`SELECT in_stock FROM products WHERE reference = $1`, reference).Scan(&inStock)

	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("WasOutOfStock.Query: %w", err)
	}

	return !inStock, nil
}

func (r *PostgresRepository) GetAllProducts(ctx context.Context) ([]Product, error) {
	rows, err := r.db.QueryContext(ctx, `
        SELECT id, reference, title, price, image_url, product_url, 
               category, event_type, event_date, in_stock 
        FROM products
        ORDER BY event_date DESC
    `)
	if err != nil {
		return nil, fmt.Errorf("GetAllProducts: %w", err)
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(
			&p.ID, &p.Reference, &p.Title, &p.Price,
			&p.ImageURL, &p.ProductURL, &p.Category,
			&p.EventType, &p.EventDate, &p.InStock,
		); err != nil {
			return nil, fmt.Errorf("GetAllProducts.Scan: %w", err)
		}
		products = append(products, p)
	}
	return products, nil
}
