package util

import (
	"fmt"
	"regexp"
)

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

//var englishLetter = []rune(AllASCII)

var PlainEnglishOnlyRegexp = regexp.MustCompile(fmt.Sprintf(`^[%s%s]+$`, AllASCII, "_"))

//
//func PlainEnglishAndNumberOnly(s string) bool {
//loop: for _, letterRune := range englishLetter {
//	for _, i2 := range s {
//		if i2 == letterRune {
//			continue loop
//		}
//	}
//}
//	return false
//}
//
//func RandStringRunes(n int) string {
//	b := make([]rune, n)
//	for i := range b {
//		b[i] = englishLetter[rand.Intn(len(englishLetter))]
//	}
//	return string(b)
//}
