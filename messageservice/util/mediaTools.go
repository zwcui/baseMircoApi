/*
@Time : 2019/3/4 下午4:33 
@Author : zwcui
@Software: GoLand
*/
package util

import (
	"net/http"
	"bytes"
	"io/ioutil"
	"encoding/json"
	"jingting_server/messageservice/models"
	"errors"
	"mime/multipart"
	"os"
	"io"
	"fmt"
	"strconv"
	"image/jpeg"
	"image"
	"image/color"
	"github.com/golang/freetype"
	"image/draw"
	"github.com/nfnt/resize"
	"github.com/satori/go.uuid"
	"strings"
	"image/png"
)

//上传文件至微信
//mediaType 1为图片（image）、2为语音（voice）、3为视频（video）、4为缩略图（thumb）
func uploadMedia(filePath string, authInfo models.AuthInfo, mediaType int) (mediaId string, err error) {
	media := ""
	if mediaType == 1 {
		media = "image"
	} else if mediaType == 2 {
		media = "voice"
	} else if mediaType == 3 {
		media = "video"
	} else if mediaType == 4 {
		media = "thumb"
	} else {
		return "", errors.New("mediaType 范围1-4")
	}

	mediaId, err = postFile(filePath, media, authInfo)
	if err != nil {
		Logger.Info("err:"+err.Error())
		return "", err
	}

	Logger.Info("uploadMedia mediaId : "+mediaId)

	return mediaId, err
}

//上传网络文件至微信
//mediaType 1为图片（image）、2为语音（voice）、3为视频（video）、4为缩略图（thumb）
func UploadOnlineMedia(fileUrl string, authInfo models.AuthInfo, mediaType int) (mediaId string, err error) {
	currentPath, _ := os.Getwd()
	Logger.Info("current path:", currentPath)
	uuid, _ := uuid.NewV4()
	commonPath := currentPath + "/media/authInfo_"+strconv.FormatInt(authInfo.Id, 10)+"/temp/"
	fileType := SubstrByLength(fileUrl, strings.LastIndex(fileUrl, ".") + 1, len(fileUrl) - strings.LastIndex(fileUrl, "."))
	fileName := uuid.String() + "." + fileType
	Logger.Info("temp filename:"+fileName)

	if !IsExist(commonPath) {
		err := os.MkdirAll(commonPath,os.ModePerm)
		// 创建文件夹
		if err != nil {
			Logger.Info("mkdir failed![%v]\n", err)
			return "", err
		} else {
			Logger.Info("mkdir success!\n")
		}
	}

	var mediaFile *os.File
	Logger.Info("create mediaFile")
	mediaFile, err = os.Create(commonPath + fileName)
	if err != nil {
		Logger.Info("os.Create err:"+err.Error())
		return "", err
	}
	resp, err := http.Get(fileUrl)
	if err != nil {
		Logger.Info("err:"+err.Error())
		return "", err
	}
	defer resp.Body.Close()
	pix, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Logger.Info("err:"+err.Error())
		return "", err
	}
	_, err = io.Copy(mediaFile, bytes.NewReader(pix))
	if err != nil {
		Logger.Info("err:"+err.Error())
		return "", err
	}

	defer mediaFile.Close()

	mediaId, err = uploadMedia(commonPath + fileName, authInfo, mediaType)

	os.Remove(commonPath + fileName)

	return mediaId, err
}


//生成裂变海报
//头像和二维码每次都重新取，防止更新
func GenerateBanner(subscriber models.Subscriber, channel string, activity models.Activity, authInfo models.AuthInfo)(mediaId string, err error){
	currentPath, _ := os.Getwd()
	Logger.Info("current path:", currentPath)

	commonPath := currentPath + "/media/authInfo_"+strconv.FormatInt(authInfo.Id, 10)+"/activity_"+strconv.FormatInt(activity.Id, 10)+"/"
	if !IsExist(commonPath) {
		err := os.MkdirAll(commonPath,os.ModePerm)
		// 创建文件夹
		if err != nil {
			Logger.Info("mkdir failed![%v]\n", err)
			return "", err
		} else {
			Logger.Info("mkdir success!\n")
		}
	}



	//获取海报背景
	bannerBackgroundUrl := activity.Banner
	bannerBackgroundFileName := ""
	uuidParam, _ := uuid.NewV4()
	if strings.HasSuffix(bannerBackgroundUrl, "png") {
		bannerBackgroundFileName = "bannerBackground_"+uuidParam.String()+".png"
	} else {
		bannerBackgroundFileName = "bannerBackground_"+uuidParam.String()+".jpeg"
	}
	Logger.Info("bannerBackground:"+bannerBackgroundUrl)

	var bannerBackgroundFile *os.File
	var bannerBackgroundFileBytes []byte
	if !IsExist(commonPath + bannerBackgroundFileName) {
		Logger.Info("create bannerBackgroundFile")
		bannerBackgroundFile, err = os.Create(commonPath + bannerBackgroundFileName)
		os.Chmod(commonPath + bannerBackgroundFileName, os.ModePerm)
		if err != nil {
			Logger.Info("os.Create err:"+err.Error())
			return "", err
		}
		resp, err := http.Get(bannerBackgroundUrl)
		if err != nil {
			Logger.Info("err:"+err.Error())
			return "", err
		}
		defer resp.Body.Close()
		bannerBackgroundFileBytes, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			Logger.Info("err:"+err.Error())
			return "", err
		}
		_, err = io.Copy(bannerBackgroundFile, bytes.NewReader(bannerBackgroundFileBytes))
		if err != nil {
			Logger.Info("err:"+err.Error())
			return "", err
		}
		bannerBackgroundFile, err = os.Open(commonPath + bannerBackgroundFileName)
		if err != nil {
			Logger.Info("err:"+err.Error())
			return "", err
		}

	} else {
		Logger.Info("open bannerBackgroundFile")
		//bannerBackgroundFile, err = os.Open(commonPath + bannerBackgroundFileName)
		bannerBackgroundFile, err = os.OpenFile(commonPath + bannerBackgroundFileName, os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			Logger.Info("err:"+err.Error())
			return "", err
		}
	}
	defer bannerBackgroundFile.Close()


	//获取微信头像暂存本地
	avatarFileName := "openid_" + subscriber.Openid + "_avatar.jpeg"

	var avatarFile *os.File
	var avatarImg image.Image
	var avatarResizedFileName string
	if subscriber.Headimgurl != "" {
		Logger.Info("create avatarFile")
		avatarFile, err = os.Create(commonPath + avatarFileName)
		if err != nil {
			Logger.Info("os.Create err:"+err.Error())
			return "", err
		}
		resp, err := http.Get(subscriber.Headimgurl)
		if err != nil {
			Logger.Info("err:"+err.Error())
			return "", err
		}
		defer resp.Body.Close()
		pix, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			Logger.Info("err:"+err.Error())
			return "", err
		}
		_, err = io.Copy(avatarFile, bytes.NewReader(pix))
		if err != nil {
			Logger.Info("err:"+err.Error())
			return "", err
		}
		avatarFile, err = os.Open(commonPath + avatarFileName)
		if err != nil {
			Logger.Info("err:"+err.Error())
			return "", err
		}

		defer avatarFile.Close()
		avatarImg, err = jpeg.Decode(bytes.NewReader(pix))
		//avatarImg, err = jpeg.Decode(avatarFile)
		if err != nil {
			Logger.Info("err:"+err.Error())
			return "", err
		}
	}


	//获取二维码暂存本地
	qrcodeqUrl := ""
	qrcodeFileName := "openid_" + subscriber.Openid + "_qrcode.jpeg"

	dataMap := map[string]string{
		"activity": strconv.FormatInt(activity.Id, 10),
		"channel": channel,
		"openid": subscriber.Openid,
	}
	qrcode, err := GenerateWechatQrcodeWithDataMap(dataMap, authInfo, 0)
	if err != nil {
		Logger.Info("GenerateQrcode err:"+err.Error())
		return "", err
	}
	qrcodeqUrl = "https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket=" + qrcode.Ticket

	qrcodeFile, err := os.Create(commonPath + qrcodeFileName)
	if err != nil {
		Logger.Info("os.Create err:"+err.Error())
		return "", err
	}
	defer qrcodeFile.Close()
	resp, err := http.Get(qrcodeqUrl)
	defer resp.Body.Close()
	pix, err := ioutil.ReadAll(resp.Body)
	_, err = io.Copy(qrcodeFile, bytes.NewReader(pix))
	if err != nil {
		Logger.Info("err:"+err.Error())
		return "", err
	}
	qrcodeFile, err = os.Open(commonPath + qrcodeFileName)
	if err != nil {
		Logger.Info("os.Open(commonPath + qrcodeFileName)  "+ commonPath + qrcodeFileName)
		Logger.Info("err:"+err.Error())
		return "", err
	}
	qrcodeImg, err := jpeg.Decode(qrcodeFile)
	if err != nil {
		Logger.Info("err:"+err.Error())
		return "", err
	}

	//二维码等比例缩放
	qrcodeResizedFileName := "openid_" + subscriber.Openid + "_resized_qrcode.jpeg"
	var qrcodeResizedFile *os.File
	//if !isExist(commonPath + qrcodeResizedFileName) {
		Logger.Info("create qrcodeResizedFile")
		qrcodeM := resize.Resize(uint(activity.BannerQrcodeSideLength), 0, qrcodeImg, resize.NearestNeighbor)
		qrcodeOut, err := os.Create(commonPath + qrcodeResizedFileName)
		if err != nil {
			Logger.Info(err)
		}
		defer qrcodeOut.Close()

		jpeg.Encode(qrcodeOut, qrcodeM, nil)

		qrcodeResizedFile, err = os.Open(commonPath + qrcodeResizedFileName)
		if err != nil {
			Logger.Info("err:"+err.Error())
			return "", err
		}
	//} else {
	//	Logger.Info("open qrcodeResizedFile")
	//	qrcodeResizedFile, err = os.Open(commonPath + qrcodeResizedFileName)
	//	if err != nil {
	//		Logger.Info("err:"+err.Error())
	//		return err
	//	}
	//}
	defer qrcodeResizedFile.Close()
	qrcodeResizedImg, err := jpeg.Decode(qrcodeResizedFile)
	if err != nil {
		Logger.Info("err:"+err.Error())
		return "", err
	}




	//生成裂变海报
	//bannerBackgroundImg, err := jpeg.Decode(bytes.NewReader(bannerBackgroundFileBytes))
	bannerBackgroundImg, err := jpeg.Decode(bannerBackgroundFile)
	if err != nil {
		Logger.Info("jpeg.Decode(bannerBackgroundFile)   "+commonPath + bannerBackgroundFileName)
		Logger.Info("err:" + err.Error())
		bannerBackgroundImg, err = png.Decode(bytes.NewReader(bannerBackgroundFileBytes))
		//bannerBackgroundImg, err = png.Decode(bannerBackgroundFile)
		if err != nil {
			Logger.Info("png.Decode(bannerBackgroundFile)  png.Decode(bytes.NewReader(bannerBackgroundFileBytes))  "+commonPath + bannerBackgroundFileName)
			Logger.Info("err:" + err.Error())
			bannerBackgroundImg, err = png.Decode(bannerBackgroundFile)
			if err != nil {
				Logger.Info("png.Decode(bannerBackgroundFile)  png.Decode(bannerBackgroundFile)  " + commonPath + bannerBackgroundFileName)
				Logger.Info("err:" + err.Error())
				return "", err
			}
		}
	}
	bannerFileName := "openid_" + subscriber.Openid + "_banner.jpeg"
	bannerFile, err := os.Create(commonPath + bannerFileName)
	if err != nil {
		Logger.Info("os.Create err:"+err.Error())
		return "", err
	}
	defer bannerFile.Close()

	jpg := image.NewRGBA(image.Rect(0, 0, bannerBackgroundImg.Bounds().Dx(), bannerBackgroundImg.Bounds().Dy()))

	draw.Draw(jpg, jpg.Bounds(), bannerBackgroundImg, bannerBackgroundImg.Bounds().Min, draw.Over)                   //首先将一个图片信息存入jpg
	//draw.DrawMask(jpg, jpg.Bounds(), avatarImg, image.ZP, &circle{image.Pt(activity.BannerAvatarX, activity.BannerAvatarY), activity.BannerAvatarRadius}, image.ZP, draw.Over)
	//draw.DrawMask(jpg, jpg.Bounds(), avatarImg, image.ZP, &circle{image.Pt(100, 100), 90}, image.ZP, draw.Over)
	//draw.Draw(jpg, jpg.Bounds(), avatarImg, image.Point{10, 50}, draw.Over)
	//draw.DrawMask(jpg, jpg.Bounds(), avatarImg, image.ZP, &circle{image.Pt(10, 100), 120}, image.ZP, draw.Over)
	//draw.DrawMask(jpg, jpg.Bounds(), avatarImg, jpg.Bounds().Min.Add(image.Pt(-50, -70)), &circle{image.Pt(90, 100), 120}, image.ZP, draw.Over)
	//draw.DrawMask(jpg, jpg.Bounds(), avatarImg, jpg.Bounds().Min.Add(image.Pt(-100, -200)), &circle{image.Pt(140, 230), 120}, image.ZP, draw.Over)
	//draw.DrawMask(jpg, jpg.Bounds(), avatarImg, jpg.Bounds().Min.Add(image.Pt(avatarX-240, avatarY-430)), &circle{image.Pt(avatarX, avatarY), avatarR}, image.ZP, draw.Over)
	//draw.DrawMask(jpg, jpg.Bounds(), avatarImg, jpg.Bounds().Min.Add(image.Pt(-100, -200)), &circle{image.Pt(avatarX, avatarY), avatarR}, image.ZP, draw.Over)
	//draw.DrawMask(jpg, jpg.Bounds(), avatarImg, jpg.Bounds().Min.Add(image.Pt(-50, -70)), &circle{jpg.Bounds().Min.Add(image.Pt(-50, -70)), 120}, image.ZP, draw.Over)
	//draw.DrawMask(jpg, jpg.Bounds(), avatarJpg, image.ZP, &circle{image.Pt(100, 100), 50}, image.ZP, draw.Over)
	//draw.DrawMask(jpg, image.Rectangle{image.Point{0, 0},image.Point{1080, 1920}}, avatarImg, image.ZP, &circle{image.Pt(10, 100), 120}, image.ZP, draw.Over)

	//if avatarResizedImg != nil {
		//draw.Draw(jpg, jpg.Bounds(), avatarResizedImg, avatarResizedImg.Bounds().Min.Sub(image.Pt(activity.BannerAvatarX, activity.BannerAvatarY)), draw.Over)
	//}

	if avatarImg != nil {
		//如果半径*2超过头像图片的长度，则图片放大
		if activity.BannerAvatarRadius * 2 > avatarImg.Bounds().Dx() {
			var avatarResizedImg image.Image
			var avatarResizedFileName string

			//头像等比例缩放
			avatarResizedFileName = "openid_" + subscriber.Openid + "_resized_avatar.jpeg"
			var avatarResizedFile *os.File
			Logger.Info("create avatarResizedFile")
			avatarM := resize.Resize(uint(avatarImg.Bounds().Dx() *  (activity.BannerAvatarRadius * 2)/avatarImg.Bounds().Dx() ), uint(avatarImg.Bounds().Dy() *  (activity.BannerAvatarRadius * 2)/avatarImg.Bounds().Dy() ), avatarImg, resize.NearestNeighbor)
			avatarOut, err := os.Create(commonPath + avatarResizedFileName)
			if err != nil {
				Logger.Info(err)
			}
			defer avatarOut.Close()

			jpeg.Encode(avatarOut, avatarM, nil)

			avatarResizedFile, err = os.Open(commonPath + avatarResizedFileName)
			if err != nil {
				Logger.Info("err:" + err.Error())
				return "", err
			}

			defer avatarResizedFile.Close()
			avatarResizedImg, err = jpeg.Decode(avatarResizedFile)
			if err != nil {
				Logger.Info("err:" + err.Error())
				return "", err
			}

			draw.DrawMask(jpg, jpg.Bounds(), avatarResizedImg, image.Point{-1*activity.BannerAvatarX, -1*activity.BannerAvatarY}, &circle{image.Point{activity.BannerAvatarX+activity.BannerAvatarRadius, activity.BannerAvatarY+activity.BannerAvatarRadius}, activity.BannerAvatarRadius}, image.Point{0, 0}, draw.Over)
		} else {
			draw.DrawMask(jpg, jpg.Bounds(), avatarImg, image.Point{-1*activity.BannerAvatarX, -1*activity.BannerAvatarY}, &circle{image.Point{activity.BannerAvatarX+activity.BannerAvatarRadius, activity.BannerAvatarY+activity.BannerAvatarRadius}, activity.BannerAvatarRadius}, image.Point{0, 0}, draw.Over)
		}
	}

	draw.Draw(jpg, jpg.Bounds(), qrcodeResizedImg, qrcodeResizedImg.Bounds().Min.Sub(image.Pt(activity.BannerQrcodeX, activity.BannerQrcodeY)), draw.Over)


	//读取字体数据
	fontBytes,err := ioutil.ReadFile(currentPath + "/media/Arial Unicode.ttf")
	if err != nil {
		Logger.Info("err:"+err.Error())
		return "", err
	}
	//载入字体数据
	font,err := freetype.ParseFont(fontBytes)
	if err != nil {
		Logger.Info("err:"+err.Error())
		return "", err
	}
	f := freetype.NewContext()
	//设置分辨率
	f.SetDPI(72)
	//设置字体
	f.SetFont(font)
	//设置尺寸
	f.SetFontSize(float64(activity.BannerNickNameFontSize))
	f.SetClip(bannerBackgroundImg.Bounds())
	//设置输出的图片
	f.SetDst(jpg)
	//设置字体颜色
	if activity.BannerNickNameColor != "" {
		hexStr := ""
		if strings.HasPrefix(activity.BannerNickNameColor, "#") {
			hexStr = activity.BannerNickNameColor[1:]
		} else {
			hexStr = activity.BannerNickNameColor
		}
		hex := HEX{hexStr}
		rgb := hex.hex2rgb()
		f.SetSrc(image.NewUniform(color.RGBA{uint8(rgb.red),uint8(rgb.green),uint8(rgb.blue),255}))
	} else {
		f.SetSrc(image.NewUniform(color.RGBA{0,0,0,255}))
	}

	//设置字体的位置
	pt := freetype.Pt(activity.BannerNickNameX, activity.BannerNickNameY + activity.BannerNickNameFontSize)	//字体是左下角坐标

	_,err = f.DrawString(subscriber.Nickname, pt)
	if err != nil {
		Logger.Info("err:"+err.Error())
		return "", err
	}

	//以png 格式写入文件
	//err = png.Encode(imgfile,img)
	//if err != nil {
	//	util.Logger.Info(err.Error())
	//}


	jpeg.Encode(bannerFile, jpg, nil)



	//上传至微信
	mediaId, err = uploadMedia(commonPath + bannerFileName, authInfo, 1)
	if err != nil {
		Logger.Info("err:"+err.Error())
		return "", err
	}

	//删除服务器文件
	os.Remove(commonPath + avatarFileName)
	if avatarResizedFileName != "" {
		os.Remove(commonPath + avatarResizedFileName)
	}
	os.Remove(commonPath + qrcodeFileName)
	os.Remove(commonPath + qrcodeResizedFileName)
	os.Remove(commonPath + bannerFileName)
	os.Remove(commonPath + bannerBackgroundFileName)

	return mediaId, nil
}

// 判断文件夹是否存在
func IsExist(path string) bool {
	_, err := os.Stat(path)    //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}


//上传文件至微信
func postFile(filePath string, mediaType string, authInfo models.AuthInfo) (mediaId string, err error){
	targetUrl := "https://api.weixin.qq.com/cgi-bin/media/upload?access_token=" + authInfo.AuthAccessToken + "&type=" + mediaType
	//Logger.Info(targetUrl)
	body_buf := bytes.NewBufferString("")
	body_writer := multipart.NewWriter(body_buf)

	// use the body_writer to write the Part headers to the buffer
	_, err = body_writer.CreateFormFile("userfile", filePath)
	if err != nil {
		Logger.Info("error writing to buffer")
		return "", err
	}

	// the file data will be the second part of the body
	fh, err := os.Open(filePath)
	if err != nil {
		Logger.Info("error opening file")
		return "", err
	}
	// need to know the boundary to properly close the part myself.
	boundary := body_writer.Boundary()
	//close_string := fmt.Sprintf("\r\n--%s--\r\n", boundary)
	close_buf := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", boundary))

	// use multi-reader to defer the reading of the file data until
	// writing to the socket buffer.
	request_reader := io.MultiReader(body_buf, fh, close_buf)
	fi, err := fh.Stat()
	if err != nil {
		Logger.Info("Error Stating file: %s", filePath)
		return "", err
	}
	req, err := http.NewRequest("POST", targetUrl, request_reader)
	if err != nil {
		return "", err
	}

	// Set headers for multipart, and Content Length
	req.Header.Add("Content-Type", "multipart/form-data; boundary="+boundary)
	req.ContentLength = fi.Size() + int64(body_buf.Len()) + int64(close_buf.Len())


	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		Logger.Info("resp, err := client.Do(r) err:" + err.Error())
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Logger.Info("body, err := ioutil.ReadAll(resp.Body) err:" + err.Error())
		return "", err
	}

	Logger.Info(string(body))

	response := models.MediaUploadJsonBody{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		Logger.Info("json.Unmarshal(body, &response) err :" + err.Error())
		return "", err
	}


	return response.MediaId, nil
}




type circle struct {
	p image.Point
	r int
}

func (c *circle) ColorModel() color.Model {
	return color.AlphaModel
}

func (c *circle) Bounds() image.Rectangle {
	return image.Rect(c.p.X-c.r, c.p.Y-c.r, c.p.X+c.r, c.p.Y+c.r)
}

func (c *circle) At(x, y int) color.Color {
	xx, yy, rr := float64(x-c.p.X)+0.5, float64(y-c.p.Y)+0.5, float64(c.r)
	if xx*xx+yy*yy < rr*rr {
		return color.Alpha{255}
	}
	return color.Alpha{0}
}



type RGB struct {
	red, green, blue int64
}

type HEX struct {
	str string
}

//如果得到的字符串只有一位，则在其左侧补上一个0
func t2x ( t int64 ) string {
	result := strconv.FormatInt(t, 16)
	if len(result) == 1{
		result = "0" + result
	}
	return result
}

func (color RGB) rgb2hex() HEX {
	r := t2x(color.red)
	g := t2x(color.green)
	b := t2x(color.blue)
	return HEX{r+g+b}
}

func (color HEX) hex2rgb() RGB {
	r, _ := strconv.ParseInt(color.str[:2], 16, 10)
	g, _ := strconv.ParseInt(color.str[2:4], 16, 18)
	b, _ := strconv.ParseInt(color.str[4:], 16, 10)
	return RGB{r,g,b}
}






