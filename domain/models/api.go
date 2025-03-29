package models

type AppService struct {
	App      App      `json:",omitempty"`
	Route    Route    `json:",omitempty"`
	Database Database `json:",omitempty"`
}

type App struct {
	Name string `json:",omitempty"`
	Port string `json:",omitempty"`
}

type Route struct {
	Methods []string `json:",omitempty"`
	Headers []string `json:",omitempty"`
	Origins []string `json:",omitempty"`
}

type Database struct {
	Read DBDetail `json:",omitempty"`
}

type DBDetail struct {
	Username     string `json:",omitempty"`
	Password     string `json:",omitempty"`
	URL          string `json:",omitempty"`
	Port         string `json:",omitempty"`
	DBName       string `json:",omitempty"`
	MaxIdleConns int    `json:",omitempty"`
	MaxOpenConns int    `json:",omitempty"`
	MaxLifeTime  int    `json:",omitempty"`
	Timeout      string `json:",omitempty"`
	SSLMode      string `json:",omitempty"`
}
