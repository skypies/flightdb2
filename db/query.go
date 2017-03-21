package db

import "fmt"

// Query is a thin skin over the datastore query API. It provides for a textual dump of the
// query, and also to paper over the two providers (cloud datastore, appengine datastore)
type Query struct {
	Kind      string
	Filters []Filter
	OrderStr   string
	LimitVal   int
}

type Filter struct {
	Field string
	Value interface{}
}

func (q *Query)String() string {
	str := fmt.Sprintf("NewQuery(%q)\n", q.Kind)
	for _,f := range q.Filters {
		str += fmt.Sprintf("  .Filter(%q, %v)\n", f.Field, f.Value)
	}
	if q.OrderStr != "" { str += fmt.Sprintf("  .Order(%q)\n", q.OrderStr) }
	if q.LimitVal != 0 { str += fmt.Sprintf("  .Limit(%d)\n", q.LimitVal) }
	return str
}

func NewQuery(kind string) *Query { return &Query{Kind:kind} }

func (q *Query)Filter(field string, val interface{}) *Query {
	q.Filters = append(q.Filters, Filter{field, val})
	return q
}

func (q *Query)Order(o string) *Query {
	q.OrderStr = o
	return q
}

func (q *Query)Limit(l int) *Query {
	q.LimitVal = l
	return q
}
