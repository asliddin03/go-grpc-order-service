package service

type Product struct {
	ProductID int64
	Price     int64
	Available bool
}

type InventoryService struct {
	products map[int64]Product
}

func NewInventoryService() *InventoryService {
	return &InventoryService{
		products: map[int64]Product{
			1: {
				ProductID: 1,
				Price:     1000,
				Available: true,
			},
			2: {
				ProductID: 2,
				Price:     500,
				Available: true,
			},
			10: {
				ProductID: 10,
				Price:     250,
				Available: true,
			},
			20: {
				ProductID: 20,
				Price:     700,
				Available: false,
			},
		},
	}
}

func (s *InventoryService) GetProducts(productIDs []int64) []Product {
	result := make([]Product, 0, len(productIDs))

	for _, productID := range productIDs {
		product, ok := s.products[productID]
		if !ok {
			continue
		}

		result = append(result, product)
	}

	return result
}
