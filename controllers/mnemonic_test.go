package controllers

import (
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"reflect"
)

var _ = Describe("Mnemonic", func() {

	const Pass = "test"
	const Mnemonic = "above pioneer library glimpse exhibit analyst monitor holiday boil art ketchup mail hunt since now pattern vacant arch museum tourist brisk come pilot devote"
	const NotEnoughWords = "above pioneer library glimpse exhibit analyst monitor holiday boil art ketchup mail hunt since now pattern vacant arch museum tourist brisk come pilot"
	const OneInvalidWord = "lnd pioneer library glimpse exhibit analyst monitor holiday boil art ketchup mail hunt since now pattern vacant arch museum tourist brisk come pilot devote"
	const ExpectedRootKey = "xprv9s21ZrQH143K2mhtoUGzSM4Nk8P4oM5CEfmhus3D5fPN6TxDPEtjT8dsBLLdbQFV7kDomWWLYB8M7w8FcAYNomJBKGKKWAtb2WEQcXrtiyY"
	const ExpectedCompressedWIF = "L39ktrwNGvuJCwdf3nMrj8h85w6yH6Z59a7CsNt3LzrJAc4mkr73"
	const ExpectedP2PKHAddress = "18TmnxkpmPsED1D23dpBHkL465s1T1Z4Jx"

	It("Should accept a valid string", func() {

		By("Accepting a string with 24 words")
		seedMnemonic, _ := initializeMnemonic(Mnemonic)
		Expect(reflect.TypeOf(seedMnemonic).Name()).To(Equal("Mnemonic"))
		Expect(len(seedMnemonic)).To(Equal(24))
		Expect(seedMnemonic[0]).To(Equal("above"))
		Expect(seedMnemonic[23]).To(Equal("devote"))

	})

	It("Should reject an invalid string", func() {

		By("Rejecting a string with 23 words")
		_, err := initializeMnemonic(NotEnoughWords)
		Expect(err.Error()).To(Equal("mnemonic contains the wrong number of words"))

		By("Rejecting a string with an invalid word")
		_, err = initializeMnemonic(OneInvalidWord)
		Expect(err.Error()).To(Equal("mnemonic contains an invalid word at index 0"))

	})

	It("Should decrypt to a cipher seed", func() {

		By("providing the password")
		seedMnemonic, err := initializeMnemonic(Mnemonic)
		Expect(err).To(Not(HaveOccurred()))
		pass := []byte(Pass)
		_, err = seedMnemonic.ToCipherSeed(pass)
		Expect(err).To(Not(HaveOccurred()))

	})

	It("Should be able to derive the expected HDKey", func() {

		By("providing network params")
		seedMnemonic, err := initializeMnemonic(Mnemonic)
		Expect(err).To(Not(HaveOccurred()))
		pass := []byte(Pass)
		cipherSeed, err := seedMnemonic.ToCipherSeed(pass)
		Expect(err).To(Not(HaveOccurred()))
		hdkey, err := hdkeychain.NewMaster(cipherSeed.Entropy[:], &chaincfg.MainNetParams)
		Expect(hdkey.String()).To(Equal(ExpectedRootKey))

	})

	It("Should derive the first expected BIP-49 address", func() {

		By("Providing the derivation path?")
		seedMnemonic, err := initializeMnemonic(Mnemonic)
		Expect(err).To(Not(HaveOccurred()))
		pass := []byte(Pass)
		cipherSeed, err := seedMnemonic.ToCipherSeed(pass)
		Expect(err).To(Not(HaveOccurred()))
		hdkey, err := hdkeychain.NewMaster(cipherSeed.Entropy[:], &chaincfg.MainNetParams)
		Expect(hdkey.String()).To(Equal(ExpectedRootKey))

		privateKey, _ := hdkey.ECPrivKey()
		privateKeyBytes, _ := btcec.PrivKeyFromBytes(privateKey.Serialize())
		btcwif, _ := btcutil.NewWIF(privateKeyBytes, &chaincfg.MainNetParams, true)
		Expect(btcwif.String()).To(Equal(ExpectedCompressedWIF))

		// This fails, I'm doing something wrong
		serializedPubKey := btcwif.SerializePubKey()
		witnessProg := btcutil.Hash160(serializedPubKey)
		addressWitnessPubKeyHash, _ := btcutil.NewAddressWitnessPubKeyHash(witnessProg, &chaincfg.MainNetParams)
		serializedScript, _ := txscript.PayToAddrScript(addressWitnessPubKeyHash)
		addressScriptHash, _ := btcutil.NewAddressScriptHash(serializedScript, &chaincfg.MainNetParams)
		segwitNested := addressScriptHash.EncodeAddress()
		Expect(segwitNested).To(Equal(ExpectedP2PKHAddress))

	})
})
