package generator

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// ColumnInfo 数据库列信息
type ColumnInfo struct {
	Name         string
	Type         string
	FullType     string // 完整类型，如 tinyint(1)
	IsNullable   bool
	DefaultValue *string
	IsPrimary    bool
	IsAutoIncr   bool
	Comment      string
	MaxLength    *int64
	NumericScale *int64
}

// TableInfo 数据库表信息
type TableInfo struct {
	Name    string
	Comment string
	Columns []*ColumnInfo
}

// MySQLAnalyzer MySQL表结构分析器
type MySQLAnalyzer struct {
	db *sql.DB
}

// NewMySQLAnalyzer 创建MySQL分析器
func NewMySQLAnalyzer(host string, port int, user, password, database string) (*MySQLAnalyzer, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, database)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %w", err)
	}

	// 测试连接
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping MySQL: %w", err)
	}

	return &MySQLAnalyzer{db: db}, nil
}

// Close 关闭数据库连接
func (a *MySQLAnalyzer) Close() error {
	if a.db != nil {
		return a.db.Close()
	}
	return nil
}

// GetTableInfo 获取表结构信息
func (a *MySQLAnalyzer) GetTableInfo(tableName string) (*TableInfo, error) {
	// 获取表注释
	tableComment, err := a.getTableComment(tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to get table comment: %w", err)
	}

	// 获取列信息
	columns, err := a.getColumnInfo(tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to get column info: %w", err)
	}

	return &TableInfo{
		Name:    tableName,
		Comment: tableComment,
		Columns: columns,
	}, nil
}

// getTableComment 获取表注释
func (a *MySQLAnalyzer) getTableComment(tableName string) (string, error) {
	query := `
		SELECT TABLE_COMMENT 
		FROM INFORMATION_SCHEMA.TABLES 
		WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?
	`

	var comment sql.NullString
	err := a.db.QueryRow(query, tableName).Scan(&comment)
	if err != nil {
		return "", err
	}

	return comment.String, nil
}

// getColumnInfo 获取列信息
func (a *MySQLAnalyzer) getColumnInfo(tableName string) ([]*ColumnInfo, error) {
	query := `
		SELECT
			COLUMN_NAME,
			DATA_TYPE,
			COLUMN_TYPE,
			IS_NULLABLE,
			COLUMN_DEFAULT,
			COLUMN_KEY,
			EXTRA,
			COLUMN_COMMENT,
			CHARACTER_MAXIMUM_LENGTH,
			NUMERIC_PRECISION,
			NUMERIC_SCALE
		FROM INFORMATION_SCHEMA.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?
		ORDER BY ORDINAL_POSITION
	`

	rows, err := a.db.Query(query, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []*ColumnInfo
	for rows.Next() {
		var (
			name             string
			dataType         string
			columnType       string
			isNullable       string
			defaultValue     sql.NullString
			columnKey        string
			extra            string
			comment          string
			maxLength        sql.NullInt64
			numericPrecision sql.NullInt64
			numericScale     sql.NullInt64
		)

		err := rows.Scan(
			&name, &dataType, &columnType, &isNullable, &defaultValue,
			&columnKey, &extra, &comment, &maxLength, &numericPrecision, &numericScale,
		)
		if err != nil {
			return nil, err
		}

		column := &ColumnInfo{
			Name:       name,
			Type:       dataType,
			FullType:   columnType,
			IsNullable: isNullable == "YES",
			IsPrimary:  columnKey == "PRI",
			IsAutoIncr: strings.Contains(extra, "auto_increment"),
			Comment:    comment,
		}

		if defaultValue.Valid {
			column.DefaultValue = &defaultValue.String
		}

		// 对于字符类型使用CHARACTER_MAXIMUM_LENGTH
		if maxLength.Valid {
			column.MaxLength = &maxLength.Int64
		}

		// 对于数值类型，如果没有字符长度，使用数值精度
		if !maxLength.Valid && numericPrecision.Valid {
			column.MaxLength = &numericPrecision.Int64
		}

		if numericScale.Valid {
			column.NumericScale = &numericScale.Int64
		}

		columns = append(columns, column)
	}

	return columns, rows.Err()
}

// ListTables 列出数据库中的所有表
func (a *MySQLAnalyzer) ListTables() ([]string, error) {
	query := `
		SELECT TABLE_NAME 
		FROM INFORMATION_SCHEMA.TABLES 
		WHERE TABLE_SCHEMA = DATABASE() AND TABLE_TYPE = 'BASE TABLE'
		ORDER BY TABLE_NAME
	`

	rows, err := a.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}
		tables = append(tables, tableName)
	}

	return tables, rows.Err()
}
