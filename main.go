package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/inhies/go-bytesize"
	"github.com/jasonlvhit/gocron"
	"github.com/net-byte/opensocks-gui/static"
	"github.com/net-byte/opensocks/client"
	"github.com/net-byte/opensocks/config"
	"github.com/net-byte/opensocks/counter"
)

var version string = "v1.5.3"

func main() {
	app := app.New()
	app.SetIcon(static.Icon)
	win := app.NewWindow(fmt.Sprintf("opensocks-gui %v", version))
	win.Resize(fyne.NewSize(400, 200))
	win.SetFixedSize(true)
	config := loadConfig()
	localAddr := widget.NewEntry()
	localAddr.Text = config.LocalAddr
	serverAddr := widget.NewEntry()
	serverAddr.Text = config.ServerAddr
	key := widget.NewPasswordEntry()
	key.Text = config.Key
	bypassRadio := widget.NewRadioGroup([]string{"Yes", "No"}, func(value string) {
		if value == "Yes" {
			config.Bypass = true
		} else {
			config.Bypass = false
		}
	})
	bypassRadio.Horizontal = true
	if config.Bypass {
		bypassRadio.SetSelected("Yes")
	} else {
		bypassRadio.SetSelected("No")
	}
	obfuscateRadio := widget.NewRadioGroup([]string{"Yes", "No"}, func(value string) {
		if value == "Yes" {
			config.Obfuscate = true
		} else {
			config.Obfuscate = false
		}
	})
	obfuscateRadio.Horizontal = true
	if config.Obfuscate {
		obfuscateRadio.SetSelected("Yes")
	} else {
		obfuscateRadio.SetSelected("No")
	}
	msg := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{})
	form := &widget.Form{
		Items: []*widget.FormItem{
			widget.NewFormItem("local addr:", localAddr),
			widget.NewFormItem("server addr:", serverAddr),
			widget.NewFormItem("key:", key),
			widget.NewFormItem("obfuscate:", obfuscateRadio),
			widget.NewFormItem("bypass:", bypassRadio),
		},
	}
	tapped := false
	connectBtn := widget.NewButtonWithIcon("connect", theme.MailSendIcon(), func() {
		var err error
		config.LocalAddr = localAddr.Text
		config.ServerAddr = serverAddr.Text
		config.Key = key.Text
		config.Scheme = "wss"
		config.Init()
		if config.LocalAddr == "" || config.ServerAddr == "" {
			msg.Text = "addr can't be empty!"
			return
		}
		_, err = net.ResolveTCPAddr("tcp", config.LocalAddr)
		if nil != err {
			msg.Text = "error local addr format!"
			return
		}
		_, err = net.ResolveTCPAddr("tcp", config.ServerAddr)
		if nil != err {
			msg.Text = "error server addr format!"
			return
		}
		if tapped {
			msg.Text = "already connected!"
			return
		}
		//start client
		go client.Start(config)
		msg.Text = "successfully connected!"
		saveConfig(config)
		tapped = true
		//start shchedule task
		gocron.Every(2).Seconds().Do(statTask, msg)
		gocron.Start()
	})

	exitBtn := widget.NewButtonWithIcon("exit", theme.CancelIcon(), func() {
		win.Close()
	})
	box := container.NewVBox(form, connectBtn, exitBtn, msg)
	win.SetContent(box)
	win.ShowAndRun()
}

func loadConfig() config.Config {
	var result config.Config
	jsonFile, err := os.Open("./config.json")
	if err != nil {
		// init default config
		result = config.Config{}
		result.LocalAddr = "127.0.0.1:1081"
		result.ServerAddr = "demo.opensocks.org:443"
		result.Key = "6w9z$C&F)J@NcRfUjXn2r4u7x!A%D*G-"
		result.Scheme = "wss"
		result.Bypass = false
		result.Obfuscate = false
		return result
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal([]byte(byteValue), &result)
	return result
}

func saveConfig(config config.Config) {
	file, _ := json.MarshalIndent(config, "", " ")
	_ = ioutil.WriteFile("./config.json", file, 0644)
}

func statTask(label *widget.Label) {
	label.Text = fmt.Sprintf("download %v upload %v", bytesize.New(float64(counter.TotalReadByte)).String(), bytesize.New(float64(counter.TotalWriteByte)).String())
	label.Refresh()
}
