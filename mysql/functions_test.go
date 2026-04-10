package mysql

import (
	"testing"

	"github.com/google/uuid"
)

func TestUUIDToBin(t *testing.T) {
	assertSerialize(t, UUID_TO_BIN(String(uuid.Nil.String())), `uuid_to_bin(?)`, uuid.Nil.String())
}

func TestJSONArrayAggExpression(t *testing.T) {
	stmt := SELECT(
		JSON_ARRAYAGG(
			table1Col1.AS("col1"),
			table1ColInt.AS("colInt"),
		).DISTINCT().
			ORDER_BY(table1Col1.DESC()).
			LIMIT(5).
			OFFSET(2).AS("json"),
	).FROM(table1)

	assertStatementSql(t, stmt, `
SELECT JSON_ARRAYAGG(DISTINCT JSON_OBJECT(
          'col1', table1.col1,
          'colInt', table1.col_int
     )
     ORDER BY table1.col1 DESC
     LIMIT ?
     OFFSET ?) AS "json"
FROM db.table1;
`, int64(5), int64(2))
}

func TestJSONArrayAggExpressionCallOrder(t *testing.T) {
	stmt := SELECT(
		JSON_ARRAYAGG(
			table1Col1.AS("col1"),
			table1ColInt.AS("colInt"),
		).ORDER_BY(table1Col1.DESC()).
			LIMIT(5).
			OFFSET(2).
			DISTINCT().AS("json"),
	).FROM(table1)

	assertStatementSql(t, stmt, `
SELECT JSON_ARRAYAGG(DISTINCT JSON_OBJECT(
          'col1', table1.col1,
          'colInt', table1.col_int
     )
     ORDER BY table1.col1 DESC
     LIMIT ?
     OFFSET ?) AS "json"
FROM db.table1;
`, int64(5), int64(2))
}
