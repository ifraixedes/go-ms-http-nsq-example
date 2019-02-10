package redis

// closer is the interface that satisfies some (if not all) of the
// go-redis/redis.Cmdable implementation types.
type closer interface {
	Close() error
}

// closer is the interface that satisfies the most of the Cmd types of the
// go-redis/redis package
type errorer interface {
	Err() error
}
