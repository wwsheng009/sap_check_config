package sap

import (
	"encoding/xml"
	"fmt"
	"iump_check/browser"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type Context struct {
	// 系统
	System string `xml:"system,omitempty"`
	// 网络
	Network string `xml:"network,omitempty"`
	// 客户端
	Client string `xml:"client,omitempty"`
	// 事务
	Transaction string `xml:"transaction,omitempty"`
	// 屏幕
	DynproName string `xml:"dynpro_name,omitempty"`
	// 屏幕编号
	DynproNum string `xml:"dynpro_num,omitempty"`
	// 权限 r/w/x 读取/写入/执行
	Permissions string `xml:"permissions"`
	// 操作 0允许 1询问 2拒绝
	Action string `xml:"action,omitempty"`
}

// 快捷方式
type SapShortcut struct {
	Name []string `xml:"name,omitempty"`
}

type CommandLine struct {
	Name []string `xml:"name,omitempty"`
}
type Contexts struct {
	Context []Context `xml:"context,omitempty"`
}
type Directories struct {
	Name []string `xml:"name,omitempty"`
}
type Files struct {
	Name []string `xml:"name,omitempty"`
}
type RegistryValues struct {
	Name []string `xml:"name,omitempty"`
}
type RegistryKeys struct {
	Name []string `xml:"name,omitempty"`
}
type Controls struct {
	Name []string `xml:"name,omitempty"`
}
type EnvVariables struct {
	Name []string `xml:"name,omitempty"`
}
type FileExtensions struct {
	Name []string `xml:"name,omitempty"`
}
type Rule struct {
	ID string `xml:"id,attr"`
	//命令行
	CommandLine *CommandLine `xml:"command_line,omitempty"`
	//目录
	Directories *Directories `xml:"directories,omitempty"`
	//文件
	Files *Files `xml:"files,omitempty"`
	//快捷方式
	SapShortcut *SapShortcut `xml:"sap_shortcut,omitempty"`
	//注册表值
	RegistryValues *RegistryValues `xml:"registry_values,omitempty"`
	//注册表键
	RegistryKeys *RegistryKeys `xml:"registry_keys,omitempty"`
	//控件
	Controls *Controls `xml:"controls,omitempty"`
	//环境变量
	EnvVariables *EnvVariables `xml:"environment_variables,omitempty"`
	//文件扩展名
	FileExtensions *FileExtensions `xml:"file_extensions,omitempty"`
	//r/w/x 读取/写入/执行
	Permissions string `xml:"permissions,omitempty"`
	//操作 3 与上下文有关 0允许 1询问 2拒绝
	Action string `xml:"action,omitempty"`

	Contexts *Contexts `xml:"contexts,omitempty"`
}
type SAP struct {
	Type      string `xml:"type"`
	Version   string `xml:"version"`
	Timestamp string `xml:"timestamp"`
	Rules     struct {
		Rules []Rule `xml:"rule"`
	} `xml:"rules"`
}

func GetShortCutRules() []Rule {
	rules := make([]Rule, 0)
	return rules
}

func UpdateShortCutRules(sap SAP, paths []string) SAP {
	for _, v := range paths {
		sapfile := path.Join(v, "*.sap")
		sap = UpdateShortCutRule(sap, sapfile)
	}
	return sap
}
func UpdateShortCutRule(sap SAP, shortcutFilePath string) SAP {
	log.Println("更新规则：" + shortcutFilePath)
	rules := make([]Rule, 0)
OuterLoop:
	for _, v := range sap.Rules.Rules {
		//删除被拒绝的
		if v.SapShortcut != nil {
			for _, k := range v.SapShortcut.Name {
				sap_sh := strings.ToLower(k)
				// 询问或是拒绝
				if strings.Contains(sap_sh, ".sap") && (v.Action == "1" || v.Action == "2") {
					continue OuterLoop
				}

			}
		}
		rules = append(rules, v)
	}
	// copy(sap.Rules.Rules[:], rules)
	sap.Rules.Rules = rules

	shortcutFilePath = strings.ReplaceAll(shortcutFilePath, "\\", "/")
	found := false
	for _, v := range sap.Rules.Rules {
		if v.SapShortcut != nil {
			for _, k := range v.SapShortcut.Name {
				if k == shortcutFilePath {
					found = true
					break
				}
				if found {
					break
				}
			}
		}

		if found {
			break
		}
	}

	id := 1
	if len(sap.Rules.Rules) > 0 {
		last := sap.Rules.Rules[len(sap.Rules.Rules)-1]

		id1, err := strconv.Atoi(last.ID)
		if err == nil {
			id = id1 + 1
		}
	}

	if !found {
		if sap.Rules.Rules == nil {
			sap.Rules.Rules = make([]Rule, 0)
		}
		sap.Rules.Rules = append(sap.Rules.Rules, Rule{
			ID: fmt.Sprintf("%d", id),
			SapShortcut: &SapShortcut{Name: []string{
				shortcutFilePath,
			}},
			Permissions: "x",
			Action:      "0",
			// Contexts: &Contexts{Context: []Context{
			// 	{
			// 		System:      "*",
			// 		Network:     "",
			// 		Client:      "800",
			// 		Transaction: "<none>",
			// 		DynproName:  "*",
			// 		DynproNum:   "*",
			// 		Permissions: "x",
			// 		Action:      "0",
			// 	},
			// }},
		})
	}

	return sap
}
func UpdateSAPConfig(filepath string) (bool, error) {
	filename := filepath

	// Check if the file exists
	_, err := os.Stat(filename)

	directory := browser.GetDirectorys()

	if os.IsNotExist(err) {
		// File does not exist, create a new one
		// fmt.Println("Creating new file...")
		log.Println("SAP客户端安全性配置文件不存在，创建新文件")
		rules := SAP{}
		rules = UpdateShortCutRules(rules, directory)
		rules.Type = "SAP object rules"
		rules.Version = "1.1"
		t := time.Now()
		format := "2006-01-02 15:04:05"
		formattedTime := t.Format(format)
		rules.Timestamp = formattedTime
		file, err := os.Create(filename)
		if err != nil {
			return false, err
		}
		defer file.Close()
		file.WriteString(xml.Header)
		xmlEncoder := xml.NewEncoder(file)
		if err := xmlEncoder.Encode(rules); err != nil {
			return false, err
		}
	} else {
		// File exists, open and load content
		// fmt.Println("File already exists...")
		// log.Println("正在更新SAP客户端安全性配置文件")
		data, err := os.ReadFile(filename) //os.OpenFile(filename, os.O_CREATE, 0644)
		if err != nil {
			if os.IsPermission(err) {
				// guilog.Println("缺少权限")
				log.Println("没有读取SAP客户端配置文件的权限")
			}
			return false, err
		}
		var rules SAP
		err = xml.Unmarshal(data, &rules)
		if err != nil {
			return false, err
		}

		rules = UpdateShortCutRules(rules, directory)

		xmlData, err := xml.Marshal(rules)
		if err != nil {
			return false, err
		}
		xmlData = []byte(xml.Header + string(xmlData))
		err = os.WriteFile(filename, xmlData, 0644)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}
