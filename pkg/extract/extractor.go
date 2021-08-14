package extract

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/actor168/m3u8-tool/log"
	"github.com/actor168/m3u8-tool/pkg"
)

type Extractor struct {
	M3U8 pkg.M3U8
}

func (e *Extractor) Extract(file string) (*string, error) {
	tmpDir := file[0:strings.LastIndexByte(file, os.PathSeparator)+1] +
		"tmp" + string(os.PathSeparator)
	lines, err := readLines(file)
	if err != nil {
		return nil, err
	}

	keyUrl := ""
	keyCrptoMethod := ""
	keyIv := ""
	for i, line := range lines {
		switch strings.Split(line, ":")[0] {
		case "#EXT-X-KEY":
			// extract key
			kvs := strings.Split(line, ",")
			for _, v := range kvs {
				kv := strings.Split(v, "=")
				switch kv[0] {
				case "URI":
					keyUrl = strings.Trim(kv[1], "\"")
				case "METHOD":
					keyCrptoMethod = kv[1]
				case "IV":
					keyIv = kv[1]
				}
			}
		case "#EXTINF":
			// save tmp file
			err := download(e.M3U8.URLPrefix, lines[i+1], tmpDir)
			if err != nil {
				return nil, err
			}
		}
	}

	e.M3U8 = pkg.M3U8{
		EncryptMethod: keyCrptoMethod,
		EncryptURL:    keyUrl,
		EncryptIV:     keyIv,
	}
	return nil, nil
}

// download
func download(url string, fileName string, path string) error {
	log.LOGGER.Debugf("m3u8 key url: %s", url)
	res, err := http.Get(url)

	if err != nil {
		return err
	}
	defer res.Body.Close()

	f, err := os.Create(path + fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	io.Copy(f, res.Body)
	return nil
}

func (e *Extractor) ToString() {
	fmt.Printf("m3u8 url: %s\nm3u8 key: %s\nm3u8 iv: %s", e.M3U8.EncryptURL,
		e.M3U8.EncryptMethod, e.M3U8.EncryptIV)
}

// readline of file
func readLines(path string) (lines []string, err error) {
	var (
		file   *os.File
		part   []byte
		prefix bool
	)

	if file, err = os.Open(path); err != nil {
		return
	}

	reader := bufio.NewReader(file)
	buffer := bytes.NewBuffer(make([]byte, 1024))

	for {
		if part, prefix, err = reader.ReadLine(); err != nil {
			break
		}
		buffer.Write(part)
		if !prefix {
			lines = append(lines, buffer.String())
			buffer.Reset()
		}
	}
	if err == io.EOF {
		err = nil
	}
	return
}
