package bulkWriter

import (
	"io"
	"os"
	"strings"

	"github.com/google/uuid"
)

func readbyte(r io.Reader) (byte, error) {
	var p [1]byte

	n, err := r.Read(p[:])
	if n > 0 {
		return p[0], nil
	}
	return 0, err
}

func createFile(curKey string, curVal string) {
	if curKey == "" && curVal == "" {
		return
	}
	var f *os.File

	f, _ = os.OpenFile("DATA.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)

	defer f.Close()
	var strBuf strings.Builder
	strBuf.Grow(100)

	if curKey != "" {
		strBuf.WriteString(curKey)
		strBuf.WriteString(":")
		strBuf.WriteString(curVal)
		strBuf.WriteString("\n")
	} else {
		strBuf.WriteString(curVal)
		strBuf.WriteString("\n")
	}
	f.WriteString((strBuf.String()))
}

func WriteJson(db string, r io.Reader) {
	var curKey, curVal string
	buf := make([]byte, 0)

	// JSON.Decodeを使わない方向でやる
	var n byte
	var err error
	var inEscape bool
	var inQuote bool // 必要ないスペース判定のために、クォートの中かどうかも判定が必要

	inEscape = false
	inQuote = false

	RootDir := "jsonRoot/" + db
	os.Mkdir(RootDir, 0777)
	os.Chdir(RootDir)

	for {
		n, err = readbyte(r)
		if err != io.EOF {
			switch n {

			case '\\':
				if inEscape {
					// すでにエスケープ中でまた￥が来た場合は、￥のエスケープ
					inEscape = false
				} else {
					// 最初のエスケープ
					inEscape = true
				}
				// エスケープは無視しない（てか出来ない）
				buf = append(buf, n)

			case '"':
				// ダブルクォートは、エスケープされていれば足す
				if inEscape {
					buf = append(buf, n)
					inEscape = false
				} else {
					// エスケープされていないので無視するが、クォート以外の空白を無視するためにInQuote処理をする
					if inQuote {
						// inEscapeではなく、inQuoteの場合はクォートの終わりと判断
						inQuote = false
						// KeyなしかKeyありかにかかわらず、クォートされている文字があった状態だが、ここでは、KeyかValか判定できない
						// Val判定は「,」か「｝」で判定するしかない
					} else {
						inQuote = true
					}
				}
			case ':':
				if inQuote {
					buf = append(buf, n)
				} else {
					// ：まできたら、bufの中身はキーと判定できる
					curKey = string(buf)
					// 足したら、バッファをクリア
					buf = make([]byte, 0)
				}
			case '\t', '\r', '\n', ' ':
				// 空白系の無視
				if inQuote {
					// クォート中なので無視せず足しこむ
					buf = append(buf, n)
				} else {
					// クォート外なので無視する
				}
			case '{', '[':
				if inQuote {
					buf = append(buf, n)
				} else {
					// フォルダを作ってチェンジディレクトリ
					if curKey == "" {
						// curKeyが無い場合は、uuidで代替
						u, _ := uuid.NewRandom()
						curKey = u.String()
						// 同じフォルダにあるDATA.txtにUUIDを書き込む
						createFile("", "-#"+curKey) // 頭に-#を付けてUUIDとする。これで出てきた順序が保存できるはず
					}
					if n == '[' {
						curKey += "[]" // 配列の印をつける
					}
					os.Mkdir(curKey, 0777)
					os.Chdir(curKey)
					curKey = ""
				}
			case '}', ',', ']':
				if inQuote {
					buf = append(buf, n)
				} else {
					if len(buf) > 0 {
						// ここまででbufがあるということはcurValの確定がされていない
						curVal = string(buf)  // まずは確定
						buf = make([]byte, 0) // bufをクリア
					} else {
						// bufが無いという事は、すでに出力済みなので、なにもしない？

					}
					// キーと値でそのディレクトリ内にファイルを出力
					createFile(curKey, curVal)
					curKey = ""
					curVal = ""
					if n == '}' || n == ']' {
						// }か]のときはディレクトリを戻る
						os.Chdir("..")
					}
				}
			default:
				// その他の文字は普通に足す
				buf = append(buf, n)
			}
		} else {
			//		EOF
			break
		}
	}

}
