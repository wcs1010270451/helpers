package pg

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

const dbTag = "select"

var validColumnName = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*(\.[a-zA-Z_][a-zA-Z0-9_]*)?$`)

// isSafeColumn 检查列名是否合法，防止 SQL 注入。
// 允许的字符: 小写字母、数字、下划线、点号。
func isSafeColumn(field string) bool {
	validColumn := regexp.MustCompile(`^[a-z0-9_\.]+$`)
	return validColumn.MatchString(field)
}

// escapeColumn 转义 SQL 查询中的列名 (PostgreSQL 版本)。
// - 如果包含 "."，表示带有表别名，例如 "u.name"
// - 否则直接转义整个字段
func escapeColumn(field string) string {
	parts := strings.Split(field, ".")
	if len(parts) == 2 {
		return fmt.Sprintf(`"%s"."%s"`, parts[0], parts[1])
	}
	return fmt.Sprintf(`"%s"`, field)
}

// DealWithWhereSafe 生成 WHERE 子句和参数列表（PostgreSQL 兼容）。
func DealWithWhereSafe(params ...Condition) (string, []interface{}, error) {
	if len(params) == 0 {
		return "", nil, nil
	}

	var (
		conditions []string
		args       []interface{}
		argIndex   = 1 // PostgreSQL 占位符从 $1 开始
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
			conditions = append(conditions, fmt.Sprintf("%s %s $%d", escapedCol, symbol, argIndex))
			args = append(args, param.Value)
			argIndex++

		case "IN":
			valSlice, ok := param.Value.([]interface{})
			if !ok || len(valSlice) == 0 {
				return "", nil, fmt.Errorf("IN requires a non-empty slice for column %s", param.Field)
			}
			placeholders := make([]string, len(valSlice))
			for i, v := range valSlice {
				placeholders[i] = fmt.Sprintf("$%d", argIndex)
				argIndex++
				args = append(args, v)
			}
			conditions = append(conditions, fmt.Sprintf("%s IN (%s)", escapedCol, strings.Join(placeholders, ",")))

		case "BETWEEN":
			valSlice, ok := param.Value.([]interface{})
			if !ok || len(valSlice) != 2 {
				return "", nil, fmt.Errorf("BETWEEN requires exactly 2 values for column %s", param.Field)
			}
			conditions = append(conditions, fmt.Sprintf("%s BETWEEN $%d AND $%d", escapedCol, argIndex, argIndex+1))
			args = append(args, valSlice[0], valSlice[1])
			argIndex += 2

		case "IS NULL":
			conditions = append(conditions, fmt.Sprintf("%s IS NULL", escapedCol))

		case "IS NOT NULL":
			conditions = append(conditions, fmt.Sprintf("%s IS NOT NULL", escapedCol))

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
