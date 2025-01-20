//go:build darwin || linux

package kiwi

/*
#cgo darwin CFLAGS: -I ./mac/include
#cgo darwin LDFLAGS: -L /Users/yangwoolee/kiwi -lkiwi
#cgo linux CFLAGS: -I ./linux/include
#cgo linux LDFLAGS: -L ./linux/lib -lkiwi
#include <stdlib.h>
#include <string.h>
#include <stdint.h>
#include <kiwi/capi.h>
*/
import "C"
import (
	"log"
	"os"
	"regexp"
	"slices"
)

func Version() {
	version := C.GoString(C.kiwi_version())
	log.Println("Kiwi version:", version)
}

type AnalyzeOption int

const (
	KIWI_MATCH_URL                  AnalyzeOption = C.KIWI_MATCH_URL
	KIWI_MATCH_EMAIL                AnalyzeOption = C.KIWI_MATCH_EMAIL
	KIWI_MATCH_HASHTAG              AnalyzeOption = C.KIWI_MATCH_HASHTAG
	KIWI_MATCH_MENTION              AnalyzeOption = C.KIWI_MATCH_MENTION
	KIWI_MATCH_ALL                  AnalyzeOption = C.KIWI_MATCH_ALL
	KIWI_MATCH_NORMALIZE_CODA       AnalyzeOption = C.KIWI_MATCH_NORMALIZE_CODA
	KIWI_MATCH_ALL_WITH_NORMALIZING AnalyzeOption = C.KIWI_MATCH_ALL_WITH_NORMALIZING
)

// BuildOption is a bitwise OR of the KiwiBuildOption values.
type BuildOption int

const (
	KIWI_BUILD_LOAD_DEFAULT_DICT   BuildOption = C.KIWI_BUILD_LOAD_DEFAULT_DICT
	KIWI_BUILD_INTEGRATE_ALLOMORPH BuildOption = C.KIWI_BUILD_INTEGRATE_ALLOMORPH
	KIWI_BUILD_DEFAULT             BuildOption = C.KIWI_BUILD_DEFAULT
)

// Kiwi is a wrapper for the kiwi C library.
type Kiwi struct {
	handler C.kiwi_h
}

type TokenResult struct {
	Tokens []TokenInfo
	Score  float32
}

// TokenInfo returns the token info for the given token(Str).
type TokenInfo struct {
	// Position is the index of this token appears in the original text.
	Position int

	// Tag represents a type of this token (e.g. VV, NNG, ...).
	Tag POSType

	// Form is the actual string of this token.
	Form string
}

func (k *Kiwi) Analyze(text string, topN int, options AnalyzeOption) ([]string, error) {
	kiwiResH := C.kiwi_analyze(k.handler, C.CString(text), C.int(topN), C.int(options), nil, nil)
	defer C.kiwi_res_close(kiwiResH)

	resSize := int(C.kiwi_res_size(kiwiResH))
	var res []string
	for i := 0; i < resSize; i++ {
		tokenCount := int(C.kiwi_res_word_num(kiwiResH, C.int(i)))
		for j := 0; j < tokenCount; j++ {
			res = append(res, C.GoString(C.kiwi_res_form(kiwiResH, C.int(i), C.int(j))))
		}
	}
	return res, nil
}
func (k *Kiwi) Analyze_Noun(text string, topN int, options AnalyzeOption) ([]string, error) {
	kiwiResH := C.kiwi_analyze(k.handler, C.CString(text), C.int(topN), C.int(options), nil, nil)
	defer C.kiwi_res_close(kiwiResH)

	resSize := int(C.kiwi_res_size(kiwiResH))
	var res []string
	for i := 0; i < resSize; i++ {
		tokenCount := int(C.kiwi_res_word_num(kiwiResH, C.int(i)))
		for j := 0; j < tokenCount; j++ {
			pos, _ := ParsePOSType(C.GoString(C.kiwi_res_tag(kiwiResH, C.int(i), C.int(j))))

			if slices.Contains([]POSType{POS_NNG, POS_NNP, POS_UNKNOWN, POS_NNB, POS_NP, POS_NR, POS_SL}, pos) {
				res = append(res, C.GoString(C.kiwi_res_form(kiwiResH, C.int(i), C.int(j))))
			}

		}
	}
	return k.RemoveDuplicates(res), nil
}
func (k *Kiwi) RemoveSpecialChars(input string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9ㄱ-ㅎㅏ-ㅣ가-힣 ]+`)
	return re.ReplaceAllString(input, " ")
}
func (k *Kiwi) RemoveDuplicates(tokens []string) []string {
	t := make(map[string]bool)
	for _, v := range tokens {
		if !t[v] {
			t[v] = true
		}
	}
	keys := make([]string, 0, len(t))
	for k2 := range t {
		keys = append(keys, k2)
	}
	return keys
}
func (k *Kiwi) Close() int {
	if k.handler != nil {
		out := int(C.kiwi_close(k.handler))
		k.handler = nil
		return out
	}
	return 0
}

// KiwiBuilder is a wrapper for the kiwi C library.
type KiwiBuilder struct {
	handler C.kiwi_builder_h
}

// NewBuilder returns a new KiwiBuilder instance.
// Don't forget to call Close after this.
func NewBuilder(modelPath string, numThread int, options BuildOption) *KiwiBuilder {
	_, err := os.Stat(modelPath)
	if os.IsNotExist(err) {
		log.Panicln("modelPath doesn't exist", modelPath)
	}
	return &KiwiBuilder{
		handler: C.kiwi_builder_init(C.CString(modelPath), C.int(numThread), C.int(options)),
	}
}

// AddWord set custom word with word, pos, score.
func (kb *KiwiBuilder) AddWord(word string, pos POSType, score float32) int {
	return int(C.kiwi_builder_add_word(kb.handler, C.CString(word), C.CString(string(pos)), C.float(score)))
}

// LoadDict loads user dict with dict file path.
func (kb *KiwiBuilder) LoadDict(dictPath string) int {
	return int(C.kiwi_builder_load_dict(kb.handler, C.CString(dictPath)))
}

// Build creates kiwi instance with user word etc.
func (kb *KiwiBuilder) Build() *Kiwi {
	h := C.kiwi_builder_build(kb.handler, nil, C.float(0))
	return &Kiwi{
		handler: h,
	}
}

func (kb *KiwiBuilder) Close() int {
	if kb.handler != nil {
		out := int(C.kiwi_builder_close(kb.handler))
		kb.handler = nil
		return out
	}
	return 0
}
