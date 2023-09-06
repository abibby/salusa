package hooks

import (
	"context"
	"reflect"

	"github.com/abibby/salusa/database/internal/helpers"
)

type BeforeSaver interface {
	BeforeSave(ctx context.Context, tx helpers.QueryExecer) error
}
type AfterSaver interface {
	AfterSave(ctx context.Context, tx helpers.QueryExecer) error
}

type AfterLoader interface {
	AfterLoad(ctx context.Context, tx helpers.QueryExecer) error
}

func BeforeSave(ctx context.Context, tx helpers.QueryExecer, model interface{}) error {
	if model, ok := model.(BeforeSaver); ok {
		err := model.BeforeSave(ctx, tx)
		if err != nil {
			return err
		}
	}
	return eachField(reflect.ValueOf(model), func(i interface{}) error {
		return BeforeSave(ctx, tx, i)
	})
}

func AfterSave(ctx context.Context, tx helpers.QueryExecer, model interface{}) error {
	if model, ok := model.(AfterSaver); ok {
		err := model.AfterSave(ctx, tx)
		if err != nil {
			return err
		}
	}
	return eachField(reflect.ValueOf(model), func(i interface{}) error {
		return AfterSave(ctx, tx, i)
	})
}

func AfterLoad(ctx context.Context, tx helpers.QueryExecer, model interface{}) error {
	if model, ok := model.(AfterLoader); ok {
		err := model.AfterLoad(ctx, tx)
		if err != nil {
			return err
		}
	}

	return eachField(reflect.ValueOf(model), func(i interface{}) error {
		return AfterLoad(ctx, tx, i)
	})
}

func eachField(v reflect.Value, callback func(model interface{}) error) error {
	if v.Kind() == reflect.Ptr {
		return eachField(v.Elem(), callback)
	}
	if v.Kind() == reflect.Struct {
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			if t.Field(i).Anonymous {
				f := v.Field(i)
				if f.Kind() != reflect.Ptr {
					f = f.Addr()
				}
				err := callback(f.Interface())
				if err != nil {
					return err
				}
			}
		}
		return nil
	}
	if v.Kind() == reflect.Slice {
		for i := 0; i < v.Len(); i++ {
			err := callback(v.Index(i).Interface())
			if err != nil {
				return err
			}
		}
		return nil
	}
	return nil
}
