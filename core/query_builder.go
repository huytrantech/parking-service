package core

import (
	"fmt"
	"parking-service/repository"
	"strings"
)

type UpdateOption struct {
	Filter         map[string]interface{}
	Skip           int
	Limit          int
	Updated        map[string]interface{}
	finalQuery     string
	Table          string
	finalDataParam []interface{}
}

func (o *UpdateOption) BuildQuery() {
	query := fmt.Sprintf("update $table$ set $update$ where $condition$")

	if len(o.Table) > 0 {
		query = strings.ReplaceAll(query, "$table$", o.Table)
	}
	count := 1
	if len(o.Updated) > 0 {
		filterArr := make([]string, 0)
		for key, value := range o.Updated {
			filterArr = append(filterArr, fmt.Sprintf("%s = $%d", key, count))
			count += 1
			o.finalDataParam = append(o.finalDataParam, value)
		}
		query = strings.ReplaceAll(query, "$update$", strings.Join(filterArr, " , "))
	}

	if len(o.Filter) > 0 {
		filterArr := make([]string, 0)
		for key, value := range o.Filter {
			filterArr = append(filterArr, fmt.Sprintf("%s = $%d", key, count))
			count += 1
			o.finalDataParam = append(o.finalDataParam, value)
		}
		query = strings.ReplaceAll(query, "$condition$", strings.Join(filterArr, " and "))
	}

	o.finalQuery = query
}

func (o *UpdateOption) GetFinalQuery() (string, []interface{}) {
	return o.finalQuery, o.finalDataParam
}

type MathOperationFilter struct {
	Operation MathOperationEnum
	Value     interface{}
}

type MathOperationEnum struct {
	Operation string
}

func Gt() MathOperationEnum {
	return MathOperationEnum{Operation: ">"}
}
func Gte() MathOperationEnum {
	return MathOperationEnum{Operation: ">="}
}
func Lt() MathOperationEnum {
	return MathOperationEnum{Operation: "<"}
}
func Lte() MathOperationEnum {
	return MathOperationEnum{Operation: "<="}
}
func Equal() MathOperationEnum {
	return MathOperationEnum{Operation: "="}
}
func NotEqual() MathOperationEnum {
	return MathOperationEnum{Operation: "!="}
}

type SelectOption struct {
	Projection     []string
	Skip           int
	Limit          int
	Sort           map[string]int
	Query          string
	FilterAdvance  map[string]MathOperationFilter
	Filter         map[string]interface{}
	finalQuery     string
	Table          string
	finalDataParam []interface{}
	IsCount        bool
}

func (o *SelectOption) GetFinalQuery() (string, []interface{}) {
	return o.finalQuery, o.finalDataParam
}

func GetFieldsByTable(table string) string {
	switch table {
	case "parking":
		return strings.Join(repository.ParkingRepositoryCol , ",")
	default:
		return "*"
	}
}
func (o *SelectOption) BuildQuery() {
	query := fmt.Sprintf("select $field$ from $table$ $condition$ $orderby$ $limit$ $offset$")
	if o.IsCount {
		query = strings.ReplaceAll(query, "$field$", "count(*)")
	} else {
		if len(o.Projection) > 0 {
			query = strings.ReplaceAll(query, "$field$", strings.Join(o.Projection, ","))
		} else {
			query = strings.ReplaceAll(query, "$field$", GetFieldsByTable(o.Table))
		}
	}

	if len(o.Table) > 0 {
		query = strings.ReplaceAll(query, "$table$", o.Table)
	}

	if len(o.Filter) > 0 {
		count := 1
		filterArr := make([]string, 0)
		for key, value := range o.Filter {
			filterArr = append(filterArr, fmt.Sprintf("%s = $%d", key, count))
			count += 1
			o.finalDataParam = append(o.finalDataParam, value)
		}
		query = strings.ReplaceAll(query, "$condition$", " where "+strings.Join(filterArr, " and "))
	} else if len(o.FilterAdvance) > 0 {
		count := 1
		filterArr := make([]string, 0)
		for keyAdvance, valueAdvance := range o.FilterAdvance {
			operation := ""
			if len(valueAdvance.Operation.Operation) == 0 {
				operation = Equal().Operation
			} else {
				operation = valueAdvance.Operation.Operation
			}
			filterArr = append(filterArr, fmt.Sprintf("%s %s $%d", keyAdvance, operation, count))
			count += 1
			o.finalDataParam = append(o.finalDataParam, valueAdvance.Value)
		}
		query = strings.ReplaceAll(query, "$condition$", " where "+strings.Join(filterArr, " and "))
	} else {
		query = strings.ReplaceAll(query, "$condition$", "")
	}

	if len(o.Sort) > 0 && o.IsCount == false{
		sortArr := make([]string, 0)
		for key, value := range o.Sort {
			if value == - 1 {
				sortArr = append(sortArr, fmt.Sprintf("%s desc", key))
			} else if value == 1 {
				sortArr = append(sortArr, fmt.Sprintf("%s asc", key))
			}
		}

		query = strings.ReplaceAll(query, "$orderby$", " order by "+strings.Join(sortArr, ","))
	} else {
		query = strings.ReplaceAll(query, "$orderby$", "")
	}

	if o.Limit > 0 {
		query = strings.ReplaceAll(query, "$limit$", fmt.Sprintf(" limit %d", o.Limit))
	} else {
		query = strings.ReplaceAll(query, "$limit$", "")
	}

	if o.Skip > 0 {
		query = strings.ReplaceAll(query, "$offset$", fmt.Sprintf(" offset %d", o.Skip))
	} else {
		query = strings.ReplaceAll(query, "$offset$", "")
	}

	o.finalQuery = query
}
