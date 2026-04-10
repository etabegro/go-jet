package mysql

import "github.com/go-jet/jet/v2/internal/jet"

func newJsonArrayExpression(projections ...Projection) jet.JsonArrayExpression {
	newExp := &jsonArrayExpressionImpl{
		projections: projections,
		distinct:    false,
		orderBy:     &jet.ClauseOrderBy{},
		limit:       &jet.ClauseLimit{Count: -1},
		offset:      &jet.ClauseOffset{},
	}
	newExp.refreshExpression()

	return newExp
}

type jsonArrayExpressionImpl struct {
	jet.Expression

	projections []Projection

	distinct bool
	orderBy  *jet.ClauseOrderBy
	limit    *jet.ClauseLimit
	offset   *jet.ClauseOffset
}

func (f *jsonArrayExpressionImpl) refreshExpression() {
	distinct := ""
	if f.distinct {
		distinct = "DISTINCT "
	}

	f.Expression = Func("JSON_ARRAYAGG", CustomExpression(
		jet.Token(distinct),
		Func("JSON_OBJECT", CustomExpression(jet.JsonObjProjectionList(f.projections))),
		f.orderBy,
		f.limit,
		f.offset,
	))
}

func (f *jsonArrayExpressionImpl) DISTINCT() jet.JsonArrayExpression {
	f.distinct = true
	f.refreshExpression()

	return f
}

func (f *jsonArrayExpressionImpl) ORDER_BY(orderBy ...OrderByClause) jet.JsonArrayExpression {
	f.orderBy.List = orderBy

	return f
}

func (f *jsonArrayExpressionImpl) LIMIT(limit int64) jet.JsonArrayExpression {
	f.limit.Count = limit

	return f
}

func (f *jsonArrayExpressionImpl) OFFSET(offset int64) jet.JsonArrayExpression {
	f.offset.Count = Int(offset)

	return f
}
