package app

import (
	"testing"

	"go.uber.org/mock/gomock"
)

func TestGetExpensiveProducts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockProductRepository(ctrl)

	testProducts := []Product{
		{ID: 1, Model: "iPhone", Company: "Apple", Price: 1000},
		{ID: 2, Model: "Galaxy", Company: "Samsung", Price: 800},
		{ID: 3, Model: "Xiaomi", Company: "Xiaomi", Price: 500},
	}

	mockRepo.EXPECT().GetAll().Return(testProducts, nil).Times(1)

	service := NewProductService(mockRepo)

	result, err := service.GetExpensiveProducts(700)

	if err != nil {
		t.Errorf("Не ожидалось ошибки, получили: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Ожидалось 2 товара, получили: %d", len(result))
	}

	t.Log("TestGetExpensiveProducts прошел успешно!")
}

func TestApplyDiscount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockProductRepository(ctrl)
	service := NewProductService(mockRepo)

	t.Run("Успешная скидка 10%", func(t *testing.T) {
		expectedProduct := &Product{ID: 1, Price: 1000}
		mockRepo.EXPECT().GetByID(1).Return(expectedProduct, nil).Times(1)
		mockRepo.EXPECT().UpdatePrice(1, 900).Return(nil).Times(1)

		err := service.ApplyDiscount(1, 10)

		if err != nil {
			t.Errorf("Не ожидалось ошибки: %v", err)
		}
	})

	t.Run("Невалидная скидка 150%", func(t *testing.T) {
		mockRepo2 := NewMockProductRepository(ctrl)
		service2 := NewProductService(mockRepo2)

		err := service2.ApplyDiscount(1, 150)

		if err == nil {
			t.Error("Ожидалась ошибка о невалидной скидке, но ее не было")
		}
	})
}

func TestGetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockProductRepository(ctrl)
	service := NewProductService(mockRepo)

	expectedProduct := &Product{ID: 1, Model: "iPhone", Company: "Apple", Price: 1000}

	mockRepo.EXPECT().GetByID(1).Return(expectedProduct, nil).Times(1)

	product, err := service.GetByID(1)

	if err != nil {
		t.Errorf("Не ожидалось ошибки: %v", err)
	}

	if product.ID != 1 {
		t.Errorf("Ожидался ID=1, получили: %d", product.ID)
	}

	t.Log("✅ TestGetByID прошел успешно!")
}

func TestCreateProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockProductRepository(ctrl)
	service := NewProductService(mockRepo)

	newProduct := &Product{Model: "Pixel 8", Company: "Google", Price: 90000}

	mockRepo.EXPECT().
		Create(gomock.Any()).
		DoAndReturn(func(p *Product) error {
			p.ID = 3
			return nil
		}).
		Times(1)

	err := service.Create(newProduct)

	if err != nil {
		t.Errorf("Не ожидалось ошибки: %v", err)
	}

	if newProduct.ID != 3 {
		t.Errorf("Ожидался ID=3 после создания, получили: %d", newProduct.ID)
	}

	t.Log("TestCreateProduct прошел успешно!")
}
