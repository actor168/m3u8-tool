package crypto

import (
	"bufio"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/actor168/m3u8-tool/pkg"
	log "github.com/actor168/m3u8-tool/pkg/log"
	"github.com/forgoer/openssl"
)

type Decryptor struct {
	M3U8 *pkg.M3U8
}

func (d *Decryptor) Decrypt(suffix *string) {
	// decode key
	key, err := d.fetchKey()
	if err != nil {
		log.LOGGER.Errorf("fetch decrpyte key error: %s", err.Error())
		return
	}
	log.LOGGER.Debugf("%0xd", key)
	d.M3U8.EncryptKey = key
	// decode video slice
	d.DecodeVideoSlice(suffix)
}

func readDir(dirName string) ([]os.DirEntry, error) {
	f, err := os.Open(dirName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dirs, err := f.ReadDir(-1)
	sort.Slice(dirs, func(i, j int) bool {
		iv, _ := strconv.ParseUint(strings.Split(dirs[i].Name(), ".")[0], 10, 64)
		jv, _ := strconv.ParseUint(strings.Split(dirs[j].Name(), ".")[0], 10, 64)
		return iv < jv
	})
	return dirs, nil
}
func (d *Decryptor) DecodeVideoSlice(suffix *string) {
	fileNames, err := readDir(d.M3U8.TmpURL)
	if err != nil {
		log.LOGGER.Error(err.Error())
		return
	}

	for _, v := range fileNames {
		// decode slice video into
		f, err := os.Open(d.M3U8.TmpURL + v.Name())
		if err != nil {
			log.LOGGER.Error(err)
			return
		}
		defer f.Close()
		src, _ := ioutil.ReadAll(f)
		dest, err := openssl.AesCBCDecrypt(src, d.M3U8.EncryptKey, d.M3U8.EncryptIV,
			openssl.PKCS7_PADDING)
		if err != nil {
			log.LOGGER.Error(err)
			return
		}

		destFile, err := os.OpenFile(d.M3U8.TmpURL+".."+string(os.PathSeparator)+*suffix,
			os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.LOGGER.Errorf("append: %s", err.Error())
			return
		}
		write := bufio.NewWriter(destFile)
		write.Write(dest)
		//Flush将缓存的文件真正写入到文件中
		write.Flush()
	}

	os.RemoveAll(d.M3U8.TmpURL)

}

func (d *Decryptor) fetchKey() ([]byte, error) {
	res, err := http.Get(d.M3U8.EncryptURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode/100 > 2 {
		return nil, errors.New("fetch url res error")
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}
