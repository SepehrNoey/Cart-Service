package request

import (
	"fmt"

	"github.com/asaskevich/govalidator"
)

type BasketCreate struct {
	Data string `json:"data,omitempty" valid:"length(0|2048),optional"`
}

func (bc BasketCreate) CreateValidate() error {
	if _, err := govalidator.ValidateStruct(bc); err != nil {
		return fmt.Errorf("create request validation failed %w", err)
	}

	return nil

}

type BasketUpdate struct {
	Data  string `json:"data,omitempty" valid:"length(0|2048),optional"`
	State string `josn:"state,omitempty" valid:"in(COMPLETED|PENDING),optional"`
}

func (bu BasketUpdate) UpdateValidate() error {
	if _, err := govalidator.ValidateStruct(bu); err != nil {
		return fmt.Errorf("update request validation failed %w", err)
	}

	return nil
}
