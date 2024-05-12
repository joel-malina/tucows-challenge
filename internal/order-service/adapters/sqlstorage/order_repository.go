package sqlstorage

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/joel-malina/tucows-challenge/internal/order-service/model"
)

type OrderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) OrderRepository {
	return OrderRepository{
		db: db,
	}
}

type orderRecord struct {
	OrderID    uuid.UUID `db:"order_id"`
	CustomerID uuid.UUID `db:"customer_id"`
	OrderDate  time.Time `db:"order_date"`
	Status     string    `db:"status"`
	TotalPrice float64   `db:"total_price"`
}

type orderItemRecord struct {
	ItemID    uuid.UUID `db:"item_id"`
	OrderID   uuid.UUID `db:"order_id"`
	ProductID uuid.UUID `db:"product_id"`
	Quantity  int       `db:"quantity"`
	Price     float64   `db:"price"`
}

func (f OrderRepository) OrderCreate(ctx context.Context, order model.Order) error {
	_, err := f.db.NamedExecContext(ctx, `
INSERT INTO orders (order_id, customer_id, order_date, status, total_price) VALUES (:order_id, :customer_id, :order_date, :status, :total_price)
`, orderRecordFromModel(order))
	if err != nil {
		return err
	}

	// iterate over the list and add the OrderItems to the 'order_items' table
	for i := range order.OrderItems {
		_, err = f.db.NamedExecContext(ctx, `
INSERT INTO order_items (item_id, order_id, product_id, quantity, price) VALUES (:order_id, :item_id, :product_id, :quantity, :price)
`, orderItemRecordFromModel(order.OrderItems[i]))
		if err != nil {
			return err
		}
	}

	return nil
}

// OrderDelete does a soft delete. The deletion can be done at a later time when the db isn't busy or the record can be recovered
func (f OrderRepository) OrderDelete(ctx context.Context, id uuid.UUID) error {
	result, err := f.db.ExecContext(ctx, "UPDATE orders SET deleted_at=CURRENT_TIMESTAMP AT TIME ZONE 'UTC' WHERE order_id=$1", id)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return model.ErrOrderNotFound
	}

	return nil
}

func (f OrderRepository) OrderUpdate(ctx context.Context, order model.Order) error {

	result, err := f.db.NamedExecContext(ctx, `
UPDATE orders SET customer_id=:customer_id, order_date=:order_date, status=:status, total_price=:total_price WHERE order_id=:order_id
`, orderRecordFromModel(order))
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return model.ErrOrderNotFound
	}

	for i := range order.OrderItems {
		result, err = f.db.NamedExecContext(ctx, `
UPDATE order_items SET product_id=:product_id, quantity=:quantity, price=:price WHERE item_id=:item_id AND order_id=:order_id
`, orderItemRecordFromModel(order.OrderItems[i]))
		if err != nil {
			return err
		}
	}

	return nil
}

func (f OrderRepository) OrderGet(ctx context.Context, id uuid.UUID) (model.Order, error) {
	var order orderRecord
	err := f.db.GetContext(ctx, &order, "SELECT order_id, customer_id, order_date, status, total_price FROM orders WHERE order_id=$1 AND deleted_at IS NULL", id.String())
	if err != nil {
		return model.Order{}, err
	}

	// also get each item from order_items that has the same order_id
	var orderItems []orderItemRecord
	err = f.db.GetContext(ctx, &orderItems, "SELECT item_id, order_id, product_id, quantity, price FROM orders_items WHERE order_id=$1 AND deleted_at IS NULL", id.String())
	if err != nil {
		return model.Order{}, err
	}

	return orderModelFromRecord(order, orderItems), nil
}

//func (f OrderRepository) OrdersGet(ctx context.Context) ([]model.Order, error) {
//	var orderRecords []orderRecord
//	err := f.db.SelectContext(ctx, &orderRecords, "SELECT id, product_name, quantity, price, status, created_at, created_at, last_update, deleted_at FROM orders WHERE deleted_at IS NULL")
//	if err != nil {
//		return nil, err
//	}
//
//	orders := pie.Map(orderRecords, orderModelFromRecord)
//
//	return orders, err
//}

func orderModelFromRecord(dbOrder orderRecord, dbOrderItems []orderItemRecord) model.Order {

	return model.Order{
		ID:         dbOrder.OrderID,
		CustomerID: dbOrder.CustomerID,
		OrderDate:  dbOrder.OrderDate,
		Status:     model.OrderStatus(dbOrder.Status),
		TotalPrice: dbOrder.TotalPrice,
		OrderItems: orderItemFromRecord(dbOrderItems),
	}
}

func orderItemFromRecord(dbOrderItems []orderItemRecord) []model.OrderItem {
	var orderItems []model.OrderItem
	for i := range dbOrderItems {
		orderItem := model.OrderItem{
			ID:        dbOrderItems[i].ItemID,
			OrderID:   dbOrderItems[i].OrderID,
			ProductID: dbOrderItems[i].ProductID,
			Quantity:  dbOrderItems[i].Quantity,
			Price:     dbOrderItems[i].Price,
		}
		orderItems = append(orderItems, orderItem)
	}

	return orderItems
}

func orderRecordFromModel(order model.Order) orderRecord {
	return orderRecord{
		OrderID:    order.ID,
		CustomerID: order.CustomerID,
		OrderDate:  order.OrderDate.UTC(),
		Status:     string(order.Status),
		TotalPrice: order.TotalPrice,
	}
}

func orderItemRecordFromModel(orderItem model.OrderItem) orderItemRecord {
	return orderItemRecord{
		ItemID:    orderItem.ID,
		OrderID:   orderItem.OrderID,
		ProductID: orderItem.ProductID,
		Quantity:  orderItem.Quantity,
		Price:     orderItem.Price,
	}
}
