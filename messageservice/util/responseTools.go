/*
@Time : 2019/5/6 上午9:21 
@Author : zwcui
@Software: GoLand
*/
package util

import (
	"encoding/base64"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"bytes"
	"errors"
	"time"
	"sort"
	"crypto/sha1"
	"io"
	"strings"
	"io/ioutil"
	"strconv"
	"encoding/xml"
	"crypto/rand"
	"net/http"
	"fmt"
	"jingting_server/messageservice/models"
)

const (
	//以下均为公众号管理后台设置项
	token          = "XXXXXXXX"
	appID          = "XXXXXXXXXX"
	encodingAESKey = "XXXXXXXXXXXXXXX"
)


var AesKey []byte


func EncodingAESKey2AESKey(encodingKey string) []byte {
	data, _ := base64.StdEncoding.DecodeString(encodingKey + "=")
	return data
}


func init() {
	//AesKey = EncodingAESKey2AESKey(models.JTThirdPartyPlatformEncodingAESKey)
	AesKey = EncodingAESKey2AESKey(models.JTGZHEncodingAESKey)
}


type TextRequestBody struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   time.Duration
	MsgType      string
	Url          string
	PicUrl       string
	MediaId      string
	ThumbMediaId string
	Content      string
	MsgId        int
	Location_X   string
	Location_Y   string
	Label        string
}


type TextResponseBody struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   CDATAText
	FromUserName CDATAText
	CreateTime   string
	MsgType      CDATAText
	Content      CDATAText
}


type EncryptRequestBody struct {
	XMLName    xml.Name `xml:"xml"`
	ToUserName string
	Encrypt    string
}


type EncryptResponseBody struct {
	XMLName      xml.Name `xml:"xml"`
	Encrypt      CDATAText
	MsgSignature CDATAText
	TimeStamp    string
	Nonce        CDATAText
}


type EncryptResponseBody1 struct {
	XMLName      xml.Name `xml:"xml"`
	Encrypt      string
	MsgSignature string
	TimeStamp    string
	Nonce        string
}


type CDATAText struct {
	Text string `xml:",innerxml"`
}


func MakeSignature(timestamp, nonce string) string {
	//sl := []string{models.JTThirdPartyPlatformSignToken, timestamp, nonce}
	sl := []string{models.JTGZHSignToken, timestamp, nonce}
	sort.Strings(sl)
	s := sha1.New()
	io.WriteString(s, strings.Join(sl, ""))
	return fmt.Sprintf("%x", s.Sum(nil))
}


func MakeMsgSignature(timestamp, nonce, msg_encrypt string) string {
	//sl := []string{models.JTThirdPartyPlatformSignToken, timestamp, nonce, msg_encrypt}
	sl := []string{models.JTGZHSignToken, timestamp, nonce, msg_encrypt}
	sort.Strings(sl)
	s := sha1.New()
	io.WriteString(s, strings.Join(sl, ""))
	return fmt.Sprintf("%x", s.Sum(nil))
}


func ValidateUrl(timestamp, nonce, signatureIn string) bool {
	signatureGen := MakeSignature(timestamp, nonce)
	if signatureGen != signatureIn {
		return false
	}
	return true
}


func ValidateMsg(timestamp, nonce, msgEncrypt, msgSignatureIn string) bool {
	msgSignatureGen := MakeMsgSignature(timestamp, nonce, msgEncrypt)
	if msgSignatureGen != msgSignatureIn {
		return false
	}
	return true
}


func ParseEncryptRequestBody(r *http.Request) *EncryptRequestBody {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil
	}
	//  mlog.AppendObj(nil, "Wechat Message Service: RequestBody--", body)
	requestBody := &EncryptRequestBody{}
	xml.Unmarshal(body, requestBody)
	return requestBody


}


func ParseTextRequestBody(r *http.Request) *TextRequestBody {
	body, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		Logger.Info(err)
		return nil
	}
	requestBody := &TextRequestBody{}
	xml.Unmarshal(body, requestBody)
	return requestBody
}


func Value2CDATA(v string) CDATAText {
	//return CDATAText{[]byte("<![CDATA[" + v + "]]>")}
	return CDATAText{"<![CDATA[" + v + "]]>"}
}


func MakeTextResponseBody(fromUserName, toUserName, content string) ([]byte, error) {
	textResponseBody := &TextResponseBody{}
	textResponseBody.FromUserName = Value2CDATA(fromUserName)
	textResponseBody.ToUserName = Value2CDATA(toUserName)
	textResponseBody.MsgType = Value2CDATA("text")
	textResponseBody.Content = Value2CDATA(content)
	//textResponseBody.CreateTime = strconv.Itoa(int(time.Duration(time.Now().Unix())))
	textResponseBody.CreateTime = strconv.FormatInt(UnixOfBeijingTime(), 10)
	return xml.MarshalIndent(textResponseBody, " ", "  ")
}
func MakeEncryptResponseBody(fromUserName, toUserName, content, nonce, timestamp string) ([]byte, error) {
	encryptBody := &EncryptResponseBody{}


	encryptXmlData, _ := MakeEncryptXmlData(fromUserName, toUserName, timestamp, content)
	encryptBody.Encrypt = Value2CDATA(encryptXmlData)
	encryptBody.MsgSignature = Value2CDATA(MakeMsgSignature(timestamp, nonce, encryptXmlData))
	encryptBody.TimeStamp = timestamp
	encryptBody.Nonce = Value2CDATA(nonce)


	return xml.MarshalIndent(encryptBody, " ", "  ")
}


func MakeEncryptXmlData(fromUserName, toUserName, timestamp, content string) (string, error) {
	textResponseBody := &TextResponseBody{}
	textResponseBody.FromUserName = Value2CDATA(fromUserName)
	textResponseBody.ToUserName = Value2CDATA(toUserName)
	textResponseBody.MsgType = Value2CDATA("text")
	textResponseBody.Content = Value2CDATA(content)
	textResponseBody.CreateTime = timestamp
	body, err := xml.MarshalIndent(textResponseBody, " ", "  ")
	if err != nil {
		return "", errors.New("xml marshal error")
	}


	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, int32(len(body)))
	if err != nil {
		Logger.Info(err, "Binary write err:", err)
	}
	bodyLength := buf.Bytes()


	randomBytes := []byte("abcdefghijklmnop")


	plainData := bytes.Join([][]byte{randomBytes, bodyLength, body, []byte(models.JTThirdPartyPlatformAppId)}, nil)
	cipherData, err := AesEncrypt(plainData, AesKey)
	if err != nil {
		return "", errors.New("AesEncrypt error")
	}
	return base64.StdEncoding.EncodeToString(cipherData), nil
}


// PadLength calculates padding length, from github.com/vgorin/cryptogo
func PadLength(slice_length, blocksize int) (padlen int) {
	padlen = blocksize - slice_length%blocksize
	if padlen == 0 {
		padlen = blocksize
	}
	return padlen
}


//from github.com/vgorin/cryptogo
func PKCS7Pad(message []byte, blocksize int) (padded []byte) {
	// block size must be bigger or equal 2
	if blocksize < 1<<1 {
		panic("block size is too small (minimum is 2 bytes)")
	}
	// block size up to 255 requires 1 byte padding
	if blocksize < 1<<8 {
		// calculate padding length
		padlen := PadLength(len(message), blocksize)


		// define PKCS7 padding block
		padding := bytes.Repeat([]byte{byte(padlen)}, padlen)


		// apply padding
		padded = append(message, padding...)
		return padded
	}
	// block size bigger or equal 256 is not currently supported
	panic("unsupported block size")
}


func AesEncrypt(plainData []byte, aesKey []byte) ([]byte, error) {
	k := len(aesKey)
	if len(plainData)%k != 0 {
		plainData = PKCS7Pad(plainData, k)
	}


	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}


	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}


	cipherData := make([]byte, len(plainData))
	blockMode := cipher.NewCBCEncrypter(block, iv)
	blockMode.CryptBlocks(cipherData, plainData)


	return cipherData, nil
}


func aesDecrypt(cipherData []byte, aesKey []byte) ([]byte, error) {
	k := len(aesKey) //PKCS#7
	if len(cipherData)%k != 0 {
		return nil, errors.New("crypto/cipher: ciphertext size is not multiple of aes key length")
	}


	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}


	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	plainData := make([]byte, len(cipherData))
	blockMode.CryptBlocks(plainData, cipherData)
	return plainData, nil
}


func ValidateAppId(id []byte) bool {
	if string(id) == models.JTThirdPartyPlatformAppId {
		return true
	}
	return false
}


func ParseEncryptTextRequestBody(plainText []byte) (*TextRequestBody, error) {


	// Read length
	buf := bytes.NewBuffer(plainText[16:20])
	var length int32
	binary.Read(buf, binary.BigEndian, &length)


	// appID validation
	appIDstart := 20 + length
	id := plainText[appIDstart : int(appIDstart)+len(models.JTThirdPartyPlatformAppId)]
	if !ValidateAppId(id) {
		Logger.Info("Wechat Message Service: appid is invalid!")
		return nil, errors.New("Appid is invalid")
	}
	Logger.Info("Wechat Message Service: appid validation is ok!")


	textRequestBody := &TextRequestBody{}
	xml.Unmarshal(plainText[20:20+length], textRequestBody)
	return textRequestBody, nil
}


func ParseEncryptResponse(responseEncryptTextBody []byte) {
	textResponseBody := &EncryptResponseBody1{}
	xml.Unmarshal(responseEncryptTextBody, textResponseBody)


	if !ValidateMsg(textResponseBody.TimeStamp, textResponseBody.Nonce, textResponseBody.Encrypt, textResponseBody.MsgSignature) {
		Logger.Info("msg signature is invalid")
		return
	}


	cipherData, err := base64.StdEncoding.DecodeString(textResponseBody.Encrypt)
	if err != nil {
		Logger.Info(err, "Wechat Message Service: Decode base64 error")
		return
	}


	plainText, err := aesDecrypt(cipherData, AesKey)
	if err != nil {
		Logger.Info(err)
		return
	}


	Logger.Info(string(plainText))
}


func DecryptWechatAppletUser(encryptedData string, session_key string, iv string) ([]byte, error) {
	ciphertext, _ := base64.StdEncoding.DecodeString(encryptedData)
	key, _ := base64.StdEncoding.DecodeString(session_key)
	keyBytes := []byte(key)
	block, err := aes.NewCipher(keyBytes) //选择加密算法
	if err != nil {
		return nil, err
	}
	iv_b, _ := base64.StdEncoding.DecodeString(iv)
	blockModel := cipher.NewCBCDecrypter(block, iv_b)
	plantText := make([]byte, len(ciphertext))
	blockModel.CryptBlocks(plantText, ciphertext)
	plantText = PKCS7UnPadding(plantText, block.BlockSize())
	return plantText, nil
}


func PKCS7UnPadding(plantText []byte, blockSize int) []byte {
	length := len(plantText)
	unpadding := int(plantText[length-1])
	return plantText[:(length - unpadding)]
}
