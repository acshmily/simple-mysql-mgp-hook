package moudle

type YamlConfig struct {
	Node struct {
		Name string `yaml:"name"`
	}
	// 数据库配置
	Mysql struct {
		User     string `yaml:"user"`
		Host     string `yaml:"host"`
		Password string `yaml:"password"`
		Port     string `yaml:"port"`
		Name     string `yaml:"name"`
	}
	// 定义日志路径
	LogPath struct {
		Path string `yaml:"path"`
	}

	// 定义间隔
	Heartbeat struct {
		Interval    int               `yaml:"interval"`
		Sql         string            `yaml:"sql"`
		Query       map[string]string `yaml:"query-key-value,flow"`
		Checkvalue  string            `yaml:"check-value"`
		Downcommand []string          `yaml:"downcommand,flow"`
		Upcommand   []string          `yaml:"upcommand,flow"`
	}
}
