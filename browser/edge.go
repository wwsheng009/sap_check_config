package browser

import (
	"encoding/json"
	"fmt"
	"iump_check/guilog"
	"iump_check/utils"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

func UpdateMSEdgeConfig() bool {
	is_running, mainwindow, e := utils.CheckProcessIsRunning("msedge.exe")
	if e != nil {
		guilog.Println("发生异常：", e.Error())
		return false
	}
	if is_running {
		// is_back, e := utils.Check_is_background("msedge.exe")
		if e == nil && mainwindow == 0 {
			utils.KillProcess("msedge.exe")
		}
		if mainwindow > 0 {
			guilog.Println("请关闭浏览器后再次执行")
			return false
		}
	}
	local, e := getEdgeConfigLocation()
	if e != nil {
		log.Println(e.Error())
		return false
	}
	ok, e := updateEdgeConfigFile(local)
	if e != nil {
		log.Println(e.Error())
		return false
	} else if !ok {
		return false
	}
	return true
}
func getEdgeConfigLocation() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	chromeDir := filepath.Join(home, "AppData", "Local", "Microsoft", "Edge", "User Data", "Default")
	prefsPath := filepath.Join(chromeDir, "Preferences")
	if _, err := os.Stat(prefsPath); os.IsNotExist(err) {
		// Preferences file does not exist
		return "", errors.New("配置不存在：" + prefsPath)
	}
	return prefsPath, nil
}

func updateEdgeConfigFile(filepath string) (bool, error) {
	// guilog.Println("浏览器配置文件地址：", filepath)
	// Let's first read the `config.json` file
	// Copy the contents of the source file to the backup file
	err := copyFile(filepath, filepath+".json")
	if err != nil {
		log.Println("无法备份浏览器配置")
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

	// userProfile := os.Getenv("USERPROFILE")
	appData := os.Getenv("APPDATA")

	tempFolder := path.Join(appData, "Local", "Temp", "MicrosoftEdgeDownloads", "*")

	AddDirectory(tempFolder)
	// download_location["edeg_tmp"] = path.Join(tempFolder, "*", "*.sap")

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
	// update_data_in_map(payload, "safebrowsing.enabled", false)
	// 增强保护
	// update_data_in_map(payload, "safebrowsing.enhanced", false)

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

func checkEdgeInstallLocation() string {
	loc := "C:\\Program Files (x86)\\Microsoft\\Edge\\Application\\msedge.exe"
	if _, err := os.Stat(loc); err != nil {
		loc = ""
	}

	// loc := utils.GetAppPath("msedge.exe")
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
