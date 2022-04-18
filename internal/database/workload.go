package database

import (
	"errors"
	"github.com/teeverr/dummy_service/internal/domain"
	"gorm.io/gorm"
)

func (c *Client) SetCpu(data *domain.Workload) error {
	return c.db.Create(data).Error
}

func (c *Client) GetLastCpu() (domain.Workload, error) {
	var model domain.Workload
	result := c.db.Last(&model)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			model.TargetCPULoad = 0
			return model, nil
		}
		return model, result.Error
	} else {
		return model, nil
	}
}
