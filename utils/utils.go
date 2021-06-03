package utils

import (
	"fmt"
	"harbored/i18n"
	"io"
	"os"
)

func CopyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}
	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}
	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()
	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func GetKeyName(key string) string {
	switch key {
	case "Left":
		return i18n.T("key:arrowLeft")
	case "Right":
		return i18n.T("key:arrowRight")
	case "Up":
		return i18n.T("key:arrowUp")
	case "Down":
		return i18n.T("key:arrowDown")
	case "Prior":
		return "Page Up"
	case "Next":
		return "Page Down"
	case "Space":
		return i18n.T("key:space")
	case "LeftShift":
		return i18n.T("key:leftShift")
	case "RightShift":
		return i18n.T("key:rightShift")
	case "Return":
		return "Enter"
	default:
		return key
	}
}

var NonprintableKeys = []string{
	"Escape",
	"Return",
	"Tab",
	"BackSpace",
	"Insert",
	"Delete",
	"Right",
	"Left",
	"Down",
	"Up",
	"Prior",
	"Next",
	"Home",
	"End",
	"F1",
	"F2",
	"F3",
	"F4",
	"F5",
	"F6",
	"F7",
	"F8",
	"F9",
	"F10",
	"F11",
	"F12",
	"CapsLock",
	"LeftShift",
	"RightShift",
	"Space",
}

var Keymap = map[string]string{
	// Russian
	"й": "Q",
	"ц": "Window",
	"у": "E",
	"к": "R",
	"е": "T",
	"н": "Y",
	"г": "U",
	"ш": "I",
	"щ": "O",
	"з": "P",
	"х": "[",
	"ъ": "]",
	"ф": "App",
	"ы": "S",
	"в": "D",
	"а": "F",
	"п": "G",
	"р": "H",
	"о": "J",
	"л": "K",
	"д": "L",
	"ж": ";",
	"э": "'",
	"я": "Z",
	"ч": "X",
	"с": "C",
	"м": "V",
	"и": "B",
	"т": "N",
	"ь": "M",
	"б": ",",
	"ю": ".",
}

// Get an index of an element in an array
func StringIndex(vs []string, t string) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}

// Check if an array contains an element
func StringArrayIncludes(vs []string, t string) bool {
	return StringIndex(vs, t) >= 0
}

// Get array elements that pass a condition check
func StringFilter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}
