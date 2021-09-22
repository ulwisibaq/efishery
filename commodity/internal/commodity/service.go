package commodity

import "github.com/ulwisibaq/efishery/commodity/internal/models"

type Service interface {
	GetCommodity() (commodities []models.Commodity, err error)
	Verify(token string) (*models.UserClaims, error)
}
