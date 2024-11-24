package service

import (
	"MapPhotoMD/internal/config"
	"MapPhotoMD/mywidget"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rwcarlsen/goexif/exif"
)

// TravelData 旅行记录结构体
type TravelData struct {
	TravelName string               //旅行名称
	TravelDate string               //旅行日期
	InputPath  string               //照片导入路径
	OutputPath string               //MD文件导出路径
	ProIndex   []*mywidget.Property //所有属性控件的指针
}

// location 经纬度结构体
type location struct {
	lat  float64 //纬度
	long float64 //经度
}

// photoData 照片相关数据的结构体
type photoData struct {
	centerLocation    location   //leaflet地图中心坐标
	rawLocation       []location //照片原始经纬度
	convertedLocation []location //高德坐标下的经纬度
	device            []string   //拍摄设备
	date              []string   //拍摄时间
	invalidPhotos     []string   //无法转换的照片
	validPhotos       []string   //可以转换的照片
}

var pData photoData

// 高德地图坐标转化系统请求头
const gaodeapiSiteHead = "https://restapi.amap.com/v3/assistant/coordinate/convert?locations="

// NewTravelData 创建照片数据结构体
func NewTravelData() *TravelData {
	return &TravelData{}
}

// GenerateMD 读取指定导入目录下的照片，在指定导出目录下按用户配置生成旅行记录MD文件夹
func (travelData *TravelData) GenerateMD(cfg *config.UserConfig) []string {
	//旅行记录文件夹根目录
	basePath := filepath.Join(travelData.OutputPath, travelData.TravelName)
	os.MkdirAll(basePath, 0755)

	//获取照片中的位置信息
	travelData.decodeEXIF(cfg)
	//创建旅行记录文件及其文件夹
	travelData.makeTravelNote(basePath, cfg)
	//创建标记点文件及其文件夹
	travelData.makeMarkers(basePath)
	//转存照片
	travelData.movePhoto(basePath, cfg)
	//删除原照片
	travelData.deletePhoto(cfg)
	//返回无法转换的照片名
	return pData.invalidPhotos
}

// makeTravelNote 创建旅行记录MD文件
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

// decodeEXIF 读取照片的EXIF信息，将定位信息转换为高德坐标，并计算地图的中心坐标
func (travelData *TravelData) decodeEXIF(cfg *config.UserConfig) {
	//读取照片的EXIF
	filepath.Walk(travelData.InputPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			fileExt := filepath.Ext(path)
			if fileExt == ".jpg" { //只读取jpg格式的文件
				file, _ := os.Open(path)
				defer file.Close()
				fileName := filepath.Base(path)

				//解码EXIF信息
				x, e := exif.Decode(file)
				if e != nil {
					pData.invalidPhotos = append(pData.invalidPhotos, fileName)
					return nil
				}

				//读取照片经纬度
				raw := location{}
				raw.lat, raw.long, e = x.LatLong()
				if e != nil || raw.lat == 0 || raw.long == 0 {
					pData.invalidPhotos = append(pData.invalidPhotos, fileName)
					return nil
				}
				pData.rawLocation = append(pData.rawLocation, raw)
				pData.validPhotos = append(pData.validPhotos, fileName)

				//读取拍摄日期
				time, e := x.DateTime()
				if e != nil {
					pData.date = append(pData.date, "")
					return nil
				}
				pData.date = append(pData.date, time.Format("2006-01-02 15:04:05"))

				//读取拍摄设备
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
	var totalLat float64
	var totalLong float64
	for _, raw := range pData.rawLocation {
		gaodeApiSite := fmt.Sprintf("%s%v,%v&coordsys=gps&output=json&key=%s", gaodeapiSiteHead, raw.long, raw.lat, cfg.Key)

		resp, err := http.Get(gaodeApiSite)
		if err != nil {
			log.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Failed to read response body: %v", err)
		}

		var respMap map[string]interface{}
		if err := json.Unmarshal(body, &respMap); err != nil {
			log.Fatalf("Failed to unmarshal JSON: %v", err)
		}

		locations, ok := respMap["locations"].(string)
		if !ok {
			log.Fatalf("locations is not a string or is nil: %v", respMap["locations"])
		}

		l := strings.Split(locations, ",")
		tempLat, _ := strconv.ParseFloat(l[1], 64)
		tempLong, _ := strconv.ParseFloat(l[0], 64)
		pData.convertedLocation = append(pData.convertedLocation, location{
			tempLat,
			tempLong,
		})
		totalLat += tempLat
		totalLong += tempLong
	}
	//计算地图中心坐标
	length := float64(len(pData.convertedLocation))
	pData.centerLocation.lat = totalLat / length
	pData.centerLocation.long = totalLong / length
}

// makeMarkers 创建标记点MD文件
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

// movePhoto 转存照片文件到指定目录下（不会删除原照片）
func (t *TravelData) movePhoto(basePath string, cfg *config.UserConfig) {
	if cfg.MovePhoto {
		var copyPath string
		_, err := os.Stat(cfg.PhotoPath)
		if err != nil { //指定目录不存在，则将转存目录改为默认目录（basePath/pictures）
			copyPath = filepath.Join(basePath, "pictures")
			os.MkdirAll(copyPath, 0755)
		} else {
			copyPath = cfg.PhotoPath
		}
		for _, fileName := range pData.validPhotos {
			source, _ := os.Open(filepath.Join(t.InputPath, fileName))
			copy, _ := os.Create(filepath.Join(copyPath, fileName))

			if cfg.PhotoQuality == 100 { //质量为100时不压缩
				io.Copy(copy, source)
				copy.Sync() //刷新缓冲区，确保成功保存
			} else { //压缩图片
				img, _, err := image.Decode(source)
				if err != nil {
					log.Fatalf("图片解码失败: %v", err)
				}

				//JPEG编码选项
				options := jpeg.Options{
					Quality: cfg.PhotoQuality,
				}

				err = jpeg.Encode(copy, img, &options)
				if err != nil {
					log.Fatalf("图片编码失败: %v", err)
				}
			}

			copy.Close()
			source.Close()
		}
	}
}

// deletePhoto 删除原照片
func (t *TravelData) deletePhoto(cfg *config.UserConfig) {
	if cfg.MovePhoto {
		if cfg.DeletePhoto {
			for _, fileName := range pData.validPhotos {
				os.Remove(filepath.Join(t.InputPath, fileName))
			}
		}
	}
}
