package main

import (
	"bytes"
	"embed"
	"io/ioutil"
	"iump_check/browser"
	"iump_check/guilog"
	"iump_check/ico"
	"iump_check/local"
	"iump_check/sap"
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

func main() {

	file, err := os.Create("检查日志.log")
	if err != nil {
		log.Fatal("Failed to create log file")
	}

	log.SetOutput(file)

	// Close the log file
	defer file.Close()
	var mw *walk.MainWindow

	// b, err := walk.NewSolidColorBrush(walk.RGB(0, 255, 0))
	if err := (MainWindow{
		// Icon:     walk.IconWinLogo,
		AssignTo: &mw,
		Title:    "SAP客户端兼容性检查，版本：20230608 v1.0",
		MinSize:  Size{Width: 400, Height: 600},
		Size:     Size{Width: 600, Height: 800},
		Layout:   VBox{},
		Font:     Font{PointSize: 12},
		Children: []Widget{
			TextEdit{AssignTo: &outTE, ReadOnly: true, VScroll: true},
			Composite{
				Layout: HBox{},
				Children: []Widget{

					PushButton{
						Background: SolidColorBrush{Color: walk.RGB(0, 255, 0)},
						// MinSize:    Size{Height: 50, Width: 100},
						Text: "重新检查并更新",
						OnClicked: func() {
							check()
						},
					},
					PushButton{
						Background: SolidColorBrush{Color: walk.RGB(0, 255, 0)},
						// MinSize:    Size{Height: 50, Width: 100},
						Text: "查看日志",
						OnClicked: func() {
							view_log()
						},
					},
					HSpacer{},
				},
			},
		},
	}.Create()); err != nil {
		log.Fatal(err)
	}
	guilog.SetControl(outTE)
	check()
	mw.SetIcon(walk.IconInformation())
	icon1, err := res.ReadFile("assets/checklist.ico")
	if err == nil {

		icoImg, err := ico.Decode(bytes.NewReader(icon1))
		if err != nil {
			panic(err)
		}
		if err == nil {
			icon1, err := walk.NewIconFromImage(icoImg)
			if err == nil {
				mw.SetIcon(icon1)
			}
			// handle error
		}
	}

	mw.Run()

}
func check() {
	guilog.Reset()
	// guilog.Println("程序用于检查系统环境是否满足使用SAP GUI客户端,")
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

	content, err := ioutil.ReadFile("检查日志.log")
	if err != nil {
		guilog.Println("无法读取日志文件:" + err.Error())
		return
	}
	lines := strings.Split(string(content), "\n")

	guilog.Reset()
	// Print each line
	for _, line := range lines {
		guilog.Println(line)
	}

}
