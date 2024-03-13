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
	rows, err := r.db.Query("SELECT o.id, o.ordered_at, o.customer_name, i.code, i.description, i.quantity FROM orders o LEFT JOIN items i ON o.id = i.order_id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make(map[int64]*core.OrderResponse)

	for rows.Next() {
		var orderID int64
		var orderedAt time.Time
		var customerName string
		var itemCode, description string
		var quantity int

		err := rows.Scan(&orderID, &orderedAt, &customerName, &itemCode, &description, &quantity)
		if err != nil {
			return nil, err
		}

		order, ok := orders[orderID]
		if !ok {
			order = &core.OrderResponse{
				ID:           orderID,
				OrderedAt:    orderedAt,
				CustomerName: customerName,
			}
			orders[orderID] = order
		}

		order.Items = append(order.Items, core.Item{
			ItemCode:    itemCode,
			Description: description,
			Quantity:    quantity,
		})
	}

	var result []core.OrderResponse
	for _, order := range orders {
		result = append(result, *order)
	}

	return result, nil
}

func (r *Repository) DeleteOrder(orderID int64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM items WHERE order_id = $1", orderID)
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM orders WHERE id = $1", orderID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateOrder(orderID int64, order core.OrderResponse) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("UPDATE orders SET ordered_at = $1, customer_name = $2 WHERE id = $3",
		order.OrderedAt, order.CustomerName, orderID)
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM items WHERE order_id = $1", orderID)
	if err != nil {
		return err
	}

	for _, item := range order.Items {
		_, err = tx.Exec("INSERT INTO items (order_id, code, description, quantity) VALUES ($1, $2, $3, $4)",
			orderID, item.ItemCode, item.Description, item.Quantity)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
