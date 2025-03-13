package config

import "flag"

var configOath string

func init() {
	flag.StringVar(&configOath, "config-path", ".env", "path to config file")
}
