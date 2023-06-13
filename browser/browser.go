package browser

import (
	"fmt"
	"io"
	"iump_check/guilog"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows/registry"
)

const FIREFOXEXE = "firefox.exe"
const CHROMEEXE = "chrome.exe"
const MSEDGEEXE = "msedge.exe"

var download_location = make(map[string]string)

func GetDirectorys() []string {
	var arr []string
	for _, v := range download_location {
		arr = append(arr, v)
	}
	return arr
}
func AddDirectory(directory string) {
	if directory == "" {
		return
	}
	directory_sap := strings.ReplaceAll(directory, "\\", "/")
	directory_sap = strings.ReplaceAll(directory_sap, "//", "/")
	if _, ok := download_location[directory_sap]; !ok {
		download_location[directory_sap] = directory_sap
		log.Println("下载目录:" + directory_sap)
	}

}
func CheckBrowser() (bool, error) {
	guilog.PrintDivideLine()

	home, err := os.UserHomeDir()
	if err == nil {
		download := filepath.Join(home, "Downloads")
		AddDirectory(download)
	} else {
		log.Println(err)
	}

	default_browser, err := getDefaultBrowser()
	if err == nil {
		browser_name := default_browser + "-不支持"
		if default_browser == "ChromeHTML" {
			browser_name = "谷歌浏览器"
		} else if default_browser == "MSEdgeHTM" {
			browser_name = "微软Edge浏览器"
		} else if default_browser == "360seURL" {
			browser_name = "360安全浏览器-不支持"
		} else if default_browser == "IE.HTTP" {
			browser_name = "IE浏览器-不支持"
		} else if default_browser == "360ChromeXURL" {
			browser_name = "360极速浏览器X-不支持"
		} else if strings.Contains(default_browser, "FirefoxURL") {
			browser_name = "火狐浏览器"
		}
		guilog.Println("默认浏览器：", browser_name)
	} else {
		guilog.Println("无法获取系统默认浏览器")
		log.Println(err)
	}

	guilog.PrintDivideLine()
	guilog.Println("谷歌浏览器:")
	if checkChromeInstallLocation() != "" {
		ok := UpdateChromeConfig()
		if ok {
			guilog.Println("配置成功")
		} else {
			guilog.Println("配置失败")
		}
	} else {
		guilog.Println("没找到")
	}

	guilog.PrintDivideLine()
	guilog.Println("微软Edge浏览器:")

	if checkEdgeInstallLocation() != "" {
		ok := UpdateMSEdgeConfig()
		if ok {
			guilog.Println("配置成功")
		} else {
			guilog.Println("配置失败")
		}
	} else {
		guilog.Println("没找到")
	}
	guilog.PrintDivideLine()
	guilog.Println("火狐浏览器:")

	if checkFireFoxInstallLocation() != "" {
		ok := UpdateFireFoxConfig()
		if ok {
			guilog.Println("配置成功")
		} else {
			guilog.Println("配置失败")
		}
	} else {
		guilog.Println("没找到")
	}

	return false, nil
}

func interfaceToString(i interface{}) string {
	if i == nil {
		return ""
	}
	return fmt.Sprintf("%v", i)
}

func copyFile(source string, target string) error {
	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}
	return nil
}
func getDataInMap(data map[string]interface{}, path string) (interface{}, bool) {
	keys := strings.Split(path, ".")
	for _, key := range keys {
		value, ok := data[key]
		if !ok {
			return nil, false
		}
		data, ok = value.(map[string]interface{})
		if !ok {
			log.Println("读取配置:" + path + fmt.Sprintf("==>%v", value))
			return value, true
		}
	}
	return nil, false
}
func updateDataInMap(data map[string]interface{}, key string, val interface{}) map[string]interface{} {
	if key == "" {
		return data
	}
	log.Println("更新配置:" + key + fmt.Sprintf("==>%v", val))
	// Find the key or path to update
	path := key //"person.address.city"
	keys := strings.Split(path, ".")
	var current interface{} = data
	for i, key := range keys {
		if current == nil {
			continue
		}
		if i == len(keys)-1 {
			// Update the value at the specified key
			current.(map[string]interface{})[key] = val
		} else {
			// Traverse the map or struct
			current = current.(map[string]interface{})[key]
		}
	}
	return data
}

func getDefaultBrowser() (string, error) {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\Shell\Associations\UrlAssociations\http\UserChoice`, registry.READ)
	if err != nil {
		return "", err
	}
	defer k.Close()

	// Read the value of the ProgId subkey
	progId, _, err := k.GetStringValue("ProgId")
	if err != nil {
		return "", err
	}

	return progId, nil
}
