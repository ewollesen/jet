package sqlbuilder

import (
	"gotest.tools/assert"
	"testing"
)

func TestCase1(t *testing.T) {
	query := CASE().
		WHEN(table3Col1.EQ(Int(1))).THEN(table3Col1.ADD(Int(1))).
		WHEN(table3Col1.EQ(Int(2))).THEN(table3Col1.ADD(Int(2)))

	queryData := &queryData{}

	err := query.serialize(select_statement, queryData)

	assert.NilError(t, err)
	assert.Equal(t, queryData.buff.String(), `(CASE WHEN table3.col1 = $1 THEN table3.col1 + $2 WHEN table3.col1 = $3 THEN table3.col1 + $4 END)`)
}

func TestCase2(t *testing.T) {
	query := CASE(table3Col1).
		WHEN(Int(1)).THEN(table3Col1.ADD(Int(1))).
		WHEN(Int(2)).THEN(table3Col1.ADD(Int(2))).
		ELSE(Int(0))

	queryData := &queryData{}

	err := query.serialize(select_statement, queryData)

	assert.NilError(t, err)
	assert.Equal(t, queryData.buff.String(), `(CASE table3.col1 WHEN $1 THEN table3.col1 + $2 WHEN $3 THEN table3.col1 + $4 ELSE $5 END)`)
}

func TestInterval(t *testing.T) {
	query := INTERVAL(`6 years 5 months 4 days 3 hours 2 minutes 1 second`)

	queryData := &queryData{}

	err := query.serialize(select_statement, queryData)

	assert.NilError(t, err)
	assert.Equal(t, queryData.buff.String(), `INTERVAL $1`)
}