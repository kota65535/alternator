package cmd

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Database struct {
	Dialect string
	Version string
	Port    int
}

type Fixture struct {
	Database
	Dir string
}

func (r Fixture) Name() string {
	dir := strings.Replace(r.Dir, RootPath+"/", "", 1)
	return fmt.Sprintf("%s_%s_%s", dir, r.Dialect, r.Version)
}

var Databases = []Database{
	{
		Dialect: "mysql",
		Version: "5",
		Port:    13306,
	},
	{
		Dialect: "mysql",
		Version: "8",
		Port:    13307,
	},
}

var RootPath = "test"

// TODO: add assertions for MySQL 5
var Skipped = []string{
	"table/check_mysql_5",
}

func getDirs(path string) []string {
	files, _ := ioutil.ReadDir(path)
	ret := []string{}
	containsDir := false
	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		if strings.HasPrefix(f.Name(), "_") {
			continue
		}
		containsDir = true
		ret = append(ret, getDirs(filepath.Join(path, f.Name()))...)
	}
	if !containsDir {
		ret = append(ret, path)
	}
	return ret
}

func getAlter(dir string, db Database) (string, error) {
	f := testFile(dir, "alter.sql", db)
	s, err := ioutil.ReadFile(f)
	return string(s), err
}

func getDiff(dir string, db Database) (string, error) {
	f := testFile(dir, "diff.txt", db)
	s, err := ioutil.ReadFile(f)
	return string(s), err
}

func prepareDb(dir string, d Database) error {
	dsn := fmt.Sprintf("root@(localhost:%d)/?multiStatements=true", d.Port)
	db, err := sql.Open("mysql", dsn)
	db.SetMaxOpenConns(10)
	defer db.Close()

	// TODO: more smartly
	_, err = db.Exec("DROP DATABASE IF EXISTS db1")
	if err != nil {
		return err
	}
	_, err = db.Exec("DROP DATABASE IF EXISTS db2")
	if err != nil {
		return err
	}
	_, err = db.Exec("DROP DATABASE IF EXISTS db3")
	if err != nil {
		return err
	}
	_, err = db.Exec("DROP DATABASE IF EXISTS db4")
	if err != nil {
		return err
	}

	f := testFile(dir, "from.sql", d)
	b, err := ioutil.ReadFile(f)
	s := string(b)
	if s == "" {
		return nil
	}
	_, err = db.Exec(string(b))
	if err != nil {
		return err
	}
	return nil
}

func testFile(dir string, name string, db Database) string {
	path := filepath.Join(dir, fmt.Sprintf("%s_%s_%s%s", strings.TrimSuffix(name, filepath.Ext(name)), db.Dialect, db.Version, filepath.Ext(name)))
	_, err := os.Stat(path)
	if err != nil {
		path = filepath.Join(dir, name)
		_, err = os.Stat(path)
		if err != nil {
			return ""
		}
	}
	return path
}
