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
)

var relationshipInterface = reflect.TypeOf((*relationship.Relationship)(nil)).Elem()

func appendColumnsAndValues(v reflect.Value, m map[string]any) {
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		if field.Anonymous && field.Type.Kind() == reflect.Struct {
			appendColumnsAndValues(v.Field(i), m)
		} else {
			tag := helpers.DBTag(field)
			if tag.Name == "-" || tag.Readonly || field.Type.Implements(relationshipInterface) {
				continue
			}
			m[tag.Name] = v.Field(i).Interface()
		}
	}
}

func columnsAndValues(v reflect.Value) map[string]any {
	m := map[string]any{}
	appendColumnsAndValues(v, m)
	return m
}

func MustSave(tx database.DB, v Model) {
	err := Save(tx, v)
	if err != nil {
		panic(err)
	}
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
func MustSaveContext(ctx context.Context, tx database.DB, v Model) {
	err := SaveContext(ctx, tx, v)
	if err != nil {
		panic(err)
	}
}
func SaveContext(ctx context.Context, tx database.DB, v Model) error {
	inDB := v.InDatabase()
	err := hooks.BeforeSave(ctx, tx, v)
	if err != nil {
		return fmt.Errorf("before save hooks: %w", err)
	}

	d := dialects.New()
	m := columnsAndValues(reflect.ValueOf(v).Elem())
	if inDB {
		err = update(ctx, tx, d, v, m)
		if err != nil {
			return fmt.Errorf("update: %w", err)
		}
	} else {
		err = insert(ctx, tx, d, v, m)
		if err != nil {
			return fmt.Errorf("insert: %w", err)
		}
	}

	err = relationship.InitializeRelationships(v)
	if err != nil {
		return fmt.Errorf("initialize relationships: %w", err)
	}

	err = hooks.AfterSave(ctx, tx, v)
	if err != nil {
		return fmt.Errorf("after save hooks: %w", err)
	}
	return nil
}

func insert(ctx context.Context, tx database.DB, d dialects.Dialect, v any, m map[string]any) error {
	rPKey, pKey, isAuto := isAutoIncrementing(v)
	if isAuto {
		delete(m, pKey)
	}
	sql, err := d.EncodeInsertQuery(&dialects.InsertQuery{
		Table:  database.GetTable(v),
		Values: []map[string]any{m},
	})
	if err != nil {
		return err
	}

	result, err := tx.ExecContext(ctx, sql.Query, sql.Bindings...)
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

func update(ctx context.Context, tx database.DB, d dialects.Dialect, v any, m map[string]any) error {
	pKey := helpers.PrimaryKey(v)

	wheres := []dialects.Condition{}
	for _, k := range pKey {
		pKeyValue, ok := helpers.GetValue(v, k)
		if !ok {
			return fmt.Errorf("no primary key found")
		}
		wheres = append(wheres, dialects.Condition{
			Column:   dialects.Column{Column: k},
			Operator: "=",
			Value:    pKeyValue,
		})
	}

	result, err := d.EncodeUpdateQuery(&dialects.UpdateQuery{
		Table:  database.GetTable(v),
		Values: m,
		Wheres: wheres,
	})
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, result.Query, result.Bindings...)
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
	maps := make([]map[string]any, len(models))
	for i, v := range models {
		maps[i] = columnsAndValues(reflect.ValueOf(v).Elem())
	}
	err := insertMany(ctx, tx, d, models[0], maps)
	if err != nil {
		return fmt.Errorf("insert: %w", err)
	}
	for _, v := range models {
		// err = relationship.InitializeRelationships(v)
		// if err != nil {
		// 	return fmt.Errorf("initialize relationships: %w", err)
		// }
		err := hooks.AfterSave(ctx, tx, v)
		if err != nil {
			return fmt.Errorf("before save hooks: %w", err)
		}
	}
	return nil
}

func insertMany(ctx context.Context, tx database.DB, d dialects.Dialect, v any, maps []map[string]any) error {

	rPKey, pKey, isAuto := isAutoIncrementing(v)
	if isAuto {
		for i, _ := range maps {
			delete(maps[i], pKey)
		}
	}
	sql, err := dialects.New().EncodeInsertQuery(&dialects.InsertQuery{
		Table:  database.GetTable(v),
		Values: maps,
	})
	if err != nil {
		return err
	}

	result, err := tx.ExecContext(ctx, sql.Query, sql.Bindings...)
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
