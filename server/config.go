package server

type Config struct {
	Debug bool
	Port  int

	RedisAddress  string
	RedisPassword string
	RedisDB       int

	SlackClientID     string
	SlackClientSecret string
}
