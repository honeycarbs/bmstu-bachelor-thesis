package psqlcli

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Client struct {
	DB *sqlx.DB
}

func NewClient(host, port, user, password, dbName, ssl string) (*Client, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		host, port, user, dbName, password, ssl)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("can't ping DB: %v", err)
	}

	return &Client{
		DB: db,
	}, nil
}

func (c *Client) Close(sqlx.DB) error {
	err := c.DB.Close()
	if err != nil {
		return err
	}
	return nil
}
