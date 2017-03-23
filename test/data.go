package main

import (
	"database/sql"
	"fmt"

	"os"

	"github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

var writeDB, readDB *sql.DB

const stmtInsertPost = "INSERT INTO posts (id, uuid, title, slug, markdown, html, featured, page, status, image, author_id, created_at, created_by, updated_at, updated_by, published_at, published_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
const stmtInsertUser = "INSERT INTO users (id, uuid, name, slug, password, email, image, cover, created_at, created_by, updated_at, updated_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
const stmtInsertRoleUser = "INSERT INTO roles_users (id, role_id, user_id) VALUES (?, ?, ?)"
const stmtInsertTag = "INSERT INTO tags (id, uuid, name, slug, created_at, created_by, updated_at, updated_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
const stmtInsertPostTag = "INSERT INTO posts_tags (id, post_id, tag_id) VALUES (?, ?, ?)"
const stmtInsertSetting = "INSERT INTO settings (id, uuid, `key`, `value`, `type`, created_at, created_by, updated_at, updated_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"

func init() {
	user := "root"
	pass := "admin"
	host := "127.0.0.1"
	port := "3306"
	db := "blog"
	database := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, db)
	logrus.Infof("DB(%s)", database)
	var err error
	writeDB, err = sql.Open("mysql", database)
	if err != nil {
		logrus.Error("Open mysql error.", err.Error())
	}
	writeDB.SetMaxOpenConns(5)
	readDB, err = sql.Open("sqlite3", "journey.db")
	if err != nil {
		logrus.Error("Open sqlite3 error.", err.Error())
	}
	readDB.SetMaxOpenConns(5)
}
func main() {
	tx, _ := writeDB.Begin()
	if re := ReadTable("posts", tx); re != nil {
		os.Exit(1)
	}
	if re := ReadTable("users", tx); re != nil {
		os.Exit(1)
	}
	if re := ReadTable("roles_users", tx); re != nil {
		os.Exit(1)
	}
	if re := ReadTable("tags", tx); re != nil {
		os.Exit(1)
	}
	if re := ReadTable("posts_tags", tx); re != nil {
		os.Exit(1)
	}
	if re := ReadTable("settings", tx); re != nil {
		os.Exit(1)
	}
	tx.Commit()
}

func ReadTable(tableName string, tx *sql.Tx) error {
	rows, err := readDB.Query("select * from " + tableName)
	if err != nil {
		logrus.Error("Select from "+tableName+" error.", err.Error())
		return err
	}
	result := createData(rows)
	if tableName == "posts" {
		for _, re := range result {
			_, err := tx.Query(stmtInsertPost, re["id"], re["uuid"], re["title"], re["slug"], re["markdown"], re["html"], re["featured"], re["page"], re["status"], re["image"], re["author_id"], re["created_at"], re["created_by"], re["updated_at"], re["updated_by"], re["published_at"], re["published_by"])
			if err != nil {
				tx.Rollback()
				logrus.Error(err)
				return err
			}
		}
	}
	if tableName == "users" {
		for _, re := range result {
			_, err := tx.Query(stmtInsertUser, re["id"], re["uuid"], re["name"], re["slug"], re["password"], re["email"], re["image"], re["cover"], re["created_at"], re["created_by"], re["updated_at"], re["updated_by"])
			if err != nil {
				tx.Rollback()
				logrus.Error(err)
				return err
			}
		}
	}
	if tableName == "roles_users" {
		for _, re := range result {
			_, err := tx.Query(stmtInsertRoleUser, re["id"], re["role_id"], re["user_id"])
			if err != nil {
				tx.Rollback()
				logrus.Error(err)
				return err
			}
		}
	}
	if tableName == "tags" {
		for _, re := range result {
			_, err := tx.Query(stmtInsertTag, re["id"], re["uuid"], re["name"], re["slug"], re["created_at"], re["created_by"], re["updated_at"], re["updated_by"])
			if err != nil {
				tx.Rollback()
				logrus.Error(err)
				return err
			}
		}
	}
	if tableName == "posts_tags" {
		for _, re := range result {
			_, err := tx.Query(stmtInsertPostTag, re["id"], re["post_id"], re["tag_id"])
			if err != nil {
				tx.Rollback()
				logrus.Error(err)
				return err
			}
		}
	}
	if tableName == "settings" {
		for _, re := range result {
			_, err := tx.Query(stmtInsertSetting, re["id"], re["uuid"], re["key"], re["value"], re["type"], re["created_at"], re["created_by"], re["updated_at"], re["updated_by"])
			if err != nil {
				tx.Rollback()
				logrus.Error(err)
				return err
			}
		}
	}
	return nil
}

func createData(rows *sql.Rows) (result []map[string]interface{}) {
	defer rows.Close()
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}
	for rows.Next() {
		record := make(map[string]interface{})
		//将行数据保存到record字典
		err := rows.Scan(scanArgs...)
		if err != nil {
			logrus.Error("查询mysql数据错误," + err.Error())
			return nil
		}
		for i, col := range values {
			if col != nil {
				if colByte, ok := col.([]byte); ok {
					record[columns[i]] = string(colByte)
				} else {
					record[columns[i]] = col
				}
			}
		}
		result = append(result, record)
	}
	return
}
