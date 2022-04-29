package stock

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	npool "github.com/NpoolPlatform/message/npool/stockmgr"
	"github.com/NpoolPlatform/stock-manager/pkg/test-init" //nolint

	constant "github.com/NpoolPlatform/stock-manager/pkg/const"

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
		InService: 0,
		Sold:      0,
	}

	schema, err := New(context.Background(), nil)
	assert.Nil(t, err)

	info, err := schema.Create(context.Background(), &stock)
	if assert.Nil(t, err) {
		assert.NotEqual(t, info.ID, uuid.UUID{}.String())
		assertStock(t, info, &stock)
	}

	stock.InService = 100
	stock.ID = info.ID

	schema, err = New(context.Background(), nil)
	assert.Nil(t, err)

	info, err = schema.Update(context.Background(), &stock)
	if assert.Nil(t, err) {
		assert.Equal(t, info, &stock)
	}

	schema, err = New(context.Background(), nil)
	assert.Nil(t, err)

	info, err = schema.Row(context.Background(), uuid.MustParse(info.ID))
	if assert.Nil(t, err) {
		assert.Equal(t, info, &stock)
	}

	schema, err = New(context.Background(), nil)
	assert.Nil(t, err)

	stock.InService = 2001
	stock.Sold = 3001

	info, err = schema.UpdateFields(context.Background(),
		uuid.MustParse(info.ID),
		map[string]interface{}{
			constant.StockFieldInService: stock.InService,
			constant.StockFieldSold:      stock.Sold,
		})
	if assert.Nil(t, err) {
		assert.Equal(t, info, &stock)
	}

	schema, err = New(context.Background(), nil)
	assert.Nil(t, err)

	stock.InService = 2002
	stock.Sold = 3002

	info, err = schema.AddFields(context.Background(),
		uuid.MustParse(info.ID),
		map[string]interface{}{
			constant.StockFieldInService: 1,
			constant.StockFieldSold:      1,
		})
	if assert.Nil(t, err) {
		assert.Equal(t, info, &stock)
	}

	assert.Nil(t, err)
	assert.NotNil(t, info)

	schema, err = New(context.Background(), nil)
	assert.Nil(t, err)

	stock.InService = 2001
	stock.Sold = 3001

	info, err = schema.SubFields(context.Background(),
		uuid.MustParse(info.ID),
		map[string]interface{}{
			constant.StockFieldInService: 1,
			constant.StockFieldSold:      1,
		})
	if assert.Nil(t, err) {
		assert.Equal(t, info, &stock)
	}

	stock1 := &npool.Stock{
		GoodID:    uuid.New().String(),
		Total:     1000,
		InService: 0,
		Sold:      0,
	}
	stock2 := &npool.Stock{
		GoodID:    uuid.New().String(),
		Total:     1000,
		InService: 0,
		Sold:      0,
	}

	schema, err = New(context.Background(), nil)
	assert.Nil(t, err)

	infos, err := schema.CreateBulk(context.Background(), []*npool.Stock{stock1, stock2})
	if assert.Nil(t, err) {
		assert.Equal(t, len(infos), 2)
		assert.NotEqual(t, infos[0].ID, uuid.UUID{}.String())
		assert.NotEqual(t, infos[1].ID, uuid.UUID{}.String())
	}
}
