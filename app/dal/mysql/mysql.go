package mysql

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

var (
	db *sqlx.DB
)

func InitDB() error {
	mysqluser := viper.GetString("mysql.user")
	mysqlpwd := viper.GetString("mysql.password")
	mysqlhost := viper.GetString("mysql.host")
	mysqlport := viper.GetInt64("mysql.port")
	mysqlname := viper.GetString("mysql.mysqlname")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", mysqluser, mysqlpwd, mysqlhost, mysqlport, mysqlname)
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		//zap.L().Error("Mysql open failed,err:%v", zap.Error(err))
		return err
	}

	err = db.Ping()
	if err != nil {
		//	zap.L().Error("Mysql ping failed,err:%v", zap.Error(err))
		return err
	}

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(20)
	return nil
}

func Close() {
	_ = db.Close()
}
