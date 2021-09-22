package service

import (
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/dgrijalva/jwt-go"
	"github.com/patrickmn/go-cache"
	"github.com/ulwisibaq/efishery/commodity/internal/commodity"
	"github.com/ulwisibaq/efishery/commodity/internal/models"
	"github.com/ulwisibaq/efishery/commodity/pkg/constants"
	"github.com/ulwisibaq/efishery/commodity/pkg/helpers"
)

type CommodityService struct {
	commodityRepo commodity.RestRepository
	cache         *cache.Cache
}

func NewCommodityService(
	commodityRepo commodity.RestRepository,
	cache *cache.Cache,
) commodity.Service {
	return &CommodityService{
		commodityRepo: commodityRepo,
		cache:         cache,
	}
}

func (cs CommodityService) GetCommodity() (commodities []models.Commodity, err error) {

	rate := float64(0)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		rate = cs.getCachedRate()

		// get exchange rate if not found on cache
		if rate == float64(0) {
			rate, err = cs.commodityRepo.GetExchangeRate(constants.UrlExchangeRateIdrUsd)
			if err != nil {
				return
			}
		}
	}()

	go func() {
		defer wg.Done()
		commodities, err = cs.commodityRepo.GetCommodity()
		if err != nil {
			return
		}
	}()

	wg.Wait()

	for i, commodity := range commodities {
		price := ""
		if commodity.Price != nil {
			priceStr := helpers.DeReferenceString(commodity.Price)

			if s, err := strconv.ParseFloat(priceStr, 64); err == nil {
				price = fmt.Sprintf("%f", s*rate)
			}

			commodities[i].Price = &price
		}

	}

	return
}

func (cs CommodityService) getCachedRate() (rate float64) {
	if v, found := cs.cache.Get("idr-to-usd"); found {
		rate = v.(float64)
		return rate
	}

	return
}

func (cs CommodityService) Verify(token string) (*models.UserClaims, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			err := errors.New(constants.ErrInvalidToken)
			return nil, err
		}
		return []byte(constants.JwtSecret), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &models.UserClaims{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		errExp := errors.New(constants.ErrExpiredToken)
		if ok && errors.Is(verr.Inner, errExp) {
			return nil, errExp
		}
		return nil, errors.New(constants.ErrInvalidToken)
	}

	payload, ok := jwtToken.Claims.(*models.UserClaims)
	if !ok {
		return nil, errors.New(constants.ErrInvalidToken)
	}

	return payload, nil
}
