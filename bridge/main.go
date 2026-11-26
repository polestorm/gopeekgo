// Command bridge is a thin C layer over package preview. It is built as an archive
// (go build -buildmode=c-archive) and linked into the Quick Look extension written in Swift, so the
// preview is made inside the extension process, without starting other programs. The declarations
// for Swift are in macos/bridge.h and must match this file. See macos/build.sh.
package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"unsafe"

	"github.com/polestorm/gopeekgo/preview"
)

// GopeekgoPreview returns the HTML page for the Markdown source of n bytes at src. The length of the
// result is stored in outLen; the caller frees the result with GopeekgoFree.
//
//export GopeekgoPreview
func GopeekgoPreview(src *C.char, n C.int, title *C.char, outLen *C.int) *C.char {
	page := preview.Page(C.GoBytes(unsafe.Pointer(src), n), C.GoString(title))
	*outLen = C.int(len(page))
	return (*C.char)(C.CBytes(page))
}

// GopeekgoFree frees a result of GopeekgoPreview.
//
//export GopeekgoFree
func GopeekgoFree(p *C.char) {
	C.free(unsafe.Pointer(p))
}

func main() {}
