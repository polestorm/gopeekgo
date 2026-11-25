// Polecenie bridge to cienka warstwa C nad pakietem preview. Budujemy je jako archiwum
// (go build -buildmode=c-archive) i łączymy z rozszerzeniem Quick Look napisanym w Swifcie,
// więc podgląd powstaje w procesie rozszerzenia, bez uruchamiania innych programów.
// Deklaracje dla Swifta są w ../macos/bridge.h i muszą zgadzać się z tym plikiem.
package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"unsafe"

	"github.com/polestorm/gopeekgo/preview"
)

// GopeekgoPreview zwraca stronę HTML dla źródła Markdown (n bajtów pod src). Długość wyniku trafia
// do outLen, a pamięć wyniku zwalnia wywołujący przez GopeekgoFree.
//
//export GopeekgoPreview
func GopeekgoPreview(src *C.char, n C.int, title *C.char, outLen *C.int) *C.char {
	page := preview.Page(C.GoBytes(unsafe.Pointer(src), n), C.GoString(title))
	*outLen = C.int(len(page))
	return (*C.char)(C.CBytes(page))
}

//export GopeekgoFree
func GopeekgoFree(p *C.char) {
	C.free(unsafe.Pointer(p))
}

func main() {}
