package parserUA

import (
	"github.com/ua-parser/uap-go/uaparser"
)

func GetUA(ua string) (*uaparser.Client,error) {
	parser, err := uaparser.New("../common/parserUA/regexes.yaml")
	if err != nil {
		return nil,err
	}
	return parser.Parse(ua), nil
}
