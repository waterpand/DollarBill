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
	"math"
	"net/http"
	"strconv"
	"strings"
)

var (
	rate                                                                                       ValCurs
	a, i                                                                                       int
	dateOfPurchase, rateValuteNow                                                              string
	rateOfPurchase, amountOfСurrency, sumOfPurchase, todayCurrency, rateOfToday, percentOfRate float64
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

func currencySelection() (a int) {
	fmt.Println("Доступные валюты:")
	for j := 0; j < 34; j++ {
		fmt.Println(j+1, "--", rate.Valute[j].CharCode, "--", rate.Valute[j].Name)
	}

	fmt.Println("введите номер валюты:")
	for i := 0; i < 1; {
		fmt.Scanln(&a)
		if a < 1 || a > 34 {
			fmt.Println("Неверное число, попробуйте ещё раз:")
		} else {
			i = 1
		}
	}
	return a - 1
}

func ratePrint(i int) { // Вывод на печать курса валюты в формате: USD -- 69,5725. Коды для валют: 10 - USD, 11 - EUR, 30 - CHF.
	fmt.Println("		  ", rate.Valute[i].CharCode, "--", rate.Valute[i].Value)
}

func stringConvert(in string) string {
	/*
		Эта функция меняет запятую на точку в данных, которые подтягиваются по xml, что их удобно было конвертировать в float64
	*/
	out := strings.Replace(in, ",", ".", -1) // Замена запятой на точку
	return out
}

func stringToFloat(in string) float64 {
	/*
		Эта функция конвертирует тип string в тип float64
	*/
	out, _ := strconv.ParseFloat(in, 8)
	return out
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

	//a = currencySelection()
	a = 10
	ratePrint(a)
	fmt.Println("Введите дату покупки в формате ДД.ММ.ГГГГ:")
	fmt.Scanln(&dateOfPurchase)

	rateValuteNow = stringConvert(rate.Valute[a].Value) // Замена запятой на точку
	rateOfToday = stringToFloat(rateValuteNow)          // конвертация в float64

	if dateOfPurchase != rate.Date {
		fmt.Println("Введите курс покупки (формат $$.$$$$):")
		fmt.Scanln(&rateOfPurchase)
	} else {

		rateOfPurchase = rateOfToday
	}

	fmt.Println("Введите количество купленной валюты:")
	fmt.Scanln(&amountOfСurrency)

	sumOfPurchase = rateOfPurchase * amountOfСurrency // Сумма покупки
	fmt.Println("Сумма покупки:", sumOfPurchase)

	todayCurrency = amountOfСurrency * rateOfToday // Стоимость по текущему курсу
	fmt.Println("Стоимость по текущему курсу:", todayCurrency)

	percentOfRate = ((rateOfToday / rateOfPurchase) - 1) * 10000
	percentOfRate = math.Round(percentOfRate) * 0.01
	fmt.Println("Текущий результат: \n", percentOfRate, "%\n", todayCurrency-sumOfPurchase, "руб.")
}
