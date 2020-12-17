package conf

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Map map[string]interface{}

var Config *Map

func LoadConf(path string) (err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	m := make(Map)
	err = yaml.Unmarshal(data, &m)
	if err != nil {
		return
	}
	Config = &m
	return
}

func Get(key string) (interface{}, error) {
	config := *Config
	if v, ok := config[key]; !ok {
		return nil, errors.New("字段错误配置不存在")
	} else {
		return v, nil
	}
}

func ToInt(v interface{}) (int, error) {
	if value, ok := v.(int); ok {
		return value, nil
	}
	return 0, errors.New("配置参数转换失败")
}

func ToString(v interface{}) (string, error) {
	if value, ok := v.(string); ok {
		return value, nil
	}
	return "", errors.New("配置参数转换失败")
}
