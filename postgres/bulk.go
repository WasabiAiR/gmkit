package postgres

import (
	"context"
)

// Bulker is an interface that defines the behavior a type needs to
// implement to be bulk insert/updated into PG.
type Bulker interface {
	Len() int
	PrepareStatement() string
	KeyedArgsAtIndex(index int) (key string, arguments []any)
	TableName() string
}

// BulkInsert is a general bulk insert func that can be used to insert any valid Bulker type.
func BulkInsert(ctx context.Context, tx CRUD, itemsToInsert Bulker) (rowsInserted int, insertErr error) {
	if itemsToInsert.Len() == 0 {
		return 0, nil
	}

	stmt, err := tx.PrepareContext(ctx, itemsToInsert.PrepareStatement())
	if err != nil {
		return 0, Err(itemsToInsert.TableName(), err)
	}

	idsMap := make(map[string]bool)
	for i := 0; i < itemsToInsert.Len(); i++ {
		key, args := itemsToInsert.KeyedArgsAtIndex(i)
		if key != "" && idsMap[key] {
			continue
		}
		idsMap[key] = true
		_, err := stmt.ExecContext(ctx, args...)
		if err != nil {
			return 0, Err(itemsToInsert.TableName(), err)
		}
	}
	idsMap = nil

	if _, err := stmt.ExecContext(ctx); err != nil {
		return 0, Err(itemsToInsert.TableName(), err)
	}

	if err := stmt.Close(); err != nil {
		return 0, Err(itemsToInsert.TableName(), err)
	}

	return itemsToInsert.Len(), nil
}
