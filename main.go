package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/actor168/m3u8-tool/pkg"
	"github.com/actor168/m3u8-tool/pkg/crypto"
	"github.com/actor168/m3u8-tool/pkg/extract"
	log "github.com/actor168/m3u8-tool/pkg/log"
)

var (
	prefix = flag.String("", "prefix", "download video slice prefix")
	file   = flag.String("", "file", "file path")
)

func main() {
	log.Init()
	m3u8 := &pkg.M3U8{
		URLPrefix:  *prefix,
		Downloaded: true,
	}
	// 文件提取解析模块
	extractor := extract.Extractor{
		M3U8: m3u8,
	}
	_, err := extractor.Extract(*file)
	extractor.ToString()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}
	// 加解密模块
	decrypter := crypto.Decryptor{
		M3U8: m3u8,
	}
	suffix := "video.mp4"
	decrypter.Decrypt(&suffix)
	// 合成模块
	// composer := compose.Composer{}
	// success := composer.Compose(*file, suffix)
	// if !success {
	// 	fmt.Println("Download video failed!")
	// 	os.Exit(1)
	// }
	fmt.Println("Download video success!")
	os.Exit(0)
}
