package storage

import (
	"errors"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ryoeuyo/demoexamen/internal/domain/order"
)

type Storage struct {
	orders []order.Order
}

func NewStorage() (*Storage, error) {
	orders := make([]order.Order, 0, 20)
	return &Storage{
		orders: orders,
	}, nil
}

func (s *Storage) Create(o *order.Order) (int64, error) {
	s.orders = append(s.orders, *o)

	return o.ID, nil
}

func (s *Storage) FetchAll() (*[]order.Order, error) {
	return &s.orders, nil
}

func (s *Storage) Update(id int64, new order.Order) error {
	for i, v := range s.orders {
		log.Print(v.ID, id)
		if id == v.ID {
			new.ID = s.orders[i].ID
			new.CreatedAt = s.orders[i].CreatedAt
			new.UpdatedAt = time.Now()
			s.orders[i] = new
			return nil
		}
	}

	return errors.New("not found")
}
