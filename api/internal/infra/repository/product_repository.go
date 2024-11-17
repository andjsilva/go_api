package repository

import (
	"api/internal/entity"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository struct {
	DB *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{DB: db}
}

func (r *ProductRepository) Create(product *entity.Product) error {
	_, err := r.DB.Exec(context.Background(), "INSERT INTO products (id, name, price) VALUES ($1, $2, $3)",
		product.ID, product.Name, product.Price)
	if err != nil {
		return err
	}
	return nil
}

func (r *ProductRepository) FindAll() ([]*entity.Product, error) {
	rows, err := r.DB.Query(context.Background(), "SELECT id, name, price FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*entity.Product
	for rows.Next() {
		var product entity.Product
		err = rows.Scan(&product.ID, &product.Name, &product.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}
	return products, nil
}
