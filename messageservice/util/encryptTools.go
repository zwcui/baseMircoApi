package util

import (
	"encoding/base64"
	"crypto/sha256"
	"fmt"
	"errors"
	"crypto/rand"
	"crypto/aes"
	"io"
	"crypto/cipher"
	"bytes"
	"encoding/binary"
	"encoding/xml"
	"jingting_server/messageservice/models"
)

//base64
const (
	Table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
)
var coder = base64.NewEncoding(Table)
func Base64Encode(src []byte) []byte {
	return []byte(coder.EncodeToString(src))
}
func Base64Decode(src []byte) ([]byte, error) {
	return coder.DecodeString(string(src))
}


//获取盐
func saltString() (salt string, error error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", errors.New(err.Error())
	}
	return fmt.Sprintf("%x", b), nil
}

//对密码进行sha256
func EncryptPasswordWithSalt(password, salt string) (hashedPwd string, error error) {
	sha_256 := sha256.New()
	_, err := sha_256.Write([]byte(password + salt))
	if err != nil {
		return "", errors.New(err.Error())
	}
	return fmt.Sprintf("%x", sha_256.Sum(nil)), nil
}

//加密密码 返回加密结果 以及使用的盐
func EncryptPassword(password string) (hashedPwd string, salt string, error error) {
	saltStr, err := saltString()
	if err != nil {
		return "", "", errors.New("server err")
	}

	hashedPd, err := EncryptPasswordWithSalt(password, saltStr)
	if err != nil {
		return "", "", errors.New("server err")
	}
	return hashedPd, saltStr, nil
}

func AesDecrypt(cipherData []byte, aesKey []byte) ([]byte, error) {
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

func ParseEncryptTextRequestBodyToComponentVerifyTicket(plainText []byte) (*models.ComponentVerifyTicketXmlBody, error) {
	//fmt.Println(string(plainText))

	// Read length
	buf := bytes.NewBuffer(plainText[16:20])
	var length int32
	binary.Read(buf, binary.BigEndian, &length)
	//fmt.Println(string(plainText[20 : 20+length]))

	// appID validation
	//appIDstart := 20 + length
	//id := plainText[appIDstart : int(appIDstart)+len(models.AppId)]
	//if !validateAppId(id) {
	//	Logger.Info("Wechat Service: appid is invalid!")
	//	return nil, errors.New("Appid is invalid")
	//}
	//Logger.Info("Wechat Service: appid validation is ok!")

	// xml Decoding
	textRequestBody := &models.ComponentVerifyTicketXmlBody{}
	xml.Unmarshal(plainText[20:20+length], textRequestBody)
	return textRequestBody, nil
}

func ParseEncryptTextRequestBodyToReceiveMessage(plainText []byte) (*models.ReceiveMessageXmlBody, error) {
	//fmt.Println(string(plainText))

	// Read length
	buf := bytes.NewBuffer(plainText[16:20])
	var length int32
	binary.Read(buf, binary.BigEndian, &length)
	//fmt.Println(string(plainText[20 : 20+length]))

	// appID validation
	//appIDstart := 20 + length
	//id := plainText[appIDstart : int(appIDstart)+len(models.AppId)]
	//if !validateAppId(id) {
	//	Logger.Info("Wechat Service: appid is invalid!")
	//	return nil, errors.New("Appid is invalid")
	//}
	//Logger.Info("Wechat Service: appid validation is ok!")

	// xml Decoding
	textRequestBody := &models.ReceiveMessageXmlBody{}
	xml.Unmarshal(plainText[20:20+length], textRequestBody)
	return textRequestBody, nil
}