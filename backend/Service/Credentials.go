package service

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/godror/godror"
)

type Credentials struct {
	Database string
	Server   string
	Port     int

	Password string
	User     string

	ShowAll bool

	db *sql.DB
}

type Table struct {
	Name        string
	ColumnNames []string
	Rows        [][]any
}

type ColumnsMetadata struct {
	Name        string
	DataType    string
	DefaultData string
	Detail      string
	IsNullable  bool
	IsIdentity  bool
	OrdPosition int
}

type ConstraintMetadata struct {
	ConstraintName     string
	ConstraintType     string
	ConstraintTypeName string
	ColumnName         string
	SearchCondition    string
	RefConstraint      string
	Position           int
}

type FunctionArgument struct {
	Name       string
	Position   int
	DataType   string
	InOut      string
	Length     int
	Precision  int
	Scale      int
	HasDefault bool
}

func CreateCredFromGin(c *gin.Context) Credentials {
	i, _ := strconv.Atoi(c.Request.PostFormValue("Port"))
	showAll := c.Request.PostFormValue("ShowAll")
	return Credentials{
		Server:   c.Request.PostFormValue("Server"),
		Port:     i,
		Database: c.Request.PostFormValue("Database"),
		User:     c.Request.PostFormValue("Username"),
		Password: c.Request.PostFormValue("Password"),
		ShowAll:  showAll == "on",
	}

}

func (c *Credentials) Connect() error {
	if c.db != nil {
		return nil
	}

	db, err := sql.Open("godror", fmt.Sprintf(`user="%s" password="%s" connectString="%s:%d/%s"`, c.User, c.Password, c.Server, c.Port, c.Database))
	if err != nil {
		return fmt.Errorf("Unable to stablish connection!\n%v", err)
	}
	c.db = db
	return nil
}

func (c *Credentials) Ping() error {
	if c.db == nil {
		if err := c.Connect(); err != nil {
			return nil
		}
	}

	return c.db.Ping()
}

func (c *Credentials) Close() error {
	if c.db != nil {
		return c.db.Close()
	}

	return nil
}

func (c *Credentials) GetDB() *sql.DB {
	return c.db
}

func (c *Credentials) Exec(query string) (int64, error) {
	db := c.GetDB()
	result, err := db.Exec(query)
	if err != nil {
		return -1, fmt.Errorf("[ERROR] %v\n", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return -1, fmt.Errorf("[ERROR] %v\n", err)
	}

	return affected, nil
}

func (c *Credentials) QueryRows(query string) (*Table, error) {
	rows, err := c.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] %v\n", err)
	}
	defer rows.Close()

	var table = Table{}
	err = fetchRows(rows, &table)

	return &table, err
}

func (c *Credentials) QueryTable(tableName string) (*Table, error) {
	rows, err := c.db.Query(fmt.Sprintf("SELECT * FROM %s", tableName))
	if err != nil {
		return nil, fmt.Errorf("[ERROR] %v\n", err)
	}
	defer rows.Close()

	var table = Table{
		Name: tableName,
	}
	err = fetchRows(rows, &table)

	return &table, err
}

func (c *Credentials) QueryTableColumnsMetadata(tableName string) ([]ColumnsMetadata, error) {
	rows, err := c.db.Query(
		`
		SELECT 
			COLUMN_NAME,
			DATA_TYPE,
			CASE 
				WHEN DATA_TYPE IN ('VARCHAR2', 'CHAR', 'NVARCHAR2', 'NCHAR')
					THEN TO_CHAR(CHAR_LENGTH)
				WHEN DATA_TYPE = 'NUMBER'
					THEN 
						CASE
							WHEN DATA_PRECISION IS NULL THEN NULL
							WHEN DATA_SCALE IS NULL THEN TO_CHAR(DATA_PRECISION)
							ELSE TO_CHAR(DATA_PRECISION) || ',' || TO_CHAR(DATA_SCALE)
						END
				WHEN DATA_TYPE = 'FLOAT'
					THEN TO_CHAR(DATA_PRECISION)
				WHEN DATA_TYPE LIKE 'TIMESTAMP%'
					THEN TO_CHAR(DATA_SCALE)
				ELSE NULL
			END AS DATA_DETAIL,
			NULLABLE,
			IDENTITY_COLUMN,
			COLUMN_ID,
			DATA_DEFAULT
		FROM ALL_TAB_COLUMNS
		WHERE TABLE_NAME = :1
		ORDER BY COLUMN_ID
		`, strings.ToUpper(tableName))
	if err != nil {
		return nil, fmt.Errorf("[ERROR] %v", err)
	}
	defer rows.Close()

	var columns []ColumnsMetadata
	for rows.Next() {
		var c ColumnsMetadata
		var nullable string
		var identity string
		var defaultData sql.NullString
		var detail sql.NullString

		if err := rows.Scan(
			&c.Name,
			&c.DataType,
			&detail,
			&nullable,
			&identity,
			&c.OrdPosition,
			&defaultData,
		); err != nil {
			return nil, fmt.Errorf("[ERROR] %v", err)
		}

		c.Detail = ""
		if detail.Valid {
			c.Detail = detail.String
		}

		c.DefaultData = ""
		if defaultData.Valid {
			c.DefaultData = defaultData.String
		}

		c.IsNullable = nullable == "Y"
		c.IsIdentity = identity == "YES"

		columns = append(columns, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("[ERROR] %v", err)
	}

	return columns, err
}

func (c *Credentials) QueryTableConstraints(tableName string) ([]ConstraintMetadata, error) {
	rows, err := c.db.Query(
		`
		SELECT
			c.constraint_name,
			c.constraint_type,
			CASE c.constraint_type
				WHEN 'P' THEN 'PRIMARY KEY'
				WHEN 'R' THEN 'FOREIGN KEY'
				WHEN 'U' THEN 'UNIQUE'
				WHEN 'C' THEN 'CHECK'
			END AS constraint_type_name,
			col.column_name,
			c.search_condition,
			c.r_constraint_name,
			col.position
		FROM ALL_CONSTRAINTS c
		JOIN ALL_CONS_COLUMNS col
		  ON c.owner = col.owner
		 AND c.constraint_name = col.constraint_name
		WHERE c.table_name = :1
		  AND c.constraint_type IN ('P','R','U','C')
		ORDER BY c.constraint_name, col.position
		`, strings.ToUpper(tableName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var constraints []ConstraintMetadata
	for rows.Next() {
		var c ConstraintMetadata
		var searchCondition sql.NullString
		var refConstraint sql.NullString

		if err := rows.Scan(
			&c.ConstraintName,
			&c.ConstraintType,
			&c.ConstraintTypeName,
			&c.ColumnName,
			&searchCondition,
			&refConstraint,
			&c.Position,
		); err != nil {
			return nil, fmt.Errorf("[ERROR] %v", err)
		}

		if searchCondition.Valid {
			c.SearchCondition = searchCondition.String
		} else {
			c.SearchCondition = ""
		}

		if refConstraint.Valid {
			c.RefConstraint = refConstraint.String
		} else {
			c.RefConstraint = ""
		}

		constraints = append(constraints, c)
	}

	return constraints, nil
}

func (c *Credentials) QueryFunctionArguments(owner string, name string) ([]FunctionArgument, error) {
	rows, err := c.db.Query(
		`
		SELECT
			NVL(argument_name, 'RETURN') AS argument_name,
			position,
			data_type,
			in_out,
			CASE 
				WHEN data_type IN ('VARCHAR2','CHAR','NVARCHAR2','NCHAR')
					THEN TO_CHAR(data_length)
				WHEN data_type = 'NUMBER'
					THEN 
						CASE
							WHEN data_precision IS NULL THEN NULL
							WHEN data_scale IS NULL THEN TO_CHAR(data_precision)
							ELSE TO_CHAR(data_precision) || ',' || TO_CHAR(data_scale)
						END
				WHEN data_type = 'FLOAT'
					THEN TO_CHAR(data_precision)
				WHEN data_type LIKE 'TIMESTAMP%'
					THEN TO_CHAR(data_scale)
				ELSE TO_CHAR(data_length)
			END AS detail,
			CASE defaulted WHEN 'Y' THEN 1 ELSE 0 END AS has_default
		FROM all_arguments
		WHERE owner = :owner
		  AND object_name = :object_name
		ORDER BY position;
		`, owner, name)
	if err != nil {
		return nil, fmt.Errorf("[ERROR]: %w", err)
	}
	defer rows.Close()

	var args []FunctionArgument
	for rows.Next() {
		var arg FunctionArgument
		var hasDefaultInt int
		if err := rows.Scan(
			&arg.Name,
			&arg.Position,
			&arg.DataType,
			&arg.InOut,
			&arg.Length,
			&arg.Precision,
			&arg.Scale,
			&hasDefaultInt,
		); err != nil {
			return nil, fmt.Errorf("[ERROR] %w", err)
		}

		arg.HasDefault = hasDefaultInt == 1
		args = append(args, arg)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("[ERROR] %w", err)
	}

	return args, nil
}

func (c *Credentials) QueryDDL(objectType string, objectName string, owner string) (string, error) {
	var ddl string
	err := c.db.QueryRow(
		`
		SELECT DBMS_METADATA.GET_DDL(:1, :2, :3) AS ddl
		FROM dual
		`, objectType, objectName, owner).Scan(&ddl)
	if err != nil {
		return "", fmt.Errorf("[ERROR] %v", err)
	}

	return ddl, nil
}

func fetchRows(rows *sql.Rows, table *Table) error {
	cols, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("[WARN] Column couldn't be read!\n\t[ERROR] Error %v", err)
	}
	colLen := len(cols)

	table.ColumnNames = cols
	table.Rows = [][]any{}

	for rows.Next() {
		columns := make([]any, colLen)
		columnPointers := make([]any, colLen)
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return fmt.Errorf("[ERROR] %v", err)
		}

		rowData := make([]any, colLen)
		for i := range cols {
			val := columns[i]

			if b, ok := val.([]byte); ok {
				rowData[i] = string(b)
			} else {
				rowData[i] = val
			}
		}
		table.Rows = append(table.Rows, rowData)
		fmt.Println()
	}

	return nil
}
