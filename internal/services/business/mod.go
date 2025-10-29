package business

import (
	"time"

	"irj/internal/config"
	queries "irj/internal/postgres/_generated"
	"irj/internal/smtp"
	"irj/pkg/api"
	"irj/pkg/framework"
	"irj/pkg/utils"
)

const defaultTimeOut = time.Second * 30

type BusinessService struct {
	stopper         *utils.AppStopper
	config          *config.Config
	smtpService     *smtp.SMTPService
	postgresService *framework.DB[*queries.Queries]
}

//nolint:lll
func NewBusinessService(stopper *utils.AppStopper, cfg *config.Config, smtpService *smtp.SMTPService, postgresService *framework.DB[*queries.Queries]) *BusinessService {
	return &BusinessService{
		stopper:         stopper,
		config:          cfg,
		smtpService:     smtpService,
		postgresService: postgresService,
	}
}

func parseLocation(
	city,
	departement,
	region,
	country []byte,
) (*api.BasicFilter, *api.BasicFilter, *api.BasicFilter, *api.BasicFilter, error) {
	var (
		c   api.BasicFilter
		d   api.BasicFilter
		r   api.BasicFilter
		p   api.BasicFilter
		err error
	)

	err = c.UnmarshalBinary(city)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	err = d.UnmarshalBinary(departement)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	err = r.UnmarshalBinary(region)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	err = p.UnmarshalBinary(country)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return &c, &d, &r, &p, nil
}

func interfaceSliceToBasicFilterSlice(input []interface{}) []*api.BasicFilter {
	result := make([]*api.BasicFilter, 0, len(input))
	for _, v := range input {
		m, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		bf := &api.BasicFilter{}

		switch id := m["id"].(type) {
		case int32:
			bf.ID = utils.PtrTo(int64(id))
		case float64:
			bf.ID = utils.PtrTo(int64(id))
		}

		if label, ok := m["name"].(string); ok {
			bf.Name = &label
		}

		result = append(result, bf)
	}

	return result
}
