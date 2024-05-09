package sqlstorage

import (
	"context"
	"database/sql"
	"time"

	"github.com/elliotchance/pie/v2"
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
	ID                 uuid.UUID         `db:"id"`
	ProductName        string            `db:"product_name"`
	Quantity           int               `db:"quantity"`
	Price              float64           `db:"price"`
	ProductDescription string            `db:"product_description"`
	Status             model.OrderStatus `db:"status"`
	CreatedAt          time.Time         `db:"created_at"`
	LastUpdate         time.Time         `db:"last_update"`
	DeletedAt          sql.NullTime      `db:"deleted_at"`
}

func (f OrderRepository) OrderCreate(ctx context.Context, order model.Order) error {
	_, err := f.db.NamedExecContext(ctx, `
INSERT INTO order-service (id, product_name, quantity, price, status, created_at, last_update) VALUES (:id, :product_name, :quantity, :name, :price, :created_at, :last_update)
`, orderRecordFromModel(order))
	if err != nil {
		return err
	}

	return nil
}

// OrderDelete does a soft delete. The deletion can be done at a later time when the db isn't busy or recovered
func (f OrderRepository) OrderDelete(ctx context.Context, id uuid.UUID) error {
	result, err := f.db.ExecContext(ctx, "UPDATE order-service SET deleted_at=CURRENT_TIMESTAMP AT TIME ZONE 'UTC' WHERE id=$1", id)
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
UPDATE order-service SET product_name=:product_name, quantity=:quantity, price=:price, last_update=:last_update WHERE id=:id
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

	return nil
}

func (f OrderRepository) OrderGet(ctx context.Context, id uuid.UUID) (model.Order, error) {

	var record orderRecord
	err := f.db.GetContext(ctx, &record, "SELECT id, product_name, quantity, price, status, created_at, created_at, last_update, deleted_at FROM order-service WHERE id=$1 AND deleted_at IS NULL", id.String())
	if err != nil {
		return model.Order{}, err
	}

	return orderModelFromRecord(record), nil
}

func (f OrderRepository) OrderGetAll(ctx context.Context) ([]model.Order, error) {
	var orderRecords []orderRecord
	err := f.db.SelectContext(ctx, &orderRecords, "SELECT id, product_name, quantity, price, status, created_at, created_at, last_update, deleted_at FROM order-service WHERE deleted_at IS NULL")
	if err != nil {
		return nil, err
	}

	orders := pie.Map(orderRecords, orderModelFromRecord)
	return orders, err
}

func orderModelFromRecord(r orderRecord) model.Order {
	return model.Order{
		ID:                 r.ID,
		ProductName:        r.ProductName,
		Quantity:           r.Quantity,
		Price:              r.Price,
		ProductDescription: r.ProductDescription,
		Status:             r.Status,
		LastUpdate:         r.LastUpdate.UTC(),
		CreatedAt:          r.CreatedAt.UTC(),
	}
}
func orderRecordFromModel(order model.Order) orderRecord {
	return orderRecord{
		ID:                 order.ID,
		ProductName:        order.ProductName,
		Quantity:           order.Quantity,
		Price:              order.Price,
		ProductDescription: order.ProductDescription,
		Status:             order.Status,
		LastUpdate:         order.LastUpdate.UTC(),
		CreatedAt:          order.CreatedAt.UTC(),
	}
}
