package main

/* Желаемый функционал приложения:
- парсинг курса доллара (позже других валют - евро и швейцарского франка) - complete
- запись значений курса за произвольный промежуток времени -- для этого нужно создать отдельную структуру, которая будет содержать данные по дням и будет записываться в файл
- запись данных о покупкх и продажах с сохранением
- фиксация даты покупки валюты и прибыльность к текущему курсу - complete
- расчет прибыли для гипотетической покупки - complete
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
	offlineRate                                                                                ValCurs
	cursOfToday                                                                                Curs
	f                                                                                          []byte
	b                                                                                          int = 10
	a, i, c                                                                                    int
	dateOfPurchase, rateValuteNow                                                              string
	rateOfPurchase, amountOfСurrency, sumOfPurchase, todayCurrency, rateOfToday, percentOfRate float64
)

// ValCurs : сгененрирована автоматически на сайте https://www.onlinetool.io/xmltogo/ по ссылке ЦБ (https://www.cbr-xml-daily.ru/daily_utf8.xml)
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

// Curs создана для получения только используемых данных из банковского xml
type Curs struct {
	Date   string
	Valute [34]struct {
		Name     string
		CharCode string
		Value    float64
	}
}

// ValCursToCurs : данная функция записывает данные полученные из xml в структуру CursOfToday
func ValCursToCurs() {
	fmt.Print("Запись полученных данных в структуру cursOfToday")
	cursOfToday.Date = rate.Date
	for i := 0; i < 34; i++ {
		cursOfToday.Valute[i].Name = rate.Valute[i].Name
		cursOfToday.Valute[i].CharCode = rate.Valute[i].CharCode
		cursOfToday.Valute[i].Value = stringToFloat(stringConvert(rate.Valute[i].Value))
	}
	fmt.Println("...complete")
	mainMenu()
}

// ValCursToCurs2 : данная функция записывает данные полученные из файла ValCurs.bin в структуру CursOfToday
func ValCursToCurs2() {
	fmt.Println()
	fmt.Print("Запись полученных данных в структуру cursOfToday")
	cursOfToday.Date = offlineRate.Date
	for i := 0; i < 34; i++ {
		cursOfToday.Valute[i].Name = offlineRate.Valute[i].Name
		cursOfToday.Valute[i].CharCode = offlineRate.Valute[i].CharCode
		cursOfToday.Valute[i].Value = stringToFloat(stringConvert(offlineRate.Valute[i].Value))
	}
	fmt.Println("...complete")
	mainMenu()
}

/*
func currencySelection() (a int) { //список доступных валют
	/*
		Вывод на экран списка всех доступных валют
	*

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
*/

/*
func currencySelection2() { //список доступных валют без запроса кода
	/*
		Вывод на экран списка всех доступных валют
	*

	fmt.Println("Доступные валюты:")
	for j := 0; j < 34; j++ {
		fmt.Println(j+1, "--", rate.Valute[j].CharCode, "--", rate.Valute[j].Name)
	}

	mainMenu()
}
*/

/*
func currencySelection3() { //список доступных валют в 4 столбца
	/*
		Вывод на экран списка всех доступных валют в четыре столбца
		Увеличить количество столбцов до 6 или 7 чтобы было удобно читать
	*
	fmt.Println()
	fmt.Println(offlineRate.Date, " Доступные валюты:")
	for j := 0; j < 10; j++ {
		if j < 4 {
			fmt.Print(j+1, " -- ", cursOfToday.Valute[j].CharCode, "			")
			fmt.Print(j+11, " -- ", cursOfToday.Valute[j+10].CharCode, "			")
			fmt.Print(j+21, " -- ", cursOfToday.Valute[j+20].CharCode, "			")
			fmt.Println(j+31, "--", cursOfToday.Valute[j+30].CharCode, "			")
		} else {
			fmt.Print(j+1, " -- ", cursOfToday.Valute[j].CharCode, "			")
			fmt.Print(j+11, " -- ", cursOfToday.Valute[j+10].CharCode, "			")
			fmt.Println(j+21, "--", cursOfToday.Valute[j+20].CharCode, "			")
		}

	}

	mainMenu()
}
*/

//список доступных валют в 7 столбцов
func currencySelection4() {
	fmt.Println()
	fmt.Println(offlineRate.Date, " Доступные валюты:")
	for j := 0; j < 5; j++ {
		if j < 4 {
			fmt.Print(j+1, "  -- ", cursOfToday.Valute[j].CharCode, "		")
			fmt.Print(j+6, "  -- ", cursOfToday.Valute[j+5].CharCode, "		")
			fmt.Print(j+11, " -- ", cursOfToday.Valute[j+10].CharCode, "		")
			fmt.Print(j+16, " -- ", cursOfToday.Valute[j+15].CharCode, "		")
			fmt.Print(j+21, " -- ", cursOfToday.Valute[j+20].CharCode, "		")
			fmt.Print(j+26, " -- ", cursOfToday.Valute[j+25].CharCode, "		")
			fmt.Println(j+31, "--", cursOfToday.Valute[j+30].CharCode, "		")
		} else {
			fmt.Print(j+1, "  -- ", cursOfToday.Valute[j].CharCode, "		")
			fmt.Print(j+6, " -- ", cursOfToday.Valute[j+5].CharCode, "		")
			fmt.Print(j+11, " -- ", cursOfToday.Valute[j+10].CharCode, "		")
			fmt.Print(j+16, " -- ", cursOfToday.Valute[j+15].CharCode, "		")
			fmt.Print(j+21, " -- ", cursOfToday.Valute[j+20].CharCode, "		")
			fmt.Println(j+26, "--", cursOfToday.Valute[j+25].CharCode, "		")
		}

	}

	mainMenu()
}

/*
func ratePrint(i int) { // Курс конкретной валюты
	/*
		Вывод на печать курса валюты в формате: USD -- 69,5725 -- Американский доллар.
		Переписать, чтобы читалось из структуры
	*
	fmt.Println("	  ", rate.Valute[i].CharCode, "--", rate.Valute[i].Value, "--", rate.Valute[i].Name)
}
*/

// ratePrint2 : Запрос номера валюты и вывод на печать курса
func ratePrint2() {
	/*
		Вывод на печать курса валюты в формате: USD -- 69,5725 -- Американский доллар.
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
	fmt.Println("	 ", cursOfToday.Valute[a-1].CharCode, "--", cursOfToday.Valute[a-1].Value, "--", cursOfToday.Valute[a-1].Name)
	mainMenu()
}

// stringConvert : меняет запятую на точку в данных, которые подтягиваются по xml, что их удобно было конвертировать в float64
func stringConvert(in string) string {

	out := strings.Replace(in, ",", ".", -1) // Замена запятой на точку. "-1" означает что будет производиться замена всех найденных символов
	return out
}

// stringToFloat : конвертирует тип string в тип float64
func stringToFloat(in string) float64 {

	out, _ := strconv.ParseFloat(in, 8)
	return out
}

/*
func httpGet() ValCurs { // вычитка из xml
	/*
		Эта функция берет данные из банковского xml-файла формирует переменную rate типа ValCurs
		Но я не имею ни малейшего понятия, как это работает
	*
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
*/

func httpGet2() ValCurs { // вычитка из xml и запись в файл
	/*
		Эта функция берет данные из банковского xml-файла формирует переменную rate типа ValCurs
		Записывает эти данные в файл ValCurs.bin
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

	file, err := os.Create("D:/_development/_projects/DollarBill/ValCurs.bin")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	file.Write(byteValue)

	err = xml.Unmarshal(byteValue, &rate)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Данные получены и записаны в файл", file.Name())
	return rate
}

// readTheFile : Чтение из файла ValCurs.bin
func readTheFile() {
	file, err := os.Open("D:/_development/_projects/DollarBill/ValCurs.bin")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	data, err := ioutil.ReadFile("D:/_development/_projects/DollarBill/ValCurs.bin")
	if err != nil {
		fmt.Println(err)
	}

	err = xml.Unmarshal(data, &offlineRate)
	if err != nil {
		log.Fatal(err)
	}

}

func mainMenu() {
	fmt.Println()
	fmt.Printf("Меню:\n1 -- Вычитать данные из xml, записать в файл ValCurs.bin и записать данные в структуру cursOfToday\n2 -- Прочитать информацию из файла и записать в структуру cursOfToday\n3 -- Вывести список доступных валют\n4 -- Посмотреть курс конкретной валюты/выбрать валюту\n5 -- Произвести рассчет\n0 -- Выход из программы\n")
	fmt.Scanln(&b)

	switch b {
	case 1:
		httpGet2()
		ValCursToCurs()
	case 2:
		readTheFile()
		ValCursToCurs2()
	case 3:
		currencySelection4()
	case 4:
		ratePrint2()
	case 5:
		rateCalculation()
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
	//if cursOfToday.Valute[a-1].Value == 0 {
	//fmt.Println("Необходимо записать данные в структуру // пункт 2 основного меню")
	//}
	fmt.Println("Валюта для расчета:", cursOfToday.Valute[a-1].CharCode, "  ", cursOfToday.Valute[a-1].Name)
	fmt.Print("Расчитать для текущей валюты - 1\n              Сменить валюту - 2 ")
	fmt.Scanln(&c)
	if c != 1 {
		mainMenu()
	}

	fmt.Println("Введите дату покупки в формате ДД.ММ.ГГГГ:")
	fmt.Scanln(&dateOfPurchase)

	rateOfToday = cursOfToday.Valute[a-1].Value

	if dateOfPurchase != cursOfToday.Date {
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
	fmt.Println("Текущий результат: \n", percentOfRate, "%\n", todayCurrency-sumOfPurchase, "руб.\n", offlineRate.Valute[a-1].Name, amountOfСurrency, "шт.")
	fmt.Println("----------------------------------------------------------------")
	mainMenu()
}

func main() {
	mainMenu()

}
