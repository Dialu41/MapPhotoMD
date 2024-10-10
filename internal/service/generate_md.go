package service

import (
	"MapPhotoMD/internal/config"
	"MapPhotoMD/mywidget"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rwcarlsen/goexif/exif"
)

type TravelData struct {
	TravelName string
	TravelDate string
	InputPath  string
	OutputPath string
	ProIndex   []*mywidget.Property
}

type location struct {
	lat  float64 //纬度
	long float64 //经度
}

type photoData struct {
	totalLocation     location   //累加坐标
	centerLocation    location   //地图中心坐标
	rawLocation       []location //原始坐标
	convertedLocation []location //转换后坐标
	device            []string   //拍摄设备
	date              []string   //拍摄时间
	invalidPhotos     []string   //无法转换的照片
	validPhotos       []string   //可以转换的照片
}

var pData photoData

// 高德地图坐标转化系统请求头
const gaodeapiSiteHead = "https://restapi.amap.com/v3/assistant/coordinate/convert?locations="

func NewTravelData() *TravelData {
	return &TravelData{}
}

func (travelData *TravelData) GenerateMD(cfg *config.UserConfig) []string {
	basePath := filepath.Join(travelData.OutputPath, travelData.TravelName)
	os.MkdirAll(basePath, 0755)

	//获取照片中的位置信息
	travelData.decodeEXIF(cfg)
	//创建旅行记录文件及其文件夹
	travelData.makeTravelNote(basePath, cfg)
	//创建标记点文件及其文件夹
	travelData.makeMarkers(basePath)
	//转存照片

	//删除原照片

	//返回无法转换的照片名
	return pData.invalidPhotos
}

// makeTravelNote 创建旅行记录md文件
func (travelData *TravelData) makeTravelNote(basePath string, cfg *config.UserConfig) {
	path := filepath.Join(basePath, travelData.TravelName+".md")
	file, _ := os.Create(path)
	defer file.Close()

	//将设置的属性写入旅行记录MD文件中
	file.WriteString("---\n")
	for _, pro := range travelData.ProIndex {
		//对不同属性类型进行解析和写入
		switch pro.ProType.Selected {
		case mywidget.ProType_List:
			file.WriteString(pro.GetPropertyName() + ": \n")
			items := strings.Split(pro.GetPropertyValue(), ",")
			for _, item := range items {
				file.WriteString("  - " + item + "\n")
			}
		default:
			file.WriteString(pro.GetPropertyName() + ": " + pro.GetPropertyValue() + "\n")
		}
	}
	file.WriteString("---\n\n")

	//写入Leaflet代码块
	file.WriteString("```leaflet\n")
	leafCode := fmt.Sprintf(`id: %s
osmLayer: false
tileServer: http://webrd0{s}.is.autonavi.com/appmaptile?lang=zh_cn&size=1&scale=1&style=8&x={x}&y={y}&z={z}
tileSubdomains: ["1", "2", "3", "4"]
lat: %v
long: %v
height: 500px
width: 100%%
defaultZoom: 16
maxzoom: 18
minzoom: 1
unit: meters
scale: 1
markerFolder: %s/%s/markers
`, travelData.TravelDate, pData.centerLocation.lat, pData.centerLocation.long, cfg.NotePath, travelData.TravelName)
	file.WriteString(leafCode)
	file.WriteString("```\n")
}

// decodeEXIF 读取照片的EXIF信息，将定位信息转换为高德坐标，计算地图的中心坐标
func (travelData *TravelData) decodeEXIF(cfg *config.UserConfig) {
	//读取照片的EXIF
	filepath.Walk(travelData.InputPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			fileExt := filepath.Ext(path)
			if fileExt == ".jpg" {
				file, _ := os.Open(path)
				defer file.Close()
				fileName := filepath.Base(path)

				x, e := exif.Decode(file)
				if e != nil {
					pData.invalidPhotos = append(pData.invalidPhotos, fileName)
					return nil
				}

				raw := location{}
				raw.lat, raw.long, e = x.LatLong()
				if e != nil || raw.lat == 0 || raw.long == 0 {
					pData.invalidPhotos = append(pData.invalidPhotos, fileName)
					return nil
				}
				pData.rawLocation = append(pData.rawLocation, raw)
				pData.validPhotos = append(pData.validPhotos, fileName)

				time, e := x.DateTime()
				if e != nil {
					pData.date = append(pData.date, "")
					return nil
				}
				pData.date = append(pData.date, time.Format("2006-01-02 15:04:05"))

				camModel, e := x.Get(exif.Model)
				if e != nil {
					pData.device = append(pData.device, "")
					return nil
				}
				pData.device = append(pData.device, strings.Trim(camModel.String(), `"`))
			}
		}
		return nil
	})
	//转换坐标
	for _, raw := range pData.rawLocation {
		gaodeApiSite := fmt.Sprintf("%s%v,%v&coordsys=gps&output=json&key=%s", gaodeapiSiteHead, raw.long, raw.lat, cfg.Key)

		resp, err := http.Get(gaodeApiSite)
		if err != nil {

		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		var respMap map[string]interface{}
		json.Unmarshal(body, &respMap)

		l := strings.Split(respMap["locations"].(string), ",")
		tempLat, _ := strconv.ParseFloat(l[1], 64)
		tempLong, _ := strconv.ParseFloat(l[0], 64)
		pData.convertedLocation = append(pData.convertedLocation, location{
			tempLat,
			tempLong,
		})
		pData.totalLocation.lat += tempLat
		pData.totalLocation.long += tempLong
	}
	//计算地图中心坐标
	length := float64(len(pData.convertedLocation))
	pData.centerLocation.lat = pData.totalLocation.lat / length
	pData.centerLocation.long = pData.totalLocation.long / length
}

func (t *TravelData) makeMarkers(basePath string) {
	markerPath := filepath.Join(basePath, "markers")
	os.MkdirAll(markerPath, 0755)
	for i := range pData.validPhotos {
		file, _ := os.Create(filepath.Join(markerPath, fmt.Sprintf("%f,%f", pData.rawLocation[i].lat, pData.rawLocation[i].long)+".md"))

		markerStr := fmt.Sprintf(`---
mapmarker: default
date: %s
device: %s
gps: [%f,%f]
gn: [%f,%f]
location: [%f,%f]
---
![[%s]]`, pData.date[i], pData.device[i],
			pData.rawLocation[i].lat, pData.rawLocation[i].long,
			pData.convertedLocation[i].lat, pData.convertedLocation[i].long,
			pData.convertedLocation[i].lat, pData.convertedLocation[i].long,
			pData.validPhotos[i])

		file.WriteString(markerStr)

		file.Close()
	}
}
