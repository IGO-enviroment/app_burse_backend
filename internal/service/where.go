package service

import (
	"fmt"
	"strconv"
	"strings"
)

type FieldDB struct {
	Name  string
	Value interface{}
}

func NewFieldDB(name string, value interface{}) *FieldDB {
	return &FieldDB{
		Name:  name,
		Value: value,
	}
}

type MapperFields struct {
	fields []FieldDB
}

func NewMapperFields(fields ...FieldDB) *MapperFields {
	return &MapperFields{
		fields: fields,
	}
}

func (m *MapperFields) GetColumnNames() string {
	var columnNames []string
	for _, field := range m.fields {
		columnNames = append(columnNames, field.Name)
	}
	return fmt.Sprintf("(%s)", strings.Join(columnNames, ", "))
}

func (m *MapperFields) GetPlaceholders() string {
	var placeholders []string
	for i := 0; i < len(m.fields); i++ {
		placeholders = append(placeholders, "$"+strconv.Itoa(i+1))
	}
	return fmt.Sprintf("(%s)", strings.Join(placeholders, ", "))
}

func (m *MapperFields) GetValues() []interface{} {
	var values []interface{}
	for _, field := range m.fields {
		values = append(values, field.Value)
	}
	return values
}
