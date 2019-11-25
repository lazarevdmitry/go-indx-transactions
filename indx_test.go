package indx

import (
	"fmt"
	"testing"
)

const INDX_LOGIN string = "login"
const INDX_PASSWORD string = "password"
const INDX_WMID string = "1234567890"
const INDX_CULTURE string = "ru-RU"

func getIndx() Indx {
	i := Indx{
		Login:    INDX_LOGIN,
		Password: INDX_PASSWORD,
		Wmid:     INDX_WMID,
		Culture:  INDX_CULTURE,
	}
	return i
}

func TestBalance(t *testing.T) {
	fmt.Println("\nIndx.Balance() method testing")
	i := getIndx()
	response, err := i.Balance()
	if err != "" {
		t.Log("Error = ", err)
		t.Fail()
	}
	fmt.Println("Test response = ", response)
	return
}

func TestTools(t *testing.T) {
	fmt.Println("\nIndx.Tools() method testing")
	i := getIndx()
	response, err := i.Tools()
	if err != "" {
		t.Log("Error = ", err)
		t.Fail()
	}
	fmt.Println("Test response = ", response)
	return
}

func TestHistoryTrading(t *testing.T) {
	fmt.Println("\nIndx.HistoryTrading() method testing")
	i := getIndx()
	response, err := i.HistoryTrading("12345678", "20190101", "20191101")
	if err != "" {
		t.Log("Error = ", err)
		t.Fail()
	}
	fmt.Println("Test response = ", response)
	return
}

func TestHistoryTransaction(t *testing.T) {
	fmt.Println("\nIndx.HistoryTransaction() method testing")
	i := getIndx()
	response, err := i.HistoryTransaction("12345678", "20190101", "20191101")
	if err != "" {
		t.Log("Error = ", err)
		t.Fail()
	}
	fmt.Println("Test response = ", response)
	return
}

func TestOfferList(t *testing.T) {
	fmt.Println("\nIndx.OfferList() method testing")
	i := getIndx()
	response, err := i.OfferList("12345678")
	if err != "" {
		t.Log("Error = ", err)
		t.Fail()
	}
	fmt.Println("Test response = ", response)
	return
}

func TestOfferAdd(t *testing.T) {
	fmt.Println("\nIndx.OfferAdd() method testing")
	i := getIndx()
	response, err := i.OfferAdd("12345678", "10", true, true, "10")
	if err != "" {
		t.Log("Error = ", err)
		t.Fail()
	}
	fmt.Println("Test response = ", response)
	return
}

func TestOfferDelete(t *testing.T) {
	fmt.Println("\nIndx.OfferDelete() method testing")
	i := getIndx()
	response, err := i.OfferDelete("12345678")
	if err != "" {
		t.Log("Error = ", err)
		t.Fail()
	}
	fmt.Println("Test response = ", response)
	return
}

func TestTick(t *testing.T) {
	fmt.Println("\nIndx.Tick() method testing")
	i := getIndx()
	response, err := i.Tick("12345678", "4")
	if err != "" {
		t.Log("Error = ", err)
		t.Fail()
	}
	fmt.Println("Test response = ", response)
	return
}

func equalitySignatures(pattern, generated string) {
	fmt.Println("Pattern = " + pattern)
	fmt.Println("Generated = " + generated)
	if generated != pattern {
		fmt.Println("Сигнатуры не совпадают")
	} else {
		fmt.Println("Сигнатуры совпадают!!!")
	}
	fmt.Println("")
}
func TestSignature(t *testing.T) {
	var signatures map[string]string = map[string]string{
		"1234567890":          "x3Xnt1ft5jDNCqERO9ECZhqziCnKUqZCKreChi8mhkY=",
		"0987654321":          "F3VjFevUe3EQNZ/HsWgXm/by3zZG/MiIvIqgXHizisE=",
		"aAbBcCdD;0123456789": "DVOgMklYv6GRKsCrTSjvzmRfJa/tII/8Pzx6hYMCJHA=",
	}

	i := Indx{}

	for src, patt := range signatures {
		equalitySignatures(patt, i.signature(src))
	}
}
