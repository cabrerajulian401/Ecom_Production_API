package orders

import (
	"context"
	"errors"
	"fmt"

	repo "github.com/cabrerajulian401/ecom/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5"
)

var (
	ErrProductNotFound = errors.New("product not found")
)

var (
	ErrProductOutOfStock = errors.New("product does not have enough stock")
)

type Service interface {
	PlaceOrder(ctx context.Context, tempOrder createOrderParams) (repo.Order, error)
}

type svc struct {
	repo *repo.Queries
	db   *pgx.Conn
}

// dependency injection for a new instance of a service
func NewService(repo *repo.Queries, db *pgx.Conn) Service {
	return &svc{
		repo: repo,
		db:   db,
	}
}

// Sending a payload of customer parametrs
func (s *svc) PlaceOrder(ctx context.Context, tempOrder createOrderParams) (repo.Order, error) {

	// Payload validation - checking if request items exist on DB

	if tempOrder.CustomerID == 0 {
		return repo.Order{}, fmt.Errorf("customer ID is required")
	}
	if len(tempOrder.Items) == 0 {
		return repo.Order{}, fmt.Errorf("at least one item is required")
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return repo.Order{}, err
	}
	defer tx.Rollback(ctx) // Rollback after exiting
	qtx := s.repo.WithTx(tx)

	// create an orders
	order, err := qtx.CreateOrder(ctx, tempOrder.CustomerID)
	if err != nil {
		return repo.Order{}, err
	}
	// look for the product if exists | if product does not exist we can rollback transaction and throw error
	for _, item := range tempOrder.Items {
		product, err := qtx.FindProductByID(ctx, item.ProductID)
		if err != nil {
			return repo.Order{}, ErrProductNotFound
		}
		if product.Quantity < item.Quantity {
			return repo.Order{}, ErrProductOutOfStock
		}
		// create order item | if product is in inventory then we can create an order item
		_, err = qtx.CreateOrderItem(ctx, repo.CreateOrderItemParams{
			OrderID:    order.ID,
			ProductID:  item.ProductID,
			Quantity:   item.Quantity,
			PriceCents: product.PriceInCents,
		})

		if err != nil {
			return repo.Order{}, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return repo.Order{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return order, nil
}
