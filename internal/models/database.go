package models

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func GetDBConfig() DBConfig {
	return DBConfig{
		Host:     "127.0.0.1",
		Port:     "5432",
		User:     "postgres",
		Password: "negr321",
		DBName:   "Negri",
		SSLMode:  "disable",
	}
}
