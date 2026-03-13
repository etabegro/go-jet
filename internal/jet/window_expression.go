package jet

type commonWindowImpl struct {
	expression Expression
	window     Window
}

func (w *commonWindowImpl) over(window ...Window) {
	if len(window) > 0 {
		w.window = window[0]
	} else {
		w.window = newWindowImpl(nil)
	}
}

func (w *commonWindowImpl) serialize(statement StatementType, out *SQLBuilder, options ...SerializeOption) {
	w.expression.serialize(statement, out)
	if w.window != nil {
		out.WriteString("OVER")
		w.window.serialize(statement, out, FallTrough(options)...)
	}
}

// --------------------------------------

type windowExpression interface {
	Expression
	OVER(window ...Window) Expression
}

func newWindowExpression(exp Expression) windowExpression {
	newExp := &windowExpressionImpl{
		Expression: exp,
	}

	newExp.commonWindowImpl.expression = exp
	exp.setRoot(newExp)

	return newExp
}

type windowExpressionImpl struct {
	Expression
	commonWindowImpl
}

func (f *windowExpressionImpl) OVER(window ...Window) Expression {
	f.commonWindowImpl.over(window...)
	return f
}

func (f *windowExpressionImpl) serialize(statement StatementType, out *SQLBuilder, options ...SerializeOption) {
	f.commonWindowImpl.serialize(statement, out, FallTrough(options)...)
}

// -----------------------------------------------------

type floatWindowExpression interface {
	FloatExpression
	OVER(window ...Window) FloatExpression
}

func newFloatWindowExpression(floatExp FloatExpression) floatWindowExpression {
	newExp := &floatWindowExpressionImpl{
		FloatExpression: floatExp,
	}

	newExp.commonWindowImpl.expression = floatExp
	floatExp.setRoot(newExp)

	return newExp
}

type floatWindowExpressionImpl struct {
	FloatExpression
	commonWindowImpl
}

func (f *floatWindowExpressionImpl) OVER(window ...Window) FloatExpression {
	f.commonWindowImpl.over(window...)
	return f
}

func (f *floatWindowExpressionImpl) serialize(statement StatementType, out *SQLBuilder, options ...SerializeOption) {
	f.commonWindowImpl.serialize(statement, out, FallTrough(options)...)
}

// ------------------------------------------------

type integerWindowExpression interface {
	IntegerExpression
	OVER(window ...Window) IntegerExpression
}

func newIntegerWindowExpression(intExp IntegerExpression) integerWindowExpression {
	newExp := &integerWindowExpressionImpl{
		IntegerExpression: intExp,
	}

	newExp.commonWindowImpl.expression = intExp
	intExp.setRoot(newExp)

	return newExp
}

type integerWindowExpressionImpl struct {
	IntegerExpression
	commonWindowImpl
}

func (f *integerWindowExpressionImpl) OVER(window ...Window) IntegerExpression {
	f.commonWindowImpl.over(window...)
	return f
}

func (f *integerWindowExpressionImpl) serialize(statement StatementType, out *SQLBuilder, options ...SerializeOption) {
	f.commonWindowImpl.serialize(statement, out, FallTrough(options)...)
}

// ------------------------------------------------

type boolWindowExpression interface {
	BoolExpression
	OVER(window ...Window) BoolExpression
}

func newBoolWindowExpression(boolExp BoolExpression) boolWindowExpression {
	newExp := &boolWindowExpressionImpl{
		BoolExpression: boolExp,
	}

	newExp.commonWindowImpl.expression = boolExp
	boolExp.setRoot(newExp)

	return newExp
}

type boolWindowExpressionImpl struct {
	BoolExpression
	commonWindowImpl
}

func (f *boolWindowExpressionImpl) OVER(window ...Window) BoolExpression {
	f.commonWindowImpl.over(window...)
	return f
}

func (f *boolWindowExpressionImpl) serialize(statement StatementType, out *SQLBuilder, options ...SerializeOption) {
	f.commonWindowImpl.serialize(statement, out, FallTrough(options)...)
}

// -----------------------------------------------------

type JsonArrayExpression interface {
	Expression

	DISTINCT() JsonArrayExpression
	ORDER_BY(orderByClauses ...OrderByClause) JsonArrayExpression
	LIMIT(limit int64) JsonArrayExpression
	OFFSET(offset int64) JsonArrayExpression
}

func newJsonArrayExpression(projections ...Projection) JsonArrayExpression {
	newExp := &jsonArrayExpressionImpl{
		projections: projections,
		distinct:    false,
		orderBy:     &ClauseOrderBy{},
		limit:       &ClauseLimit{Count: -1},
		offset:      &ClauseOffset{},
	}
	newExp.Expression = newExpression(newExp)

	return newExp
}

type jsonArrayExpressionImpl struct {
	Expression

	projections []Projection

	distinct bool
	orderBy  *ClauseOrderBy
	limit    *ClauseLimit
	offset   *ClauseOffset
}

func (f *jsonArrayExpressionImpl) DISTINCT() JsonArrayExpression {
	f.distinct = true

	return f
}

func (f *jsonArrayExpressionImpl) ORDER_BY(orderBy ...OrderByClause) JsonArrayExpression {
	f.orderBy.List = orderBy

	return f
}

func (f *jsonArrayExpressionImpl) LIMIT(limit int64) JsonArrayExpression {
	f.limit.Count = limit

	return f
}

func (f *jsonArrayExpressionImpl) OFFSET(offset int64) JsonArrayExpression {
	f.offset.Count = Int(offset)

	return f
}

func (f *jsonArrayExpressionImpl) serialize(statement StatementType, out *SQLBuilder, options ...SerializeOption) {
	distinct := ""
	if f.distinct {
		distinct = "DISTINCT "
	}

	jsonProjection := Func("JSON_ARRAYAGG", CustomExpression(
		Token(distinct),
		JSON_OBJECT(f.projections...),
		f.orderBy,
		f.limit,
		f.offset,
	))

	jsonProjection.serialize(statement, out, options...)
}
