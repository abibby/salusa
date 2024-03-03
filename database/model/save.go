package model

import (
	"context"
	"fmt"
	"reflect"

	"github.com/abibby/salusa/database"
	"github.com/abibby/salusa/database/dialects"
	"github.com/abibby/salusa/database/hooks"
	"github.com/abibby/salusa/internal/helpers"
	"github.com/abibby/salusa/internal/relationship"
	"github.com/abibby/salusa/slices"
)

var relationshipInterface = reflect.TypeOf((*relationship.Relationship)(nil)).Elem()

func columnsAndValues(v reflect.Value) ([]string, []any) {
	t := v.Type()
	numFields := t.NumField()
	columns := make([]string, 0, numFields)
	values := make([]any, 0, numFields)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		if field.Anonymous {
			subColumns, subValues := columnsAndValues(v.Field(i))
			columns = append(columns, subColumns...)
			values = append(values, subValues...)
		} else {
			tag := helpers.DBTag(field)
			if tag.Name == "-" || tag.Readonly || field.Type.Implements(relationshipInterface) {
				continue
			}
			columns = append(columns, tag.Name)
			values = append(values, v.Field(i).Interface())
		}
	}
	return columns, values
}

func Save(tx database.DB, v Model) error {
	ctx := context.Background()
	if v, ok := v.(Contexter); ok {
		modelCtx := v.Context()
		if modelCtx != nil {
			ctx = modelCtx
		}
	}
	return SaveContext(ctx, tx, v)
}
func SaveContext(ctx context.Context, tx database.DB, v Model) error {
	err := hooks.BeforeSave(ctx, tx, v)
	if err != nil {
		return fmt.Errorf("before save hooks: %w", err)
	}

	d := dialects.New()
	columns, values := columnsAndValues(reflect.ValueOf(v).Elem())
	if v.InDatabase() {
		err = update(ctx, tx, d, v, columns, values)
		if err != nil {
			return fmt.Errorf("update: %w", err)
		}
	} else {
		err = insert(ctx, tx, d, v, columns, values)
		if err != nil {
			return fmt.Errorf("insert: %w", err)
		}
	}

	if err != nil {
		return fmt.Errorf("initialize relationships: %w", err)
	}
	err = hooks.AfterSave(ctx, tx, v)
	if err != nil {
		return fmt.Errorf("after save hooks: %w", err)
	}
	return nil
}

func insert(ctx context.Context, tx database.DB, d dialects.Dialect, v any, columns []string, values []any) error {
	rPKey, pKey, isAuto := isAutoIncrementing(v)
	if isAuto {
		newColumns := make([]string, 0, len(columns))
		newValues := make([]any, 0, len(values))
		for i, column := range columns {
			if column != pKey {
				newColumns = append(newColumns, column)
				newValues = append(newValues, values[i])
			}
		}
		columns = newColumns
		values = newValues
	}
	r := helpers.Result().
		AddString("INSERT INTO").
		Add(helpers.Identifier(database.GetTable(v))).
		Add(
			helpers.Group(
				helpers.Join(
					helpers.IdentifierList(columns),
					", ",
				),
			),
		).
		AddString("VALUES").
		Add(
			helpers.Group(
				helpers.Join(
					helpers.LiteralList(values),
					", ",
				),
			),
		)

	q, bindings, err := r.ToSQL(d)
	if err != nil {
		return fmt.Errorf("failed to generate sql: %w", err)
	}

	result, err := tx.ExecContext(ctx, q, bindings...)
	if err != nil {
		return fmt.Errorf("failed to insert model: %w", err)
	}

	if isAuto {
		id, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("could not get last insert id: %w", err)
		}
		rPKey.SetInt(id)
	}
	return nil
}

func isAutoIncrementing(v any) (reflect.Value, string, bool) {
	pKeys := helpers.PrimaryKey(v)
	if len(pKeys) != 1 {
		return reflect.Value{}, "", false
	}

	pKey := pKeys[0]
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return reflect.Value{}, "", false
	}
	var pKeyTag *helpers.Tag
	var rPKey reflect.Value
	errFound := fmt.Errorf("found")
	err := helpers.EachField(rv, func(sf reflect.StructField, fv reflect.Value) error {
		tag := helpers.DBTag(sf)
		if tag.Name == pKey {
			pKeyTag = tag
			rPKey = fv
			if !rPKey.IsZero() {
				return nil
			}
			return errFound
		}
		return nil
	})
	if err != errFound {
		return reflect.Value{}, "", false
	}
	if pKeyTag != nil && !pKeyTag.AutoIncrement {
		return reflect.Value{}, "", false
	}
	return rPKey, pKey, true
}

func update(ctx context.Context, tx database.DB, d dialects.Dialect, v any, columns []string, values []any) error {
	pKey := helpers.PrimaryKey(v)
	r := helpers.Result().
		AddString("UPDATE").
		Add(helpers.Identifier(database.GetTable(v))).
		AddString("SET")

	for i, column := range columns {
		if i != 0 {
			r.AddString(",")
		}
		r.Add(helpers.Identifier(column))
		r.AddString("=")
		r.Add(helpers.Literal(values[i]))
	}

	r.AddString("WHERE")

	for i, k := range pKey {
		pKeyValue, ok := helpers.GetValue(v, k)
		if !ok {
			return fmt.Errorf("no primary key found")
		}

		if i != 0 {
			r.AddString("AND")
		}

		r.Add(helpers.Identifier(k)).
			AddString("=").
			Add(helpers.Literal(pKeyValue))
	}

	q, bindings, err := r.ToSQL(d)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, q, bindings...)
	if err != nil {
		return err
	}
	return nil
}

func InsertMany[T Model](tx database.DB, models []T) error {
	return InsertManyContext(context.Background(), tx, models)
}
func InsertManyContext[T Model](ctx context.Context, tx database.DB, models []T) error {
	for _, v := range models {
		err := hooks.BeforeSave(ctx, tx, v)
		if err != nil {
			return fmt.Errorf("before save hooks: %w", err)
		}
	}

	d := dialects.New()
	var columns []string
	values := make([][]any, len(models))
	for i, v := range models {

		c, v := columnsAndValues(reflect.ValueOf(v).Elem())
		if columns == nil {
			columns = c
		}
		values[i] = v
	}
	err := insertMany(ctx, tx, d, models[0], columns, values)
	if err != nil {
		return fmt.Errorf("insert: %w", err)
	}
	for _, v := range models {
		if err != nil {
			return fmt.Errorf("initialize relationships: %w", err)
		}
		err := hooks.AfterSave(ctx, tx, v)
		if err != nil {
			return fmt.Errorf("before save hooks: %w", err)
		}
	}
	return nil
}

func insertMany(ctx context.Context, tx database.DB, d dialects.Dialect, v any, columns []string, values [][]any) error {
	_, pKey, isAuto := isAutoIncrementing(v)
	pKeyIndex := -1
	if isAuto {
		newColumns := make([]string, 0, len(columns))
		for i, column := range columns {
			if column != pKey {
				newColumns = append(newColumns, column)
			} else {
				pKeyIndex = i
			}
		}
		columns = newColumns
	}
	r := helpers.Result().
		AddString("INSERT INTO").
		Add(helpers.Identifier(database.GetTable(v))).
		Add(
			helpers.Group(
				helpers.Join(
					helpers.IdentifierList(columns),
					", ",
				),
			),
		).
		AddString("VALUES").
		Add(
			helpers.Join(
				slices.Map(values, func(v []any) helpers.ToSQLer {
					newValues := make([]any, 0, len(columns))
					for i, val := range v {
						if i != pKeyIndex {
							newValues = append(newValues, val)
						}
					}
					return helpers.Group(
						helpers.Join(
							helpers.LiteralList(newValues),
							", ",
						),
					)
				}),
				", ",
			),
		)

	q, bindings, err := r.ToSQL(d)
	if err != nil {
		return fmt.Errorf("failed to generate sql: %w", err)
	}

	_, err = tx.ExecContext(ctx, q, bindings...)
	if err != nil {
		return fmt.Errorf("failed to insert model: %w", err)
	}

	return nil
}
