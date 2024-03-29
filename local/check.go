package local

import (
	"errors"
	"fmt"
	"iump_check/guilog"
	"log"
	"net"
	"os"
	"runtime"

	"golang.org/x/sys/windows/registry"
)

func CheckOS() (bool, error) {
	guilog.PrintDivideLine()

	log.Println("操作系统:", runtime.GOOS)
	log.Println("架构:", runtime.GOARCH)
	hostname, err := os.Hostname()
	if err == nil {
		guilog.Println("主机名:", hostname)
	}
	username := os.Getenv("USER")
	if username != "" {
		guilog.Println("用户名:", username)
	}
	addrs, err := net.InterfaceAddrs()
	if err == nil {
		log.Println("IP地址：")
		for _, addr := range addrs {
			// check if the IP address is not a loopback and is IPv4
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
				log.Println(ipnet.IP.String())
			}
		}
	}
	guilog.PrintDivideLine()
	if runtime.GOOS != "windows" {
		return false, errors.New("不支持非Windows操作系统")
	} else {
		k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
		if err != nil {
			return false, err
		}
		defer k.Close()

		pn, _, _ := k.GetStringValue("ProductName")

		maj, _, _ := k.GetStringValue("EditionID")
		if maj != "" {
			maj = fmt.Sprintf(" [%s]", maj)
		}
		cv, _, _ := k.GetStringValue("DisplayVersion")
		if cv != "" {
			cv = fmt.Sprintf(" [%s]", cv)
		}
		guilog.Println(fmt.Sprintf("操作系统版本: %s%s%s", pn, maj, cv))
	}

	return true, nil
}
