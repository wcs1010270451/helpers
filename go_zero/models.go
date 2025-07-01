package go_zero

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

const dbTag = "select"

var validColumnName = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*(\.[a-zA-Z_][a-zA-Z0-9_]*)?$`)

// isSafeColumn 检查列名是否有效
// 该函数使用正则表达式匹配输入的列名字符串，以确保列名符合预期的格式
// 参数:
//
//	col: 待检查的列名字符串
//
// 返回值:
//
//	bool: 如果列名有效则返回true，否则返回false
func isSafeColumn(col string) bool {
	return validColumnName.MatchString(col)
}

// escapeColumn 用于转义 SQL 查询中的列名。
// 当字段名包含"."时，表明它是带有表别名的列名，此时会分别转义表别名和列名；
// 否则，将直接转义整个字段名。
// 这对于防止 SQL 注入攻击和处理特殊字符非常有用。
//
// 参数:
//
//	field: 待转义的字段名或表别名.字段名。
//
// 返回值:
//
//	转义后的列名字符串。
func escapeColumn(field string) string {
	// 将字段名按"."分割，以处理可能的表别名情况。
	parts := strings.Split(field, ".")
	// 如果字段名被分割为两部分，说明它包含了表别名和列名。
	if len(parts) == 2 {
		// 分别转义表别名和列名，并重新组合。
		return fmt.Sprintf("`%s`.`%s`", parts[0], parts[1])
	}
	// 如果字段名未被分割，直接转义整个字段名。
	return fmt.Sprintf("`%s`", field)
}

func DealWithWhereSafe(params ...Condition) (string, []interface{}, error) {
	if len(params) == 0 {
		return "", nil, nil
	}

	var (
		conditions []string
		args       []interface{}
	)

	for _, param := range params {
		if !isSafeColumn(param.Field) {
			return "", nil, fmt.Errorf("invalid column name: %s", param.Field)
		}

		symbol := strings.ToUpper(strings.TrimSpace(param.Symbol))
		if symbol == "" {
			symbol = "="
		}

		escapedCol := escapeColumn(param.Field)

		switch symbol {
		case "=", "!=", ">", ">=", "<", "<=", "LIKE":
			conditions = append(conditions, fmt.Sprintf("%s %s ?", escapedCol, symbol))
			args = append(args, param.Value)

		case "IN":
			valSlice, ok := param.Value.([]interface{})
			if !ok || len(valSlice) == 0 {
				return "", nil, fmt.Errorf("IN symbol requires a non-empty slice for column %s", param.Field)
			}
			placeholders := make([]string, len(valSlice))
			for i, v := range valSlice {
				placeholders[i] = "?"
				args = append(args, v)
			}
			conditions = append(conditions, fmt.Sprintf("`%s` IN (%s)", escapedCol, strings.Join(placeholders, ",")))

		default:
			return "", nil, fmt.Errorf("unsupported symbol: %s", symbol)
		}
	}

	whereClause := "WHERE " + strings.Join(conditions, " AND ")
	return whereClause, args, nil
}

func RawFieldNames(in any, postgreSql ...bool) []string {
	out := make([]string, 0)
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	var pg bool
	if len(postgreSql) > 0 {
		pg = postgreSql[0]
	}

	if v.Kind() != reflect.Struct {
		panic(fmt.Errorf("RawFieldNames only accepts structs; got %T", v))
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fi := typ.Field(i)
		tagv := fi.Tag.Get(dbTag)

		if tagv == "-" {
			continue
		}

		// 解析 db:"xxx,option1,option2"
		if strings.Contains(tagv, ",") {
			tagv = strings.TrimSpace(strings.Split(tagv, ",")[0])
		}
		if tagv == "" {
			tagv = fi.Name
		}
		if tagv == "-" {
			continue
		}

		// 处理字段格式化
		if pg {
			out = append(out, tagv)
		} else {
			out = append(out, escapeColumn(tagv))
		}
	}
	return out
}

func GetOrderBy(sorts []Sort, query string) string {
	for i, s := range sorts {
		filed := escapeColumn(s.Filed)
		if i > 0 {
			query += fmt.Sprintf(", %s %s", filed, s.Order)
		} else {
			query += fmt.Sprintf(" order by %s %s", filed, s.Order)
		}
	}
	return query
}

type (
	Total struct {
		Number int64 `db:"number"`
	}
	Pages struct {
		Page int64
		Size int64
	}
	Sort struct {
		Filed string //排序字段
		Order string //顺序：asc 倒叙:desc
	}
	Condition struct {
		Field  string
		Symbol string //符号 可不传，不传默认 =
		Value  interface{}
	}
	ListConditions struct {
		Pages
		Conditions []Condition
		Sorts      []Sort
	}
)
