package config

// App defines the environment variable configuration for the whole app
type App struct {
	// Port is the port of the HTTP server.
	Port string `envconfig:"port" default:"1001"`

	// LogLevel defines the log level.
	LogLevel string `envconfig:"log_level" default:"info"`

	// JWTSignKey defines the key used for signing JWTs.
	JWTSignKey string `envconfig:"jwt_sign_key" default:"info"`

	// Mongo is the configuration of the MongoDB server.
	Mongo Mongo `envconfig:"mongo"`

	// Mongo is the configuration of the MongoDB server.
	Webhook Webhook `envconfig:"webhook"`
}

// Mongo defines the environment variable configuration for MongoDB
type Mongo struct {
	// URI is the MongoDB server URI.
	URI string `envconfig:"uri" default:"mongodb://localhost:27017"`
}

// Webhook defines the environment variable configuration for webhook
type Webhook struct {

	// URL is called before/after every request.
	URL string `envconfig:"url" default:""`

	// URL is called before/after every request.
	Headers string `envconfig:"headers" default:""`
}
