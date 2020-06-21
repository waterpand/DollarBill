package main

/* Желаемый функционал приложения:
- парсинг курса доллара (позже других валют - евро и швейцарского франка) - complete
- запись значений курса за произвольный промежуток времени -- для этого нужно создать отдельную структуру, которая будет содержать данные по дням и будет записываться в файл
- запись данных о покупкх и продажах с сохранением
- фиксация даты покупки валюты и прибыльность к текущему курсу - complete
- расчет прибыли для гипотетической покупки - complete
*/

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	rate                                                                                       ValCurs
	offlineRate                                                                                ValCurs
	cursOfToday                                                                                Curs
	op                                                                                         Order
	f                                                                                          []byte
	b, e                                                                                       int = 10, 10
	a, i, c, d                                                                                 int
	dateOfPurchase, rateValuteNow                                                              string
	rateOfPurchase, amountOfСurrency, sumOfPurchase, todayCurrency, rateOfToday, percentOfRate float64
	temp                                                                                       Transact
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

// Order : структура для записи проведенных операция покупки/продажи валюты
type Order struct {
	CharCode    string     // название валюты
	Transaction []Transact // структура Transact

}

// Transact : структура для записи и сохранения операций купли/продажи
type Transact struct { // все параметры операции
	IdOpp     int     // ID операции для навигации по срезу
	CharCode  string  // название валюты
	Operation bool    // 1-покупка 0-продажа
	Price     float64 // курс
	Date      string  // дата операции
	Quantity  float64 // кол-во валюты
	Flag      bool    // возможность исключить операцию из расчета
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
// вычитка из xml и запись в файл
func httpGet2() ValCurs {
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

	file, err := os.Create("D:/_development/_projects/DollarBill/ValCurs.bin") // создание файла
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
	T := time.Date(2020, time.June, 21, 8, 5, 2, 0, time.Local)
	fmt.Println()
	fmt.Printf(T.Format("_2.1.2006"))
	fmt.Print(" // ", cursOfToday.Valute[a-1].CharCode, " -- ", cursOfToday.Valute[a-1].Value)
	fmt.Printf(" // Меню:\n1 -- Вычитать данные из xml, записать в файл ValCurs.bin и записать данные в структуру cursOfToday\n2 -- Вывести список доступных валют\n3 -- Сменить валюту\n4 -- Произвести рассчет\n5 -- Прочитать из файла историю операций\n6 -- Показать историю операций\n7 -- Записать историю операций в файл\n8 -- Техническое меню\n9 -- func main() \n0 -- Выход из программы\n")
	fmt.Scanln(&b)

	switch b {
	case 1:
		httpGet2()
		ValCursToCurs()
	case 2:
		currencySelection4()
	case 3:
		ratePrint2()
	case 4:
		rateCalculation()
	case 5:
		readTheFile2()
	case 6:
		fmt.Println(op)
		mainMenu()
	case 7:
		WriteTheFile(op)
	case 8:
		techMenu()
	case 9:
		fmt.Println("Выход из меню")
	case 0:
		fmt.Println("Выход")
		os.Exit(0)

	default:
		fmt.Println("Введено неверное значение")
		mainMenu()
	}

}

func techMenu() {
	fmt.Printf(" // Техническое меню:\n1 -- \n2 -- Прочитать информацию из файла и записать в структуру cursOfToday\n3 -- \n4 -- \n5 -- \n8 -- Возврвт в основное меню - mainMenu\n9 -- Выход в func main()\n0 -- Выход из программы\n")
	fmt.Scanln(&b)

	switch b {
	case 1:

	case 2:
		readTheFile()
		ValCursToCurs2()
	case 3:

	case 4:

	case 5:

	case 8:
		mainMenu()
	case 9:
		fmt.Println("Выход из меню")
	case 0:
		fmt.Println("Выход")
		os.Exit(0)

	default:
		fmt.Println("Введено неверное значение")
		techMenu()
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

	fmt.Println("Запомнить результат? для сохранения - 1, для сброса - любое число")
	fmt.Scanln(&d)
	if c == 1 {
		SafeOperation(cursOfToday.Valute[a-1].CharCode, dateOfPurchase, true, true, rateOfPurchase, amountOfСurrency)
	}

	mainMenu()
}

// SafeOperation : сохранение операции по покупке валюты / запись в структуру Order
func SafeOperation(ChC, D string, Opp, Fl bool, Pr, Q float64) {

	/*
		Если произвести несколько операций с одной валютой, а потом заменить валюту и произвести еще одну операцию, то все опвалюта всех операций изменится
		Варианты:
		-- либо копать в сторону карт, чтобы для каждой валюты операции записывались отдельно
		-- (пока выбран этот вариант) либо внести обозначение валюты внутрь структуры, чтобы каждую операцию можно было идентифицировать по валюте
	*/

	fmt.Println("Запись в структуру op (type Order):")
	fmt.Println("Название валюты", ChC)
	fmt.Println("Покупка или продажа (true - покупка)", Opp)
	fmt.Println("цена покупки (курс)", Pr)
	fmt.Println("Дата покупки", D)
	fmt.Println("Количество валюты", Q)
	fmt.Println("Учет операции (true - учитывать)", Fl)
	fmt.Println()

	temp.IdOpp = len(op.Transaction)
	temp.CharCode = ChC
	temp.Date = D
	temp.Operation = Opp
	temp.Price = Pr
	temp.Quantity = Q
	temp.Flag = Fl
	op.Transaction = append(op.Transaction, temp)

	fmt.Println("Все операции:", op)
	fmt.Println()
	fmt.Println("Текущая операция:", temp)

	fmt.Println("1 - Записать в файл\n2 - выйти в меню")
	fmt.Scanln(&e)
	if e == 1 {
		WriteTheFile(op)
	}
	mainMenu()
}

// readTheFile2 : Чтение из файла ValCurs.bin
func readTheFile2() {
	file, err := os.Open("D:/_development/_projects/DollarBill/OperationDamp.json")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	data, err := ioutil.ReadFile("D:/_development/_projects/DollarBill/OperationDamp.json")
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(data, &op)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(op)

	mainMenu()
}

// WriteTheFile : данная функция записывает в файл данные из структуры с записями всех операций
func WriteTheFile(op Order) {

	byteValue, err := json.Marshal(op)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create("D:/_development/_projects/DollarBill/OperationDamp.json") // создание файла
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	file.Write(byteValue) //  cannot use op (type Order) as type []byte in argument to file.Write

	fmt.Println("Данные записаны в файл", file.Name())
	mainMenu()

}

func main() {
	a = 11
	readTheFile()
	ValCursToCurs2()
	defer mainMenu()

	fmt.Println("func main")
	T := time.Date(2020, time.June, 21, 8, 5, 2, 0, time.Local)
	fmt.Println(T.Format("_2.1.2006"))
	fmt.Println(cursOfToday)
}
