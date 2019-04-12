package main

import (
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"fmt"
)

func main() {

	//生成私钥
	pri, _ := rsa.GenerateKey(rand.Reader, 1024)

	//生成公钥
	pub := &pri.PublicKey

	plainTxt := []byte("hello world，你好")

	//对原文进行hash散列
	h := md5.New()
	h.Write(plainTxt)
	hashed := h.Sum(nil)

	opts := rsa.PSSOptions{rsa.PSSSaltLengthAuto, crypto.MD5}

	//实现签名

	sign, _ := rsa.SignPSS(rand.Reader, pri, crypto.MD5, hashed, &opts)

	fmt.Println(hex.EncodeToString(sign))

	//通过公钥实现验签
	err := rsa.VerifyPSS(pub, crypto.MD5, hashed, sign, &opts)

	//err 为空 及验签成功
	fmt.Println("err:", err)

}
