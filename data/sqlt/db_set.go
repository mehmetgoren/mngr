package sqlt

import (
	"gorm.io/gorm"
)

type DbSet[T any] struct {
	db *gorm.DB
}

func (d *DbSet[T]) GetByQuery(query map[string]interface{}) []*T {
	entities := make([]*T, 0)
	d.db.Where(query).Find(&entities)
	return entities
}

func (d *DbSet[T]) GetGormDb() *gorm.DB {
	return d.db
}
