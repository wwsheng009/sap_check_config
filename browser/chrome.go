package browser

import (
	"encoding/json"
	"fmt"
	"iump_check/guilog"
	"iump_check/utils"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

func UpdateChromeConfig() bool {

	is_running, mainwindow, err := utils.CheckProcessIsRunning("chrome.exe")
	if err != nil {
		guilog.Println("配置谷歌浏览器异常：", err.Error())
		return false
	}
	if is_running {
		if err == nil && mainwindow == 0 {
			_, err := utils.KillProcess("chrome.exe")
			if err != nil {
				log.Println(err)
				return false
			}
		}
		if mainwindow > 0 {
			guilog.Println("请关闭浏览器后再次执行")
			return false
		}
	}
	local, err := getChromeConfigLocation()
	if err != nil {
		guilog.Println(err.Error())
		return false
	}
	ok, err := updateChromeConfig(local)
	if err != nil {
		log.Println(err)
		return false
	} else if !ok {
		return false
	}
	return true
}
func getChromeConfigLocation() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	chromeDir := filepath.Join(home, "AppData", "Local", "Google", "Chrome", "User Data", "Default")
	prefsPath := filepath.Join(chromeDir, "Preferences")
	if _, err := os.Stat(prefsPath); os.IsNotExist(err) {
		// Preferences file does not exist
		return "", errors.New("配置不存在" + prefsPath)
	}
	return prefsPath, nil
}
func updateChromeConfig(filepath string) (bool, error) {
	// guilog.Println("谷歌浏览器配置文件地址：", filepath)
	// Let's first read the `config.json` file
	err := copyFile(filepath, filepath+".json")
	if err != nil {
		println("无法备份浏览器配置")
		return false, err
	}
	content, err := os.ReadFile(filepath)
	if err != nil {
		return false, nil
	}

	// Now let's unmarshall the data into `payload`
	var payload map[string]interface{}
	err = json.Unmarshal(content, &payload)
	if err != nil {
		return false, nil
	}
	//默认下载地址
	default_dir, _ := getDataInMap(payload, "download.default_directory")
	if default_dir != nil {
		AddDirectory(interfaceToString(default_dir))
	}
	default_dir, _ = getDataInMap(payload, "savefile.default_directory")
	if default_dir != nil {
		AddDirectory(interfaceToString(default_dir))
	}
	default_dir, _ = getDataInMap(payload, "selectfile.last_directory")
	if default_dir != nil {
		AddDirectory(interfaceToString(default_dir))
	}

	//关闭每次提示下载
	log.Println("关闭选项：每次下载都询问我该做些什么")
	updateDataInMap(payload, "download.prompt_for_download", false)

	ext_p, _ := getDataInMap(payload, "download.extensions_to_open")
	ext := interfaceToString(ext_p)
	if ext != "" && !strings.Contains(ext, "sap") {
		ext = fmt.Sprintf("%s:sap", ext)
	} else {
		ext = "sap"
	}
	//增加sap后缀
	log.Println("设置自动打开sap文件")
	updateDataInMap(payload, "download.extensions_to_open", ext)
	//关闭保护
	log.Println("关闭选项：每次下载扫描文件")
	updateDataInMap(payload, "safebrowsing.enabled", false)
	// 增强保护
	log.Println("关闭选项：增强型保护")
	updateDataInMap(payload, "safebrowsing.enhanced", false)

	//关闭选项：下载完成后显示下载内容
	log.Println("关闭选项：下载完成后显示下载内容")
	updateDataInMap(payload, "download_bubble.partial_view_enabled", false)

	// Marshal the updated data into JSON format
	updatedData, err := json.Marshal(payload)
	if err != nil {
		return false, err
	}

	// Write the updated JSON data back to the file
	err = os.WriteFile(filepath, updatedData, 0644)
	if err != nil {
		// panic(err)
		return false, err
	}

	// guilog.Println("浏览器配置文件更新成功!")
	return true, nil
}

func checkChromeInstallLocation() string {
	loc := utils.GetAppPath("chrome.exe")
	if loc == "" {
		localAppData := os.Getenv("LOCALAPPDATA")
		locations := []string{
			"C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe",
			localAppData + "\\Google\\Chrome\\Application\\chrome.exe",
			"C:\\Program Files(x86)\\Google\\Chrome\\Application\\chrome.exe",
			"C:\\Program Files (x86)\\Google\\Application\\chrome.exe",
		}
		for _, v := range locations {
			if _, err := os.Stat(v); err == nil {
				break
			}
		}
	}

	if loc != "" {
		log.Println("安装路径:", loc)
		ver, err := utils.GetExeVersion(loc)
		if err != nil {
			log.Println(err)
		}
		if ver != "" {
			log.Println("版本:", ver)
		}
	}
	return loc
}
