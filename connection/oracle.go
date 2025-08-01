package connection

import (
	"database/sql"
	"fmt"
	"log"
	"context"
	_ "github.com/godror/godror"
)

// OracleDB manages Oracle DB connection
type OracleDB struct {
	Conn *sql.DB
}

// NewOracleDB สร้าง connection ใหม่
func NewOracleDB(username, password, host string, port int, serviceName string) *OracleDB {
	dsn := fmt.Sprintf("%s/%s@%s:%d/%s", username, password, host, port, serviceName)

	db, err := sql.Open("godror", dsn)
	if err != nil {
		log.Fatal("Error opening Oracle DB:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Cannot connect to Oracle DB:", err)
	}

	log.Println("Connected to Oracle DB")
	return &OracleDB{Conn: db}
}

func (o *OracleDB) Close() {
	o.Conn.Close()
}

// QuerySQL รับ SQL query string และ args แล้ว return *sql.Rows และ error
func (o *OracleDB) QuerySQL(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return o.Conn.QueryContext(ctx, query, args...)
}

// ExecSQL สำหรับคำสั่ง INSERT, UPDATE, DELETE
func (o *OracleDB) ExecSQL(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return o.Conn.ExecContext(ctx, query, args...)
}
