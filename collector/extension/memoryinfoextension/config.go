package memoryinfoextension

type Config struct {

	// the extension should be configured as enabled,
	// to run it's loop and emit any memory info.
	Enabled bool `mapstructure:"enabled"`

	// how often to emit memory info.
	// format: duration string (20ms, 1s, etc).
	// default is 100ms (10 times per second).
	LoopInterval string `mapstructure:"loopInterval"`
}
