package secretsharing

import (
	"crypto/rand"
	"log"
	"math/big"
)

type Share struct {
	Public, Private, Modulus *big.Int
}

func (s Share) String() string {
	return s.Public.String() + " " + s.Private.String() + " " + s.Modulus.String()
}

type Polynome struct {
	Coefficients []*big.Int
	Modulus      *big.Int
}

func (p Polynome) CalcShare(public *big.Int) Share {
	sum := new(big.Int)
	for power, coefficient := range p.Coefficients {
		tmp := new(big.Int)
		tmp.Exp(public, big.NewInt(int64(power)), p.Modulus)
		tmp.Mul(tmp, coefficient)
		sum.Add(sum, tmp)
		sum.Mod(sum, p.Modulus)
	}
	return Share{public, sum, p.Modulus}
}

func (s Share) CalcPrivateComponent(publics []*big.Int) big.Int {
	prod := new(big.Int)
	prod.Set(s.Private)

	for _, publicxj := range publics {
		if publicxj == s.Public {
			continue
		}
		otherPub := new(big.Int)
		otherPub.Set(publicxj)

		bottom := new(big.Int)
		bottom.Sub(s.Public, otherPub) // x_i - x_j
		bottom.Mod(bottom, s.Modulus)

		bottom.ModInverse(bottom, s.Modulus) // (x_i - x_j)^-1

		otherPub.Neg(otherPub)
		otherPub.Mul(otherPub, bottom)
		otherPub.Mod(otherPub, s.Modulus)

		prod.Mul(prod, otherPub)
		prod.Mod(prod, s.Modulus)
	}

	return *prod
}

func CalcSecret(shares []Share) big.Int {
	publicList := make([]*big.Int, 0)
	for _, share := range shares {
		publicList = append(publicList, share.Public)
	}
	sum := new(big.Int)
	for _, share := range shares {
		wi := share.CalcPrivateComponent(publicList)
		sum.Add(sum, &wi)
	}
	sum.Mod(sum, shares[0].Modulus)
	return *sum
}

func MakePolynome(secret *big.Int, bitsize int, degree int) Polynome {
	p := new(Polynome)
	p.Modulus = new(big.Int)

	rndprime, err := rand.Prime(rand.Reader, bitsize)
	if err != nil {
		log.Fatal(err)
	}
	p.Modulus.Set(rndprime)

	p.Coefficients = append(p.Coefficients, secret)

	for x := 1; x < degree; x++ {
		rnd, err := rand.Int(rand.Reader, p.Modulus)
		if err != nil {
			log.Fatal(err)
		}
		p.Coefficients = append(p.Coefficients, rnd)
	}
	return *p
}
