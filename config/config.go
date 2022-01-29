package config

import (
	"gopkg.in/ini.v1"
	"log"
	"net"
)

type Config struct {
	Iface         string
	WolMac        net.HardwareAddr
	HostIp        net.IP
	Port          int
	Origins       []string
	Auth0Url      string
	Auth0Audience string
}

func NewConfig(ini *ini.File) (Config, error) {
	var err error
	var config Config
	sct := ini.Section("WAKEUP")
	config.Iface = sct.Key("WAKEUP_IFACE").String()

	config.WolMac, err = net.ParseMAC(sct.Key("WOL_MAC").String())
	if err != nil {
		log.Fatalf("Failed to parse MAC: %v", err)
	}

	config.HostIp = net.ParseIP(sct.Key("BACKEND_IP").String())

	config.Port, err = sct.Key("WAKEUP_PORT").Int()
	if err != nil {
		log.Fatalf("Failed to parse port: %v", err)
	}

	sct = ini.Section("BACKEND")
	config.Origins = sct.Key("ORIGINS").Strings(",")

	sct = ini.Section("AUTH0")
	config.Auth0Url = sct.Key("URL").String()
	config.Auth0Audience = sct.Key("API_AUDIENCE").String()

	return config, nil
}
