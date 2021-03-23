package config

type Properties struct {
	Port            string `env:"MY_APP_PORT" env-default:"1323"`
	Host            string `env:"HOST" env-default:"localhost"`
	DBHost          string `env:"DB_HOST" env-default:"localhost"`
	DBPort          string `env:"DB_PORT" env-default:"27017"`
	DBName          string `env:"DB_NAME" env-default:"blog"`
	PostsCollection string `env:"PRODUCTS_COLLECTION" env-default:"posts"`
	UsersCollection string `env:"USERS_COLLECTION" env-default:"users"`
	JwtTokenSecret  string `env:"JWT_SECRET" env-default:"abrakadabra"`
}
