package repository

import (
	"assigntment2/core"
	"database/sql"
	"time"
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

func (r *Repository) GetAllOrders() ([]core.OrderResponse, error) {
	// Query database
	rows, err := r.db.Query("SELECT o.id, o.ordered_at, o.customer_name, i.code, i.description, i.quantity FROM orders o LEFT JOIN items i ON o.id = i.order_id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Initialize map to store orders
	orders := make(map[int64]*core.OrderResponse)

	// Iterate through rows
	for rows.Next() {
		var orderID int64
		var orderedAt time.Time
		var customerName string
		var itemCode, description string
		var quantity int

		// Scan row data
		err := rows.Scan(&orderID, &orderedAt, &customerName, &itemCode, &description, &quantity)
		if err != nil {
			return nil, err
		}

		// Check if order already exists, if not, create new order
		order, ok := orders[orderID]
		if !ok {
			order = &core.OrderResponse{
				ID:           orderID,
				OrderedAt:    orderedAt,
				CustomerName: customerName,
			}
			orders[orderID] = order
		}

		// Append item to order's items
		order.Items = append(order.Items, core.Item{
			ItemCode:    itemCode,
			Description: description,
			Quantity:    quantity,
		})
	}

	// Convert map to slice
	var result []core.OrderResponse
	for _, order := range orders {
		result = append(result, *order)
	}

	return result, nil
}
