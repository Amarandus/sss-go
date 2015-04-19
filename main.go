package main

import (
	"fmt"
	"github.com/Amarandus/sss-go/secretsharing"
	"log"
	"math/big"
	"os"
	"strconv"
)

func printUsage() {
	fmt.Println("Usage: " + os.Args[0] + " [join|gen] n [k]")
}

func gen(n int, k int, secret string) {
	byteArray := []byte(secret)
	bits := len(byteArray)*8 + 1

	secretBigint := new(big.Int)
	secretBigint.SetBytes(byteArray)

	p := secretsharing.MakePolynome(secretBigint, bits, k)

	for x := 1; x <= n; x++ {
		pub := *big.NewInt(int64(x))
		s := p.CalcShare(&pub)
		fmt.Println(s)
	}
}

func join(n int) {
	shares := make([]secretsharing.Share, 0)
	for i := 0; i < n; i++ {
		var pub, priv, mod string
		_, err := fmt.Scanf("%s %s %s", &pub, &priv, &mod)
		if err != nil {
			log.Fatal(err)
		}
		pubInt := new(big.Int)
		privInt := new(big.Int)
		modInt := new(big.Int)

		pubInt.SetString(pub, 10)
		privInt.SetString(priv, 10)
		modInt.SetString(mod, 10)

		shares = append(shares, secretsharing.Share{pubInt, privInt, modInt})
	}
	secret := secretsharing.CalcSecret(shares)

	fmt.Println(string(secret.Bytes()))
}

func main() {
	if len(os.Args) <= 2 {
		printUsage()
		return
	}

	if os.Args[1] == "gen" {
		if len(os.Args) != 4 {
			printUsage()
			return
		}

		n, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}

		k, err := strconv.Atoi(os.Args[3])
		if err != nil {
			log.Fatal(err)
		}

		var secret string
		_, err = fmt.Scanf("%s", &secret)
		if err != nil {
			log.Fatal(err)
		}
		gen(n, k, secret)
	}
	if os.Args[1] == "join" {
		if len(os.Args) != 3 {
			printUsage()
			return
		}

		n, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
		join(n)

	}
}
