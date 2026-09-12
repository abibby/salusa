package dialects

import "gosalusa.com/extra/sets"

type DataType struct {
	Name string
	Size int
}

var (
	DataTypeBlob   = DataType{Name: "blob"}
	DataTypeString = DataType{Name: "string"}
	DataTypeText   = DataType{Name: "text"}
	DataTypeEnum   = DataType{Name: "enum"}

	DataTypeBoolean = DataType{Name: "bool"}

	DataTypeDate     = DataType{Name: "date"}
	DataTypeDateTime = DataType{Name: "date-time"}

	DataTypeFloat32 = DataType{Name: "float32"}
	DataTypeFloat64 = DataType{Name: "float64"}

	DataTypeInt8  = DataType{Name: "int8"}
	DataTypeInt16 = DataType{Name: "int16"}
	DataTypeInt32 = DataType{Name: "int32"}
	DataTypeInt64 = DataType{Name: "int64"}

	DataTypeUInt8  = DataType{Name: "uint8"}
	DataTypeUInt16 = DataType{Name: "uint16"}
	DataTypeUInt32 = DataType{Name: "uint32"}
	DataTypeUInt64 = DataType{Name: "uint64"}

	DataTypeJSON = DataType{Name: "json"}
)

var dataTypes = sets.New(
	DataTypeBlob.Name,
	DataTypeString.Name,
	DataTypeText.Name,
	DataTypeEnum.Name,

	DataTypeBoolean.Name,

	DataTypeDate.Name,
	DataTypeDateTime.Name,

	DataTypeFloat32.Name,
	DataTypeFloat64.Name,

	DataTypeInt8.Name,
	DataTypeInt16.Name,
	DataTypeInt32.Name,
	DataTypeInt64.Name,

	DataTypeUInt8.Name,
	DataTypeUInt16.Name,
	DataTypeUInt32.Name,
	DataTypeUInt64.Name,

	DataTypeJSON.Name,
)

func (d DataType) IsValid() bool {
	return dataTypes.Has(d.Name)
}

// DataTyper must not be implemented on an interface
type DataTyper interface {
	DataType() DataType
}
