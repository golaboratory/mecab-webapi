package mecab

// #cgo CFLAGS: -I/usr/local/include
// #cgo LDFLAGS: -L/usr/local/lib -lmecab -lstdc++
// #include <mecab.h>
// #include <stdio.h>
import "C"

import "strings"

/* --------------------------------- */
/* --- error define                  */
/* --------------------------------- */
type mecabError struct {
	message string
}

func (err mecabError) Error() string {
	return err.message
}

/* --------------------------------- */
/* --- mecab parse                   */
/* --------------------------------- */
func Parse(line string, arg string) (result []string, err error) {
	var parse []string
	sentence := C.CString(line)

	model := C.mecab_model_new2(C.CString(arg))
	if model == nil {
		return parse, mecabError{"mecab model is not created."}
	}
	mecab := C.mecab_model_new_tagger(model)
	if mecab == nil {
		return parse, mecabError{"mecab tagger is not created."}
	}
	lattice := C.mecab_model_new_lattice(model)
	if lattice == nil {
		return parse, mecabError{"mecab lattice is not created."}
	}

	C.mecab_lattice_set_sentence(lattice, sentence)
	C.mecab_parse_lattice(mecab, lattice)
	
	lines := strings.Split(C.GoString(C.mecab_lattice_tostr(lattice)), "\n")

	for _, item := range lines {
		if strings.Index(item, "EOS") != 0  {
			if len(item) > 1 {
				parse = append(parse, item)
			}
		}
	}

	C.mecab_destroy(mecab)
	C.mecab_lattice_destroy(lattice)
	C.mecab_model_destroy(model)

	return parse, nil
}

