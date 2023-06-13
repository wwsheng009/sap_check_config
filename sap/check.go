package sap

import (
	"iump_check/guilog"
	"iump_check/utils"
	"log"
	"os"
	"path"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

const SAPLOGONEXE = "saplogon.exe"

func Check() (bool, error) {
	guilog.PrintDivideLine()
	guilog.Println("SAP客户端:")
	loc := checksapLogonExeLocation()
	if loc == "" {
		guilog.Println("SAP客户端未安装")
		return false, nil
	}

	is_running, _, err := utils.CheckProcessIsRunning("saplogon.exe")
	if err != nil {
		guilog.Println("配置SAP客户端异常：", err.Error())
		return false, err
	}
	if is_running {
		guilog.Println("SAP客户端正在运行,请关闭SAP客户端再次执行")
		return false, nil
	}

	// 处理sap关联文件
	ok := checkAndFixSAPAssoiaction()
	if !ok {
		guilog.Println("无法恢复.sap文件打方式")
	}

	// 处理客户端安全性设置
	appDataRoaming := os.Getenv("APPDATA")
	saprulesfile := filepath.Join(appDataRoaming, "SAP", "Common", "saprules.xml")

	ok, err = UpdateSAPConfig(saprulesfile)
	if ok {
		guilog.Println("配置成功")
	} else {
		guilog.Println("配置失败")
		guilog.Println(err.Error())
	}
	return ok, err
}
func GetSapLogonExe() (string, error) {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\SAP\SAP Shared`, registry.READ)
	if err != nil {
		return "", err
	}
	defer key.Close()
	value, _, err := key.GetStringValue("SAPsysdir")
	if err != nil {
		return "", err
	}
	loc := path.Join(value, "saplogon.exe")
	return loc, nil
}

func checksapLogonExeLocation() string {
	loc := utils.GetAppPath(SAPLOGONEXE)
	if loc == "" {
		loc1, err := GetSapLogonExe()
		if err != nil {
			log.Println(err)
		}
		loc = loc1

	}
	if loc == "" {
		loc = "C:\\Program Files (x86)\\SAP\\FrontEnd\\SAPgui\\saplogon.exe"
	}
	if _, err := os.Stat(loc); err != nil {
		loc = ""
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

func checkAndFixSAPAssoiaction() bool {
	has_error := false

	has_custom_config := false
	key, err := registry.OpenKey(registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Explorer\FileExts\.sap`, registry.QUERY_VALUE)
	if err == nil {
		log.Println(".sap文件关联存在用户自定义设置:")
		has_custom_config = true
		key.Close()
	}

	if has_custom_config {
		need_delete := true
		key, err = registry.OpenKey(registry.CURRENT_USER,
			`Software\Microsoft\Windows\CurrentVersion\Explorer\FileExts\.sap\UserChoice`, registry.QUERY_VALUE)
		if err == nil {
			// 用户修改过.sap关联
			value, _, err := key.GetStringValue("ProgId")
			key.Close()
			log.Println(".sap文件关联当前的设置:", value)
			if err == nil && value == "SAPGui.Shortcut.File" {
				need_delete = false
			}
		}
		if need_delete {
			keyToDelete := `Software\Microsoft\Windows\CurrentVersion\Explorer\FileExts\.sap`
			ok, err := deleteUserRegistryKey(keyToDelete)
			if err != nil {
				has_error = true
			} else if ok {
				has_custom_config = false
			}
		}
	}

	// 配置默认项
	need_reset := false
	// 已经删除了用户自定义的配置，需要配置SAP默认的配置
	if !has_custom_config {
		key, err = registry.OpenKey(registry.CLASSES_ROOT, `.sap`, registry.QUERY_VALUE)
		if err == nil {
			value, _, err := key.GetStringValue("")
			key.Close()
			if err == nil {
				log.Println("读取.sap文件默认关联设置成功:", value)
				if value != "SAPGui.Shortcut.File" {
					need_reset = true
				}
			} else {
				log.Println("读取SAP文件默认设置失败:", err)
				has_error = true
			}
		} else {
			log.Println("读取.sap文件默认关联失败:", err)
			has_error = true
		}
	}

	if need_reset {
		key, err := registry.OpenKey(registry.CLASSES_ROOT, `.sap`, registry.ALL_ACCESS)
		if err == nil {
			if err := key.SetStringValue("", "SAPGui.Shortcut.File"); err != nil {
				log.Println("恢复.sap文件默认关联失败", err)
				has_error = true
			} else {
				log.Println("恢复.sap文件默认关联成功")
			}
			key.Close()
		} else {
			log.Println("无法打开注册表", err)
			guilog.Println("请使用管理员身份运行")
			has_error = true
		}
	}
	if has_error {
		return false
	}
	return true
	// defer key.Close()
}

func deleteUserRegistryKey(node string) (bool, error) {
	// Open the registry key
	key, err := registry.OpenKey(registry.CURRENT_USER, node, registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		return false, err
	}
	// Get the subkey names
	subkeys, err := key.ReadSubKeyNames(-1)
	if err != nil {
		return false, err
	}
	key.Close()

	// Print the subkey names
	for _, subkey := range subkeys {
		key1 := node + "\\" + subkey
		err := registry.DeleteKey(registry.CURRENT_USER, key1)
		log.Println("删除注册表配置:", key1)
		if err != nil {
			log.Println("删除失败:", err)
			return false, err
		} else {
			log.Println("删除成功:")
		}
	}

	err = registry.DeleteKey(registry.CURRENT_USER, node)
	log.Println("删除注册表配置:", node)

	if err != nil {
		log.Println("删除失败:", err)
		return false, err
	} else {
		log.Println("删除成功:")
	}
	return true, nil
}
