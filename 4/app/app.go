package app

import (
	"database/sql"
	"fmt"
)

type Product struct {
	ID      int    `json:"id"`
	Model   string `json:"model"`
	Company string `json:"company"`
	Price   int    `json:"price"`
}

type ProductRepository interface {
	GetAll() ([]Product, error)
	GetByID(id int) (*Product, error)
	Create(product *Product) error
	UpdatePrice(id int, newPrice int) error
	Delete(id int) error
}

type ProductService struct {
	repo ProductRepository
}

func NewProductService(repo ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetAll() ([]Product, error) {
	return s.repo.GetAll()
}

func (s *ProductService) GetByID(id int) (*Product, error) {
	return s.repo.GetByID(id)
}

func (s *ProductService) Create(product *Product) error {
	return s.repo.Create(product)
}

func (s *ProductService) UpdatePrice(id int, newPrice int) error {
	if newPrice <= 0 {
		return fmt.Errorf("цена должна быть положительной")
	}
	return s.repo.UpdatePrice(id, newPrice)
}

func (s *ProductService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *ProductService) GetExpensiveProducts(minPrice int) ([]Product, error) {
	all, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	var result []Product
	for _, p := range all {
		if p.Price > minPrice {
			result = append(result, p)
		}
	}
	return result, nil
}

func (s *ProductService) ApplyDiscount(id int, discountPercent int) error {
	if discountPercent < 0 || discountPercent > 100 {
		return fmt.Errorf("скидка должна быть от 0 до 100")
	}

	product, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	newPrice := product.Price * (100 - discountPercent) / 100
	return s.repo.UpdatePrice(id, newPrice)
}

type ProductRepositoryDB struct {
	db *sql.DB
}

func NewProductRepositoryDB(db *sql.DB) *ProductRepositoryDB {
	return &ProductRepositoryDB{db: db}
}

func (r *ProductRepositoryDB) GetAll() ([]Product, error) {
	rows, err := r.db.Query("SELECT id, model, company, price FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Model, &p.Company, &p.Price); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *ProductRepositoryDB) GetByID(id int) (*Product, error) {
	var p Product
	err := r.db.QueryRow("SELECT id, model, company, price FROM products WHERE id = ?", id).
		Scan(&p.ID, &p.Model, &p.Company, &p.Price)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepositoryDB) Create(product *Product) error {
	result, err := r.db.Exec("INSERT INTO products (model, company, price) VALUES (?, ?, ?)",
		product.Model, product.Company, product.Price)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	product.ID = int(id)
	return nil
}

func (r *ProductRepositoryDB) UpdatePrice(id int, newPrice int) error {
	_, err := r.db.Exec("UPDATE products SET price = ? WHERE id = ?", newPrice, id)
	return err
}

func (r *ProductRepositoryDB) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM products WHERE id = ?", id)
	return err
}
