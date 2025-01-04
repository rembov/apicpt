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
        Host:     "localhost",
        Port:     "5432",
        User:     "postgres",
        Password: "AuX2313gdr",
        DBName:   "Negri",
        SSLMode:  "disable",
    }
}
