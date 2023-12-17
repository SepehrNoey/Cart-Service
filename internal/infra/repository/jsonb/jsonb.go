package jsonb

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/SepehrNoey/Cart-Service/internal/domain/model"
)

type JSONB map[string]interface{}

func (jf JSONB) Value() (driver.Value, error) {
	return json.Marshal(jf)
}

func (jf *JSONB) Scan(value interface{}) error {
	data, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(data, &jf)
}

func BasketToJSON(basket model.Basket) JSONB {
	jsonMap := make(JSONB)
	jsonMap["id"] = strconv.Itoa(int(basket.ID))
	jsonMap["created_at"] = basket.CreatedAt.String()
	jsonMap["updated_at"] = basket.UpdatedAt.String()
	jsonMap["state"] = basket.State
	jsonMap["data"] = basket.Data

	return jsonMap
}
