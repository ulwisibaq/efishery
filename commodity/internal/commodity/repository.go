package commodity

import "github.com/ulwisibaq/efishery/commodity/internal/models"

type RestRepository interface {
	GetCommodity() (commodities []models.Commodity, err error)
	GetExchangeRate(url string) (rate float64, err error)
}
