# easychars

Based on [saintfish/chardet](https://github.com/saintfish/chardet) and [golang.org/x/text/encoding/](golang.org/x/text/encoding/) , easychars makes it convient to detect the charsets and convert content to UTF-8 encoded.

## Detect Support charset

- Unicode: UTF-8, UTF-16-LE, UTF-16-BE, UTF-32-LE, UTF-32-BE

- Simplified Chinese: GB2312, GBK, GB18030(include GB2312 and GBK)

- Tranditional Chinese: Big5, EUC-TW

- Janpanese: EUC-JP, Shift_JIS, ISO-2022-JP

- Korean: EUC-KR, ISO-2022-KR

- Others: ISO-8859-1, ISO-8859-2, ISO-8859-5, ISO-8859-6, ISO-8859-7, ISO-8859-9, Windows-1250, Windows-1251, Windows-1254, Windows-1255, Windows-1256 ...

for other charsets, try `easychars.ToUtf8WithCharsetName` to test whether it's supported

## Example

```
package main

import (
    "fmt"
    "github.com/HeapStackTree/easychars"
    "os"
)

func ReadAndConvertFile(path string, charsetName string) (contentInUtf8 []byte, res *charset.Result, err error) {
    res = &charset.Result{
        Charset:     "unknown",
        Language:    "unknown",
        Confidence:  0,
        Convertible: false,
    }

    content, err := os.ReadFile(path)
    if err != nil {
        return
    }
    if charsetName == "" {
        contentInUtf8, res, err = easychars.DetectAndConvertToUtf8(content)
    } else {
        contentInUtf8, err = easychars.ToUtf8WithCharsetName(content, charsetName)
        if err == nil {
            res.Charset = charsetName
            res.Confidence = 100
            res.Convertible = true
        }
    }
    return
}

func main() {
    path := "tests/GB2312/_mozilla_bug171813_text.html"

    // use charset name if you are sure about it
    content, res, err := ReadAndConvertFile(path, "")
    if err != nil {
        return
    }

    // jump ascii parts
    var gbkLoc int
    for i, v := range content {
        if v >= 0x7F {
            gbkLoc = i
            break
        }
    }

    fmt.Printf("Path: %s\nCharset: %s\nLanguage: %s\nConfidence: %d\nConvetible: %t\nContent: %s\n", path, res.Charset, res.Language, res.Confidence, res.Convertible, content[gbkLoc:])
    // Ouput should be:
    // Charset: GB-18030
    // Language: zh
    // Confidence: 100
    // Convetible: true
    // Content: 搜狐在线</b></font></a></div> ...
}

```

Check [godoc](https://pkg.go.dev/github.com/HeapStackTree/easychars) for other methods.
