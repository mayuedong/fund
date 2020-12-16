package util

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"fund/entity"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path"
	"path/filepath"
	"sync"
)

type Sqlite struct {
	db    *sql.DB
	table string
	used  bool
}

const (
	cutset            = ` |：`
	rateTable         = "rate"
	fundJsTable       = "fundJs"
	fundHtmlTable     = "fundHtml"
	currencyListTable = "currencyList"
	indexListTable    = "indexList"
	mixListTable      = "mixList"
	historyTable      = "history"
	//indexTopicTable   = "indexTopic"
)

var (
	g_sqlPool  map[string]*Sqlite
	g_sqlMutex sync.Mutex
	g_sqlCond  *sync.Cond
	g_tables   = []string{rateTable, fundJsTable, fundHtmlTable, currencyListTable, indexListTable, mixListTable, historyTable}
)

func init() {
	g_sqlCond = sync.NewCond(&g_sqlMutex)
	dir, _ := filepath.Abs("data")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}

	g_sqlPool = make(map[string]*Sqlite)
	for _, table := range g_tables {
		g_sqlPool[table] = newSqlite(table)
		g_sqlPool[table].createTable()
	}
}

func GetSqlite(table string) *Sqlite {
	g_sqlMutex.Lock()
	defer g_sqlMutex.Unlock()
	db := g_sqlPool[table]
	for db.used {
		g_sqlCond.Wait()
	}
	if err := db.db.Ping(); nil != err {
		g_sqlPool[table] = db
		db = newSqlite(table)
	}
	db.used = true
	return db
}
func CLOSE() {
	for _, table := range g_tables {
		db := GetSqlite(table)
		db.db.Close()
	}
}
func (this *Sqlite) CLOSE() {
	g_sqlMutex.Lock()
	defer g_sqlMutex.Unlock()
	this.used = false
	g_sqlCond.Signal()
}

func newSqlite(table string) (ret *Sqlite) {
	dir, _ := filepath.Abs("data")
	fpath := path.Join(dir, "sqliteKV.db")
	db, err := sql.Open("sqlite3", fpath)
	if nil != err {
		entity.GetLog().Fatal(err)
	}
	ret = new(Sqlite)
	ret.table = table
	ret.db = db
	ret.used = false
	return ret
}

func (this *Sqlite) createTable() {
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS '%s'('key' VARCHAR(255) PRIMARY KEY, 'val' TEXT);", this.table)
	_, err := this.db.Exec(sql)
	if nil != err {
		entity.GetLog().Fatal(err)
	}
}

func (this *Sqlite) SET(key string, value interface{}) error {
	b, err := json.Marshal(value)
	if nil != err {
		entity.GetLog().Println(err)
	}
	sql := fmt.Sprintf("REPLACE INTO %s(key,val) values(?,?);", this.table)
	stmt, err := this.db.Prepare(sql)
	if nil != err {
		entity.GetLog().Println(sql, err)
		return err
	}
	_, err = stmt.Exec(key, string(b))
	if nil != err {
		entity.GetLog().Println(sql, err)
		return err
	}
	stmt.Close()
	return nil
}

func (this *Sqlite) GET(key string, value interface{}) error {
	sql := fmt.Sprintf("SELECT val FROM %s WHERE key='%s';", this.table, key)
	row := this.db.QueryRow(sql)
	var tmp string
	if err := row.Scan(&tmp); nil != err {
		return err
	}
	if err := json.Unmarshal([]byte(tmp), value); nil != err {
		entity.GetLog().Println(sql, err)
		return err
	}
	return nil
}

func (this *Sqlite) DELETE(key string) {
	sql := fmt.Sprintf("DELETE FROM %s WHERE key='%s';", this.table, key)
	if _, err := this.db.Exec(sql); nil != err {
		entity.GetLog().Println(sql, err)
	}
}

func (m *Sqlite) KEYS() []string {
	keys := make([]string, 0)
	sql := fmt.Sprintf("SELECT key FROM %s", m.table)
	row, err := m.db.Query(sql)
	if nil != err {
		entity.GetLog().Println(sql, err)
		return nil
	}
	defer row.Close()
	for row.Next() {
		var val string
		err = row.Scan(&val)
		if nil != err {
			entity.GetLog().Println(sql, err)
			return nil
		} else {
			keys = append(keys, val)
		}
	}
	return keys
}
func (m *Sqlite) Clean() error {
	sql := fmt.Sprintf("DELETE FROM %s", m.table)
	_, err := m.db.Exec(sql)
	return err
}
