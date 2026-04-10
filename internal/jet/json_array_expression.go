package jet

// JsonArrayExpression interface for JSON aggregate functions.
type JsonArrayExpression interface {
	Expression

	DISTINCT() JsonArrayExpression
	ORDER_BY(orderByClauses ...OrderByClause) JsonArrayExpression
	LIMIT(limit int64) JsonArrayExpression
	OFFSET(offset int64) JsonArrayExpression
}
