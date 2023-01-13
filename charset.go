package easychars

import (
	"bytes"
	"errors"
	"github.com/saintfish/chardet"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/htmlindex"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/encoding/unicode/utf32"
	"golang.org/x/text/transform"
	"io"
	"strings"
)

// Result contains all the information that charset detecfr gives.
type Result struct {
	// IANA name of the detected charset.
	Charset string
	// IANA name of the detected language. It may be empty for some charsets.
	Language string
	// Confidence of the Result. Scale from 1 to 100. The bigger, the more confident.
	Confidence int
	// a Decoder which can convert the Result.Charset to utf-8, default encoding.Nop.NewDecoder() which won't try to convert the charset.
	Decoder transform.Transformer
	// Whether the charset can be converted by this package
	Convertible bool
}

// alias for transform.Transformer
type Decoder interface {
	transform.Transformer
}

var (
	errInvalidName  = errors.New("easychars: invalid encoding name")
	errUnknown      = errors.New("easychars: unknown Encoding")
	errUnsupported  = errors.New("easychars: this encoding is not supported")
	errWrongDecoder = errors.New("easychars: wrong decoder")
)

// DetectAll returns all chardet.Results which have non-zero Confidence. The Results are sorted by Confidence in descending order.
//
// Same as saintfish/chardet - chardet.NewTextDetector().DetectAll()
// but save matched Decoder in result
func DetectAll(content []byte) (results []*Result, err error) {
	ress, err := chardet.NewTextDetector().DetectAll(content)
	for _, res := range ress {
		result := &Result{
			Charset:     res.Charset,
			Language:    res.Language,
			Confidence:  res.Confidence,
			Decoder:     encoding.Nop.NewDecoder(),
			Convertible: false,
		}
		charset := res.Charset
		if decoder, err := GetDecoderFromCharsetName(charset); err == nil {
			result.Decoder = decoder
			result.Convertible = true
		}
		results = append(results, result)
	}
	return
}

// DetectEncoding return the Result with highest Confidence.
func DetectEncoding(content []byte) (result *Result, err error) {
	if res, err := DetectAll(content); err == nil {
		result = res[0]
	}
	return
}

// Detect and convert content to UTF-8 encoded.
func DetectAndConvertToUtf8(content []byte) (convertedContent []byte, res *Result, err error) {
	convertedContent = content
	err = nil
	res, err = DetectEncoding(content)
	if err != nil {
		return
	}
	if !res.Convertible {
		return
	}

	charsetLower := strings.ToLower(res.Charset)
	switch charsetLower {
	case "", "unknown":
		return
	case "utf-8", "utf8":
		return
	case "gb18030", "gb-18030", "gb 18030", "gbk", "gb2312":
		// Check whether it's valid under GBK rule
		if !isValidGBK(content) {
			res.Confidence = 20
		}
	}

	convertedContent, err = ToUtf8WithDecoder(content, res.Decoder)
	if err != nil {
		return
	}
	return
}

// Get UTF-8 encoded []byte with encoding.Encoding.
func ToUtf8WithEncoding(content []byte, e encoding.Encoding) ([]byte, error) {
	return ToUtf8WithDecoder(content, e.NewDecoder())
}

// Get UTF-8 encoded []byte with Decoder.
func ToUtf8WithDecoder(content []byte, d Decoder) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(content), d)
	decoded, err := io.ReadAll(reader)
	if err != nil {
		// seems that it seldom happens even if the decoder is not correspond to content
		err = errWrongDecoder
		return nil, err
	}
	return decoded, nil
}

// Get UTF-8 encoded []byte with charset name.
//
// # It will return errInvalidName if there is charset name is not valid
//
// or errWrongDecoder if content can't decoded by the correspond Decoder
//
// Charset name reference:
//
// https://encoding.spec.whatwg.org/#names-and-labels
//
// http://www.iana.org/assignments/character-sets/character-sets.xhtml
func ToUtf8WithCharsetName(content []byte, charsetName string) ([]byte, error) {
	decoder, err := GetDecoderFromCharsetName(charsetName)
	if err != nil {
		return content, err
	}
	return ToUtf8WithDecoder(content, decoder)
}

// GetEncodingFromCharsetName return encoding.Encoding for given charset name (case insensitive).
//
// It will return errInvalidName if the package can't find correspond encoding.Encoding.
//
// Charset name reference:
//
// https://encoding.spec.whatwg.org/#names-and-labels
//
// http://www.iana.org/assignments/character-sets/character-sets.xhtml
func GetEncodingFromCharsetName(name string) (e encoding.Encoding, err error) {
	name = strings.TrimSpace(name)
	name = strings.ToLower(name)
	switch name {
	// only gb18030 is valid name in htmlindex and ianaindex
	case "gb-18030", "gb_18030", "gb 18030":
		name = "gb18030"

	// UTF-32 is not listed in ianaindex and html encodings,
	// so manually return correspond encoding.Encoding
	case "utf-32-le", "utf_32_le", "utf-32_le", "utf_32-le", "utf32le", "utf-32le", "utf32-le", "utf_32le", "utf32_le":
		return utf32.UTF32(utf32.LittleEndian, utf32.IgnoreBOM), nil

	case "utf-32-be", "utf_32_be", "utf-32_be", "utf_32-be", "utf32be", "utf-32be", "utf32-be", "utf_32be", "utf32_be":
		return utf32.UTF32(utf32.BigEndian, utf32.IgnoreBOM), nil

	}
	e, err = htmlindex.Get(name)
	if err != nil || e == nil {
		e, err = ianaindex.IANA.Encoding(name)
	}
	if err != nil || e == nil {
		err = errInvalidName
	}
	return
}

// GetDecoderFromCharsetName return Decoder for given charset name (case insensitive).
//
// It will return errInvalidName if the package can't find correspond Decoder.
//
// Reference: http://www.iana.org/assignments/character-sets/character-sets.xhtml
//
// and http://www.iana.org/assignments/character-sets/character-sets.xhtml.
func GetDecoderFromCharsetName(charsetName string) (decoder Decoder, err error) {
	encod, err := GetEncodingFromCharsetName(charsetName)
	if err != nil {
		return
	}
	decoder = encod.NewDecoder()
	return
}

// GetCharsetNameFromEncoding reports the canonical name of the given Encoding.
//
// # It will return errUnknown if e is not associated with a known encoding scheme
//
// or errUnsupported if e is not listed in the mibMap of htmlindex and ianaindex.
//
// Reference: http://www.iana.org/assignments/character-sets/character-sets.xhtml.
func getCharsetNameFromEncoding(e encoding.Encoding) (name string, err error) {
	// in golang.org/x/text/encoding v0.6.0
	// htmlindex and iana index name return "", errUnknown for utf-32
	name, err = ianaindex.IANA.Name(e)
	if err != nil {
		name, err = htmlindex.Name(e)
	}
	if err != nil {
		if strings.Contains(err.Error(), "not supported") {
			err = errUnsupported
		} else {
			err = errUnknown
		}
	}
	return
}
