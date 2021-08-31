package bulkWriter

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
)

func WriteArray2Disk(ii map[string]interface{}) (string, string) {
	log.Println("trace: Enter WriteArray2Disk")
	var tmpVal, tmpProp string
	var str, prop string

	for k, v := range ii {
		switch vv := v.(type) {
		case float64:
			switch k[:1] {
			case "-":
				prop += fmt.Sprintf(" %s = &#34;%v&#34;", k[1:], vv)
			case "#":
				str += fmt.Sprintf("%v", v)
			default:
				// float64はエスケープの必要が無いはず
				str += fmt.Sprintf("<%s>%v</%s>", k, v, k)
			}
		case string:
			switch k[:1] {
			case "-":
				prop += fmt.Sprintf(" %s = &#34;%v&#34;", k[1:], vv)
			case "#":
				str += fmt.Sprintf("%v", v)
			default:
				// 普通の文字だけど、タグ内なのでエスケープが必要
				str += fmt.Sprintf("<%s>%v</%s>", k, vv, k)
			}
		case []interface{}:
			tmpVal = ProcArray(vv)
			// str += fmt.Sprintf("<%v>%v</%v>", k, tmpVal, k)
		case map[string]interface{}:
			tmpVal, tmpProp = WriteArray2Disk(vv)
			if tmpProp != "" {
				str += fmt.Sprintf("<%v%v>%v</%v>", k, tmpProp, tmpVal, k)
			} else {
				str += fmt.Sprintf("<%v>%v</%v>", k, tmpVal, k)
			}
		default:
		}
	}
	log.Println("trace: Leave WriteArray2Disk")
	return str, prop
}

func ProcArray(ii []interface{}) string {
	log.Println("trace: Enter ProcArray")
	var str string
	var tmpVal string

	for k, v := range ii {
		switch vv := v.(type) {
		case float64:
			str += fmt.Sprintf("<%v>%v</%v>", k, vv, k)
		case string:
			str += fmt.Sprintf("<%v>%v</%v>", k, vv, k)
		case []interface{}:
			tmpVal = ProcArray(vv)
			str += fmt.Sprintf("<%v>%s</%v>", k, tmpVal, k)
		case map[string]interface{}:
			tmpVal, _ = WriteArray2Disk(vv)
			str += tmpVal
		default:

		}
	}
	log.Println("trace: Leave ProcArray")
	return str
}

func JSONArray2DiskSystem(ii map[string]interface{}) string {
	log.Println("trace: Enter JSONArray2DiskSystem")
	var str string
	str = ""
	for k, v := range ii {
		switch vv := v.(type) {
		case string:
			// 現ディレクトリでファイルを作成し、中に "#omit-xml-declaration": "yes"を記述するように修正する
			str += "let $" + k + " := <VAL>" + vv + "</VAL>\n"
		case float64:
			floatstr := fmt.Sprintf("%f", vv)
			str += "let $" + k + " := <VAL>" + floatstr + "</VAL>\n"
		case []interface{}:
			str += "let $" + k + " := <VAL>" + ProcArray(vv) + "</VAL>\n"
		case map[string]interface{}:
			tmpstr := JSONArray2DiskSystem(vv)
			str = fmt.Sprintf("<%s>%s</%s>", k, tmpstr, k)
		default:
		}
	}
	return str
}

func WriteJson(reader io.Reader) string {
	var iface interface{}
	var str string
	json.NewDecoder(reader).Decode(&iface)
	if iface != nil {
		m := iface.(map[string]interface{})
		str = JSONArray2DiskSystem(m)
	}
	return str
}
