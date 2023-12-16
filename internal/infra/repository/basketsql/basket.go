package basketsql

import (
	"context"
	"errors"

	"github.com/SepehrNoey/Cart-Service/internal/domain/model"
	"github.com/SepehrNoey/Cart-Service/internal/domain/repository/basketrepo"
	"gorm.io/gorm"
)

type BasketDTO struct {
	model.Basket
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(ctx context.Context, basket model.Basket) error {
	result := r.db.WithContext(ctx).Create(&BasketDTO{Basket: basket})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return basketrepo.ErrBasketIDDuplicate
		}

		return result.Error
	}

	return nil
}

func (r *Repository) Get(_ context.Context, cmd basketrepo.GetCommand) []model.Basket {
	var basketDTOs []BasketDTO

	var dto BasketDTO
	var conditions []string

	if cmd.ID != nil {
		dto.ID = *cmd.ID
		conditions = append(conditions, "ID")
	}
	if cmd.CreatedAt != nil {
		dto.CreatedAt = *cmd.CreatedAt
		conditions = append(conditions, "CreatedAt")
	}
	if cmd.UpdatedAt != nil {
		dto.UpdatedAt = *cmd.UpdatedAt
		conditions = append(conditions, "UpdatedAt")
	}
	if cmd.State != nil {
		dto.State = *cmd.State
		conditions = append(conditions, "State")
	}

	if len(conditions) == 0 {
		if err := r.db.Find(&basketDTOs); err.Error != nil {
			return nil
		}
	} else {
		if err := r.db.Where(&dto, conditions).Find(&basketDTOs); err.Error != nil {
			// we may need to change here
			return nil
		}
	}

	baskets := make([]model.Basket, len(basketDTOs))
	for i, dto := range basketDTOs {
		baskets[i] = dto.Basket
	}

	return baskets

}

func (r *Repository) Update(ctx context.Context, basket model.Basket) error {
	result := r.db.WithContext(ctx).Save(&BasketDTO{Basket: basket})
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, cmd basketrepo.GetCommand) error {
	result := r.db.WithContext(ctx).Delete(&BasketDTO{Basket: model.Basket{ID: *cmd.ID}})
	if result.Error != nil {
		return result.Error
	}

	return nil
}
