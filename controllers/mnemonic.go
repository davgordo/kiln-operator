package controllers

import (
	"errors"
	"fmt"
	"github.com/lightningnetwork/lnd/aezeed"
	"k8s.io/utils/strings/slices"
	"strings"
)

func initializeMnemonic(mnemonicStr string) (aezeed.Mnemonic, error) {
	mnemonic := aezeed.Mnemonic{}
	mnemonicSlice := strings.Fields(mnemonicStr)

	if len(mnemonicSlice) != 24 {
		return mnemonic, errors.New("mnemonic contains the wrong number of words")
	}

	for i, word := range mnemonicSlice {
		if !slices.Contains(aezeed.DefaultWordList, word) {
			return mnemonic, errors.New(fmt.Sprintf("mnemonic contains an invalid word at index %d", i))
		}
	}

	copy(mnemonic[:], mnemonicSlice)
	return mnemonic, nil
}
