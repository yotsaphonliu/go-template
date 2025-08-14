package postgresql

import (
	"fmt"
	"strings"
	"time"
)

func concatSql(query, value string) string {
	if len(value) > 0 {
		if len(query) > 0 {
			query += "AND "
		} else {
			query += "WHERE "
		}

		query += value
	}

	return query
}

func concatContainSql(query, field, value string) string {
	if len(value) > 0 {
		return concatMatchSql(query, field, "%"+value+"%")
	}

	return query
}

func concatExactlyMatchedSQL(query, field, value string) string {
	if len(value) > 0 {
		return concatMatchSql(query, field, value)
	}

	return query
}

func concatContainSqlOr(query string, fields []string, value string) string {
	if len(value) > 0 {
		return concatMatchSqlOr(query, fields, "%"+value+"%")
	}

	return query
}

func concatMatchSqlOr(query string, fields []string, value string) string {
	if len(value) > 0 || len(fields) <= 1 {
		if len(query) > 0 {
			query += "AND "
		} else {
			query += "WHERE "
		}

		query += " ("

		for index, field := range fields {
			query += fmt.Sprintf(" LOWER(%s) LIKE LOWER('%v') ", field, value)

			if index != len(fields)-1 {
				query += " OR "
			}
		}

		query += ") "
	}

	return query
}

func concatMatchSql(query, field string, value string) string {
	if len(value) > 0 {
		if len(query) > 0 {
			query += "AND "
		} else {
			query += "WHERE "
		}

		query += fmt.Sprintf(`
			LOWER(%s) LIKE LOWER('%v')
		`, field, value)
	}

	return query
}

func concatMatchSqlInt(query, field string, value int) string {
	if value != 0 {
		if len(query) > 0 {
			query += "AND "
		} else {
			query += "WHERE "
		}

		query += fmt.Sprintf(`
			%s = %v
		`, field, value)
	}

	return query
}

func concatMatchSqlInt64(query, field string, value int64) string {
	if value != 0 {
		if len(query) > 0 {
			query += "AND "
		} else {
			query += "WHERE "
		}

		query += fmt.Sprintf(`
			%s = %v
		`, field, value)
	}

	return query
}

func concatMatchSqlInt64Or(query, field1, field2 string, value int64) string {

	if value != 0 {
		if len(query) > 0 {
			query += "AND "
		} else {
			query += "WHERE "
		}

		query += fmt.Sprintf(`
			(%s = %v OR %s = %v)
		`, field1, value, field2, value)
	}

	return query
}

func concatNotMatchSqlInt64(query, field string, value int64) string {

	if len(query) > 0 {
		query += "AND "
	} else {
		query += "WHERE "
	}

	query += fmt.Sprintf(`
			%s <> %v
		`, field, value)

	return query
}

func concatMatchSqlFloat64(query, field string, value float64) string {
	if value != 0 {
		if len(query) > 0 {
			query += "AND "
		} else {
			query += "WHERE "
		}

		query += fmt.Sprintf(`
			%s = %v
		`, field, value)
	}

	return query
}

func concatMatchSqlBool(query, field string, value *bool) string {
	if value != nil {
		if len(query) > 0 {
			query += "AND "
		} else {
			query += "WHERE "
		}

		query += fmt.Sprintf(`
			%s = %t
		`, field, *value)
	}

	return query
}

func concatNotMatchSql(query, field, value string) string {
	if len(value) > 0 {
		if len(query) > 0 {
			query += "AND "
		} else {
			query += "WHERE "
		}

		query += fmt.Sprintf(`
			%s NOT LIKE '%s'
		`, field, value)
	}

	return query
}

func concatMatchInListSql(tempQuery, field string, value []string) string {
	valueIN := ""
	for index, v := range value {
		if index != len(value)-1 {
			valueIN += fmt.Sprintf(`'%s',`, v)
		} else {
			valueIN += fmt.Sprintf(`'%s'`, v)
		}
	}

	if len(value) > 0 {
		if len(tempQuery) > 0 {
			tempQuery += "AND "
		} else {
			tempQuery += "WHERE "
		}

		tempQuery += fmt.Sprintf(`
			%s IN (%s)
		`, field, valueIN)
	}

	return tempQuery
}

func concatMatchInInt64ListSql(tempQuery, field string, value []int64) string {
	valueIN := make([]string, 0)
	for _, v := range value {
		valueIN = append(valueIN, fmt.Sprintf("%d", v))
	}

	if len(value) > 0 {
		if len(tempQuery) > 0 {
			tempQuery += "AND "
		} else {
			tempQuery += "WHERE "
		}

		tempQuery += fmt.Sprintf(`
			%s IN (%s)
		`, field, strings.Join(valueIN, ","))
	}

	return tempQuery
}

func concatContainsAllSql(tempQuery, field string, value []string) string {
	containsAll := ""
	for index, v := range value {
		if index != len(value)-1 {
			containsAll += fmt.Sprintf(`'%s',`, v)
		} else {
			containsAll += fmt.Sprintf(`'%s'`, v)
		}
	}

	if len(value) > 0 {
		if len(tempQuery) > 0 {
			tempQuery += "AND "
		} else {
			tempQuery += "WHERE "
		}

		tempQuery += fmt.Sprintf(`
			%s @> ARRAY[%s]
		`, field, containsAll)
	}

	return tempQuery
}

func concatGreaterThanEqualDatetime(query, field string, t *time.Time) string {
	if t != nil {
		if len(query) > 0 {
			query += "AND "
		} else {
			query += "WHERE "
		}

		query += fmt.Sprintf(`
			DATE(%s) >= DATE('%v')
		`, field, t.Format(time.RFC3339))
	}

	return query
}

func concatLessThanDatetime(query, field string, t *time.Time) string {
	if t != nil {
		if len(query) > 0 {
			query += "AND "
		} else {
			query += "WHERE "
		}

		query += fmt.Sprintf(`
			DATE(%s) <= DATE('%v')
		`, field, t.Format(time.RFC3339))
	}

	return query
}

func concatMatchDate(query, field string, t *time.Time) string {
	if t != nil {
		if len(query) > 0 {
			query += "AND "
		} else {
			query += "WHERE "
		}

		query += fmt.Sprintf(`
			CAST(%s as DATE) = CAST('%v' as DATE)
		`, field, t.Format(time.RFC3339))
	}

	return query
}

func concatMatchInList(query, field, value string) string {
	if len(value) != 0 {
		if len(query) > 0 {
			query += "AND "
		} else {
			query += "WHERE "
		}

		query += fmt.Sprintf(`
			'%s' = ANY(%s)
		`, value, field)
	}

	return query
}

func concatContainInList(query, field, value string) string {
	if len(value) != 0 {
		if len(query) > 0 {
			query += "AND "
		} else {
			query += "WHERE "
		}

		value = "%" + value + "%"

		query += fmt.Sprintf(`
			EXISTS (
                 SELECT 1
                 FROM unnest(%s) AS cn
                 WHERE cn LIKE '%s'
             )
		`, field, value)
	}

	return query
}
