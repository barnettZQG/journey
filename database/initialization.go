package database

import (
	"database/sql"
	"time"

	"os"

	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/barnettzqg/journey/structure"
	_ "github.com/go-sql-driver/mysql"
	"github.com/twinj/uuid"
)

// Handler for read access
var readDB *sql.DB
var currentTime = time.Now()
var stmtInitialization = []string{"CREATE TABLE IF NOT EXISTS posts( `id` INT(15) NOT NULL AUTO_INCREMENT, uuid varchar(36) NOT NULL, title VARCHAR(150) NOT NULL, slug VARCHAR(150) NOT NULL, markdown TEXT, html TEXT, image TEXT, featured TINYINT NOT NULL DEFAULT '0', PAGE TINYINT NOT NULL DEFAULT '0', STATUS VARCHAR (150) NOT NULL DEFAULT 'draft', LANGUAGE VARCHAR(6) NOT NULL DEFAULT 'en_US', meta_title VARCHAR(150), meta_description VARCHAR(200), author_id INTEGER NOT NULL, created_at DATETIME NOT NULL, created_by INTEGER NOT NULL, updated_at DATETIME, updated_by INTEGER, published_at DATETIME, published_by INTEGER,PRIMARY KEY (`id`))DEFAULT CHARSET=utf8;",
	" CREATE TABLE IF NOT EXISTS users( `id` INT(15) NOT NULL AUTO_INCREMENT, uuid varchar(36) NOT NULL, `name` VARCHAR(150) NOT NULL, slug VARCHAR(150) NOT NULL, PASSWORD VARCHAR(60) NOT NULL, email VARCHAR(254) NOT NULL, image TEXT, cover TEXT, bio VARCHAR(200), website TEXT, location TEXT, accessibility TEXT, STATUS VARCHAR(150) NOT NULL DEFAULT 'active', LANGUAGE VARCHAR(6) NOT NULL DEFAULT 'en_US', meta_title VARCHAR(150), meta_description VARCHAR(200), last_login DATETIME, created_at DATETIME NOT NULL, created_by INTEGER NOT NULL, updated_at DATETIME, updated_by INTEGER,PRIMARY KEY (`id`))DEFAULT CHARSET=utf8;",
	" CREATE TABLE IF NOT EXISTS tags ( `id` INT(15) NOT NULL AUTO_INCREMENT, uuid varchar(36) NOT NULL, `name` varchar(150) NOT NULL, slug varchar(150) NOT NULL, description varchar(200), parent_id integer, meta_title varchar(150), meta_description varchar(200), created_at datetime NOT NULL, created_by integer NOT NULL, updated_at datetime, updated_by integer,PRIMARY KEY (`id`))DEFAULT CHARSET=utf8;",
	" CREATE TABLE IF NOT EXISTS posts_tags ( `id` INT(15) NOT NULL AUTO_INCREMENT, post_id integer NOT NULL, tag_id integer NOT NULL,PRIMARY KEY (`id`))DEFAULT CHARSET=utf8;",
	" CREATE TABLE IF NOT EXISTS settings( `id` INT(15) NOT NULL AUTO_INCREMENT, UUID varchar(36), `key` VARCHAR(150), `value` TEXT, `type` VARCHAR(150) DEFAULT 'core', created_at DATETIME NOT NULL, created_by INTEGER NOT NULL, updated_at DATETIME, updated_by INTEGER,PRIMARY KEY (`id`))DEFAULT CHARSET=utf8;",
	" CREATE TABLE IF NOT EXISTS roles ( `id` INT(15) NOT NULL AUTO_INCREMENT, uuid varchar(36) NOT NULL, `name` varchar(150) NOT NULL, description varchar(200), created_at datetime NOT NULL, created_by integer NOT NULL, updated_at datetime, updated_by integer ,PRIMARY KEY (`id`))DEFAULT CHARSET=utf8;",
	" CREATE TABLE IF NOT EXISTS roles_users ( `id` INT(15) NOT NULL AUTO_INCREMENT, role_id integer NOT NULL, user_id integer NOT NULL ,PRIMARY KEY (`id`))DEFAULT CHARSET=utf8;"}
var insertSQL = []string{
	"INSERT  INTO settings (id, uuid, `key`, `value`, `type`, created_at, created_by, updated_at, updated_by) VALUES (1, ?, 'title', 'My Blog', 'blog', ?, 1, ?, 1);",
	"INSERT  INTO settings (id, uuid, `key`, `value`, `type`, created_at, created_by, updated_at, updated_by) VALUES (2, ?, 'description', 'Just another Blog', 'blog', ?, 1, ?, 1);",
	"INSERT  INTO settings (id, uuid, `key`, `value`, `type`, created_at, created_by, updated_at, updated_by) VALUES (3, ?, 'email', '', 'blog', ?, 1, ?, 1);",
	"INSERT  INTO settings (id, uuid, `key`, `value`, `type`, created_at, created_by, updated_at, updated_by) VALUES (4, ?, 'logo', '/public/images/blog-logo.jpg', 'blog', ?, 1, ?, 1);",
	"INSERT  INTO settings (id, uuid, `key`, `value`, `type`, created_at, created_by, updated_at, updated_by) VALUES (5, ?, 'cover', '/public/images/blog-cover.jpg', 'blog', ?, 1, ?, 1);",
	"INSERT  INTO settings (id, uuid, `key`, `value`, `type`, created_at, created_by, updated_at, updated_by) VALUES (6, ?, 'postsPerPage', 5, 'blog', ?, 1, ?, 1);",
	"INSERT  INTO settings (id, uuid, `key`, `value`, `type`, created_at, created_by, updated_at, updated_by) VALUES (7, ?, 'activeTheme', 'promenade', 'theme', ?, 1, ?, 1);",
	"INSERT  INTO settings (id, uuid, `key`, `value`, `type`, created_at, created_by, updated_at, updated_by) VALUES (8, ?, 'navigation', '[{\"label\":\"Home\", \"url\":\"/\"}]', 'blog', ?, 1, ?, 1);",
	"INSERT  INTO roles (id, uuid, name, description, created_at, created_by, updated_at, updated_by) VALUES (1, ?, 'Administrator', 'Administrators', ?, 1, ?, 1);",
	"INSERT  INTO roles (id, uuid, name, description, created_at, created_by, updated_at, updated_by) VALUES (2, ?, 'Editor', 'Editors', ?, 1, ?, 1);",
	"INSERT  INTO roles (id, uuid, name, description, created_at, created_by, updated_at, updated_by) VALUES (3, ?, 'Author', 'Authors', ?, 1, ?, 1);",
	"INSERT  INTO roles (id, uuid, name, description, created_at, created_by, updated_at, updated_by) VALUES (4, ?, 'Owner', 'Blog Owner', ?, 1, ?, 1);",
}

//Initialize 初始化
func Initialize() error {

	var err error
	user := os.Getenv("MYSQL_USER")
	pass := os.Getenv("MYSQL_PASS")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	db := os.Getenv("MYSQL_DATABASE")
	if user == "" {
		user = "root"
	}
	if pass == "" {
		pass = "barnettblog"
	}
	if host == "" {
		host = "barnettblogmysql"
	}
	if port == "" {
		port = "3306"
	}
	if db == "" {
		db = "blog"
	}
	database := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, db)
	logrus.Infof("DB(%s)", database)
	readDB, err = sql.Open("mysql", database)
	if err != nil {
		logrus.Error("Open mysql error.", err.Error())
		return err
	}
	readDB.SetMaxIdleConns(256) // TODO: is this enough?
	err = readDB.Ping()
	if err != nil {
		logrus.Error("Ping mysql error.", err.Error())
		return err
	}
	rows, serr := readDB.Query("SELECT * FROM roles ")
	if serr == nil && rows.Next() {
		return nil
	}
	tx, _ := readDB.Begin()
	for i := 0; i < len(stmtInitialization); i++ {
		logrus.Info(stmtInitialization[i])
		_, err = tx.Exec(stmtInitialization[i])
		// TODO: Is Commit()/Rollback() needed for DB.Exec()?
		if err != nil {
			logrus.Error(err)
			tx.Rollback()
			return err
		}
	}
	for i := 0; i < len(insertSQL); i++ {
		logrus.Info(insertSQL[i])
		_, err = tx.Exec(insertSQL[i], uuid.Formatter(uuid.NewV4(), uuid.FormatCanonical), currentTime, currentTime)
		// TODO: Is Commit()/Rollback() needed for DB.Exec()?
		if err != nil {
			logrus.Error(err)
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	err = checkBlogSettings()
	if err != nil {
		return err
	}
	return nil
}

// Function to check and "insert any missing blog settings into the database (settings could be missing if migrating from Ghost).
func checkBlogSettings() error {
	tempBlog := structure.Blog{}
	// Check for title
	row := readDB.QueryRow(stmtRetrieveBlog, "title")
	err := row.Scan(&tempBlog.Title)
	if err != nil {
		// "Insert title
		err = insertSettingString("title", "My Blog", "blog", time.Now(), 1)
		if err != nil {
			return err
		}
	}
	// Check for description
	row = readDB.QueryRow(stmtRetrieveBlog, "description")
	err = row.Scan(&tempBlog.Description)
	if err != nil {
		// Insert description
		err = insertSettingString("description", "Just another Blog", "blog", time.Now(), 1)
		if err != nil {
			return err
		}
	}
	// Check for email
	var email []byte
	row = readDB.QueryRow(stmtRetrieveBlog, "email")
	err = row.Scan(&email)
	if err != nil {
		// Insert email
		err = insertSettingString("email", "", "blog", time.Now(), 1)
		if err != nil {
			return err
		}
	}
	// Check for logo
	row = readDB.QueryRow(stmtRetrieveBlog, "logo")
	err = row.Scan(&tempBlog.Logo)
	if err != nil {
		// Insert logo
		err = insertSettingString("logo", "/public/images/blog-logo.jpg", "blog", time.Now(), 1)
		if err != nil {
			return err
		}
	}
	// Check for cover
	row = readDB.QueryRow(stmtRetrieveBlog, "cover")
	err = row.Scan(&tempBlog.Cover)
	if err != nil {
		// Insert cover
		err = insertSettingString("cover", "/public/images/blog-cover.jpg", "blog", time.Now(), 1)
		if err != nil {
			return err
		}
	}
	// Check for postsPerPage
	row = readDB.QueryRow(stmtRetrieveBlog, "postsPerPage")
	err = row.Scan(&tempBlog.PostsPerPage)
	if err != nil {
		// Insert postsPerPage
		err = insertSettingInt64("postsPerPage", 5, "blog", time.Now(), 1)
		if err != nil {
			return err
		}
	}
	// Check for activeTheme
	row = readDB.QueryRow(stmtRetrieveBlog, "activeTheme")
	err = row.Scan(&tempBlog.ActiveTheme)
	if err != nil {
		// Insert activeTheme
		err = insertSettingString("activeTheme", "promenade", "theme", time.Now(), 1)
		if err != nil {
			return err
		}
	}
	// Check for navigation
	var navigation []byte
	row = readDB.QueryRow(stmtRetrieveBlog, "navigation")
	err = row.Scan(&navigation)
	if err != nil {
		// Insert navigation
		err = insertSettingString("navigation", "[{\"label\":\"Home\", \"url\":\"/\"}]", "blog", time.Now(), 1)
		if err != nil {
			return err
		}
	}
	return nil
}
