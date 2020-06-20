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

var (
	rate ValCurs
	a, i int
)

type ValCurs struct { //эта структура сгененрирована автоматически на сайте https://www.onlinetool.io/xmltogo/ по ссылке ЦБ (https://www.cbr-xml-daily.ru/daily_utf8.xml)
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

func currencySelection() {
	fmt.Println("Доступные валюты:")
	for j := 0; j < 34; j++ {
		fmt.Println(j+1, "--", rate.Valute[j].CharCode, "--", rate.Valute[j].Name)
	}

}

func ratePrint(i int) { // Вывод на печать курса валюты в формате: USD -- 69,5725. Коды для валют: 10 - USD, 11 - EUR, 30 - CHF.
	fmt.Println("		  ", rate.Valute[i].CharCode, "--", rate.Valute[i].Value)
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

	err = xml.Unmarshal(byteValue, &rate)
	if err != nil {
		log.Fatal(err)
	}
	// не имею ни малейшего понятия, как работает весь предыдущий кусок и можно ли его убрать в отдельную функцию...

	currencySelection()
	fmt.Println("введите номер валюты:")
	for i := 0; i < 1; {
		fmt.Scanln(&a)
		if a < 1 || a > 34 {
			fmt.Println("Неверное число, попробуйте ещё раз:")
		} else {
			i = 1
		}
	}
	ratePrint(a - 1)

}
