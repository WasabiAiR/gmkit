package postgres

import (
	"database/sql"
	"fmt"
	"strings"
	
	"github.com/lib/pq"
	"github.com/reiver/go-pqerror"
)

// DatabaseErr provides meaningful behavior to the postgres error.
type DatabaseErr struct {
	resource string
	msg      string
	notFound bool
	err      *pq.Error
}

// Err creates a new PG error from a provided error.
func Err(resource string, err error) error {
	if err == nil {
		return nil
	}

	pg := &DatabaseErr{
		resource: resource,
	}

	if e, ok := err.(*pq.Error); ok {
		pg.err = e
		return pg
	}

	if err == sql.ErrNoRows {
		pg.notFound = true
		return pg
	}

	pg.msg = err.Error()
	return pg
}

// Error returns the error string.
func (e *DatabaseErr) Error() string {
	msgParts := []string{
		fmt.Sprintf("resource=%q", e.resource),
	}

	if msg := e.msg; msg == "" && e.notFound {
		msgParts = append(msgParts, fmt.Sprintf("%s not found", e.resource))
	} else if msg != "" {
		msgParts = append(msgParts, fmt.Sprintf("errMsg=%q", msg))
	}

	if e := e.err; e != nil {
		if e.Severity != "" {
			msgParts = append(msgParts, fmt.Sprintf("severity=%q", e.Severity))
		}

		if err := e.Message; err != "" {
			msgParts = append(msgParts, fmt.Sprintf("err=%q", err))
		}

		if code := e.Code; code != "" {
			msgParts = append(msgParts, fmt.Sprintf("code=%q", code))
		}

		if constraint := e.Constraint; constraint != "" {
			msgParts = append(msgParts, fmt.Sprintf("constraint=%q", constraint))
		}

		if column := e.Column; column != "" {
			msgParts = append(msgParts, fmt.Sprintf("column=%q", column))
		}

		if position := e.Position; position != "" {
			msgParts = append(msgParts, fmt.Sprintf("position=%q", position))
		}

		if table := e.Table; table != "" {
			msgParts = append(msgParts, fmt.Sprintf("table=%q", table))
		}

		if hint := e.Hint; hint != "" {
			msgParts = append(msgParts, fmt.Sprintf("hint=%q", hint))
		}

		if detail := e.Detail; detail != "" {
			msgParts = append(msgParts, fmt.Sprintf("detail=%q", detail))
		}

		if intQuery := e.InternalQuery; intQuery != "" {
			msgParts = append(msgParts, fmt.Sprintf("internal_query=%q", intQuery))
		}

		if dataType := e.DataTypeName; dataType != "" {
			msgParts = append(msgParts, fmt.Sprintf("data_type_name=%q", dataType))
		}

		if where := e.Where; where != "" {
			msgParts = append(msgParts, fmt.Sprintf("where=%q", where))
		}

		if schema := e.Schema; schema != "" {
			msgParts = append(msgParts, fmt.Sprintf("schema=%q", schema))
		}
	}

	return strings.Join(msgParts, " ")
}

// NotFound returns whether this error refers to the behavior of a resource that was not found.
func (e *DatabaseErr) NotFound() bool {
	if e.err == nil {
		return e.notFound
	}
	return e.notFound || e.err.Code == pqerror.CodeCaseNotFound
}

// Exists returns whether this error refers to the behavior of a resource that already exists in the
// datastore and a violation is thrown. This will be true for a unique key violation in the PG store,
// but can be expanded in the future.
func (e *DatabaseErr) Exists() bool {
	if e.err == nil {
		return false
	}
	return e.err.Code == pqerror.CodeIntegrityConstraintViolationUniqueViolation
}

// Conflict returns where this error refers to the behavior of a resource that conflicts. At this
// time the conflict is determiend from a foreign key violation, but can be expanded in the future.
func (e *DatabaseErr) Conflict() bool {
	if e.err == nil {
		return false
	}
	return e.err.Code == pqerror.CodeIntegrityConstraintViolationForeignKeyViolation
}
