package app

import "fmt"

type Product struct {
	ID      int
	Model   string
	Company string
	Price   int
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
