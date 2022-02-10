package models

import "encoding/xml"

type OneObmen struct {
	XMLName xml.Name       `xml:"rates"`
	Rates   []OneObmenItem `xml:"item"`
}

type OneObmenItem struct {
	XMLName   xml.Name `xml:"item"`
	From      string   `xml:"from"`
	To        string   `xml:"to"`
	In        float64  `xml:"in"`
	Out       float64  `xml:"out"`
	Amount    float64  `xml:"amount"`
	MinAmount string   `xml:"minamount"`
	MaxAmount string   `xml:"maxamount"`
}

type Coin struct {
	ID        int
	Name      string
	ShortName string
	Fiat      bool
}

type Direction struct {
	ID   int
	From string
	To   string
}

var COINS = []*Coin{
	{ID: 1, Name: "Bitcoin", ShortName: "BTC", Fiat: false},
	{ID: 2, Name: "Ethereum", ShortName: "ETH", Fiat: false},
	{ID: 3, Name: "Рубль", ShortName: "SBERRUB", Fiat: true},
	{ID: 4, Name: "Tether ERC-20", ShortName: "USDT", Fiat: false},
	{ID: 5, Name: "Advanced Cash", ShortName: "ADVCUSD", Fiat: false},
	{ID: 6, Name: "Perfect Money", ShortName: "PMUSD", Fiat: false},
	{ID: 7, Name: "Payeer", ShortName: "PRUSD", Fiat: false},
	{ID: 8, Name: "Гривна", ShortName: "MONOBUAH", Fiat: true},
}

var DIRECTIONS = []*Direction{
	{ID: 1, From: "SBERRUB", To: "BTC"},
	{ID: 2, From: "SBERRUB", To: "USDT"},

	{ID: 3, From: "MONOBUAH", To: "BTC"},
	{ID: 4, From: "MONOBUAH", To: "USDT"},

	{ID: 5, From: "BTC", To: "MONOBUAH"},

	{ID: 6, From: "ETH", To: "SBERRUB"},

	{ID: 7, From: "ADVCUSD", To: "SBERRUB"},

	{ID: 8, From: "PRUSD", To: "SBERRUB"},

	{ID: 9, From: "PMUSD", To: "SBERRUB"},
}
