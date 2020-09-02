package serverconfig

// See worldfiles/server.yaml for example of this configuration
type Config struct {
	WorldFilesDir string `yaml:"worldFilesDir"`
	Log           struct {
		File  string
		Level string
	}
	ServerPort int `yaml:"serverPort"`
	WebPort    int `yaml:"webPort"`
	DB         struct {
		UseSSH bool `yaml:"useSSH"`
		SSH    struct {
			User    string
			Host    string
			Port    int
			KeyFile string `yaml:"keyfile"`
		}
		User     string
		Password string
		Host     string
		Port     int
		Name     string
	}
}
