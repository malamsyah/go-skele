package repository

import "gorm.io/gorm"

type Repository[T any] interface {
	Create(entity *T) error
	FindByID(id uint) (*T, error)
	FindAll() ([]T, error)
	Update(entity *T) (*T, error)
	Delete(id uint) error
}

type gormRepository[T any] struct {
	db *gorm.DB
}

func NewRepository[T any](db *gorm.DB) Repository[T] {
	return &gormRepository[T]{db: db}
}

func (r *gormRepository[T]) Create(entity *T) error {
	return r.db.Create(entity).Error
}

func (r *gormRepository[T]) FindByID(id uint) (*T, error) {
	var entity T
	if err := r.db.First(&entity, id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *gormRepository[T]) FindAll() ([]T, error) {
	var entities []T
	if err := r.db.Find(&entities).Error; err != nil {
		return nil, err
	}

	return entities, nil
}

func (r *gormRepository[T]) Update(entity *T) (*T, error) {
	if err := r.db.Save(entity).Error; err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *gormRepository[T]) Delete(id uint) error {
	return r.db.Delete(new(T), id).Error
}
