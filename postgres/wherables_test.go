package postgres_test

import (
	"testing"

	"github.com/graymeta/gmkit/postgres"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWherables(t *testing.T) {
	tests := []struct {
		name          string
		expectedQuery string
		expectedArgs  []any
		whereable     postgres.Whereable
	}{
		{
			"All",
			"",
			nil,
			postgres.All(),
		},
		{
			"ByFieldEquals",
			"WHERE test = $1",
			[]any{"test"},
			postgres.ByFieldEquals("test", "test"),
		},
		{
			"ByItemID",
			"WHERE item_id = $1",
			[]any{"test"},
			postgres.ByItemID("test"),
		},
		{
			"ByID",
			"WHERE id = $1",
			[]any{"test"},
			postgres.ByID("test"),
		},
		{
			"ByLocationID",
			"WHERE location_id = $1",
			[]any{"test"},
			postgres.ByLocationID("test"),
		},
		{
			"In",
			"WHERE test IN ($1, $2, $3)",
			[]any{"one", "two", "three"},
			postgres.ByFieldIn("test", "one", "two", "three"),
		},
		{
			"ByIDsIn",
			"WHERE id IN ($1, $2, $3)",
			[]any{"one", "two", "three"},
			postgres.ByIDsIn("one", "two", "three"),
		},
		{
			"ByItemIDsIn",
			"WHERE item_id IN ($1, $2, $3)",
			[]any{"one", "two", "three"},
			postgres.ByItemIDsIn("one", "two", "three"),
		},
		{
			"Like",
			"WHERE test LIKE $1",
			[]any{"test"},
			postgres.Like("test", "test"),
		},
		{
			"And",
			"WHERE (id = $1 AND item_id = $2)",
			[]any{"one", "two"},
			postgres.And(
				postgres.ByID("one"),
				postgres.ByItemID("two"),
			),
		},
		{
			"Or",
			"WHERE (id = $1 OR item_id = $2)",
			[]any{"one", "two"},
			postgres.Or(
				postgres.ByID("one"),
				postgres.ByItemID("two"),
			),
		},
		{
			"And and Or",
			"WHERE (id = $1 AND (item_id = $2 OR item_id = $3))",
			[]any{"one", "two", "three"},
			postgres.And(
				postgres.ByID("one"),
				postgres.Or(
					postgres.ByItemID("two"),
					postgres.ByItemID("three"),
				),
			),
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			query, args, err := postgres.ToWhereClause(&rebinder{}, test.whereable)

			require.NoError(t, err)
			assert.Equal(t, test.expectedQuery, query)
			assert.Equal(t, test.expectedArgs, args)
		})
	}
}

func TestQueryOption(t *testing.T) {
	tests := []struct {
		name          string
		expectedQuery string
		queryOption   postgres.QueryOption
	}{
		{
			"Limit",
			"LIMIT 123",
			postgres.Limit(123),
		},
		{
			"Offset",
			"OFFSET 456",
			postgres.Offset(456),
		},
		{
			"OrderBy",
			"ORDER BY one ASC, two DESC",
			postgres.OrderBy(
				postgres.OrderPair{Col: "one", Desc: false},
				postgres.OrderPair{Col: "two", Desc: true},
			),
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expectedQuery, test.queryOption.Clause())
		})
	}
}

func TestQueryOptions(t *testing.T) {
	expectedQuery := "ORDER BY test ASC LIMIT 123 OFFSET 456"
	actualQuery := postgres.QueryOptions{
		postgres.Offset(456),
		postgres.Limit(123),
		postgres.OrderBy(postgres.OrderPair{Col: "test", Desc: false}),
	}.Clause()

	assert.Equal(t, expectedQuery, actualQuery)
}

type rebinder struct{}

func (rebinder) Rebind(in string) string {
	return sqlx.Rebind(sqlx.DOLLAR, in)
}
