package redisgreen

// A Server is a single RedisGreen server.
type Server struct {
	// RedisGreen id of the server.
	Id string `json:"id"`
	// User-visible name of the server.
	Name string `json:"name"`
	// Redis connection URL, including password, host, and port information.
	URL string `json:"url"`
	// Ids for all RedisGreen slaves of this server.
	Slaves []string `json:"slaves"`
}
