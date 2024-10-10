package config

import (
	"MapPhotoMD/mywidget"
	"encoding/json"
	"io"
	"os"

	"fyne.io/fyne/v2"
)

// UserConfig 用户配置数据
type UserConfig struct {
	Key            string                   `json:"key"`             //高德key
	NotePath       string                   `json:"note_path"`       //ob库路径
	MovePhoto      bool                     `json:"move_photo"`      //是否转存照片
	PhotoPath      string                   `json:"photo_path"`      //转存路径
	DeletePhoto    bool                     `json:"delete_Photo"`    //是否删除原照片
	SaveProperties bool                     `json:"save_properties"` //是否保存YAML属性
	Properties     []*mywidget.PropertyData `json:"properties"`      //旅行记录YAML属性
}

func NewUserConfig() *UserConfig {
	return &UserConfig{}
}

// ReadConfigFile 用于读取配置文件，如成功则将数据保存到config，否则发送错误提醒
func (config *UserConfig) ReadConfigFile(ap fyne.App) {
	//打开配置文件
	file, err := os.Open("config.json")
	if err != nil {
		if os.IsNotExist(err) {
			ap.SendNotification(&fyne.Notification{
				Title:   "提示",
				Content: "未找到配置文件，请先在设置中完成相关设置并保存",
			})
			return
		} else {
			ap.SendNotification(&fyne.Notification{
				Title:   "错误",
				Content: "打开配置文件时出错，请联系开发者",
			})
			return
		}
	}
	defer file.Close()
	//读取配置文件
	data, err := io.ReadAll(file)
	if err != nil {
		ap.SendNotification(&fyne.Notification{
			Title:   "错误",
			Content: "读取配置文件时出错，请联系开发者",
		})
		return
	}
	//解析为JSON
	err = json.Unmarshal(data, &config)
	if err != nil {
		ap.SendNotification(&fyne.Notification{
			Title:   "错误",
			Content: "解析JSON时出错，请联系开发者",
		})
		return
	}
}

// SaveConfigFile 保存设置到config.json
func (config *UserConfig) SaveConfigFile(ap fyne.App) {
	jsonData, err := json.Marshal(config)
	if err != nil {
		ap.SendNotification(&fyne.Notification{
			Title:   "错误",
			Content: "序列化配置文件时出错，请联系开发者",
		})
		return
	}
	err = os.WriteFile("config.json", jsonData, 0644)
	if err != nil {
		ap.SendNotification(&fyne.Notification{
			Title:   "错误",
			Content: "保存配置文件时出错，请联系开发者",
		})
		return
	}
}
