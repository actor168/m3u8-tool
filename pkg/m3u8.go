package pkg

type M3U8 struct {
	EncryptMethod string
	EncryptURL    string
	EncryptIV     []byte
	EncryptKey    []byte
	TmpURL        string
	URLPrefix     string
	Downloaded    bool
}
