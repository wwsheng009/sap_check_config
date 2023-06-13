package main

import (
	"bytes"
	"embed"
	"fmt"
	"iump_check/browser"
	"iump_check/guilog"
	"iump_check/ico"
	"iump_check/local"
	"iump_check/sap"
	"iump_check/utils"
	"log"
	"os"
	"strings"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var (
	//go:embed assets
	res embed.FS
)

var outTE *walk.TextEdit

type MyMainWindow struct {
	*walk.MainWindow
	sbi *walk.StatusBarItem
}

func getIcon(iconName string) *walk.Icon {
	iconBytes, err := res.ReadFile(fmt.Sprintf("assets/%s.ico", iconName))
	if err != nil {
		log.Fatal(err)
	}
	icoImg, err := ico.Decode(bytes.NewReader(iconBytes))
	if err != nil {
		log.Fatal(err)
	}
	icon, err := walk.NewIconFromImageForDPI(icoImg, 96)
	if err != nil {
		log.Fatal(err)
	}
	return icon
}
func main() {
	file, err := os.Create("配置日志.log")
	if err != nil {
		log.Fatal("Failed to create log file")
	}
	log.SetOutput(file)

	// Close the log file
	defer file.Close()

	mw := new(MyMainWindow)
	var checkAction, forceCheckAction, showAboutBoxAction *walk.Action

	// var recentMenu *walk.Menu
	// var toggleSpecialModePB *walk.PushButton

	// b, err := walk.NewSolidColorBrush(walk.RGB(0, 255, 0))

	iconChecklist := getIcon("checklist")
	iconMain := getIcon("main")
	iconCheckRed := getIcon("check_red")
	iconLog := getIcon("log")
	if err := (MainWindow{
		// Icon:     walk.IconWinLogo,
		Icon:     iconMain,
		AssignTo: &mw.MainWindow,
		Title:    "SAP环境自动配置，版本：20230608 v1.0",
		MenuItems: []MenuItem{
			Menu{
				Text: "配置",
				Items: []MenuItem{
					// Action{
					// 	AssignTo: &openAction,
					// 	Text:     "&Open",
					// 	// Image:       "../img/open.png",
					// 	Enabled:     Bind("enabledCB.Checked"),
					// 	Visible:     Bind("!openHiddenCB.Checked"),
					// 	Shortcut:    Shortcut{walk.ModControl, walk.KeyO},
					// 	OnTriggered: mw.checkAction_Triggered,
					// },
					// Menu{
					// 	AssignTo: &recentMenu,
					// 	Text:     "Recent",
					// },
					Action{
						AssignTo: &checkAction,
						// Image:       iconChecklist,
						Text:        "自动配置",
						OnTriggered: mw.checkAction_Triggered,
					},
					Action{
						// Image:       iconCheckRed,
						AssignTo:    &forceCheckAction,
						Text:        "强制关闭SAP/浏览器并自动配置",
						OnTriggered: mw.forceCheckAction_Triggered,
					},
					Separator{},
					Action{
						// Image:       walk.IconWarning(),
						Text:        "退出",
						OnTriggered: func() { mw.Close() },
					},
				},
			},
			Menu{
				Text: "&查看",
				Items: []MenuItem{
					// Action{
					// 	Text:    "Open / Special Enabled",
					// 	Checked: Bind("enabledCB.Visible"),
					// },
					// Action{
					// 	Text:    "Open Hidden",
					// 	Checked: Bind("openHiddenCB.Visible"),
					// },
					Action{
						Text:        "查看日志",
						OnTriggered: func() { view_log() },
					},
				},
			},
			Menu{
				Text: "&帮助",
				Items: []MenuItem{
					Action{
						// AssignTo:    &showAboutBoxAction,
						Text:        "使用说明",
						OnTriggered: mw.noteAction_Triggered,
					},
					Action{
						AssignTo:    &showAboutBoxAction,
						Text:        "关于",
						OnTriggered: mw.showAboutBoxAction_Triggered,
					},
				},
			},
		},
		ToolBar: ToolBar{
			ButtonStyle: ToolBarButtonImageBeforeText,
			Items: []MenuItem{
				// ActionRef{&checkAction},
				Action{
					Text:  "自动配置",
					Image: iconChecklist,
					// Image:       "../img/system-shutdown.png",
					// Enabled:     Bind("isSpecialMode && enabledCB.Checked"),
					OnTriggered: mw.checkAction_Triggered,
				},
				Action{
					Image: iconCheckRed,
					Text:  "强制关闭SAP/浏览器并自动配置",
					// Image:       "../img/system-shutdown.png",
					// Enabled:     Bind("isSpecialMode && enabledCB.Checked"),
					OnTriggered: mw.forceCheckAction_Triggered,
				},
				// Menu{
				// 	Text: "重新自动配置",
				// 	// Image: "../img/document-new.png",
				// 	OnTriggered: mw.checkAction_Triggered,
				// },
				Separator{},
				// Menu{
				// 	Text: "View",
				// 	// Image: "../img/document-properties.png",
				// 	Items: []MenuItem{
				// 		Action{
				// 			Text:        "X",
				// 			OnTriggered: mw.changeViewAction_Triggered,
				// 		},
				// 		Action{
				// 			Text:        "Y",
				// 			OnTriggered: mw.changeViewAction_Triggered,
				// 		},
				// 		Action{
				// 			Text:        "Z",
				// 			OnTriggered: mw.changeViewAction_Triggered,
				// 		},
				// 	},
				// },
				Separator{},
				Action{
					Image: iconLog,
					Text:  "查看日志",
					// Image:       "../img/system-shutdown.png",
					// Enabled:     Bind("isSpecialMode && enabledCB.Checked"),
					OnTriggered: mw.viewLogAction_Triggered,
				},
			},
		},
		MinSize: Size{Width: 800, Height: 600},
		Size:    Size{Width: 800, Height: 600},
		Layout:  VBox{},
		Font:    Font{PointSize: 14},
		Children: []Widget{
			TextEdit{AssignTo: &outTE, ReadOnly: true, VScroll: true},
			// Composite{
			// 	Layout: HBox{},
			// 	Children: []Widget{
			// 		PushButton{
			// 			Background: SolidColorBrush{Color: walk.RGB(0, 255, 0)},
			// 			// MinSize:    Size{Height: 50, Width: 100},
			// 			Text:      "重新配置并更新",
			// 			OnClicked: mw.checkAction_Triggered,
			// 		},
			// 		PushButton{
			// 			Background: SolidColorBrush{Color: walk.RGB(0, 255, 0)},
			// 			// MinSize:    Size{Height: 50, Width: 100},
			// 			Text:      "强制退出浏览器/SAP并更新",
			// 			OnClicked: mw.forceCheckAction_Triggered,
			// 		},
			// 		PushButton{
			// 			Background: SolidColorBrush{Color: walk.RGB(0, 255, 0)},
			// 			// MinSize:    Size{Height: 50, Width: 100},
			// 			Text:      "查看日志",
			// 			OnClicked: mw.viewLogAction_Triggered,
			// 		},
			// 		HSpacer{},
			// 	},
			// },
		},
		StatusBarItems: []StatusBarItem{
			{
				AssignTo: &mw.sbi,
				// Icon:     icon1,
				// Text:        "状态",
				Width:       180,
				ToolTipText: "配置状态",
				// OnClicked: func() {
				// 	if sbi.Text() == "click" {
				// 		sbi.SetText("again")
				// 		// sbi.SetIcon(icon2)
				// 	} else {
				// 		sbi.SetText("click")
				// 		// sbi.SetIcon(icon1)
				// 	}
				// },
			},
			// StatusBarItem{
			// 	Text:        "left",
			// 	ToolTipText: "no tooltip for me",
			// },
			// StatusBarItem{
			// 	Text: "\tcenter",
			// },
			// StatusBarItem{
			// 	Text: "\t\tright",
			// },
			// StatusBarItem{
			// 	// Icon:        icon1,
			// 	ToolTipText: "An icon with a tooltip",
			// },
		},
	}.Create()); err != nil {
		log.Fatal(err)
	}
	guilog.SetControl(outTE)
	guilog.Reset()
	mw.checkAction_Triggered()
	// mw.noteAction_Triggered()
	// guilog.Println(`    ☝`)
	// guilog.Println(`    ☝`)
	// guilog.Println(`请点击上方的按钮进行自动配置SAP运行环境`)
	// check()
	// mw.SetIcon(walk.IconInformation())
	// mw.SetIcon(icon2)
	// walk.MsgBox(mw, "配置", "自动配置完成，请查看日志。", walk.MsgBoxIconInformation)
	mw.Run()

}

func (mw *MyMainWindow) noteAction_Triggered() {
	// walk.MsgBox(mw, "Open", "Pretend to open a file...", walk.MsgBoxIconInformation)
	guilog.Reset()
	guilog.Println(`    ☝`)
	guilog.Println(`    ☝`)
	snote := `功能介绍：

SAP环境自动配置：

1，取消浏览器的下载文件时的提示。
2，设置浏览器自动打开SAP客户端。
3，取消SAP客户端弹出的安全提示。

自动配置只支持以下的浏览器：
1，谷歌浏览器。
2，微软Edge浏览器。
3，火狐浏览器。

在点击上方的按钮之前请先 >关闭浏览器与SAP客户端<。

如果日志中提示 >配置失败<，请按提示操作。
`
	lines := strings.Split(snote, "\n")
	result := strings.Join(lines, "\r\n")
	guilog.Println(result)
}

func (mw *MyMainWindow) checkAction_Triggered() {
	// walk.MsgBox(mw, "New", "Newing something up... or not.", walk.MsgBoxIconInformation)
	// walk.MsgBox(mw, "New", "Newing something up... or not.", walk.MsgBoxIconInformation)
	guilog.Reset()
	mw.sbi.SetText("")
	check()
	mw.sbi.SetText("自动配置完成")
	walk.MsgBox(mw, "配置", "自动配置完成，请查看日志。", walk.MsgBoxIconInformation)
}
func (mw *MyMainWindow) forceCheckAction_Triggered() {
	// walk.MsgBox(mw, "New", "Newing something up... or not.", walk.MsgBoxIconInformation)
	guilog.Reset()
	mw.sbi.SetText("")
	res := walk.MsgBox(mw, "配置", "强制关闭浏览器/SAP，是否继续", walk.MsgBoxYesNo)
	if res != 6 {
		return
	}
	force_kill()
	mw.sbi.SetText("自动配置完成")
	walk.MsgBox(mw, "配置", "自动配置完成，请查看日志。", walk.MsgBoxIconInformation)
}

// func (mw *MyMainWindow) changeViewAction_Triggered() {
// 	walk.MsgBox(mw, "Change View", "By now you may have guessed it. Nothing changed.", walk.MsgBoxIconInformation)
// }

func (mw *MyMainWindow) showAboutBoxAction_Triggered() {
	walk.MsgBox(mw, "关于", "此程序用于自动优化SAP环境配置。", walk.MsgBoxIconInformation)
}

func (mw *MyMainWindow) viewLogAction_Triggered() {
	// walk.MsgBox(mw, "Special", "Nothing to see here.", walk.MsgBoxIconInformation)
	view_log()
}

func check() {
	// guilog.Println("程序用于配置系统环境是否满足使用SAP GUI客户端,")
	guilog.PrintDivideLine()
	currentTime := time.Now()
	guilog.Println("当前时间：" + currentTime.Format("2006-01-02 15:04:05"))
	ok, err := local.CheckOS()
	if err != nil {
		guilog.Println(err.Error())
		return
	}
	if ok {
		_, err = browser.CheckBrowser()
		if err != nil {
			guilog.Println(err.Error())
			return
		}
	}
	directorys := browser.GetDirectorys()
	if len(directorys) > 0 {
		_, err = sap.Check()
		if err != nil {
			guilog.Println(err.Error())
			return
		}
	}
}
func view_log() {
	guilog.Reset()
	content, err := os.ReadFile("配置日志.log")
	if err != nil {
		guilog.Println("无法读取日志文件:" + err.Error())
		return
	}
	lines := strings.Split(string(content), "\n")

	// Print each line
	for _, line := range lines {
		guilog.Println(line)
	}

}
func force_kill() {

	guilog.PrintDivideLine()
	guilog.Println("强制关闭:")
	process := []string{"chrome.exe", "msedge.exe", "firefox.exe", "saplogon.exe"}
	for _, v := range process {
		ok, err := utils.CheckProcessIsInTasklist(v)
		if err != nil {
			log.Println(err)
		} else if ok {
			_, err := utils.KillProcess(v)
			if err != nil {
				log.Println(err)
			} else {
				guilog.Println("关闭:", v)
			}
		}
	}
	check()
}
