package util

import (
	"database/sql"
	"fund/entity"
	_ "github.com/go-sql-driver/mysql"
)

type MysqlClient struct {
	Ip       string
	Port     int
	User     string
	Password string
	Database string
	Db       *sql.DB
}

var g_mysqlClient *MysqlClient

func GetFundMysql() *MysqlClient {
	if nil == g_mysqlClient {
		g_mysqlClient = new(MysqlClient)
		g_mysqlClient.Connect()
	}
	return g_mysqlClient
}

func (this *MysqlClient) Connect() {
	var err error
	this.Db, err = sql.Open("mysql", `root:@tcp(127.0.0.1:3306)/fund`)
	if err != nil {
		entity.GetLog().Print(err)
	}
}

func (this *MysqlClient) Close() {
	this.Db.Close()
}
