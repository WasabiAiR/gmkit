package postgres

import (
	"fmt"
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"
)

// Whereable is a way of creating a functional query statement for different lookups.
type Whereable interface {
	Clause(rebinder Rebinder) (query string, args []interface{}, err error)
}

// All provides a nil that indicates no Where clause is provided.
func All() Whereable {
	return nil
}

// Equals is a struct used to represent a part of a where clause
type Equals struct {
	Key string
	Val interface{}
}

// Clause prints out a sql ready statement for the Rebinder to repackage.
func (e Equals) Clause(rb Rebinder) (string, []interface{}, error) {
	query := fmt.Sprintf("%s = %s", e.Key, rb.Rebind("?"))
	return query, []interface{}{e.Val}, nil
}

// ToWhereClause adds WHERE prefix to a whereable's clause output if the whereable is not nil
func ToWhereClause(rebinder Rebinder, wherable Whereable) (clause string, args []interface{}, err error) {
	if wherable == nil {
		return "", nil, nil
	}
	whereClause, whereArgs, err := wherable.Clause(rebinder)
	if err != nil {
		return "", nil, err
	}
	if whereClause != "" {
		whereClause = fmt.Sprintf("WHERE %s", whereClause)
	}

	return whereClause, whereArgs, nil
}

// ByItemID creates a `id = ?` argument for a Where clause.
func ByItemID(id interface{}) Whereable {
	return Equals{
		Key: "item_id",
		Val: id,
	}
}

// ByID creates a `id = ?` argument for a Where clause.
func ByID(id interface{}) Whereable {
	return Equals{
		Key: "id",
		Val: id,
	}
}

// ByLocationID creates a `location_id = ?` arugment for a Where clause.
func ByLocationID(id interface{}) Whereable {
	return Equals{
		Key: "location_id",
		Val: id,
	}
}

type is struct {
	key string
	val interface{}
}

func (i is) Clause(rb Rebinder) (string, []interface{}, error) {
	query := fmt.Sprintf("%s IS %s", i.key, i.val)
	return query, nil, nil
}

// IDsNotNull returns all any non Null row by id.
func IDsNotNull() Whereable {
	return is{
		key: "id",
		val: "NOT NULL",
	}
}

// In is a struct used to represent an IN clause
type In struct {
	Key  string
	Vals []interface{}
}

// Clause prints out a sql ready statement for the Rebinder to repackage.
func (i In) Clause(rb Rebinder) (string, []interface{}, error) {
	query := fmt.Sprintf("%s IN (?)", i.Key)
	query, args, err := sqlx.In(query, i.Vals)
	if err != nil {
		return "", nil, err
	}
	return rb.Rebind(query), args, nil
}

// ByIDsIn provides a where clause with the IN syntax, matching many rows potentially.
func ByIDsIn(ids ...string) Whereable {
	var vals []interface{}
	for _, id := range ids {
		vals = append(vals, id)
	}
	return In{
		Key:  "id",
		Vals: vals,
	}
}

// ByItemIDsIn provides a where clause with the IN syntax, matching many rows potentially.
func ByItemIDsIn(ids ...string) Whereable {
	var vals []interface{}
	for _, id := range ids {
		vals = append(vals, id)
	}
	return In{
		Key:  "item_id",
		Vals: vals,
	}
}

type like struct {
	key string
	val interface{}
}

func (l like) Clause(rb Rebinder) (string, []interface{}, error) {
	query := fmt.Sprintf("%s LIKE %s", l.key, rb.Rebind("?"))
	return query, []interface{}{l.val}, nil
}

// Like creates a where clause that uses a like matcher.
func Like(column string, val string) Whereable {
	return like{
		key: column,
		val: val,
	}
}

// QueryOptions provides a means to turn a collection of Query Options into a sql safe clause.
type QueryOptions []QueryOption

// Clause prints out a sql ready statement for the Rebinder to repackage.
func (q QueryOptions) Clause() string {
	sort.Slice(q, func(i, j int) bool {
		return q[i].order() < q[j].order()
	})

	var opts []string
	for _, v := range q {
		if !v.valid() {
			continue
		}
		opts = append(opts, v.Clause())
	}
	return strings.Join(opts, " ")
}

// QueryOption is an interface around query options that can be of and order By, Limit or Offset.
type QueryOption interface {
	Clause() string
	order() int
	valid() bool
}

type limit struct {
	val int64
}

// Clause prints out a sql ready statement for A LIMIT.
func (l limit) Clause() string {
	return fmt.Sprintf("LIMIT %d", l.val)
}

func (l limit) order() int {
	return 5
}

func (l limit) valid() bool {
	return l.val > 0
}

// Limit returns a QueryOption that applies a LIMIT to the query.
func Limit(lim int64) QueryOption {
	return limit{
		val: lim,
	}
}

type offset struct {
	val int64
}

// Clause prints out a sql ready statement for an OFFSET.
func (o offset) Clause() string {
	return fmt.Sprintf("OFFSET %d", o.val)
}

func (o offset) order() int {
	return 6
}

func (o offset) valid() bool {
	return o.val > 0
}

// Offset returns a QueryOption that applies an OFFSET to the query.
func Offset(o int64) QueryOption {
	return offset{
		val: o,
	}
}

type orderBy struct {
	vals []string
}

// Clause prints out a sql ready statement for an ORDER BY.
func (o orderBy) Clause() string {
	return fmt.Sprintf("ORDER BY %s", strings.Join(o.vals, ", "))
}

func (o orderBy) order() int {
	return 4
}

func (o orderBy) valid() bool {
	return len(o.vals) > 0
}

// OrderBy returns a QueryOption that applies an ORDER BY to the query.
func OrderBy(cols ...OrderPair) QueryOption {
	var ss []string
	for _, v := range cols {
		if !v.valid() {
			continue
		}
		ss = append(ss, v.format())
	}
	return orderBy{
		vals: ss,
	}
}

// OrderPair is a type to declare a column to be ordered by and if it should be ASC or DESC.
type OrderPair struct {
	Col  string
	Desc bool
}

func (o OrderPair) format() string {
	if o.Desc {
		return fmt.Sprintf("%s DESC", o.Col)
	}

	return fmt.Sprintf("%s ASC", o.Col)
}

func (o OrderPair) valid() bool {
	return o.Col != ""
}
