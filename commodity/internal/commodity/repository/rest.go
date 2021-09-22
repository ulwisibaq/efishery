package repository

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/ulwisibaq/efishery/commodity/internal/commodity"
	"github.com/ulwisibaq/efishery/commodity/internal/models"
)

type CommodityRepository struct {
	url string
}

func NewCommodityRepository(
	url string,
) commodity.RestRepository {
	return &CommodityRepository{
		url: url,
	}
}

func (cr CommodityRepository) GetCommodity() (commodities []models.Commodity, err error) {
	url := cr.url

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	res, err := client.Do(req)
	if err != nil {
		return
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &commodities)
	if err != nil {
		return nil, err
	}

	return
}

func (cr CommodityRepository) GetExchangeRate(url string) (rate float64, err error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return
	}

	res, err := client.Do(req)
	if err != nil {
		return
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	var rateResp models.IdrToUsdRate

	err = json.Unmarshal(body, &rateResp)
	if err != nil {
		return
	}

	rate = rateResp.IdrUsd

	return
}
