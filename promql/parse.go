package promql

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/pkg/value"
	"github.com/prometheus/prometheus/util/strutil"
)

type parser struct {
	lex		*lexer
	token		[3]item
	peekCount	int
}
type ParseErr struct {
	Line, Pos	int
	Err		error
}

func (e *ParseErr) Error() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if e.Line == 0 {
		return fmt.Sprintf("parse error at char %d: %s", e.Pos, e.Err)
	}
	return fmt.Sprintf("parse error at line %d, char %d: %s", e.Line, e.Pos, e.Err)
}
func ParseExpr(input string) (Expr, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	p := newParser(input)
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	err = p.typecheck(expr)
	return expr, err
}
func ParseMetric(input string) (m labels.Labels, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	p := newParser(input)
	defer p.recover(&err)
	m = p.metric()
	if p.peek().typ != itemEOF {
		p.errorf("could not parse remaining input %.15q...", p.lex.input[p.lex.lastPos:])
	}
	return m, nil
}
func ParseMetricSelector(input string) (m []*labels.Matcher, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	p := newParser(input)
	defer p.recover(&err)
	name := ""
	if t := p.peek().typ; t == itemMetricIdentifier || t == itemIdentifier {
		name = p.next().val
	}
	vs := p.VectorSelector(name)
	if p.peek().typ != itemEOF {
		p.errorf("could not parse remaining input %.15q...", p.lex.input[p.lex.lastPos:])
	}
	return vs.LabelMatchers, nil
}
func newParser(input string) *parser {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	p := &parser{lex: lex(input)}
	return p
}
func (p *parser) parseExpr() (expr Expr, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	defer p.recover(&err)
	for p.peek().typ != itemEOF {
		if p.peek().typ == itemComment {
			continue
		}
		if expr != nil {
			p.errorf("could not parse remaining input %.15q...", p.lex.input[p.lex.lastPos:])
		}
		expr = p.expr()
	}
	if expr == nil {
		p.errorf("no expression found in input")
	}
	return
}

type sequenceValue struct {
	value	float64
	omitted	bool
}

func (v sequenceValue) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if v.omitted {
		return "_"
	}
	return fmt.Sprintf("%f", v.value)
}
func parseSeriesDesc(input string) (labels.Labels, []sequenceValue, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	p := newParser(input)
	p.lex.seriesDesc = true
	return p.parseSeriesDesc()
}
func (p *parser) parseSeriesDesc() (m labels.Labels, vals []sequenceValue, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	defer p.recover(&err)
	m = p.metric()
	const ctx = "series values"
	for {
		for p.peek().typ == itemSpace {
			p.next()
		}
		if p.peek().typ == itemEOF {
			break
		}
		if p.peek().typ == itemBlank {
			p.next()
			times := uint64(1)
			if p.peek().typ == itemTimes {
				p.next()
				times, err = strconv.ParseUint(p.expect(itemNumber, ctx).val, 10, 64)
				if err != nil {
					p.errorf("invalid repetition in %s: %s", ctx, err)
				}
			}
			for i := uint64(0); i < times; i++ {
				vals = append(vals, sequenceValue{omitted: true})
			}
			if t := p.expectOneOf(itemSpace, itemEOF, ctx).typ; t == itemEOF {
				break
			}
			continue
		}
		sign := 1.0
		if t := p.peek().typ; t == itemSUB || t == itemADD {
			if p.next().typ == itemSUB {
				sign = -1
			}
		}
		var k float64
		if t := p.peek().typ; t == itemNumber {
			k = sign * p.number(p.expect(itemNumber, ctx).val)
		} else if t == itemIdentifier && p.peek().val == "stale" {
			p.next()
			k = math.Float64frombits(value.StaleNaN)
		} else {
			p.errorf("expected number or 'stale' in %s but got %s (value: %s)", ctx, t.desc(), p.peek())
		}
		vals = append(vals, sequenceValue{value: k})
		if t := p.peek(); t.typ == itemSpace {
			continue
		} else if t.typ == itemEOF {
			break
		} else if t.typ != itemADD && t.typ != itemSUB {
			p.errorf("expected next value or relative expansion in %s but got %s (value: %s)", ctx, t.desc(), p.peek())
		}
		sign = 1.0
		if p.next().typ == itemSUB {
			sign = -1.0
		}
		offset := sign * p.number(p.expect(itemNumber, ctx).val)
		p.expect(itemTimes, ctx)
		times, err := strconv.ParseUint(p.expect(itemNumber, ctx).val, 10, 64)
		if err != nil {
			p.errorf("invalid repetition in %s: %s", ctx, err)
		}
		for i := uint64(0); i < times; i++ {
			k += offset
			vals = append(vals, sequenceValue{value: k})
		}
		if t := p.expectOneOf(itemSpace, itemEOF, ctx).typ; t == itemEOF {
			break
		}
	}
	return m, vals, nil
}
func (p *parser) typecheck(node Node) (err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	defer p.recover(&err)
	p.checkType(node)
	return nil
}
func (p *parser) next() item {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if p.peekCount > 0 {
		p.peekCount--
	} else {
		t := p.lex.nextItem()
		for t.typ == itemComment {
			t = p.lex.nextItem()
		}
		p.token[0] = t
	}
	if p.token[p.peekCount].typ == itemError {
		p.errorf("%s", p.token[p.peekCount].val)
	}
	return p.token[p.peekCount]
}
func (p *parser) peek() item {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if p.peekCount > 0 {
		return p.token[p.peekCount-1]
	}
	p.peekCount = 1
	t := p.lex.nextItem()
	for t.typ == itemComment {
		t = p.lex.nextItem()
	}
	p.token[0] = t
	return p.token[0]
}
func (p *parser) backup() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	p.peekCount++
}
func (p *parser) errorf(format string, args ...interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	p.error(fmt.Errorf(format, args...))
}
func (p *parser) error(err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	perr := &ParseErr{Line: p.lex.lineNumber(), Pos: p.lex.linePosition(), Err: err}
	if strings.Count(strings.TrimSpace(p.lex.input), "\n") == 0 {
		perr.Line = 0
	}
	panic(perr)
}
func (p *parser) expect(exp ItemType, context string) item {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	token := p.next()
	if token.typ != exp {
		p.errorf("unexpected %s in %s, expected %s", token.desc(), context, exp.desc())
	}
	return token
}
func (p *parser) expectOneOf(exp1, exp2 ItemType, context string) item {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	token := p.next()
	if token.typ != exp1 && token.typ != exp2 {
		p.errorf("unexpected %s in %s, expected %s or %s", token.desc(), context, exp1.desc(), exp2.desc())
	}
	return token
}

var errUnexpected = fmt.Errorf("unexpected error")

func (p *parser) recover(errp *error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	e := recover()
	if e != nil {
		if _, ok := e.(runtime.Error); ok {
			buf := make([]byte, 64<<10)
			buf = buf[:runtime.Stack(buf, false)]
			fmt.Fprintf(os.Stderr, "parser panic: %v\n%s", e, buf)
			*errp = errUnexpected
		} else {
			*errp = e.(error)
		}
	}
	p.lex.close()
}
func (p *parser) expr() Expr {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	expr := p.unaryExpr()
	for {
		op := p.peek().typ
		if !op.isOperator() {
			if op == itemLeftBracket {
				expr = p.subqueryOrRangeSelector(expr, false)
				if s, ok := expr.(*SubqueryExpr); ok {
					if p.peek().typ == itemOffset {
						offset := p.offset()
						s.Offset = offset
					}
				}
			}
			return expr
		}
		p.next()
		vecMatching := &VectorMatching{Card: CardOneToOne}
		if op.isSetOperator() {
			vecMatching.Card = CardManyToMany
		}
		returnBool := false
		if p.peek().typ == itemBool {
			if !op.isComparisonOperator() {
				p.errorf("bool modifier can only be used on comparison operators")
			}
			p.next()
			returnBool = true
		}
		if p.peek().typ == itemOn || p.peek().typ == itemIgnoring {
			if p.peek().typ == itemOn {
				vecMatching.On = true
			}
			p.next()
			vecMatching.MatchingLabels = p.labels()
			if t := p.peek().typ; t == itemGroupLeft || t == itemGroupRight {
				p.next()
				if t == itemGroupLeft {
					vecMatching.Card = CardManyToOne
				} else {
					vecMatching.Card = CardOneToMany
				}
				if p.peek().typ == itemLeftParen {
					vecMatching.Include = p.labels()
				}
			}
		}
		for _, ln := range vecMatching.MatchingLabels {
			for _, ln2 := range vecMatching.Include {
				if ln == ln2 && vecMatching.On {
					p.errorf("label %q must not occur in ON and GROUP clause at once", ln)
				}
			}
		}
		rhs := p.unaryExpr()
		expr = p.balance(expr, op, rhs, vecMatching, returnBool)
	}
}
func (p *parser) balance(lhs Expr, op ItemType, rhs Expr, vecMatching *VectorMatching, returnBool bool) *BinaryExpr {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if lhsBE, ok := lhs.(*BinaryExpr); ok {
		precd := lhsBE.Op.precedence() - op.precedence()
		if (precd < 0) || (precd == 0 && op.isRightAssociative()) {
			balanced := p.balance(lhsBE.RHS, op, rhs, vecMatching, returnBool)
			if lhsBE.Op.isComparisonOperator() && !lhsBE.ReturnBool && balanced.Type() == ValueTypeScalar && lhsBE.LHS.Type() == ValueTypeScalar {
				p.errorf("comparisons between scalars must use BOOL modifier")
			}
			return &BinaryExpr{Op: lhsBE.Op, LHS: lhsBE.LHS, RHS: balanced, VectorMatching: lhsBE.VectorMatching, ReturnBool: lhsBE.ReturnBool}
		}
	}
	if op.isComparisonOperator() && !returnBool && rhs.Type() == ValueTypeScalar && lhs.Type() == ValueTypeScalar {
		p.errorf("comparisons between scalars must use BOOL modifier")
	}
	return &BinaryExpr{Op: op, LHS: lhs, RHS: rhs, VectorMatching: vecMatching, ReturnBool: returnBool}
}
func (p *parser) unaryExpr() Expr {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch t := p.peek(); t.typ {
	case itemADD, itemSUB:
		p.next()
		e := p.unaryExpr()
		if nl, ok := e.(*NumberLiteral); ok {
			if t.typ == itemSUB {
				nl.Val *= -1
			}
			return nl
		}
		return &UnaryExpr{Op: t.typ, Expr: e}
	case itemLeftParen:
		p.next()
		e := p.expr()
		p.expect(itemRightParen, "paren expression")
		return &ParenExpr{Expr: e}
	}
	e := p.primaryExpr()
	if p.peek().typ == itemLeftBracket {
		e = p.subqueryOrRangeSelector(e, true)
	}
	if p.peek().typ == itemOffset {
		offset := p.offset()
		switch s := e.(type) {
		case *VectorSelector:
			s.Offset = offset
		case *MatrixSelector:
			s.Offset = offset
		case *SubqueryExpr:
			s.Offset = offset
		default:
			p.errorf("offset modifier must be preceded by an instant or range selector, but follows a %T instead", e)
		}
	}
	return e
}
func (p *parser) subqueryOrRangeSelector(expr Expr, checkRange bool) Expr {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	ctx := "subquery selector"
	if checkRange {
		ctx = "range/subquery selector"
	}
	p.next()
	var erange time.Duration
	var err error
	erangeStr := p.expect(itemDuration, ctx).val
	erange, err = parseDuration(erangeStr)
	if err != nil {
		p.error(err)
	}
	var itm item
	if checkRange {
		itm = p.expectOneOf(itemRightBracket, itemColon, ctx)
		if itm.typ == itemRightBracket {
			vs, ok := expr.(*VectorSelector)
			if !ok {
				p.errorf("range specification must be preceded by a metric selector, but follows a %T instead", expr)
			}
			return &MatrixSelector{Name: vs.Name, LabelMatchers: vs.LabelMatchers, Range: erange}
		}
	} else {
		itm = p.expect(itemColon, ctx)
	}
	var estep time.Duration
	itm = p.expectOneOf(itemRightBracket, itemDuration, ctx)
	if itm.typ == itemDuration {
		estepStr := itm.val
		estep, err = parseDuration(estepStr)
		if err != nil {
			p.error(err)
		}
		p.expect(itemRightBracket, ctx)
	}
	return &SubqueryExpr{Expr: expr, Range: erange, Step: estep}
}
func (p *parser) number(val string) float64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	n, err := strconv.ParseInt(val, 0, 64)
	f := float64(n)
	if err != nil {
		f, err = strconv.ParseFloat(val, 64)
	}
	if err != nil {
		p.errorf("error parsing number: %s", err)
	}
	return f
}
func (p *parser) primaryExpr() Expr {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch t := p.next(); {
	case t.typ == itemNumber:
		f := p.number(t.val)
		return &NumberLiteral{f}
	case t.typ == itemString:
		return &StringLiteral{p.unquoteString(t.val)}
	case t.typ == itemLeftBrace:
		p.backup()
		return p.VectorSelector("")
	case t.typ == itemIdentifier:
		if p.peek().typ == itemLeftParen {
			return p.call(t.val)
		}
		fallthrough
	case t.typ == itemMetricIdentifier:
		return p.VectorSelector(t.val)
	case t.typ.isAggregator():
		p.backup()
		return p.aggrExpr()
	default:
		p.errorf("no valid expression found")
	}
	return nil
}
func (p *parser) labels() []string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	const ctx = "grouping opts"
	p.expect(itemLeftParen, ctx)
	labels := []string{}
	if p.peek().typ != itemRightParen {
		for {
			id := p.next()
			if !isLabel(id.val) {
				p.errorf("unexpected %s in %s, expected label", id.desc(), ctx)
			}
			labels = append(labels, id.val)
			if p.peek().typ != itemComma {
				break
			}
			p.next()
		}
	}
	p.expect(itemRightParen, ctx)
	return labels
}
func (p *parser) aggrExpr() *AggregateExpr {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	const ctx = "aggregation"
	agop := p.next()
	if !agop.typ.isAggregator() {
		p.errorf("expected aggregation operator but got %s", agop)
	}
	var grouping []string
	var without bool
	modifiersFirst := false
	if t := p.peek().typ; t == itemBy || t == itemWithout {
		if t == itemWithout {
			without = true
		}
		p.next()
		grouping = p.labels()
		modifiersFirst = true
	}
	p.expect(itemLeftParen, ctx)
	var param Expr
	if agop.typ.isAggregatorWithParam() {
		param = p.expr()
		p.expect(itemComma, ctx)
	}
	e := p.expr()
	p.expect(itemRightParen, ctx)
	if !modifiersFirst {
		if t := p.peek().typ; t == itemBy || t == itemWithout {
			if len(grouping) > 0 {
				p.errorf("aggregation must only contain one grouping clause")
			}
			if t == itemWithout {
				without = true
			}
			p.next()
			grouping = p.labels()
		}
	}
	return &AggregateExpr{Op: agop.typ, Expr: e, Param: param, Grouping: grouping, Without: without}
}
func (p *parser) call(name string) *Call {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	const ctx = "function call"
	fn, exist := getFunction(name)
	if !exist {
		p.errorf("unknown function with name %q", name)
	}
	p.expect(itemLeftParen, ctx)
	if p.peek().typ == itemRightParen {
		p.next()
		return &Call{fn, nil}
	}
	var args []Expr
	for {
		e := p.expr()
		args = append(args, e)
		if p.peek().typ != itemComma {
			break
		}
		p.next()
	}
	p.expect(itemRightParen, ctx)
	return &Call{Func: fn, Args: args}
}
func (p *parser) labelSet() labels.Labels {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	set := []labels.Label{}
	for _, lm := range p.labelMatchers(itemEQL) {
		set = append(set, labels.Label{Name: lm.Name, Value: lm.Value})
	}
	return labels.New(set...)
}
func (p *parser) labelMatchers(operators ...ItemType) []*labels.Matcher {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	const ctx = "label matching"
	matchers := []*labels.Matcher{}
	p.expect(itemLeftBrace, ctx)
	if p.peek().typ == itemRightBrace {
		p.next()
		return matchers
	}
	for {
		label := p.expect(itemIdentifier, ctx)
		op := p.next().typ
		if !op.isOperator() {
			p.errorf("expected label matching operator but got %s", op)
		}
		var validOp = false
		for _, allowedOp := range operators {
			if op == allowedOp {
				validOp = true
			}
		}
		if !validOp {
			p.errorf("operator must be one of %q, is %q", operators, op)
		}
		val := p.unquoteString(p.expect(itemString, ctx).val)
		var matchType labels.MatchType
		switch op {
		case itemEQL:
			matchType = labels.MatchEqual
		case itemNEQ:
			matchType = labels.MatchNotEqual
		case itemEQLRegex:
			matchType = labels.MatchRegexp
		case itemNEQRegex:
			matchType = labels.MatchNotRegexp
		default:
			p.errorf("item %q is not a metric match type", op)
		}
		m, err := labels.NewMatcher(matchType, label.val, val)
		if err != nil {
			p.error(err)
		}
		matchers = append(matchers, m)
		if p.peek().typ == itemIdentifier {
			p.errorf("missing comma before next identifier %q", p.peek().val)
		}
		if p.peek().typ != itemComma {
			break
		}
		p.next()
		if p.peek().typ == itemRightBrace {
			break
		}
	}
	p.expect(itemRightBrace, ctx)
	return matchers
}
func (p *parser) metric() labels.Labels {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	name := ""
	var m labels.Labels
	t := p.peek().typ
	if t == itemIdentifier || t == itemMetricIdentifier {
		name = p.next().val
		t = p.peek().typ
	}
	if t != itemLeftBrace && name == "" {
		p.errorf("missing metric name or metric selector")
	}
	if t == itemLeftBrace {
		m = p.labelSet()
	}
	if name != "" {
		m = append(m, labels.Label{Name: labels.MetricName, Value: name})
		sort.Sort(m)
	}
	return m
}
func (p *parser) offset() time.Duration {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	const ctx = "offset"
	p.next()
	offi := p.expect(itemDuration, ctx)
	offset, err := parseDuration(offi.val)
	if err != nil {
		p.error(err)
	}
	return offset
}
func (p *parser) VectorSelector(name string) *VectorSelector {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var matchers []*labels.Matcher
	if t := p.peek(); t.typ == itemLeftBrace {
		matchers = p.labelMatchers(itemEQL, itemNEQ, itemEQLRegex, itemNEQRegex)
	}
	if name != "" {
		for _, m := range matchers {
			if m.Name == labels.MetricName {
				p.errorf("metric name must not be set twice: %q or %q", name, m.Value)
			}
		}
		m, err := labels.NewMatcher(labels.MatchEqual, labels.MetricName, name)
		if err != nil {
			panic(err)
		}
		matchers = append(matchers, m)
	}
	if len(matchers) == 0 {
		p.errorf("vector selector must contain label matchers or metric name")
	}
	notEmpty := false
	for _, lm := range matchers {
		if !lm.Matches("") {
			notEmpty = true
			break
		}
	}
	if !notEmpty {
		p.errorf("vector selector must contain at least one non-empty matcher")
	}
	return &VectorSelector{Name: name, LabelMatchers: matchers}
}
func (p *parser) expectType(node Node, want ValueType, context string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	t := p.checkType(node)
	if t != want {
		p.errorf("expected type %s in %s, got %s", documentedType(want), context, documentedType(t))
	}
}
func (p *parser) checkType(node Node) (typ ValueType) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch n := node.(type) {
	case Expressions:
		typ = ValueTypeNone
	case Expr:
		typ = n.Type()
	default:
		p.errorf("unknown node type: %T", node)
	}
	switch n := node.(type) {
	case *EvalStmt:
		ty := p.checkType(n.Expr)
		if ty == ValueTypeNone {
			p.errorf("evaluation statement must have a valid expression type but got %s", documentedType(ty))
		}
	case Expressions:
		for _, e := range n {
			ty := p.checkType(e)
			if ty == ValueTypeNone {
				p.errorf("expression must have a valid expression type but got %s", documentedType(ty))
			}
		}
	case *AggregateExpr:
		if !n.Op.isAggregator() {
			p.errorf("aggregation operator expected in aggregation expression but got %q", n.Op)
		}
		p.expectType(n.Expr, ValueTypeVector, "aggregation expression")
		if n.Op == itemTopK || n.Op == itemBottomK || n.Op == itemQuantile {
			p.expectType(n.Param, ValueTypeScalar, "aggregation parameter")
		}
		if n.Op == itemCountValues {
			p.expectType(n.Param, ValueTypeString, "aggregation parameter")
		}
	case *BinaryExpr:
		lt := p.checkType(n.LHS)
		rt := p.checkType(n.RHS)
		if !n.Op.isOperator() {
			p.errorf("binary expression does not support operator %q", n.Op)
		}
		if (lt != ValueTypeScalar && lt != ValueTypeVector) || (rt != ValueTypeScalar && rt != ValueTypeVector) {
			p.errorf("binary expression must contain only scalar and instant vector types")
		}
		if (lt != ValueTypeVector || rt != ValueTypeVector) && n.VectorMatching != nil {
			if len(n.VectorMatching.MatchingLabels) > 0 {
				p.errorf("vector matching only allowed between instant vectors")
			}
			n.VectorMatching = nil
		} else {
			if n.Op.isSetOperator() {
				if n.VectorMatching.Card == CardOneToMany || n.VectorMatching.Card == CardManyToOne {
					p.errorf("no grouping allowed for %q operation", n.Op)
				}
				if n.VectorMatching.Card != CardManyToMany {
					p.errorf("set operations must always be many-to-many")
				}
			}
		}
		if (lt == ValueTypeScalar || rt == ValueTypeScalar) && n.Op.isSetOperator() {
			p.errorf("set operator %q not allowed in binary scalar expression", n.Op)
		}
	case *Call:
		nargs := len(n.Func.ArgTypes)
		if n.Func.Variadic == 0 {
			if nargs != len(n.Args) {
				p.errorf("expected %d argument(s) in call to %q, got %d", nargs, n.Func.Name, len(n.Args))
			}
		} else {
			na := nargs - 1
			if na > len(n.Args) {
				p.errorf("expected at least %d argument(s) in call to %q, got %d", na, n.Func.Name, len(n.Args))
			} else if nargsmax := na + n.Func.Variadic; n.Func.Variadic > 0 && nargsmax < len(n.Args) {
				p.errorf("expected at most %d argument(s) in call to %q, got %d", nargsmax, n.Func.Name, len(n.Args))
			}
		}
		for i, arg := range n.Args {
			if i >= len(n.Func.ArgTypes) {
				i = len(n.Func.ArgTypes) - 1
			}
			p.expectType(arg, n.Func.ArgTypes[i], fmt.Sprintf("call to function %q", n.Func.Name))
		}
	case *ParenExpr:
		p.checkType(n.Expr)
	case *UnaryExpr:
		if n.Op != itemADD && n.Op != itemSUB {
			p.errorf("only + and - operators allowed for unary expressions")
		}
		if t := p.checkType(n.Expr); t != ValueTypeScalar && t != ValueTypeVector {
			p.errorf("unary expression only allowed on expressions of type scalar or instant vector, got %q", documentedType(t))
		}
	case *SubqueryExpr:
		ty := p.checkType(n.Expr)
		if ty != ValueTypeVector {
			p.errorf("subquery is only allowed on instant vector, got %s in %q instead", ty, n.String())
		}
	case *NumberLiteral, *MatrixSelector, *StringLiteral, *VectorSelector:
	default:
		p.errorf("unknown node type: %T", node)
	}
	return
}
func (p *parser) unquoteString(s string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	unquoted, err := strutil.Unquote(s)
	if err != nil {
		p.errorf("error unquoting string %q: %s", s, err)
	}
	return unquoted
}
func parseDuration(ds string) (time.Duration, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	dur, err := model.ParseDuration(ds)
	if err != nil {
		return 0, err
	}
	if dur == 0 {
		return 0, fmt.Errorf("duration must be greater than 0")
	}
	return time.Duration(dur), nil
}
