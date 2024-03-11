package repository

import (
	"assigntment2/core"
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) InsertOrder(order core.Order) (int64, error) {
	// Insert the main order
	sqlOrder := `
		INSERT INTO orders (ordered_at, customer_name)
		VALUES ($1, $2)
		RETURNING id
	`

	var orderID int64
	err := r.db.QueryRow(sqlOrder, order.OrderedAt, order.CustomerName).Scan(&orderID)
	if err != nil {
		return 0, err
	}

	// Insert the items associated with the order
	for _, item := range order.Items {
		sqlItem := `
			INSERT INTO items (code, description, quantity, order_id)
			VALUES ($1, $2, $3, $4)
		`
		_, err := r.db.Exec(sqlItem, item.ItemCode, item.Description, item.Quantity, orderID)
		if err != nil {
			return 0, err
		}
	}

	return orderID, nil
}
