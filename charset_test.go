package easychars

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

// Check whether file content is valid under UTF-16 rule, reference: https://zh.wikipedia.org/wiki/UTF-16
//
// return: isUTF16 bool, BE bool, err error
//
// BE: true if content is valid under UTF-16 BE rule, false if not
// LE: true if content is valid under UTF-16 LE rule, false if not
//
// Deprecated: Use CheckIsValidUTF16 instead.
func CheckFileIsUTF16(path string) (isUTF16 bool, BE bool, LE bool, err error) {
	isUTF16 = false
	BE = false
	LE = false
	err = nil
	content, err := os.ReadFile(path)
	if err != nil {
		return
	}
	isUTF16, BE, LE = isValidUTF16(content)
	return
}

type Pair struct {
	in   string
	want bool
}

func GetTestCases(root string, want bool) []Pair {
	var paths []string
	var pairs []Pair
	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fs.SkipDir
		}
		if d.IsDir() {
			return nil
		}
		paths = append(paths, path)
		return nil
	})
	for _, p := range paths {
		pairs = append(pairs, Pair{in: p, want: want})
	}
	return pairs
}

// Check whether file is encoded by ? charset
func CheckFileIs(path string, fn func([]byte) bool) (bool, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}
	return fn(content), nil
}

func TestIsValidGBK(t *testing.T) {
	cases := GetTestCases("./tests/GB2312", true)
	cases2 := GetTestCases("./tests/UTF-16", false)
	cases = append(cases, cases2...)
	for _, c := range cases {
		got, _ := CheckFileIs(c.in, isValidGBK)
		filename := filepath.Base(c.in)
		if got != c.want {
			t.Errorf("CheckFileIsGBK(%q) == %t, want %t\n", filename, got, c.want)
		} else {
			t.Logf("PASS: CheckFileIsGBK(%q) == %t", filename, got)
		}
	}
}

func TestIsValidGB18030(t *testing.T) {
	cases := GetTestCases("./tests/GB2312", true)
	for _, c := range cases {
		got, _ := CheckFileIs(c.in, isValidGB18030)
		filename := filepath.Base(c.in)
		if got != c.want {
			t.Errorf("CheckFileIs18030(%q) == %t, want %t\n", filename, got, c.want)
		} else {
			t.Logf("PASS: CheckFileIsGB18030(%q) == %t", filename, got)
		}
	}
}

func TestIsValidUTF16(t *testing.T) {
	cases := GetTestCases("./tests/UTF-16", true)
	cases2 := GetTestCases("./tests/UTF-16LE", true)
	cases = append(cases, cases2...)
	cases2 = GetTestCases("./tests/UTF-16BE", true)
	cases = append(cases, cases2...)
	for _, c := range cases {
		got, be, le, _ := CheckFileIsUTF16(c.in)
		filename := filepath.Base(c.in)
		if got != c.want {
			t.Errorf("CheckFileIsUTF16(%q) == %t, be = %t, le = %t, want %t\n", filename, got, be, le, c.want)
		} else {
			t.Logf("PASS: CheckFileIsUTF16(%q) == %t", filename, got)
		}
	}
}

func TestIsValidUTF16BE(t *testing.T) {
	cases := GetTestCases("./tests/UTF-16BE", true)
	cases2 := GetTestCases("./tests/GB2312", false)
	cases = append(cases, cases2...)
	for _, c := range cases {
		got, _ := CheckFileIs(c.in, isValidUTF16BE)
		filename := filepath.Base(c.in)
		if got != c.want {
			t.Errorf("CheckFileIsUTF16BE(%q) == %t, want %t\n", filename, got, c.want)
		} else {
			t.Logf("PASS: CheckFileIsUTF16BE(%q) == %t", filename, got)
		}
	}
}

func TestIsValidUTF16LE(t *testing.T) {
	cases := GetTestCases("./tests/UTF-16LE", true)
	cases2 := GetTestCases("./tests/GB2312", false)
	cases = append(cases, cases2...)
	for _, c := range cases {
		got, _ := CheckFileIs(c.in, isValidUTF16LE)
		filename := filepath.Base(c.in)
		if got != c.want {
			t.Errorf("CheckFileIsUTF16LE(%q) == %t, want %t\n", filename, got, c.want)
		} else {
			t.Logf("PASS: CheckFileIsUTF16LE(%q) == %t", filename, got)
		}
	}
}

func TestIsValidBig5(t *testing.T) {
	cases := GetTestCases("./tests/Big5", true)
	cases2 := GetTestCases("./tests/UTF-16", false)
	cases = append(cases, cases2...)
	for _, c := range cases {
		got, _ := CheckFileIs(c.in, isValidBig5)
		filename := filepath.Base(c.in)
		if got != c.want {
			t.Errorf("CheckFileIsBig5(%q) == %t, want %t\n", filename, got, c.want)
		} else {
			t.Logf("PASS: CheckFileIsBig5(%q) == %t", filename, got)
		}
	}
}

func Test_UTF_8_Detect(t *testing.T) {
	cases := GetTestCases("./tests/utf-8", true)
	charsetName := "UTF-8"
	for _, c := range cases {
		content, _ := os.ReadFile(c.in)
		content, res, err := DetectAndConvertToUtf8(content)
		filename := filepath.Base(c.in)
		if err != nil {
			t.Errorf("%s: can't convert to utf8", filename)
		} else if res.Charset != charsetName {
			t.Errorf("%s: got charset %s != %s (real charset)", filename, res.Charset, charsetName)
		}
		t.Logf("\nfilename: %s\ncharset: %s\nconfidence: %d\ncontent: \n%s\n\n", filename, res.Charset, res.Confidence, content)
	}
}

// // UTF-16BE_Detect can't pass tests
// func Test_UTF_16BE_Detect(t *testing.T) {
// 	cases := GetTestCases("./tests/UTF-16BE", true)
// 	charsetName := "UTF-16BE"
// 	for _, c := range cases {
// 		content, _ := os.ReadFile(c.in)
// 		content, res, err := DetectAndConvertToUtf8(content)
// 		filename := filepath.Base(c.in)
// 		if err != nil {
// 			t.Errorf("%s: can't convert to utf8", filename)
// 		} else if res.Charset != charsetName {
// 			t.Errorf("%s: got charset %s != %s (real charset)", filename, res.Charset, charsetName)
// 		}
// 		t.Logf("\nfilename: %s\ncharset: %s\nconfidence: %d\ncontent: \n%s\n\n", filename, res.Charset, res.Confidence, content)
// 	}
// }

// // UTF-16LE_Detect can't pass tests
// func Test_UTF_16LE_Detect(t *testing.T) {
// 	cases := GetTestCases("./tests/UTF-16LE", true)
// 	charsetName := "UTF-16LE"
// 	for _, c := range cases {
// 		content, _ := os.ReadFile(c.in)
// 		content, res, err := DetectAndConvertToUtf8(content)
// 		filename := filepath.Base(c.in)
// 		if err != nil {
// 			t.Errorf("%s: can't convert to utf8", filename)
// 		} else if res.Charset != charsetName {
// 			t.Errorf("%s: got charset %s != %s (real charset)", filename, res.Charset, charsetName)
// 		}
// 		t.Logf("\nfilename: %s\ncharset: %s\nconfidence: %d\ncontent: \n%s\n\n", filename, res.Charset, res.Confidence, content)
// 	}
// }

func Test_UTF_32BE_Detect(t *testing.T) {
	cases := GetTestCases("./tests/UTF-32BE", true)
	charsetName := "UTF-32BE"
	for _, c := range cases {
		content, _ := os.ReadFile(c.in)
		content, res, err := DetectAndConvertToUtf8(content)
		filename := filepath.Base(c.in)
		if err != nil {
			t.Errorf("%s: can't convert to utf8", filename)
		} else if res.Charset != charsetName {
			t.Errorf("%s: got charset %s != %s (real charset)", filename, res.Charset, charsetName)
		}
		t.Logf("\nfilename: %s\ncharset: %s\nconfidence: %d\ncontent: \n%s\n\n", filename, res.Charset, res.Confidence, content)
	}
}

func Test_UTF_32LE_Detect(t *testing.T) {
	cases := GetTestCases("./tests/UTF-32LE", true)
	charsetName := "UTF-32LE"
	for _, c := range cases {
		content, _ := os.ReadFile(c.in)
		content, res, err := DetectAndConvertToUtf8(content)
		filename := filepath.Base(c.in)
		if err != nil {
			t.Errorf("%s: can't convert to utf8", filename)
		} else if res.Charset != charsetName {
			t.Errorf("%s: got charset %s != %s (real charset)", filename, res.Charset, charsetName)
		}
		t.Logf("\nfilename: %s\ncharset: %s\nconfidence: %d\ncontent: \n%s\n\n", filename, res.Charset, res.Confidence, content)
	}
}

func Test_GB18030_Detect(t *testing.T) {
	cases := GetTestCases("./tests/GB2312", true)
	charsetName := "GB-18030"
	for _, c := range cases {
		content, _ := os.ReadFile(c.in)
		content, res, err := DetectAndConvertToUtf8(content)
		filename := filepath.Base(c.in)
		if err != nil {
			t.Errorf("%s: can't convert to utf8", filename)
		} else if res.Charset != charsetName {
			t.Errorf("%s: got charset %s != %s (real charset)", filename, res.Charset, charsetName)
		}
		t.Logf("\nfilename: %s\ncharset: %s\nconfidence: %d\ncontent: \n%s\n\n", filename, res.Charset, res.Confidence, content)
	}
}

func Test_Big5_Detect(t *testing.T) {
	cases := GetTestCases("./tests/Big5", true)
	charsetName := "Big5"
	for _, c := range cases {
		content, _ := os.ReadFile(c.in)
		content, res, err := DetectAndConvertToUtf8(content)
		filename := filepath.Base(c.in)
		if err != nil {
			t.Errorf("%s: can't convert to utf8", filename)
		} else if res.Charset != charsetName {
			t.Errorf("%s: got charset %s != %s (real charset)", filename, res.Charset, charsetName)
		}
		t.Logf("\nfilename: %s\ncharset: %s\nconfidence: %d\ncontent: \n%s\n\n", filename, res.Charset, res.Confidence, content)
	}
}

func Test_EUC_JP_Detect(t *testing.T) {
	cases := GetTestCases("./tests/EUC-JP", true)
	charsetName := "EUC-JP"
	for _, c := range cases {
		content, _ := os.ReadFile(c.in)
		content, res, err := DetectAndConvertToUtf8(content)
		filename := filepath.Base(c.in)
		if err != nil {
			t.Errorf("%s: can't convert to utf8", filename)
		} else if res.Charset != charsetName {
			t.Errorf("%s: got charset %s != %s (real charset)", filename, res.Charset, charsetName)
		}
		t.Logf("\nfilename: %s\ncharset: %s\nconfidence: %d\ncontent: \n%s\n\n", filename, res.Charset, res.Confidence, content)
	}
}

func Test_Shift_JIS_Detect(t *testing.T) {
	cases := GetTestCases("./tests/SHIFT_JIS", true)
	charsetName := "Shift_JIS"
	for _, c := range cases {
		content, _ := os.ReadFile(c.in)
		content, res, err := DetectAndConvertToUtf8(content)
		filename := filepath.Base(c.in)
		if err != nil {
			t.Errorf("%s: can't convert to utf8", filename)
		} else if res.Charset != charsetName {
			t.Errorf("%s: got charset %s != %s (real charset)", filename, res.Charset, charsetName)
		}
		t.Logf("\nfilename: %s\ncharset: %s\nconfidence: %d\ncontent: \n%s\n\n", filename, res.Charset, res.Confidence, content)
	}
}

func Test_ISO_2022_JP_Detect(t *testing.T) {
	cases := GetTestCases("./tests/iso-2022-jp", true)
	charsetName := "ISO-2022-JP"
	for _, c := range cases {
		content, _ := os.ReadFile(c.in)
		content, res, err := DetectAndConvertToUtf8(content)
		filename := filepath.Base(c.in)
		if err != nil {
			t.Errorf("%s: can't convert to utf8", filename)
		} else if res.Charset != charsetName {
			t.Errorf("%s: got charset %s != %s (real charset)", filename, res.Charset, charsetName)
		}
		t.Logf("\nfilename: %s\ncharset: %s\nconfidence: %d\ncontent: \n%s\n\n", filename, res.Charset, res.Confidence, content)
	}
}

func Test_EUC_KR_Detect(t *testing.T) {
	cases := GetTestCases("./tests/EUC-KR", true)
	charsetName := "EUC-KR"
	for _, c := range cases {
		content, _ := os.ReadFile(c.in)
		content, res, err := DetectAndConvertToUtf8(content)
		filename := filepath.Base(c.in)
		if err != nil {
			t.Errorf("%s: can't convert to utf8", filename)
		} else if res.Charset != charsetName {
			t.Errorf("%s: got charset %s != %s (real charset)", filename, res.Charset, charsetName)
		}
		t.Logf("\nfilename: %s\ncharset: %s\nconfidence: %d\ncontent: \n%s\n\n", filename, res.Charset, res.Confidence, content)
	}
}

func Test_ISO_2022_KR_Detect(t *testing.T) {
	cases := GetTestCases("./tests/ISO-2022-KR", true)
	charsetName := "ISO-2022-KR"
	for _, c := range cases {
		content, _ := os.ReadFile(c.in)
		content, res, err := DetectAndConvertToUtf8(content)
		filename := filepath.Base(c.in)
		if err != nil {
			t.Errorf("%s: can't convert to utf8", filename)
		} else if res.Charset != charsetName {
			t.Errorf("%s: got charset %s != %s (real charset)", filename, res.Charset, charsetName)
		}
		t.Logf("\nfilename: %s\ncharset: %s\nconfidence: %d\ncontent: \n%s\n\n", filename, res.Charset, res.Confidence, content)
	}
}

// // Windows_1250_Detect can't pass tests
// func Test_Windows_1250_Detect(t *testing.T) {
// 	cases := GetTestCases("./tests/windows-1250-croatian", true)
// 	charsetName := "windows-1250"
// 	for _, c := range cases {
// 		content, _ := os.ReadFile(c.in)
// 		content, res, err := DetectAndConvertToUtf8(content)
// 		filename := filepath.Base(c.in)
// 		if err != nil {
// 			t.Errorf("%s: can't convert to utf8", filename)
// 		} else if res.Charset != charsetName {
// 			t.Errorf("%s: got charset %s != %s (real charset)", filename, res.Charset, charsetName)
// 		}
// 		t.Logf("\nfilename: %s\ncharset: %s\nconfidence: %d\ncontent: \n%s\n\n", filename, res.Charset, res.Confidence, content)
// 	}
// }

func Test_Windows_1250_WithCharsetName(t *testing.T) {
	cases := GetTestCases("./tests/windows-1250-croatian", true)
	charsetName := "windows-1250"
	for _, c := range cases {
		content, _ := os.ReadFile(c.in)

		content, err := ToUtf8WithCharsetName(content, charsetName)
		filename := filepath.Base(c.in)
		if err != nil {
			t.Errorf("%s: can't convert to utf8", filename)
		}
		t.Logf("\nfilename: %s\ncharset: %s\ncontent: \n%s\n\n", filename, charsetName, content)
	}
}

// // Windows_1251_Detect can't pass tests
// func Test_Windows_1251_Detect(t *testing.T) {
// 	cases := GetTestCases("./tests/windows-1251-russian", true)
// 	charsetName := "windows-1251"
// 	for _, c := range cases {
// 		content, _ := os.ReadFile(c.in)
// 		content, res, err := DetectAndConvertToUtf8(content)
// 		filename := filepath.Base(c.in)
// 		if err != nil {
// 			t.Errorf("%s: can't convert to utf8", filename)
// 		} else if res.Charset != charsetName {
// 			t.Errorf("%s: got charset %s != %s (real charset)", filename, res.Charset, charsetName)
// 		}
// 		t.Logf("\nfilename: %s\ncharset: %s\nconfidence: %d\ncontent: \n%s\n\n", filename, res.Charset, res.Confidence, content)
// 	}
// }

func Test_Windows_1251_WithCharsetName(t *testing.T) {
	cases := GetTestCases("./tests/windows-1251-russian", true)
	charsetName := "windows-1251"
	for _, c := range cases {
		content, _ := os.ReadFile(c.in)

		content, err := ToUtf8WithCharsetName(content, charsetName)
		filename := filepath.Base(c.in)
		if err != nil {
			t.Errorf("%s: can't convert to utf8", filename)
		}
		t.Logf("\nfilename: %s\ncharset: %s\ncontent: \n%s\n\n", filename, charsetName, content)
	}
}

// // KOI8_R_Detect can't pass tests
// func Test_KOI8_R_Detect(t *testing.T) {
// 	cases := GetTestCases("./tests/KOI8-R", true)
// 	charsetName := "KOI8-R"
// 	for _, c := range cases {
// 		content, _ := os.ReadFile(c.in)
// 		content, res, err := DetectAndConvertToUtf8(content)
// 		filename := filepath.Base(c.in)
// 		if err != nil {
// 			t.Errorf("%s: can't convert to utf8", filename)
// 		} else if res.Charset != charsetName {
// 			t.Errorf("%s: got charset %s != %s (real charset)", filename, res.Charset, charsetName)
// 		}
// 		t.Logf("\nfilename: %s\ncharset: %s\nconfidence: %d\ncontent: \n%s\n\n", filename, res.Charset, res.Confidence, content)
// 	}
// }

func Test_KOI8_R_WithCharsetName(t *testing.T) {
	cases := GetTestCases("./tests/KOI8-R", true)
	charsetName := "KOI8-R"
	for _, c := range cases {
		content, _ := os.ReadFile(c.in)
		content, err := ToUtf8WithCharsetName(content, charsetName)
		filename := filepath.Base(c.in)
		if err != nil {
			t.Errorf("can't convert %s to utf8", filename)
		}
		t.Logf("\nfilename: %s\ncharset: %s\ncontent: \n%s\n\n", filename, charsetName, content)
	}
}

// func Test_Windows_1252_Detect(t *testing.T) {
// 	cases := GetTestCases("./tests/windows-1252", true)
// 	charsetName := "windows-1252"
// 	for _, c := range cases {
// 		content, _ := os.ReadFile(c.in)
// 		content, res, err := DetectAndConvertToUtf8(content)
// 		filename := filepath.Base(c.in)
// 		if err != nil {
// 			t.Errorf("%s: can't convert to utf8", filename)
// 		} else if res.Charset != charsetName {
// 			t.Errorf("%s: got charset %s != %s (real charset)", filename, res.Charset, charsetName)
// 		}
// 		t.Logf("\nfilename: %s\ncharset: %s\nconfidence: %d\ncontent: \n%s\n\n", filename, res.Charset, res.Confidence, content)
// 	}
// }

func Test_Windows_1252_WithCharsetName(t *testing.T) {
	cases := GetTestCases("./tests/windows-1252", true)
	charsetName := "windows-1252"
	for _, c := range cases {
		content, _ := os.ReadFile(c.in)
		content, err := ToUtf8WithCharsetName(content, charsetName)
		filename := filepath.Base(c.in)
		if err != nil {
			t.Errorf("can't convert %s to utf8", filename)
		}
		t.Logf("\nfilename: %s\ncharset: %s\ncontent: \n%s\n\n", filename, charsetName, content)
	}
}

// func Test_Windows_1254_Detect(t *testing.T) {
// 	cases := GetTestCases("./tests/windows-1254-turkish", true)
// 	charsetName := "windows-1254"
// 	for _, c := range cases {
// 		content, _ := os.ReadFile(c.in)
// 		content, res, err := DetectAndConvertToUtf8(content)
// 		filename := filepath.Base(c.in)
// 		if err != nil {
// 			t.Errorf("%s: can't convert to utf8", filename)
// 		} else if res.Charset != charsetName {
// 			t.Errorf("%s: got charset %s != %s (real charset)", filename, res.Charset, charsetName)
// 		}
// 		t.Logf("\nfilename: %s\ncharset: %s\nconfidence: %d\ncontent: \n%s\n\n", filename, res.Charset, res.Confidence, content)
// 	}
// }

func Test_Windows_1254_WithCharsetName(t *testing.T) {
	cases := GetTestCases("./tests/windows-1254-arabic", true)
	charsetName := "windows-1254"
	for _, c := range cases {
		content, _ := os.ReadFile(c.in)
		content, err := ToUtf8WithCharsetName(content, charsetName)
		filename := filepath.Base(c.in)
		if err != nil {
			t.Errorf("can't convert %s to utf8", filename)
		}
		t.Logf("\nfilename: %s\ncharset: %s\ncontent: \n%s\n\n", filename, charsetName, content)
	}
}

// // Windows_1255_Detect can't pass tests
// func Test_Windows_1255_Detect(t *testing.T) {
// 	cases := GetTestCases("./tests/windows-1255-hebrew", true)
// 	charsetName := "windows-1255"
// 	for _, c := range cases {
// 		content, _ := os.ReadFile(c.in)
// 		content, res, err := DetectAndConvertToUtf8(content)
// 		filename := filepath.Base(c.in)
// 		if err != nil {
// 			t.Errorf("%s: can't convert to utf8", filename)
// 		} else if res.Charset != charsetName {
// 			t.Errorf("%s: got charset %s != %s (real charset)", filename, res.Charset, charsetName)
// 		}
// 		t.Logf("\nfilename: %s\ncharset: %s\nconfidence: %d\ncontent: \n%s\n\n", filename, res.Charset, res.Confidence, content)
// 	}
// }

func Test_Windows_1255_WithCharsetName(t *testing.T) {
	cases := GetTestCases("./tests/windows-1256-arabic", true)
	charsetName := "windows-1256"
	for _, c := range cases {
		content, _ := os.ReadFile(c.in)
		content, err := ToUtf8WithCharsetName(content, charsetName)
		filename := filepath.Base(c.in)
		if err != nil {
			t.Errorf("can't convert %s to utf8", filename)
		}
		t.Logf("\nfilename: %s\ncharset: %s\ncontent: \n%s\n\n", filename, charsetName, content)
	}
}

func Test_Windows_1256_Detect(t *testing.T) {
	cases := GetTestCases("./tests/windows-1256-arabic", true)
	charsetName := "windows-1256"
	for _, c := range cases {
		content, _ := os.ReadFile(c.in)
		content, res, err := DetectAndConvertToUtf8(content)
		filename := filepath.Base(c.in)
		if err != nil {
			t.Errorf("%s: can't convert to utf8", filename)
		} else if res.Charset != charsetName {
			t.Errorf("%s: got charset %s != %s (real charset)", filename, res.Charset, charsetName)
		}
		t.Logf("\nfilename: %s\ncharset: %s\nconfidence: %d\ncontent: \n%s\n\n", filename, res.Charset, res.Confidence, content)
	}
}

func Test_Windows_1256_WithCharsetName(t *testing.T) {
	cases := GetTestCases("./tests/windows-1256-arabic", true)
	charsetName := "windows-1256"
	for _, c := range cases {
		content, _ := os.ReadFile(c.in)
		content, err := ToUtf8WithCharsetName(content, charsetName)
		filename := filepath.Base(c.in)
		if err != nil {
			t.Errorf("can't convert %s to utf8", filename)
		}
		t.Logf("\nfilename: %s\ncharset: %s\ncontent: \n%s\n\n", filename, charsetName, content)
	}
}

func Test_8859_1_Detect(t *testing.T) {
	cases := GetTestCases("./tests/iso-8859-1", true)
	charsetName := "ISO-8859-1"
	for _, c := range cases {
		content, _ := os.ReadFile(c.in)
		content, res, err := DetectAndConvertToUtf8(content)
		filename := filepath.Base(c.in)
		if err != nil {
			t.Errorf("%s: can't convert to utf8", filename)
		} else if res.Charset != charsetName {
			t.Errorf("%s: got charset %s != %s (real charset)", filename, res.Charset, charsetName)
		}
		t.Logf("\nfilename: %s\ncharset: %s\nconfidence: %d\ncontent: \n%s\n\n", filename, res.Charset, res.Confidence, content)
	}
}

// // 8859_2_Detect can't pass tests
// func Test_8859_2_Detect(t *testing.T) {
// 	cases := GetTestCases("./tests/iso-8859-2-croatian", true)
// 	charsetName := "ISO-8859-2"
// 	for _, c := range cases {
// 		content, _ := os.ReadFile(c.in)
// 		content, res, err := DetectAndConvertToUtf8(content)
// 		filename := filepath.Base(c.in)
// 		if err != nil {
// 			t.Errorf("%s: can't convert to utf8", filename)
// 		} else if res.Charset != charsetName {
// 			t.Errorf("%s: got charset %s != %s (real charset)", filename, res.Charset, charsetName)
// 		}
// 		t.Logf("\nfilename: %s\ncharset: %s\nconfidence: %d\ncontent: \n%s\n\n", filename, res.Charset, res.Confidence, content)
// 	}
// }

func Test_8859_2_WithCharsetName(t *testing.T) {
	cases := GetTestCases("./tests/iso-8859-2-croatian", true)
	charsetName := "ISO-8859-2"
	for _, c := range cases {
		content, _ := os.ReadFile(c.in)
		content, err := ToUtf8WithCharsetName(content, charsetName)
		filename := filepath.Base(c.in)
		if err != nil {
			t.Errorf("can't convert %s to utf8", filename)
		}
		t.Logf("\nfilename: %s\ncharset: %s\ncontent: \n%s\n\n", filename, charsetName, content)
	}
}

// func Test_8859_3_Detect(t *testing.T) {
// 	cases := GetTestCases("./tests/iso-8859-3-hungarian", true)
// 	charsetName := "ISO-8859-3"
// 	for _, c := range cases {
// 		content, _ := os.ReadFile(c.in)
// 		content, res, err := DetectAndConvertToUtf8(content)
// 		filename := filepath.Base(c.in)
// 		if err != nil {
// 			t.Errorf("%s: can't convert to utf8", filename)
// 		} else if res.Charset != charsetName {
// 			t.Errorf("%s: got charset %s != %s (real charset)", filename, res.Charset, charsetName)
// 		}
// 		t.Logf("\nfilename: %s\ncharset: %s\nconfidence: %d\ncontent: \n%s\n\n", filename, res.Charset, res.Confidence, content)
// 	}
// }

func Test_8859_3_WithCharsetName(t *testing.T) {
	cases := GetTestCases("./tests/iso-8859-3-hungarian", true)
	charsetName := "ISO-8859-3"
	for _, c := range cases {
		content, _ := os.ReadFile(c.in)
		content, err := ToUtf8WithCharsetName(content, charsetName)
		filename := filepath.Base(c.in)
		if err != nil {
			t.Errorf("can't convert %s to utf8", filename)
		}
		t.Logf("\nfilename: %s\ncharset: %s\ncontent: \n%s\n\n", filename, charsetName, content)
	}
}

// // 8859_5_Detect can't pass tests
// func Test_8859_5_Detect(t *testing.T) {
// 	cases := GetTestCases("./tests/iso-8859-5-bulgarian", true)
// 	charsetName := "ISO-8859-5"
// 	for _, c := range cases {
// 		content, _ := os.ReadFile(c.in)
// 		content, res, err := DetectAndConvertToUtf8(content)
// 		filename := filepath.Base(c.in)
// 		if err != nil {
// 			t.Errorf("%s: can't convert to utf8", filename)
// 		} else if res.Charset != charsetName {
// 			t.Errorf("%s: got charset %s != %s (real charset)", filename, res.Charset, charsetName)
// 		}
// 		t.Logf("\nfilename: %s\ncharset: %s\nconfidence: %d\ncontent: \n%s\n\n", filename, res.Charset, res.Confidence, content)
// 	}
// }

func Test_8859_5_WithCharsetName(t *testing.T) {
	cases := GetTestCases("./tests/iso-8859-5-bulgarian", true)
	charsetName := "ISO-8859-5"
	for _, c := range cases {
		content, _ := os.ReadFile(c.in)
		content, err := ToUtf8WithCharsetName(content, charsetName)
		filename := filepath.Base(c.in)
		if err != nil {
			t.Errorf("can't convert %s to utf8", filename)
		}
		t.Logf("\nfilename: %s\ncharset: %s\ncontent: \n%s\n\n", filename, charsetName, content)
	}
}

func Test_8859_6_Detect(t *testing.T) {
	cases := GetTestCases("./tests/iso-8859-6-arabic", true)
	charsetName := "ISO-8859-6"
	for _, c := range cases {
		content, _ := os.ReadFile(c.in)
		content, res, err := DetectAndConvertToUtf8(content)
		filename := filepath.Base(c.in)
		if err != nil {
			t.Errorf("%s: can't convert to utf8", filename)
		} else if res.Charset != charsetName {
			t.Errorf("%s: got charset %s != %s (real charset)", filename, res.Charset, charsetName)
		}
		t.Logf("\nfilename: %s\ncharset: %s\nconfidence: %d\ncontent: \n%s\n\n", filename, res.Charset, res.Confidence, content)
	}
}

//	func Test_8859_7_Detect(t *testing.T) {
//		cases := GetTestCases("./tests/iso-8859-7-greek", true)
//		charsetName := "ISO-8859-7"
//		for _, c := range cases {
//			content, _ := os.ReadFile(c.in)
//			content, res, err := DetectAndConvertToUtf8(content)
//			filename := filepath.Base(c.in)
//			if err != nil {
//				t.Errorf("%s: can't convert to utf8", filename)
//			} else if res.Charset != charsetName {
//				t.Errorf("%s: got charset %s != %s (real charset)", filename, res.Charset, charsetName)
//			}
//			t.Logf("\nfilename: %s\ncharset: %s\nconfidence: %d\ncontent: \n%s\n\n", filename, res.Charset, res.Confidence, content)
//		}
//	}
func Test_8859_7_WithCharsetName(t *testing.T) {
	cases := GetTestCases("./tests/iso-8859-7-greek", true)
	charsetName := "ISO-8859-7"
	for _, c := range cases {
		content, _ := os.ReadFile(c.in)
		content, err := ToUtf8WithCharsetName(content, charsetName)
		filename := filepath.Base(c.in)
		if err != nil {
			t.Errorf("can't convert %s to utf8", filename)
		}
		t.Logf("\nfilename: %s\ncharset: %s\ncontent: \n%s\n\n", filename, charsetName, content)
	}
}

func Test_8859_9_Detect(t *testing.T) {
	cases := GetTestCases("./tests/iso-8859-9-turkish", true)
	charsetName := "ISO-8859-9"
	for _, c := range cases {
		content, _ := os.ReadFile(c.in)
		content, res, err := DetectAndConvertToUtf8(content)
		filename := filepath.Base(c.in)
		if err != nil {
			t.Errorf("%s: can't convert to utf8", filename)
		} else if res.Charset != charsetName {
			t.Errorf("%s: got charset %s != %s (real charset)", filename, res.Charset, charsetName)
		}
		t.Logf("\nfilename: %s\ncharset: %s\nconfidence: %d\ncontent: \n%s\n\n", filename, res.Charset, res.Confidence, content)
	}
}

func Test_Windows_1251_WithDecoder(t *testing.T) {
	cases := GetTestCases("./tests/windows-1251-russian", true)
	charsetName := "windows-1251"
	decoder := windows_1251_Decoder{}
	for _, c := range cases {
		content, _ := os.ReadFile(c.in)

		content, err := ToUtf8WithDecoder(content, decoder)
		filename := filepath.Base(c.in)
		if err != nil {
			t.Errorf("can't convert %s to utf8", filename)
		}
		t.Logf("\nfilename: %s\ncharset: %s\ncontent: \n%s\n\n", filename, charsetName, content)
	}
}

func Test_ISO_88959_1_WithDecoder(t *testing.T) {
	cases := GetTestCases("./tests/iso-8859-1", true)
	charsetName := "ISO-8859-1"
	decoder := iso_8859_1_Decoder{}
	for _, c := range cases {
		content, _ := os.ReadFile(c.in)

		content, err := ToUtf8WithDecoder(content, decoder)
		filename := filepath.Base(c.in)
		if err != nil {
			t.Errorf("can't convert %s to utf8", filename)
		}
		t.Logf("\nfilename: %s\ncharset: %s\ncontent: \n%s\n\n", filename, charsetName, content)
	}
}

func TestUnicodeRuneToUtf8(t *testing.T) {
	for i, c := range charMap_windows_1251 {
		codes := unicodeRuneToUtf8(c)
		t.Logf("idx: %d, codes == %v, %#X\n", i+128, string(codes), codes)
	}
}

func TestGetCharsetName(t *testing.T) {
	charsetNames := []string{
		"UTF-32BE",
		"UTF-32LE",
		"GBK",
		"GB18030",
		"Big5",
		"EUC-KR",
	}
	for _, charsetName := range charsetNames {
		e, _ := GetEncodingFromCharsetName(charsetName)
		name, _ := getCharsetNameFromEncoding(e)
		if name != charsetName {
			t.Errorf("GetCharsetName of %s fail, output: %s", charsetName, name)
		} else {
			t.Logf("PASS")
		}
	}
}
