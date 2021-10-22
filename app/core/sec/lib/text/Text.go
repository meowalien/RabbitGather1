package text

import (
	"fmt"
	"math/rand"
	"regexp"
	"time"
	"unsafe"
)

const AllASCII = " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLM_NOPQRSTUVWXYZ[\\]^`abcdefghijklmnopqrstuvwxyz{|}~"
const AllASCIICount = 95

const (
	Reset ColorCode = iota
	Bold
	Faint
	Italic
	Underline
	BlinkSlow
	BlinkRapid
	ReverseVideo
	Concealed
	CrossedOut
)

// Foreground text colors
const (
	FgBlack ColorCode = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite
)

type ColorCode int

// Foreground Hi-Intensity text colors
const (
	FgHiBlack ColorCode = iota + 90
	FgHiRed
	FgHiGreen
	FgHiYellow
	FgHiBlue
	FgHiMagenta
	FgHiCyan
	FgHiWhite
)

// Background text colors
const (
	BgBlack ColorCode = iota + 40
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgMagenta
	BgCyan
	BgWhite
)

// Background Hi-Intensity text colors
const (
	BgHiBlack ColorCode = iota + 100
	BgHiRed
	BgHiGreen
	BgHiYellow
	BgHiBlue
	BgHiMagenta
	BgHiCyan
	BgHiWhite
)

// ColorSting Color the given string to the given color
func ColorSting(s string, color ColorCode) string {
	return fmt.Sprintf("\033[%dm%s\033[00m", color, s)
}

// ColorSting Color the given string to the given color
func ColorByteSting(s []byte, color ColorCode) []byte {
	return []byte(fmt.Sprintf("\033[%dm%s\033[00m", color, s))
}

var PlainEnglishOnlyRegexp = regexp.MustCompile(fmt.Sprintf(`^[%s%s]+$`, AllASCII, "_"))

var src = rand.NewSource(time.Now().UnixNano())

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

const letterBytes = "ABCDEFGHIJKLM_NOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_"

func randStringBytesMaskImprSrcUnsafe(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

func RandomString(i int) string {
	return randStringBytesMaskImprSrcUnsafe(i)
}
