// Code generated by entc, DO NOT EDIT.

package ent

import (
	"github.com/NpoolPlatform/stock-manager/pkg/db/ent/schema"
	"github.com/NpoolPlatform/stock-manager/pkg/db/ent/stock"
	"github.com/google/uuid"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	stockFields := schema.Stock{}.Fields()
	_ = stockFields
	// stockDescID is the schema descriptor for id field.
	stockDescID := stockFields[0].Descriptor()
	// stock.DefaultID holds the default value on creation for the id field.
	stock.DefaultID = stockDescID.Default.(func() uuid.UUID)
}
