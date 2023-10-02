package goutils

import (
	"errors"
	"fmt"

	"github.com/Muruyung/go-utilities/logger"
)

// =================================================

// direction sorting direction
type direction struct {
	dir string
}

var (
	// DirAsc sorting ASC
	DirAsc = direction{
		dir: "ASC",
	}

	// DirDesc sorting DESC
	DirDesc = direction{
		dir: "DESC",
	}

	// Direction map direction
	Direction = map[string]direction{
		"ASC":  DirAsc,
		"DESC": DirDesc,
	}
)

// =================================================

// QueryBuilderInteractor query builder interactor
type QueryBuilderInteractor interface {
	GetQuery(tablename string, aliases string) (query string, values []interface{}, err error)
	AddSelection(selection string)
	AddSum(column string, aliases string)
	AddCount(column string, aliases string)
	AddPagination(pagination *paginationOption)
	AddSort(direction direction, sortBy ...string)
	AddWhere(attribute string, operation string, value interface{})
	AddRawWhere(listWhere map[string]interface{})
	AddJoin(joinType joinType, tableName, aliases, on string)
	AddGroup(group ...string)
	AddKey(key ...interface{})
	RemoveKey()
	GetKey() string
}

type queryBuilder struct {
	selection  *[]string
	sort       *[]string
	where      *map[string]interface{}
	pagination *paginationOption
	join       *[]join
	group      *[]string
	key        string
}

// JoinType type of join table
type joinType struct {
	join string
}

var (
	// InnerJoin inner join type
	InnerJoin = joinType{
		join: "INNER",
	}

	// LeftJoin left join type
	LeftJoin = joinType{
		join: "LEFT",
	}

	// RightJoin right join type
	RightJoin = joinType{
		join: "RIGHT",
	}
)

type join struct {
	joinType  string
	tableName string
	aliases   string
	on        string
}

// NewQueryBuilder create new query builder
func NewQueryBuilder() QueryBuilderInteractor {
	return &queryBuilder{
		selection:  nil,
		sort:       nil,
		where:      nil,
		pagination: nil,
	}
}

// AddKey add cache key
func (q *queryBuilder) AddKey(key ...interface{}) {
	for _, val := range key {
		if q.key != "" {
			q.key = fmt.Sprintf("%s-%v", q.key, val)
		} else {
			q.key = fmt.Sprintf("%v", val)
		}
	}
}

// RemoveKey remove cache key
func (q *queryBuilder) RemoveKey() {
	q.key = ""
}

// GetKey get cache key
func (q *queryBuilder) GetKey() string {
	return q.key
}

// GetQuery parse query
func (q *queryBuilder) GetQuery(tablename string, aliases string) (query string, values []interface{}, err error) {
	logger.InitLogger("local", "")
	query = `SELECT`
	if q.selection == nil {
		query = fmt.Sprintf(`%s *`, query)
	} else {
		for key, val := range *q.selection {
			if key > 0 {
				query = fmt.Sprintf(`%s,`, query)
			}
			query = fmt.Sprintf(`%s %s`, query, val)
		}
	}

	if aliases != "" {
		aliases = fmt.Sprintf(` %s`, aliases)
	}

	query = fmt.Sprintf(`%s FROM %s%s`, query, tablename, aliases)

	var (
		join       string
		where      string
		sort       string
		pagination string
		group      string
	)
	if q.join != nil {
		join = parseJoin(*q.join...)
		query = fmt.Sprintf(`%s %s`, query, join)
	}

	if q.where != nil {
		where, values, err = parseWhere(*q.where)
		if err != nil {
			logger.Logger.Error(err)
			return
		}
		query = fmt.Sprintf(`%s WHERE %s`, query, where)
	}

	if q.group != nil {
		group, err = parseGroup(*q.group)
		if err != nil {
			logger.Logger.Error(err)
			return
		}
		query = fmt.Sprintf(`%s GROUP BY %s`, query, group)
	}

	if q.sort != nil {
		sort, err = parseSort(*q.sort)
		if err != nil {
			logger.Logger.Error(err)
			return
		}
		query = fmt.Sprintf(`%s ORDER BY %s`, query, sort)
	}

	if q.pagination != nil {
		pagination = parsePagination(*q.pagination)
		query = fmt.Sprintf(`%s %s`, query, pagination)
	}

	return
}

// =================================================

// AddSelection add selection query
func (q *queryBuilder) AddSelection(selection string) {
	arrSelection := make([]string, 0)
	if q.selection != nil {
		arrSelection = append(arrSelection, *q.selection...)
	}

	arrSelection = append(arrSelection, selection)
	q.selection = &arrSelection
}

// AddSum add sum query
func (q *queryBuilder) AddSum(column string, aliases string) {
	selection := fmt.Sprintf(`SUM(%s) %s`, column, aliases)
	q.AddSelection(selection)
}

// AddCount add count query
func (q *queryBuilder) AddCount(column string, aliases string) {
	selection := fmt.Sprintf(`COUNT(DISTINCT %s) %s`, column, aliases)
	q.AddSelection(selection)
}

// =================================================

// AddPagination add pagination query
func (q *queryBuilder) AddPagination(pagination *paginationOption) {
	q.pagination = pagination
	q.AddKey("limit", pagination.limit)
	q.AddKey("offset", pagination.offset)
}

func parsePagination(pagination paginationOption) (query string) {
	query = fmt.Sprintf(`LIMIT %d OFFSET %d`, pagination.limit, pagination.offset)
	return
}

// =================================================

// AddGroup add group query
func (q *queryBuilder) AddGroup(group ...string) {
	arrGroup := make([]string, 0)
	if q.group != nil {
		arrGroup = append(arrGroup, *q.group...)
	}

	arrGroup = append(arrGroup, group...)
	q.group = &arrGroup
}

func parseGroup(groups []string) (query string, err error) {
	for _, val := range groups {
		if query != "" {
			query = fmt.Sprintf("%s, %s", query, val)
		} else {
			query = val
		}
	}
	return
}

// =================================================

// AddSort add sort query
func (q *queryBuilder) AddSort(direction direction, sortBy ...string) {
	sort := make([]string, 0)
	if q.sort != nil {
		sort = *q.sort
	}

	for _, val := range sortBy {
		sort = append(sort, fmt.Sprintf("%s %s", val, direction.dir))
		q.AddKey(val, direction.dir)
	}
	q.sort = &sort
}

func parseSort(sorts []string) (query string, err error) {
	for _, val := range sorts {
		if query != "" {
			query = fmt.Sprintf("%s, %s", query, val)
		} else {
			query = val
		}
	}
	return
}

// =================================================

// AddJoin add join query
func (q *queryBuilder) AddJoin(joinType joinType, tableName string, aliases string, on string) {
	tmpJoin := make([]join, 0)
	if q.join != nil {
		tmpJoin = append(tmpJoin, *q.join...)
	}

	tmpJoin = append(tmpJoin, join{
		joinType:  string(joinType.join),
		tableName: tableName,
		aliases:   aliases,
		on:        on,
	})
	q.join = &tmpJoin
}

func parseJoin(arrJoin ...join) (query string) {
	for _, join := range arrJoin {
		if join.aliases != "" {
			join.aliases = fmt.Sprintf(" %s", join.aliases)
		}
		query = fmt.Sprintf("%s %s JOIN %s%s ON %s",
			query, join.joinType, join.tableName, join.aliases, join.on,
		)
	}

	return
}

// =================================================

// AddWhere add where query
func (q *queryBuilder) AddWhere(attribute string, operation string, value interface{}) {
	where := make(map[string]interface{})
	if q.where != nil {
		where = *q.where
	}

	if operation == "" || operation == "eq" || operation == "=" {
		where[attribute] = value
		q.AddKey(attribute, value)
	} else {
		where[operation] = map[string]interface{}{
			attribute: value,
		}
		q.AddKey(attribute, operation, value)
	}

	q.where = &where
}

// AddRawWhere add raw where query
func (q *queryBuilder) AddRawWhere(listWhere map[string]interface{}) {
	where := make(map[string]interface{})
	if q.where != nil {
		where = *q.where
	}

	for key, val := range listWhere {
		where[key] = val
	}
	q.where = &where
}

func parseWhere(where map[string]interface{}) (query string, values []interface{}, err error) {
	query = ""
	for key, val := range where {
		switch key {
		case "AND", "OR", "NOT":
			switch value := val.(type) {
			case map[string]interface{}:
				q, v, err := parseBoolOperator(key, value)
				if err != nil {
					logger.Logger.Error(err)
					return query, values, err
				}

				if query == "" {
					query = q
				} else {
					query = fmt.Sprintf("(%s AND %s)", query, q)
				}

				values = append(values, v...)
			case []map[string]interface{}:
				for _, arrVal := range value {
					q, v, err := parseBoolOperator(key, arrVal)
					if err != nil {
						return query, values, err
					}

					if query == "" {
						query = q
					} else {
						query = fmt.Sprintf("(%s %s %s)", query, key, q)
					}

					values = append(values, v...)
				}
			default:
				err = fmt.Errorf("invalid value for %v", value)
				logger.Logger.Error(err)
				return
			}
		case "BETWEEN":
			switch value := val.(type) {
			case map[string]interface{}:
				for k, v := range value {
					switch dateVal := v.(type) {
					case map[string]interface{}:
						for k2, v2 := range dateVal {
							q := fmt.Sprintf("(%s BETWEEN ? AND ?)", k)

							if query == "" {
								query = q
							} else {
								query = fmt.Sprintf("(%s AND %s)", query, q)
							}

							values = append(values, k2, v2)
						}
					default:
						err = fmt.Errorf("invalid value for %v", v)
						logger.Logger.Error(err)
						return
					}
				}
			default:
				err = fmt.Errorf("invalid value for %v", value)
				logger.Logger.Error(err)
				return
			}
		default:
			q, value, err := parseValueOperator(key, val)
			if err != nil {
				logger.Logger.Error(err)
				return query, values, err
			}

			if query == "" {
				query = q
			} else {
				query = fmt.Sprintf("(%s AND %s)", query, q)
			}

			values = append(values, value...)

		}
	}

	return
}

func getOperation(key string, op string, value interface{}) (res string) {
	switch op {
	case "lte", "<=":
		res = fmt.Sprintf("%s <= ?", key)
	case "lt", "<":
		res = fmt.Sprintf("%s < ?", key)
	case "gte", ">=":
		res = fmt.Sprintf("%s >= ?", key)
	case "gt", ">":
		res = fmt.Sprintf("%s > ?", key)
	case "eq", "=":
		if value == nil {
			res = fmt.Sprintf("%s IS NULL", key)
		} else {
			res = fmt.Sprintf("%s = ?", key)
		}
	case "is":
		if value == nil {
			res = fmt.Sprintf("%s IS NULL", key)
		} else {
			res = fmt.Sprintf("%s IS ?", key)
		}
	case "neq", "!=":
		if value == nil {
			res = fmt.Sprintf("%s IS NOT NULL", key)
		} else {
			res = fmt.Sprintf("%s != ?", key)
		}
	case "is_not":
		if value == nil {
			res = fmt.Sprintf("%s IS NOT NULL", key)
		} else {
			res = fmt.Sprintf("%s IS NOT ?", key)
		}
	case "in":
		res = fmt.Sprintf("%s = ANY(?)", key)
	case "not_in":
		res = fmt.Sprintf("NOT %s = ANY(?)", key)
	case "between":
		res = key + " BETWEEN ? AND ?"
	case "not_between":
		res = key + " NOT BETWEEN ? AND ?"
	case "starts_with":
		res = key + " LIKE ? || '%'"
	case "ends_with":
		res = key + " LIKE '%' || ?"
	case "substring":
		res = key + " LIKE '%' || ? || '%'"
	case "i_starts_with":
		res = key + " ILIKE ? || '%'"
	case "i_ends_with":
		res = key + " ILIKE '%' || ?"
	case "i_substring":
		res = key + " ILIKE '%' || ? || '%'"
	case "like":
		res = key + " LIKE ?"
	case "i_like":
		res = key + " ILIKE ?"
	case "not_like":
		res = key + " NOT LIKE ?"
	case "not_ilike":
		res = key + " NOT ILIKE ?"
	case "regexp":
		res = key + " REGEXP/~ ?"
	case "not_regexp":
		res = key + " NOT REGEXP/~ ?"
	case "i_regexp":
		res = key + " NOT ~* ?"
	case "not_i_regexp":
		res = key + " NOT !~* ?"
	default:
		if value == nil {
			res = fmt.Sprintf("%s IS NULL", key)
		} else {
			res = fmt.Sprintf("%s = ?", key)
		}
	}
	return
}

func parseValueOperator(attribute string, val interface{}) (query string, values []interface{}, err error) {
	switch value := val.(type) {
	case map[string]interface{}:
		for key, v := range value {
			if _, ok := v.([]interface{}); ok {
				for _, arrVal := range v.([]interface{}) {
					q := getOperation(key, attribute, arrVal)
					if err != nil {
						return query, values, err
					}

					if query == "" {
						query = q
					} else {
						query = fmt.Sprintf("(%s AND %s)", query, q)
					}

					if v != nil {
						values = append(values, v)
					}
				}
			} else {
				q := getOperation(key, attribute, v)
				if err != nil {
					return query, values, err
				}

				if query == "" {
					query = q
				} else {
					query = fmt.Sprintf("(%s AND %s)", query, q)
				}

				if v != nil {
					values = append(values, v)
				}
			}
		}
	default:
		q := getOperation(attribute, "=", val)
		if err != nil {
			return query, values, err
		}
		if query == "" {
			query = q
		} else {
			query = fmt.Sprintf("(%s AND %s)", query, q)
		}

		if val != nil {
			values = append(values, val)
		}
	}

	return
}

func parseBoolOperator(operator string, items map[string]interface{}) (query string, values []interface{}, err error) {
	switch operator {
	case "NOT":
		for key, item := range items {
			var q string
			var v []interface{}
			parseMap := make(map[string]interface{})
			parseMap[key] = item
			q, v, err = parseWhere(parseMap)
			if err != nil {
				return
			}

			if query == "" {
				query = q
			} else {
				query = fmt.Sprintf("(%s OR %s)", query, q)
			}

			values = append(values, v...)
		}
		query = fmt.Sprintf("NOT (%s)", query)
	case "AND":
		for key, item := range items {
			var q string
			var v []interface{}
			parseMap := make(map[string]interface{})
			parseMap[key] = item
			q, v, err = parseWhere(parseMap)
			if err != nil {
				return
			}

			if query == "" {
				query = q
			} else {
				query = fmt.Sprintf("(%s AND %s)", query, q)
			}

			values = append(values, v...)
		}
	case "OR":
		for key, item := range items {
			var q string
			var v []interface{}
			parseMap := make(map[string]interface{})
			parseMap[key] = item
			q, v, err = parseWhere(parseMap)
			if err != nil {
				return
			}

			if query == "" {
				query = q
			} else {
				query = fmt.Sprintf("(%s OR %s)", query, q)
			}

			values = append(values, v...)
		}
	default:
		err = errors.New("invalid operator type")
		return
	}
	return
}
