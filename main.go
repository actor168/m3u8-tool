package main

import (
	"fmt"
	"os"

	"github.com/actor168/m3u8-tool/pkg/compose"
	"github.com/actor168/m3u8-tool/pkg/crypto"
	"github.com/actor168/m3u8-tool/pkg/extract"
	"github.com/actor168/m3u8-tool/pkg/log"
)

func main() {
	log.Init()
	file := "/mnt/c/Users/chenc/Downloads/1583995470399.m3u8"
	// 文件提取解析模块
	extractor := extract.Extractor{}
	content, err := extractor.Extract(file)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}
	// 加解密模块
	decrypter := crypto.Decryptor{}
	decrypter.Decrypt(content)
	// 合成模块
	suffix := ".dts"
	composer := compose.Composer{}
	success := composer.Compose(file, suffix)
	if !success {
		fmt.Println("Download video failed!")
		os.Exit(1)
	}
	fmt.Println("Download video success!")
	os.Exit(0)
}
