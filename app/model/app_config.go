package model

type AppConfig struct {
	Debug        bool
	DbUser       string
	DbName       string
	DbPort       string
	DbPublicPort string
	DbHost       string
	DbPassword   string
	BindHost     string
	BindPort     string
	Per          int
	JwtSecretKey string
}
