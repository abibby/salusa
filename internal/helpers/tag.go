package helpers

import (
	"reflect"
	"strconv"
	"strings"
)

type Tag struct {
	Name          string
	Primary       bool
	Nullable      bool
	AutoIncrement bool
	Readonly      bool
	Index         bool
	Unique        bool
	Type          string
	Size          int
}

func DBTag(f reflect.StructField) *Tag {
	dbTag, ok := f.Tag.Lookup("db")
	if !ok {
		return &Tag{
			Name: f.Name,
		}
	}
	parts := strings.Split(dbTag, ",")
	tag := &Tag{
		Name: parts[0],
	}

	tagValue := reflect.ValueOf(tag).Elem()
	tagType := tagValue.Type()
	for _, p := range parts[1:] {
		typePrefix := "type:"
		if strings.HasPrefix(p, typePrefix) {
			tag.Type = p[len(typePrefix):]
			continue
		}
		sizePrefix := "size:"
		if strings.HasPrefix(p, sizePrefix) {
			s, err := strconv.Atoi(p[len(sizePrefix):])
			if err == nil {
				tag.Size = s
			}
			continue
		}
		for i := 0; i < tagType.NumField(); i++ {
			f := tagType.Field(i)
			if strings.ToLower(f.Name) == p && f.Type.Kind() == reflect.Bool {
				tagValue.Field(i).SetBool(true)
			}
		}
	}

	return tag
}
