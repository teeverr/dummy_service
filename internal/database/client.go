package database

import (
	"fmt"
	"github.com/teeverr/dummy_service/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Client struct {
	db *gorm.DB
}

func NewClient(cfg *domain.Config) (*Client, error) {
	var err error
	pg := &Client{}
	pg.db, err = gorm.Open(
		postgres.Open(fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=%s", cfg.Database.Address,
			cfg.Database.Username, cfg.Database.Password, cfg.Database.Name, cfg.Database.Port,
			cfg.Database.TimeZone)),
		&gorm.Config{NamingStrategy: schema.NamingStrategy{TablePrefix: fmt.Sprintf("%s_", cfg.Database.Prefix)}},
	)
	if err != nil {
		return &Client{}, err
	}

	err = pg.Migrate(domain.Workload{})
	if err != nil {
		return &Client{}, err
	}
	return pg, nil
}
func (c *Client) Migrate(model interface{}) error {
	return c.db.AutoMigrate(model)
}
