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

// ByFieldEquals creates a `{field} = ?` argument for a Where clause.
func ByFieldEquals(field string, value interface{}) Whereable {
	return Equals{
		Key: field,
		Val: value,
	}
}

// notEquals is a struct used to represent a part of a where clause
type notEquals struct {
	Key string
	Val interface{}
}

// Clause prints out a sql ready statement for the Rebinder to repackage.
func (ne notEquals) Clause(rb Rebinder) (string, []interface{}, error) {
	query := fmt.Sprintf("%s != %s", ne.Key, rb.Rebind("?"))
	return query, []interface{}{ne.Val}, nil
}

// ByFieldNotEquals creates a `{field} = ?` argument for a Where clause.
func ByFieldNotEquals(field string, value interface{}) Whereable {
	return notEquals{
		Key: field,
		Val: value,
	}
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

type group struct {
	op     string
	wheres []Whereable
}

// Rebind makes group into a Rebinder to call the Clause() method on other
// Whereables and maintain order abiguity (keep the ?).
func (group) Rebind(q string) string {
	return sqlx.Rebind(sqlx.QUESTION, q)
}

func (g group) Clause(rb Rebinder) (string, []interface{}, error) {
	var qParts []string
	var args []interface{}
	for _, w := range g.wheres {
		if w == nil {
			continue
		}
		q, a, err := w.Clause(g)
		if err != nil {
			return "", nil, err
		}
		qParts = append(qParts, q)
		args = append(args, a...)
	}
	query := rb.Rebind("(" + strings.Join(qParts, " "+g.op+" ") + ")")

	return query, args, nil
}

// And joins multiple Whereable clauses with AND operations.
func And(first Whereable, second Whereable, rest ...Whereable) Whereable {
	wheres := []Whereable{first, second}
	if len(rest) > 0 {
		wheres = append(wheres, rest...)
	}
	return group{
		op:     "AND",
		wheres: wheres,
	}
}

// Or joins multiple Whereable clauses with OR operations.
func Or(first Whereable, second Whereable, rest ...Whereable) Whereable {
	wheres := []Whereable{first, second}
	if len(rest) > 0 {
		wheres = append(wheres, rest...)
	}
	return group{
		op:     "OR",
		wheres: wheres,
	}
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

// ByFieldIn creates a `{field} = (?, ?, ...)` argument for a Where clause.
func ByFieldIn(field string, values ...interface{}) Whereable {
	return In{
		Key:  field,
		Vals: values,
	}
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

// QueryOptionTypeLimit describes a limit option type
const QueryOptionTypeLimit = "limit"

// QueryOptionTypeOffset describes an offset option type
const QueryOptionTypeOffset = "offset"

// QueryOptionTypeOrderBy describes an order by option type
const QueryOptionTypeOrderBy = "orderBy"

// QueryOption is an interface around query options that can be of and order By, Limit or Offset.
type QueryOption interface {
	Clause() string
	Type() string
	order() int
	valid() bool
}

type limit struct {
	val int64
}

// Clause prints out a sql ready statement for A LIMIT.
func (l limit) Clause() string {
	if l.val == 0 {
		return ""
	}
	return fmt.Sprintf("LIMIT %d", l.val)
}

func (l limit) Type() string {
	return QueryOptionTypeLimit
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
	if o.val == 0 {
		return ""
	}
	return fmt.Sprintf("OFFSET %d", o.val)
}

func (o offset) Type() string {
	return QueryOptionTypeOffset
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

func (o orderBy) Type() string {
	return QueryOptionTypeOrderBy
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
