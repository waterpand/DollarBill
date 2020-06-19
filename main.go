package main

/* Желаемый функционал приложения:
- парсинг курса доллара (позже других валют - евро и швейцарского франка)
- запись значений курса за произвольный промежуток времени
- фиксация даты покупки валюты и прибыльность к текущему курсу
- расчет прибыли для гипотетической покупки.
*/

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Text    string   `xml:",chardata"`
	Date    string   `xml:"Date,attr"`
	Name    string   `xml:"name,attr"`
	Valute  []struct {
		Text     string `xml:",chardata"`
		ID       string `xml:"ID,attr"`
		NumCode  string `xml:"NumCode"`
		CharCode string `xml:"CharCode"`
		Nominal  string `xml:"Nominal"`
		Name     string `xml:"Name"`
		Value    string `xml:"Value"`
	} `xml:"Valute"`
}

func main() {
	responce, err := http.Get("https://www.cbr-xml-daily.ru/daily_utf8.xml")
	if err != nil {
		log.Fatal(err)
	}
	defer responce.Body.Close()

	byteValue, err := ioutil.ReadAll(responce.Body)
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println(string(byteValue))

	var (
		rate ValCurs
	)
	err = xml.Unmarshal(byteValue, &rate)
	if err != nil {
		log.Fatal(err)
	}

	currency := rate.Valute[10].Name
	date := rate.Date // выдает дату, на которую даётся курс
	fmt.Println(currency)
	fmt.Println(date)
}
