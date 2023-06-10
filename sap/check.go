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
		guilog.Println("SAP客户端发生异常：", err.Error())
		return false, err
	}
	if is_running {
		guilog.Println("SAP客户端正在运行,请关闭SAP客户端再次执行")
		return false, nil
	}

	appDataRoaming := os.Getenv("APPDATA")
	saprulesfile := filepath.Join(appDataRoaming, "SAP", "Common", "saprules.xml")

	ok, err := UpdateSAPConfig(saprulesfile)
	if ok {
		guilog.Println("更新成功")
	} else {
		guilog.Println("更新失败")
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
