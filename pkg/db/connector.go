package db

import (
	"fmt"
	"log"
	"net/url"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

//defaule for postgresql
type Config struct {
	Addr            string            `json:"server" env:"DB_ADDR" envDefault:"127.0.0.1"`
	Port            string            `json:"port" env:"DB_PORT" envDefault:"5432"`
	Database        string            `json:"database" env:"DB_USER" envDefault:"user"` //TODO: to do for production
	Username        string            `json:"username" env:"DB_USER" envDefault:"user"`
	Password        string            `json:"password" env:"DB_PASSWORD" envDefault:"password"`
	Parameter       map[string]string `json:"parameter" env:"DB_PARAMETER" envDefault:""`
	MaxIdleConns    int               `json:"maxidleconns" env:"MAX_IDLE_CONNS" envDefault:"5"`
	MaxOpenConns    int               `json:"maxopenconns" env:"MAX_OPEN_CONNS" envDefault:"5"`
	ConnMaxLifetime int               `json:"connmaxlifetime" env:"CONN_MAX_LIFETIME" envDefault:"90"`
	Automigrate     bool              `json:"automigrate" env:"AUTOMIGRATE" envDefault:"1"` //TODO: 為了測試預設為1
}

type Logger struct {
	Logger logger.Interface
}

func ConnectPostgres(cfg *Config, l *Logger) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Taipei",
		cfg.Addr, cfg.Port, cfg.Username, cfg.Password, cfg.Database)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	// db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
	// 	Logger: logger.Default.LogMode(logger.Info),
	// }) //TODO
	if err != nil {
		return nil, err
	}

	return db, nil
}

func ConnectMySqlWithoutDB(c *Config, l *Logger) (*gorm.DB, error) {

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/",
		c.Username, c.Password,
		c.Addr, c.Port)

	gormCfg := &gorm.Config{}
	if l != nil {
		gormCfg.Logger = l.Logger
	}
	return gorm.Open(mysql.Open(dataSourceName), gormCfg)

}
func ConnectMySql(c *Config, l *Logger) (*gorm.DB, error) {

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		c.Username, c.Password,
		c.Addr, c.Port, c.Database)

	baseURL, err := url.Parse(dataSourceName)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	baseURL.RawQuery = MapToQueryString(c.Parameter)
	dataSource := fmt.Sprint(baseURL)

	gormCfg := &gorm.Config{}
	if l != nil {
		gormCfg.Logger = l.Logger
	}
	return gorm.Open(mysql.Open(dataSource), gormCfg)

}

func MapToQueryString(m map[string]string) string {

	params := url.Values{}
	for k, v := range m {
		params.Add(k, v)
	}

	return params.Encode()
}

func GormOpen(cfg *Config, l *Logger) (*gorm.DB, error) {

	//db, err := ConnectMySql(cfg, l)
	db, err := ConnectPostgres(cfg, l)

	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)

	return db, nil
}

func GormOpenWithoutDB(cfg *Config, l *Logger) (*gorm.DB, error) {
	db, err := ConnectMySqlWithoutDB(cfg, l)
	return db, err
}

func GormCreateDB(cfg *Config, l *Logger) error {
	db, err := GormOpenWithoutDB(cfg, l)

	if err != nil {
		return err
	}

	s := fmt.Sprintf("CREATE DATABASE %s CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;",
		cfg.Database)

	dbc := db.Exec(s)
	return dbc.Error
}
