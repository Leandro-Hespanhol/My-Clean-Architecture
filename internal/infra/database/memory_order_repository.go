package database

import (
	"sync"

	"MyCleanArchitecture/internal/entity"
)

type MemoryOrderRepository struct {
	orders []entity.Order
	mutex  sync.RWMutex
}

func NewMemoryOrderRepository() *MemoryOrderRepository {
	return &MemoryOrderRepository{
		orders: make([]entity.Order, 0),
	}
}

func (r *MemoryOrderRepository) Save(order *entity.Order) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	r.orders = append(r.orders, *order)
	return nil
}

func (r *MemoryOrderRepository) FindAll() ([]*entity.Order, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	result := make([]*entity.Order, len(r.orders))
	for i := range r.orders {
		result[i] = &r.orders[i]
	}
	return result, nil
}

func (r *MemoryOrderRepository) GetTotal() (int, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	return len(r.orders), nil
}
