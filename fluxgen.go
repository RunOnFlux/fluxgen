package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/RunOnFlux/fluxgen/fluxcrypto"
)

func main() {
	//	var networkId fluxcrypto.NetworkId
	boolPtr := flag.Bool("test", false, "generate a testnet wallet")
	nPtr := flag.Int("n", 1, "Number of addresses to generate up to 100")
	boolPtr3 := flag.Bool("o", false, "enable output to file outputfluxgen.txt")
	flag.Parse()

	var output bool = *boolPtr3

	// Generate the wallet
	wallet, err := fluxcrypto.CreateWallet(!(*boolPtr), *nPtr)

	if err != nil {
		log.Panicln(err.Error())
	}

	log.Println("Wallet generated!")
	fmt.Println("Mnemonic:", wallet.Mnemonic)
	fmt.Println("Address\t\t\t\tPrivate key")

	file, err := os.OpenFile("outputfluxgen.txt", os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil && output {
		fmt.Println("File does not exists or cannot be created")
		os.Exit(1)
	}

	w := bufio.NewWriter(file)

	if output {
		fmt.Fprintln(w, "Mnemonic:", wallet.Mnemonic)
		fmt.Fprintln(w, "Address\t\t\t\t\t\t\t\tPrivate key")
		w.Flush()
	}

	for i := 0; i <= len(wallet.Addresses)-1; i++ {
		fmt.Println(wallet.Addresses[i].Value, wallet.Addresses[i].PrivateKey)
		if output {
			fmt.Fprintln(w, wallet.Addresses[i].Value, wallet.Addresses[i].PrivateKey)
			w.Flush()
		}
	}
}
