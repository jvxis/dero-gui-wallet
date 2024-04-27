package settings

import (
	"encoding/json"
	"os"
	"path/filepath"

	"gioui.org/app"
	sysTheme "gioui.org/x/pref/theme"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/theme"
	"github.com/g45t345rt/g45w/utils"
)

var (
	MainTabBarsToken = "tokens"
	MainTabBarsTxs   = "txs"
	FolderLayoutGrid = "grid"
	FolderLayoutList = "list"
)

type AppSettings struct {
	Language                string `json:"language"`
	HideBalance             bool   `json:"hide_balance"`
	SendRingSize            int    `json:"send_ring_size"`
	MainTabBars             string `json:"main_tab_bars"`
	Theme                   string `json:"theme"`
	FolderLayout            string `json:"folder_layout"`
	NodeSelect              string `json:"node_select"`
	Testnet                 bool   `json:"testnet"`
	MobileBackgroundService bool   `json:"mobile_background_service"`
}

var (
	AppDir            string
	IntegratedNodeDir string
	WalletsDir        string
	CacheDir          string
)

var App AppSettings

var Name = "deroW"
var DonationAddress = "dero1qys28f7w548hcg36c5t6jwvrr7w7ce8c5h79lxue7xwsp7wyk05kqqqq32jf7"

// vars below are replaced by -ldflags during build
var Version = "development"
var BuildTime = ""
var GitVersion = "development"

var settingsPath = ""

func Load() error {
	dataDir, err := app.DataDir()
	if err != nil {
		return err
	}

	appDir := filepath.Join(dataDir, "g45w")
	err = utils.CreateFolderPath(appDir)
	if err != nil {
		return err
	}

	settingsPath = filepath.Join(appDir, "settings.json")

	// default values
	appSettings := AppSettings{
		Language:                "en",
		HideBalance:             false,
		SendRingSize:            16,
		NodeSelect:              "",
		MainTabBars:             MainTabBarsTxs,
		FolderLayout:            FolderLayoutGrid,
		Testnet:                 false,
		MobileBackgroundService: false,
	}

	_, err = os.Stat(settingsPath)
	if err == nil {
		data, err := os.ReadFile(settingsPath)
		if err != nil {
			return err
		}

		err = json.Unmarshal(data, &appSettings)
		if err != nil {
			return err
		}
	}

	language := lang.Get(appSettings.Language)
	if language != nil {
		lang.Current = language.Key
	} else {
		lang.Current = "en"
		appSettings.Language = "en"
	}

	currentTheme := theme.Get(appSettings.Theme)
	if currentTheme != nil {
		theme.Current = currentTheme
	} else {
		// check system user theme preference
		isDark, _ := sysTheme.IsDarkMode()
		if isDark {
			appSettings.Theme = "dark"
		} else {
			appSettings.Theme = "light"
		}

		theme.Current = theme.Get(appSettings.Theme)
	}

	App = appSettings

	if App.Testnet {
		appDir = filepath.Join(appDir, "testnet")
		err = utils.CreateFolderPath(appDir)
		if err != nil {
			return err
		}
	}

	AppDir = appDir
	IntegratedNodeDir = filepath.Join(appDir, "node")
	WalletsDir = filepath.Join(appDir, "wallets")
	CacheDir = filepath.Join(appDir, "cache")

	return nil
}

func Save() error {
	data, err := json.MarshalIndent(App, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(settingsPath, data, os.ModePerm)
}
