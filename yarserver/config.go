package yarserver

const Version = "1.1.0"
const packageName = "yarserver"
const headerSize int = 90
const magicNum uint32 = 0x80DFEC60
const packager = "JSON"

type header struct {
	Id       uint32 `pack:"N"`
	Version  uint16 `pack:"n"`
	MagicNum uint32 `pack:"N"`
	Reserved uint32 `pack:"N"`
	Provider string `pack:"a32"`
	Token    string `pack:"a32"`
	BodyLen  uint32 `pack:"N"`
	Packager string `pack:"a8"`
}

type body struct {
	Id     uint32        `json:"i"`
	Method string        `json:"m"`
	Param  []interface{} `json:"p"`
}
