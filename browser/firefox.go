package browser

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"iump_check/guilog"
	"iump_check/utils"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func UpdateFireFoxConfig() bool {
	is_running, mainwindow, e := utils.CheckProcessIsRunning("firefox.exe")
	if e != nil {
		log.Println("配置火狐浏览器异常：", e.Error())
		return false
	}
	if is_running {
		// is_back, e := utils.Check_is_background("msedge.exe")
		if e == nil && mainwindow == 0 {
			utils.KillProcess("firefox.exe")
		}
		if mainwindow > 0 {
			guilog.Println("请关闭浏览器后再次执行")
			return false
		}
	}

	profile_loc := getFirefoxProfileLocation()
	if profile_loc == "" {
		return false
	}

	prefsjs := path.Join(profile_loc, "prefs.js")
	if _, err := os.Stat(prefsjs); !os.IsNotExist(err) {
		err := copyFile(prefsjs, prefsjs+".back")
		if err != nil {
			log.Println("无法备份配置文件:" + "handlers.json")
			log.Println(err)
			return false
		}
	}

	download_dir, err := getFireFoxConfig(prefsjs, "browser.download.dir")
	if err != nil {
		log.Println(err)
	}
	AddDirectory(download_dir)
	download_dir, err = getFireFoxConfig(prefsjs, "browser.download.lastDir")
	if err != nil {
		log.Println(err)
	}
	AddDirectory(download_dir)

	//用户配置
	userjs := path.Join(profile_loc, "user.js")
	if _, err := os.Stat(userjs); !os.IsNotExist(err) {
		download_dir, err = getFireFoxConfig(userjs, "browser.download.dir")
		if err != nil {
			log.Println(err)
			return false
		}
		AddDirectory(download_dir)
		download_dir, err = getFireFoxConfig(userjs, "browser.download.lastDir")
		if err != nil {
			log.Println(err)
			return false
		}
		AddDirectory(download_dir)
	}

	//更新sap默认打开方式
	_, err = updateFireFoxHandlder(profile_loc)
	if err != nil {
		log.Println(err)
		return false
	}
	//默认使用下载目录
	_, err = updateFireFoxConfig(prefsjs, "browser.download.useDownloadDir", true)
	if err != nil {
		log.Println(err)
		return false
	}
	//不需要显示下载对话框
	_, err = updateFireFoxConfig(prefsjs, "browser.download.panel.shown", false)
	if err != nil {
		log.Println(err)
		return false
	}
	//不需要显示下载对话框
	_, err = updateFireFoxConfig(prefsjs, "browser.download.alwaysOpenPanel", false)
	if err != nil {
		log.Println(err)
		return false
	}

	_, err = updateFireFoxConfig(prefsjs, "browser.download.enable_spam_prevention", false)
	if err != nil {
		log.Println(err)
		return false
	}
	//自动隐藏下载按钮
	// updateFireFoxConfig(prefsjs, "browser.download.autohideButton", true)
	return true
}

func getFirefoxProfileLocation() string {
	appData := os.Getenv("APPDATA")
	firefoxPath := filepath.Join(appData, "Mozilla", "Firefox")
	profilesPath := filepath.Join(firefoxPath, "Profiles")

	// List all directories in the Firefox profiles folder
	profiles, err := os.ReadDir(profilesPath)
	if err != nil {
		// fmt.Println(err)
		return ""
	}

	// Find the directory with the ".default" suffix
	for _, profile := range profiles {
		if profile.IsDir() {
			filePath := filepath.Join(profilesPath, profile.Name(), "handlers.json")
			_, err := os.Stat(filePath)
			if !os.IsNotExist(err) {
				profilePath := filepath.Join(profilesPath, profile.Name())
				// fmt.Println("File found!")
				return profilePath
			}
		}
	}
	return ""
}

// 自动打开SAP文件的配置
func updateFireFoxHandlder(profile string) (bool, error) {
	filePath := filepath.Join(profile, "handlers.json")
	content, err := os.ReadFile(filePath)
	if err != nil {
		return false, err
	}

	err = copyFile(filePath, filePath+".back")
	if err != nil {
		log.Println("无法备份浏览器配置:" + "handlers.json")
		return false, err
	}

	// Now let's unmarshall the data into `payload`
	var payload map[string]interface{}
	err = json.Unmarshal(content, &payload)
	if err != nil {
		return false, err
	}
	updateDataInMap(payload, "mimeTypes.application/x-sapshortcut",
		map[string]interface{}{
			"action":     4,
			"extensions": []string{"sap"},
		})
	updatedData, err := json.Marshal(payload)
	if err != nil {
		return false, err
	}

	// Write the updated JSON data back to the file
	err = os.WriteFile(filePath, updatedData, 0644)
	if err != nil {
		// panic(err)
		return false, err
	}
	return true, nil
}

func getFireFoxConfig(firefoxConfigPath string, configKey string) (string, error) {
	configBytes, err := os.ReadFile(firefoxConfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", err
		}
		// log.Println("无法打开配置文件:", firefoxConfigPath, err)
		return "", err
	}
	configString := string(configBytes)
	configLines := strings.Split(configString, "\r\n")
	// const configKey = "network.proxy.http"
	var configValue string

	for _, line := range configLines {
		if strings.HasPrefix(line, "user_pref(\""+configKey+"\",") {
			configValue = strings.TrimSuffix(strings.TrimPrefix(line, "user_pref(\""+configKey+"\", \""), "\");")
			break
		}
	}
	log.Println("读取配置:" + configKey + fmt.Sprintf("==>%v", configValue))

	return configValue, nil
}

func updateFireFoxConfig(firefoxConfigPath string, configKey string, configValue interface{}) (bool, error) {
	if configKey == "" {
		return false, errors.New("不能配置空值")
	}
	log.Println("更新配置:" + configKey + fmt.Sprintf("==>%v", configValue))

	// Open the prefs.js file for writing
	file, err := os.OpenFile(firefoxConfigPath, os.O_CREATE, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(firefoxConfigPath)
			if err != nil {
				log.Println(err)
				return false, err
			}
		} else {
			return false, err
		}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	newFile := ""
	key := fmt.Sprintf("user_pref(\"%s\"", configKey)
	val := fmt.Sprintf("user_pref(\"%s\", %v);\r\n", configKey, configValue)
	// Loop through each line of the prefs.js file
	is_new := true
	for scanner.Scan() {
		line := scanner.Text()

		// Check if the line contains the special key you want to update
		if len(line) > len(key) && line[:len(key)] == key {
			is_new = false
			newFile += val
		} else {
			newFile += line + "\r\n"
		}
	}
	// add new line
	if is_new {
		newFile += val
	}
	err = file.Truncate(0)
	if err != nil {
		return false, err
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		return false, err
	}
	_, err = file.Write([]byte(newFile))
	if err != nil {
		return false, err
	}
	return true, nil
}
func checkFireFoxInstallLocation() string {
	loc := utils.GetAppPath("firefox.exe")
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
