package stock

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	npool "github.com/NpoolPlatform/message/npool/stockmgr"
	"github.com/NpoolPlatform/stock-manager/pkg/test-init" //nolint

	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
)

func init() {
	if runByGithubAction, err := strconv.ParseBool(os.Getenv("RUN_BY_GITHUB_ACTION")); err == nil && runByGithubAction {
		return
	}
	if err := testinit.Init(); err != nil {
		fmt.Printf("cannot init test stub: %v\n", err)
	}
}

func assertStock(t *testing.T, actual, expected *npool.Stock) {
	assert.Equal(t, actual.GoodID, expected.GoodID)
	assert.Equal(t, actual.Total, expected.Total)
	assert.Equal(t, actual.InService, expected.InService)
	assert.Equal(t, actual.Sold, expected.Sold)
}

func TestCRUD(t *testing.T) {
	stock := npool.Stock{
		GoodID:    uuid.New().String(),
		Total:     1000,
		InService: 20,
		Sold:      10000,
	}

	schema, err := New(context.Background(), nil)
	assert.Nil(t, err)

	info, err := schema.Create(context.Background(), &stock)
	if assert.Nil(t, err) {
		assert.NotEqual(t, info.ID, uuid.UUID{}.String())
		assertStock(t, info, &stock)
	}
}
