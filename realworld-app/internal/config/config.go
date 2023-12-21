package config

type Config struct {
	Auth       Auth       `yaml:"auth"`
	DataSource DataSource `yaml:"dataSource"`
}

type Auth struct {
	TokenSecret string `yaml:"tokenSecret"`
	TokenTTL    int64  `yaml:"tokenTTL"`
}

type DataSource struct {
	Host            string `yaml:"host" env:"DB_HOST"`
	Port            string `yaml:"port" env:"DB_PORT"`
	User            string `yaml:"user" env:"DB_USER"`
	Password        string `yaml:"password" env:"DB_PASSWORD"`
	MigrateUser     string `yaml:"migrateUser" env:"DB_MIGRATE_USER"`
	MigratePassword string `yaml:"migratePassword" env:"DB_MIGRATE_PASSWORD"`
	Name            string `yaml:"name" env:"DB_NAME"`
}
