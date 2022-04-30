package stock

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
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
		if assert.NotEqual(t, info.ID, uuid.UUID{}.String()) {
			stock.ID = info.ID
		}
		assert.Equal(t, info, &stock)
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

	infos, total, err := schema.Rows(context.Background(), map[string]*cruder.Cond{
		constant.FieldID: {
			Op:  cruder.EQ,
			Val: info.ID,
		},
	}, 0, 0)
	if assert.Nil(t, err) {
		assert.Equal(t, total, 1)
		assert.Equal(t, infos[0], &stock)
	}

	schema, err = New(context.Background(), nil)
	assert.Nil(t, err)

	exist, err := schema.ExistConds(context.Background(), map[string]*cruder.Cond{
		constant.FieldID: {
			Op:  cruder.EQ,
			Val: info.ID,
		},
	})
	if assert.Nil(t, err) {
		assert.Equal(t, exist, true)
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

	infos, err = schema.CreateBulk(context.Background(), []*npool.Stock{stock1, stock2})
	if assert.Nil(t, err) {
		assert.Equal(t, len(infos), 2)
		assert.NotEqual(t, infos[0].ID, uuid.UUID{}.String())
		assert.NotEqual(t, infos[1].ID, uuid.UUID{}.String())
	}

	schema, err = New(context.Background(), nil)
	assert.Nil(t, err)

	count, err := schema.Count(context.Background(), map[string]*cruder.Cond{
		constant.FieldID: {
			Op:  cruder.EQ,
			Val: info.ID,
		},
	})
	if assert.Nil(t, err) {
		assert.Equal(t, count, uint32(1))
	}

	schema, err = New(context.Background(), nil)
	assert.Nil(t, err)

	info, err = schema.Delete(context.Background(), uuid.MustParse(info.ID))
	if assert.Nil(t, err) {
		assert.Equal(t, info, &stock)
	}

	schema, err = New(context.Background(), nil)
	assert.Nil(t, err)

	count, err = schema.Count(context.Background(), map[string]*cruder.Cond{
		constant.FieldID: {
			Op:  cruder.EQ,
			Val: info.ID,
		},
	})
	if assert.Nil(t, err) {
		assert.Equal(t, count, uint32(0))
	}

	schema, err = New(context.Background(), nil)
	assert.Nil(t, err)

	_, err = schema.Row(context.Background(), uuid.MustParse(info.ID))
	assert.NotNil(t, err)
}
