package activecampaign

import (
	"fmt"
	"net/http"
	"net/url"

	a_types "github.com/leapforce-libraries/go_activecampaign/types"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	go_types "github.com/leapforce-libraries/go_types"
)

type DealStages struct {
	DealStages []DealStage `json:"dealStages"`
	Meta       Meta        `json:"meta"`
}

type DealStage struct {
	GroupId     go_types.Int64String           `json:"group"`
	Title       string                         `json:"title"`
	Color       string                         `json:"color"`
	Order       go_types.Int64String           `json:"order"`
	Width       go_types.Int64String           `json:"width"`
	DealOrder   string                         `json:"dealOrder"`
	CardRegion1 *go_types.String               `json:"cardRegion1"`
	CardRegion2 *go_types.String               `json:"cardRegion2"`
	CardRegion3 *go_types.String               `json:"cardRegion3"`
	CardRegion4 *go_types.String               `json:"cardRegion4"`
	CardRegion5 *go_types.String               `json:"cardRegion5"`
	CreatedDate a_types.DateTimeTimezoneString `json:"cdate"`
	UpdatedDate a_types.DateTimeTimezoneString `json:"udate"`
	Links       *Links                         `json:"links"`
	Id          go_types.Int64String           `json:"id"`
}

type GetDealStagesConfig struct {
	Limit        *uint64
	Offset       *uint64
	Title        *string
	GroupId      *int64
	OrderByTitle *OrderByDirection
}

func (service *Service) GetDealStages(getDealStagesConfig *GetDealStagesConfig) (*DealStages, bool, *errortools.Error) {
	params := url.Values{}

	dealStages := DealStages{}
	rowCount := uint64(0)
	limit := defaultLimit

	if getDealStagesConfig != nil {
		if getDealStagesConfig.Limit != nil {
			limit = getLimit(*getDealStagesConfig.Limit)
		}
		if getDealStagesConfig.Offset != nil {
			service.nextOffsets.DealStage = *getDealStagesConfig.Offset
		}
		if getDealStagesConfig.Title != nil {
			params.Add("filters[title]", *getDealStagesConfig.Title)
		}
		if getDealStagesConfig.GroupId != nil {
			params.Add("filters[d_groupid]", fmt.Sprintf("%v", *getDealStagesConfig.GroupId))
		}
		if getDealStagesConfig.OrderByTitle != nil {
			params.Add("orders[title]", string(*getDealStagesConfig.OrderByTitle))
		}
	}
	params.Add("limit", fmt.Sprintf("%v", limit))

	for {
		params.Set("offset", fmt.Sprintf("%v", service.nextOffsets.DealStage))

		dealStagesBatch := DealStages{}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.url(fmt.Sprintf("dealStages?%s", params.Encode())),
			ResponseModel: &dealStagesBatch,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, false, e
		}

		if len(dealStagesBatch.DealStages) < int(limit) {
			service.nextOffsets.DealStage = 0
			break
		}

		service.nextOffsets.DealStage += limit
		rowCount += limit

		if rowCount >= service.maxRowCount {
			return &dealStages, true, nil
		}
	}

	return &dealStages, false, nil
}
