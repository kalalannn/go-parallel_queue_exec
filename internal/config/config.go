package config

type Config struct {
	App struct {
		Host                      string `yaml:"host"`
		Port                      int    `yaml:"port"`
		StaticEndpoint            string `yaml:"static_endpoint"`
		PublicFolder              string `yaml:"public_folder"`
		ViewsFolder               string `yaml:"views_folder"`
		TemplatesExt              string `yaml:"templates_ext"`
		FiberShutdownTimeoutMs    int    `yaml:"fiber_shutdown_timeout_ms"`
		ExecutorShutdownTimeoutMs int    `yaml:"executor_shutdown_timeout_ms"`
	} `yaml:"app"`

	ExecutorService struct {
		WorkersLimit int `yaml:"workers_limit"`
	} `yaml:"executor_service"`
}
