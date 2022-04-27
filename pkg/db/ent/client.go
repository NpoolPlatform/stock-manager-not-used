// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"log"

	"github.com/NpoolPlatform/stock-manager/pkg/db/ent/migrate"
	"github.com/google/uuid"

	"github.com/NpoolPlatform/stock-manager/pkg/db/ent/stock"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
)

// Client is the client that holds all ent builders.
type Client struct {
	config
	// Schema is the client for creating, migrating and dropping schema.
	Schema *migrate.Schema
	// Stock is the client for interacting with the Stock builders.
	Stock *StockClient
}

// NewClient creates a new client configured with the given options.
func NewClient(opts ...Option) *Client {
	cfg := config{log: log.Println, hooks: &hooks{}}
	cfg.options(opts...)
	client := &Client{config: cfg}
	client.init()
	return client
}

func (c *Client) init() {
	c.Schema = migrate.NewSchema(c.driver)
	c.Stock = NewStockClient(c.config)
}

// Open opens a database/sql.DB specified by the driver name and
// the data source name, and returns a new client attached to it.
// Optional parameters can be added for configuring the client.
func Open(driverName, dataSourceName string, options ...Option) (*Client, error) {
	switch driverName {
	case dialect.MySQL, dialect.Postgres, dialect.SQLite:
		drv, err := sql.Open(driverName, dataSourceName)
		if err != nil {
			return nil, err
		}
		return NewClient(append(options, Driver(drv))...), nil
	default:
		return nil, fmt.Errorf("unsupported driver: %q", driverName)
	}
}

// Tx returns a new transactional client. The provided context
// is used until the transaction is committed or rolled back.
func (c *Client) Tx(ctx context.Context) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, fmt.Errorf("ent: cannot start a transaction within a transaction")
	}
	tx, err := newTx(ctx, c.driver)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = tx
	return &Tx{
		ctx:    ctx,
		config: cfg,
		Stock:  NewStockClient(cfg),
	}, nil
}

// BeginTx returns a transactional client with specified options.
func (c *Client) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, fmt.Errorf("ent: cannot start a transaction within a transaction")
	}
	tx, err := c.driver.(interface {
		BeginTx(context.Context, *sql.TxOptions) (dialect.Tx, error)
	}).BeginTx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = &txDriver{tx: tx, drv: c.driver}
	return &Tx{
		ctx:    ctx,
		config: cfg,
		Stock:  NewStockClient(cfg),
	}, nil
}

// Debug returns a new debug-client. It's used to get verbose logging on specific operations.
//
//	client.Debug().
//		Stock.
//		Query().
//		Count(ctx)
//
func (c *Client) Debug() *Client {
	if c.debug {
		return c
	}
	cfg := c.config
	cfg.driver = dialect.Debug(c.driver, c.log)
	client := &Client{config: cfg}
	client.init()
	return client
}

// Close closes the database connection and prevents new queries from starting.
func (c *Client) Close() error {
	return c.driver.Close()
}

// Use adds the mutation hooks to all the entity clients.
// In order to add hooks to a specific client, call: `client.Node.Use(...)`.
func (c *Client) Use(hooks ...Hook) {
	c.Stock.Use(hooks...)
}

// StockClient is a client for the Stock schema.
type StockClient struct {
	config
}

// NewStockClient returns a client for the Stock from the given config.
func NewStockClient(c config) *StockClient {
	return &StockClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `stock.Hooks(f(g(h())))`.
func (c *StockClient) Use(hooks ...Hook) {
	c.hooks.Stock = append(c.hooks.Stock, hooks...)
}

// Create returns a create builder for Stock.
func (c *StockClient) Create() *StockCreate {
	mutation := newStockMutation(c.config, OpCreate)
	return &StockCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Stock entities.
func (c *StockClient) CreateBulk(builders ...*StockCreate) *StockCreateBulk {
	return &StockCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Stock.
func (c *StockClient) Update() *StockUpdate {
	mutation := newStockMutation(c.config, OpUpdate)
	return &StockUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *StockClient) UpdateOne(s *Stock) *StockUpdateOne {
	mutation := newStockMutation(c.config, OpUpdateOne, withStock(s))
	return &StockUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *StockClient) UpdateOneID(id uuid.UUID) *StockUpdateOne {
	mutation := newStockMutation(c.config, OpUpdateOne, withStockID(id))
	return &StockUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Stock.
func (c *StockClient) Delete() *StockDelete {
	mutation := newStockMutation(c.config, OpDelete)
	return &StockDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *StockClient) DeleteOne(s *Stock) *StockDeleteOne {
	return c.DeleteOneID(s.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *StockClient) DeleteOneID(id uuid.UUID) *StockDeleteOne {
	builder := c.Delete().Where(stock.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &StockDeleteOne{builder}
}

// Query returns a query builder for Stock.
func (c *StockClient) Query() *StockQuery {
	return &StockQuery{
		config: c.config,
	}
}

// Get returns a Stock entity by its id.
func (c *StockClient) Get(ctx context.Context, id uuid.UUID) (*Stock, error) {
	return c.Query().Where(stock.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *StockClient) GetX(ctx context.Context, id uuid.UUID) *Stock {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// Hooks returns the client hooks.
func (c *StockClient) Hooks() []Hook {
	return c.hooks.Stock
}
