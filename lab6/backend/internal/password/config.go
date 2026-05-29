package password

type PasswordConfig struct {
	SaltSize   uint32 `env:"SALT_SIZE"`
	Iterations uint32 `env:"ITERATIONS"`
	Memory     uint32 `env:"MEMORY"`
	Threads    uint8  `env:"THREADS"`
	KeyLength  uint32 `env:"KEY_LENGTH"`
}
