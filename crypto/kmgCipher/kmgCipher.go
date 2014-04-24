package kmgCipher

import (
	"crypto/aes"
	"crypto/sha512"
	//"crypto/sha256"
	"bytes"
	"crypto/cipher"
	"encoding/base64"
	//"crypto/rand"
	"fmt"
)

/*
对称加密,  安全性没有这么简单
 key为任意长度,使用简单
 data为任意长度,使用简单
 使用aes,cbc,32位密码
 输入密码hash使用sha384
 数据padding使用PKCS5Padding
 不会修改输入的数据
*/
func Encrypt(key []byte, data []byte) (output []byte, err error) {
	keyHash := sha512.Sum384(key)
	aseKey := keyHash[:32]
	cbcIv := keyHash[32:]
	block, err := aes.NewCipher(aseKey)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	paddingSize := blockSize - len(data)%blockSize
	afterCbcSize := paddingSize + len(data)
	output = make([]byte, afterCbcSize)
	copy(output, data)
	copy(output[len(data):], bytes.Repeat([]byte{byte(paddingSize)}, paddingSize))
	blockmode := cipher.NewCBCEncrypter(block, cbcIv)
	blockmode.CryptBlocks(output, output)
	return output, nil
}

/*
对称解密,
	不会修改输入的数据
*/
func Decrypt(key []byte, data []byte) (output []byte, err error) {
	keyHash := sha512.Sum384(key)
	aseKey := keyHash[:32]
	cbcIv := keyHash[32:]
	block, err := aes.NewCipher(aseKey)
	if err != nil {
		return nil, err
	}
	if len(data)%block.BlockSize() != 0 {
		return nil, fmt.Errorf("[kmgCipher.Decrypt] input not full blocks")
	}
	output = make([]byte, len(data))
	blockmode := cipher.NewCBCDecrypter(block, cbcIv)
	blockmode.CryptBlocks(output, data)
	paddingSize := int(output[len(output)-1])
	if paddingSize > block.BlockSize() {
		return nil, fmt.Errorf("[kmgCipher.Decrypt] paddingSize out of range")
	}
	beforeCbcSize := len(data) - paddingSize
	return output[:beforeCbcSize], nil
}

func EncryptString(key string, data []byte) (output string, err error) {
	outputByte, err := Encrypt([]byte(key), []byte(data))
	if err != nil {
		return
	}
	return base64.URLEncoding.EncodeToString(outputByte), nil
}
func DecryptString(key string, data string) (output []byte, err error) {
	dataByte, err := base64.URLEncoding.DecodeString(data)
	if err != nil {
		return
	}
	outputByte, err := Decrypt([]byte(key), dataByte)
	if err != nil {
		return
	}
	return outputByte, nil
}

/*
对称加密
 TODO 了解加密实际需求
 key为任意长度,使用简单
 data为任意长度,使用简单
 使用aes,cbc,32位密码
 数据内容验证使用sha256,避免截断攻击
 输入密码hash使用sha384
 数据padding使用PKCS5Padding
 在数据前部添加一定随机性,相同输入,不同输出.避免已知原文攻击
 不会修改输入的数据
*/
/*
func Encrypt(key []byte,data []byte)(output []byte,err error){
	simpleData:=make([]byte,len(data)+40)
	//8个字节随机信息
	_,err = rand.Reader.Read(simpleData[:8])
	if err!=nil{
		return
	}
	//32个字节sha256数据验证
	sumOutput:=sha256.Sum256(data)
	copy(simpleData[8:40],sumOutput[:])
	copy(simpleData[40:],data)
	return SimpleEncrypt(key,simpleData)
}
*/
/*
对称解密,
	不会修改输入的数据
*/
/*
func Decrypt(key []byte,data []byte)(output []byte,err error){
	simpleData,err:=SimpleDecrypt(key,data)
	if err!=nil{
		return
	}
	//sha256数据验证
	sumOutput:=sha256.Sum256(simpleData[40:])
	if bytes.Compare(sumOutput[:],simpleData[8:40])!=0{
		return nil,fmt.Errorf("[kmgCipher.Decrypt] checksum fail!")
	}
	return simpleData[40:],nil
}
*/
