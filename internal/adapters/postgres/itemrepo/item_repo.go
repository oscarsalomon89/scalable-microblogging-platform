package itemrepo

import (
	"github.com/oscarsalomon89/go-hexagonal/internal/application/item"
	"github.com/oscarsalomon89/go-hexagonal/internal/platform/db"
)

type ItemRepository struct {
	db db.Connections
}

func NewItemRepository(db db.Connections) *ItemRepository {
	return &ItemRepository{db: db}
}

func (r *ItemRepository) Save(item item.Item) error {
	return r.db.MasterConn.Create(&item).Error
}
