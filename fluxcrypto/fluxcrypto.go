// Standalone library to generate Flux addresses
package fluxcrypto

import (
	"encoding/hex"
	"fmt"

	"github.com/RunOnFlux/fluxgen/base58"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/btcsuite/golangcrypto/ripemd160"
	"github.com/tyler-smith/go-bip39"
)

type FluxWallet struct {
	Mnemonic  string        `json:"mnemonic"`
	HexSeed   string        `json:"hexSeed"`
	Addresses []FluxAddress `json:"addresses"`
	RequestId string        `json:"requestId"`
}

type FluxAddress struct {
	Value      string `json:"value"`
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
}

func getExtendedKeyFromMnemonic(mainnet bool, mnemonicWords string) (*hdkeychain.ExtendedKey, error) {
	var networkCfg chaincfg.Params

	// Switch depending on mainnet or testnet
	if mainnet == true {
		networkCfg = chaincfg.MainNetParams

	} else {
		networkCfg = chaincfg.TestNet3Params
	}

	seed := bip39.NewSeed(mnemonicWords, "")
	hexSeed := hex.EncodeToString(seed)

	hexValue, err := hex.DecodeString(hexSeed)

	if err != nil {
		return nil, err
	}

	masterKey, err := hdkeychain.NewMaster(hexValue, &networkCfg)
	if err != nil {
		return nil, err
	}

	// get m/0'/0/0
	// Hardened key for account 0. ie 0'
	acct0, err := masterKey.Derive(hdkeychain.HardenedKeyStart + 0)
	if err != nil {
		return nil, err
	}

	// External account for 0'
	extAcct0, err := acct0.Derive(0)
	if err != nil {
		return nil, err
	}

	return extAcct0, nil
}

func getAddressFromMnemonic(mainnet bool, mnemonicWords string, position uint32) (FluxAddress, error) {
	var returnValue FluxAddress
	var networkId NetworkId
	var networkCfg chaincfg.Params

	// Switch depending on mainnet or testnet
	if mainnet == true {
		networkId = MainnnetId
		networkCfg = chaincfg.MainNetParams

	} else {
		networkId = TestnetId
		networkCfg = chaincfg.TestNet3Params
	}

	extendedKey, err := getExtendedKeyFromMnemonic(mainnet, mnemonicWords)
	if err != nil {
		return returnValue, err
	}

	key, err := extendedKey.Derive(uint32(position))
	if err != nil {
		return returnValue, err
	}

	// Serialize to compressed key bytes and pkhash
	pk, err := key.ECPubKey()
	if err != nil {
		return returnValue, err
	}
	pkSerialized := pk.SerializeCompressed()
	pkHash := btcutil.Hash160(pkSerialized)

	encodedAddress := base58.CheckEncode(pkHash[:ripemd160.Size], networkId)

	// Get the pubkey and serialise the compressed public key
	privKey, err := key.ECPrivKey()
	if err != nil {
		return returnValue, err
	}

	returnValue.Value = fmt.Sprintf("%s", encodedAddress)
	wif, err := btcutil.NewWIF(privKey, &networkCfg, true)

	if err != nil {
		return returnValue, err
	}

	returnValue.PrivateKey = wif.String()
	returnValue.PublicKey = hex.EncodeToString(privKey.PubKey().SerializeCompressed())

	return returnValue, nil
}

func CreateWallet(mainnet bool, numberOfAddressesToGenerate int) (FluxWallet, error) {
	var wallet FluxWallet
	var numAddresses int
	var networkId NetworkId
	var networkCfg chaincfg.Params

	// Switch depending on mainnet or testnet
	if mainnet == true {
		networkId = MainnnetId
		networkCfg = chaincfg.MainNetParams

	} else {
		networkId = TestnetId
		networkCfg = chaincfg.TestNet3Params
	}

	if numberOfAddressesToGenerate <= 0 {
		numAddresses = 20
	} else if numberOfAddressesToGenerate > 100 {
		numAddresses = 100
	} else {
		numAddresses = numberOfAddressesToGenerate
	}

	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return wallet, err
	}

	m, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return wallet, err
	}

	seed := bip39.NewSeed(m, "")

	wallet.Mnemonic = m
	wallet.HexSeed = hex.EncodeToString(seed)

	extendedKey, err := getExtendedKeyFromMnemonic(mainnet, wallet.Mnemonic)
	if err != nil {
		return wallet, err
	}

	// Derive extended key (repeat this from 0 to number of addresses-1)
	for i := 0; i <= numAddresses-1; i++ {
		var address FluxAddress

		key, err := extendedKey.Derive(uint32(i))
		if err != nil {
			return wallet, err
		}

		// Serialize to compressed key bytes and pkhash
		pk, err := key.ECPubKey()
		if err != nil {
			return wallet, err
		}
		pkSerialized := pk.SerializeCompressed()
		pkHash := btcutil.Hash160(pkSerialized)

		encodedAddress := base58.CheckEncode(pkHash[:ripemd160.Size], networkId)

		// Get the pubkey and serialise the compressed public key
		privKey, err := key.ECPrivKey()
		if err != nil {
			return wallet, err
		}

		address.Value = fmt.Sprintf("%s", encodedAddress)
		wif, err := btcutil.NewWIF(privKey, &networkCfg, true)

		if err != nil {
			return wallet, err
		}

		address.PrivateKey = wif.String()
		address.PublicKey = hex.EncodeToString(privKey.PubKey().SerializeCompressed())

		wallet.Addresses = append(wallet.Addresses, address)
	}

	return wallet, nil
}

func GetWalletFromMnemonic(mainnet bool, mnemonicWords string, position uint32) (FluxWallet, error) {
	var result FluxWallet
	var address FluxAddress

	address, err := getAddressFromMnemonic(mainnet, mnemonicWords, position)

	if err != nil {
		return result, err
	}

	result.Mnemonic = mnemonicWords
	result.Addresses = append(result.Addresses, address)

	return result, nil
}
