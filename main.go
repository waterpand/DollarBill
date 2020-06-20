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
	"os"
	"strconv"
	"strings"
)

var (
	rate                                                                                       ValCurs
	CursOfToday                                                                                Curs
	b                                                                                          int = 10
	a, i, c                                                                                    int
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

type Curs struct {
	Date   string
	Valute [34]struct {
		Name     string
		CharCode string
		Value    float64
	}
}

func ValCursToCurs() { //xml --> struct
	fmt.Print("Запись полученных данных в структуру CursOfToday")
	CursOfToday.Date = rate.Date
	for i := 0; i < 34; i++ {
		CursOfToday.Valute[i].Name = rate.Valute[i].Name
		CursOfToday.Valute[i].CharCode = rate.Valute[i].CharCode
		CursOfToday.Valute[i].Value = stringToFloat(stringConvert(rate.Valute[i].Value))
	}
	fmt.Println("...complete")
	mainMenu()
}

func currencySelection() (a int) { //список доступных валют
	/*
		Вывод на экран списка всех доступных валют
	*/

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
	mainMenu()
	return a - 1
}

func currencySelection2() { //список доступных валют без запроса кода
	/*
		Вывод на экран списка всех доступных валют
	*/

	fmt.Println("Доступные валюты:")
	for j := 0; j < 34; j++ {
		fmt.Println(j+1, "--", rate.Valute[j].CharCode, "--", rate.Valute[j].Name)
	}

	mainMenu()
}

func currencySelection3() { //список доступных валют в 4 столбца
	/*
		Вывод на экран списка всех доступных валют
	*/

	fmt.Println("Доступные валюты:")
	for j := 0; j < 10; j++ {
		if j < 4 {
			fmt.Print(j+1, " -- ", CursOfToday.Valute[j].CharCode, "			")
			fmt.Print(j+11, " -- ", CursOfToday.Valute[j+10].CharCode, "			")
			fmt.Print(j+21, " -- ", CursOfToday.Valute[j+20].CharCode, "			")
			fmt.Println(j+31, "--", CursOfToday.Valute[j+30].CharCode, "			")
		} else {
			fmt.Print(j+1, " -- ", CursOfToday.Valute[j].CharCode, "			")
			fmt.Print(j+11, " -- ", CursOfToday.Valute[j+10].CharCode, "			")
			fmt.Println(j+21, "--", CursOfToday.Valute[j+20].CharCode, "			")
		}

	}

	mainMenu()
}

func ratePrint(i int) { // Курс конкретной валюты
	/*
		Вывод на печать курса валюты в формате: USD -- 69,5725 -- Американский доллар.
		Переписать, чтобы читалось из структуры
	*/
	fmt.Println("	  ", rate.Valute[i].CharCode, "--", rate.Valute[i].Value, "--", rate.Valute[i].Name)
}

func ratePrint2() { // Запрос номера валюты и вывод на печать курса
	/*
		Вывод на печать курса валюты в формате: USD -- 69,5725 -- Американский доллар.
		Переписать, чтобы читалось из структуры
	*/
	fmt.Println("Введите номер валюты:")
	for i := 0; i < 1; {
		fmt.Scanln(&a)
		if a < 1 || a > 34 {
			fmt.Println("Неверное число, попробуйте ещё раз:")
		} else {
			i = 1
		}
	}
	fmt.Println("	 ", CursOfToday.Valute[a-1].CharCode, "--", CursOfToday.Valute[a-1].Value, "--", CursOfToday.Valute[a-1].Name)
	mainMenu()
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

func httpGet() ValCurs { // вычитка из xml
	/*
		Эта функция берет данные из банковского xml-файла формирует переменную rate типа ValCurs
		Но я не имею ни малейшего понятия, как это работает
	*/
	fmt.Println("Запрос...https://www.cbr-xml-daily.ru/daily_utf8.xml")
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
	fmt.Println("Данные получены")
	mainMenu()
	return rate
}

func mainMenu() {
	fmt.Printf("Меню:\n1 -- Вычитать данные из xml\n2 -- Записать данные в структуру\n3 -- Вывести список доступных валют\n4 -- Посмотреть курс конкретной валюты/выбрать валюту\n5 -- Произвести рассчет\n0 -- Выход из программы\n")
	fmt.Scanln(&b)

	switch b {
	case 1:
		httpGet()
	case 2:
		ValCursToCurs()
	case 3:
		currencySelection3()
	case 4:
		ratePrint2()
	case 5:
		rateCalculation()
	case 6:
		httpGet()
		ValCursToCurs()
	case 0:
		fmt.Println("Выход")
		os.Exit(0)
	default:
		fmt.Println("Введено неверное значение")
		mainMenu()
	}

}

func rateCalculation() { // расчет по выбранной валюте
	/*
		Добавить учет разряда валют, например, если курс установлен за 10 крон...
	*/
	//if CursOfToday.Valute[a-1].Value == 0 {
	//fmt.Println("Необходимо записать данные в структуру // пункт 2 основного меню")
	//}
	fmt.Println("Валюта для расчета:", CursOfToday.Valute[a-1].CharCode, "  ", CursOfToday.Valute[a-1].Name)
	fmt.Print("Расчитать для текущей валюты - 1\n              Сменить валюту - 2 ")
	fmt.Scanln(&c)
	if c != 1 {
		mainMenu()
	}

	fmt.Println("Введите дату покупки в формате ДД.ММ.ГГГГ:")
	fmt.Scanln(&dateOfPurchase)

	rateOfToday = CursOfToday.Valute[a-1].Value

	if dateOfPurchase != CursOfToday.Date {
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
	fmt.Println("Текущий результат: \n", percentOfRate, "%\n", todayCurrency-sumOfPurchase, "руб.\n", rate.Valute[a-1].Name, amountOfСurrency, "шт.")
	fmt.Println("----------------------------------------------------------------")
	mainMenu()
}

func main() {
	mainMenu()

}
