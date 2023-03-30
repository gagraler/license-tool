package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

// Config 存储全局配置信息
type Config struct {
	config map[string]interface{}
	mutex  sync.Mutex
}

// LoadConfig 从指定路径加载配置文件
func (c *Config) LoadConfig(filename string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	config := make(map[string]interface{})
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "#") && strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if b, err := strconv.ParseBool(value); err == nil {
				config[key] = b
			} else if f, err := strconv.ParseFloat(value, 64); err == nil {
				config[key] = f
			} else {
				config[key] = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	c.config = config
	return nil
}

// GetString 获取指定键的字符串值
func (c *Config) GetString(key string) string {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if value, ok := c.config[key].(string); ok {
		return value
	}

	return ""
}

// GetInt 获取指定键的整型值
func (c *Config) GetInt(key string) (int, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if value, ok := c.config[key].(float64); ok {
		return int(value), nil
	}

	return 0, fmt.Errorf("key %s not found in configuration", key)
}

// GetFloat 获取指定键的浮点型值
func (c *Config) GetFloat(key string) (float64, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if value, ok := c.config[key].(float64); ok {
		return value, nil
	}

	return 0.0, fmt.Errorf("key %s not found in configuration", key)
}

// GetBool 获取指定键的布尔型值
func (c *Config) GetBool(key string) (bool, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if value, ok := c.config[key].(bool); ok {
		return value, nil
	}

	return false, fmt.Errorf("key %s not found in configuration", key)
}

var (
	configOnce   sync.Once
	globalConfig *Config
)

func getConfig() *Config {
	configOnce.Do(func() {
		globalConfig = &Config{}
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = globalConfig.LoadConfig(filepath.Join(dir, "config.ini"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	})
	return globalConfig
}

func main() {
	fmt.Println(getConfig().GetString("database.host"))
	fmt.Println(getConfig().GetString("database.port"))

	port, err := getConfig().GetInt("server.port")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(port)
	}

	floatValue, err := getConfig().GetFloat("some.float.value")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(floatValue)
	}

	boolValue, err := getConfig().GetBool("some.bool.value")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(boolValue)
	}
}
