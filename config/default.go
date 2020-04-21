package config

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type tomlConfig struct {
	ContentDir string
	LayoutDir  string
	Title      string
	Owner      ownerInfo
	DB         database `toml:"database"`
	Servers    map[string]server
	Clients    clients
	Xorm       xorm
}

type xorm struct {
	User          string
	Passwd        string
	Database      string
	SecurePivFile string `toml:"secure_piv_file"`
}

type ownerInfo struct {
	Name string
	Org  string `toml:"organization"`
	Bio  string
	DOB  time.Time
}

type database struct {
	Server  string
	Ports   []int
	ConnMax int `toml:"connection_max"`
	Enabled bool
}

type server struct {
	IP string
	DC string
}

type clients struct {
	Data  [][]interface{}
	Hosts []string
}

// Config global config
var Config tomlConfig

func init() {
	viper.SetDefault("ContentDir", "content")
	viper.SetDefault("LayoutDir", "layouts")

	viper.SetConfigName("config") // name of config file (without extension)
	// viper.SetConfigType("toml")   // REQUIRED if the config file does not have the extension in the name
	// viper.AddConfigPath("/etc/appname/")  // path to look for the config file in
	// viper.AddConfigPath("$HOME/.appname") // call multiple times to add many search paths
	// viper.AddConfigPath("./") // optionally look for config in the working directory
	dir, _ := filepath.Abs("./config")
	// fmt.Println(dir)
	viper.AddConfigPath(dir)    // optionally look for config in the working directory
	err := viper.ReadInConfig() // Find and read the config file
	// if err != nil {             // Handle errors reading the config file
	// 	panic(fmt.Errorf("Fatal error config file: %s \n", err))
	// }
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			fmt.Println("config file not found")
		} else {
			// Config file was found but another error was produced
			panic("Fatal error config file")
		}
	}
	// fmt.Println(viper.Get("clients"))
	// fmt.Println(viper.Get("clients.hosts"))
	err = viper.Unmarshal(&Config)
	if err != nil {
		panic(fmt.Errorf("unable to decode into struct, %v", err))
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		// 不能直接 unmarshal
		// err = viper.Unmarshal(&Config)
		// if err != nil {
		// 	panic(fmt.Errorf("unable to decode into struct, %v", err))
		// }
		// 可以考虑加锁操作
		var newConfig tomlConfig
		err = viper.Unmarshal(&newConfig)
		if err != nil {
			panic(fmt.Errorf("unable to decode into struct, %v", err))
		}
		Config = newConfig
	})
}
