package request

import (
	"fmt"

	"github.com/asaskevich/govalidator"
)

type BasketCreate struct {
	Data []byte `json:"data,omitempty" valid:"length(0|2048),optional"`
}

func (bc BasketCreate) Validate() error {
	if _, err := govalidator.ValidateStruct(bc); err != nil {
		return fmt.Errorf("create request validation failed %w", err)
	}

	return nil

}
