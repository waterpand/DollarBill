package main

/* Желаемый функционал приложения:
- парсинг курса доллара (позже других валют - евро и швейцарского франка) - complete
- запись значений курса за произвольный промежуток времени -- для этого нужно создать отдельную структуру, которая будет содержать данные по дням и будет записываться в файл
- запись данных о покупкх и продажах с сохранением
- фиксация даты покупки валюты и прибыльность к текущему курсу - complete
- расчет прибыли для гипотетической покупки - complete
- бумажная доходность к текущему курсу
- сохранение архивных курсов в срез и файл
- При запросе текущего курса сохранять его в архивный файл
- считывать архивный файл при запуске приложения
- добавить в функцию конвертации string другие форматы записи числа 01/01/2020 01.01.2020 01,01,2020 01:01:2020
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
	rate, rateOld, offlineRate                                                                 ValCurs
	cursOfToday, cursOfOldDay                                                                  Curs2
	op                                                                                         Order
	archiveCurses                                                                              []Curs2
	f                                                                                          []byte
	a, b, e                                                                                    int = 11, 10, 10
	i, c, d                                                                                    int
	dateOfPurchase, rateValuteNow                                                              string
	rateOfPurchase, amountOfСurrency, sumOfPurchase, todayCurrency, rateOfToday, percentOfRate float64
	temp                                                                                       Transact
	buy                                                                                        bool
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

// Curs : получение только используемых данных из банковского xml
type Curs struct {
	Date   string
	Valute [34]struct {
		Name     string
		CharCode string
		Value    float64
	}
}

// Curs2 : получение только используемых данных из банковского xml
type Curs2 struct {
	DD     int
	MM     int
	YYYY   int
	Valute [34]struct {
		Name     string
		CharCode string
		Value    float64
	}
}

// Order : структура для записи проведенных операция покупки/продажи валюты
type Order struct {
	Fresh       string
	Transaction []Transact // структура Transact

}

// Transact : структура для записи и сохранения операций купли/продажи
type Transact struct { // все параметры операции
	Date      string  // дата операции
	Quantity  float64 // кол-во валюты
	CharCode  string  // название валюты
	Price     float64 // курс
	Operation bool    // 1-покупка 0-продажа
	Flag      bool    // возможность исключить операцию из расчета
}

// ValCursToCurs : данная функция записывает данные полученные из xml в структуру CursOfToday
func ValCursToCurs(ret bool) {
	fmt.Println(rate.Date)
	fmt.Print("Запись полученных данных в структуру cursOfToday")
	//cursOfToday.Date = rate.Date
	DD, MM, YYYY := stringDateToInt(rate.Date)
	cursOfToday.YYYY = YYYY
	cursOfToday.MM = MM
	cursOfToday.DD = DD
	for i := 0; i < 34; i++ {
		cursOfToday.Valute[i].Name = rate.Valute[i].Name
		cursOfToday.Valute[i].CharCode = rate.Valute[i].CharCode
		cursOfToday.Valute[i].Value = stringToFloat(stringConvert(rate.Valute[i].Value))
	}
	fmt.Println("...complete")
	returnMenu(ret)
}

// ValCursToCurs2 : данная функция записывает данные полученные из файла ValCurs.bin в структуру CursOfToday
func ValCursToCurs2(print, ret bool) { // print(true) - печатать структуру, ret(true) - возврат в основное меню
	fmt.Println()
	fmt.Print("Запись полученных данных в структуру cursOfToday")
	//cursOfToday.Date = offlineRate.Date
	DD, MM, YYYY := stringDateToInt(offlineRate.Date)
	cursOfToday.YYYY = YYYY
	cursOfToday.MM = MM
	cursOfToday.DD = DD
	for i := 0; i < 34; i++ {
		cursOfToday.Valute[i].Name = offlineRate.Valute[i].Name
		cursOfToday.Valute[i].CharCode = offlineRate.Valute[i].CharCode
		cursOfToday.Valute[i].Value = stringToFloat(stringConvert(offlineRate.Valute[i].Value))
	}
	fmt.Println("...complete")
	if print == true {
		fmt.Println(cursOfToday)
	}
	returnMenu(ret)
}

// ValCursToCurs3 : данная функция записывает данные полученные из фрхивного xml в структуру cursOfOldDay для последующей записи в []Curs
func ValCursToCurs3(ret bool) {
	/*
	   Далее: записывать archiveCurses в файл, а перед append считывать файл и циклом проверять, нет ли там уже этих чисел (или не перед добавлением в срез, а сразу после запроса, чтобы лишний раз не парсить...)
	*/

	fmt.Println(rateOld.Date)
	fmt.Print("Запись полученных данных в структуру cursOfOldDay")
	//cursOfOldDay.Date = rateOld.Date
	DD, MM, YYYY := stringDateToInt(rateOld.Date)
	cursOfOldDay.YYYY = YYYY
	cursOfOldDay.MM = MM
	cursOfOldDay.DD = DD

	for i := 0; i < 34; i++ {
		cursOfOldDay.Valute[i].Name = rateOld.Valute[i].Name
		cursOfOldDay.Valute[i].CharCode = rateOld.Valute[i].CharCode
		cursOfOldDay.Valute[i].Value = stringToFloat(stringConvert(rateOld.Valute[i].Value))

	}
	fmt.Println("...complete")

	fmt.Println(cursOfOldDay)

	fmt.Println("Добавление данных в срез archiveCurses")
	archiveCurses = append(archiveCurses, cursOfOldDay)
	fmt.Println(archiveCurses)

	WriteFileValCursArchive(archiveCurses, ret) //когда заработает - сделать return archiveCurses и далее его в функцию Write File....
}

//список доступных валют в 7 столбцов
func currencySelection4(ret bool) {
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

	returnMenu(ret)
}

// ratePrint2 : Запрос номера валюты и вывод на печать курса
func ratePrint2(ret bool) {
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
	returnMenu(ret)
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

// stringDateToInt : конвертирует string дату в три int
func stringDateToInt(date string) (int, int, int) {

	dd, err := strconv.Atoi(date[:2])
	if err != nil {
		log.Fatal(err)
	}

	mm, err := strconv.Atoi(date[3:5])
	if err != nil {
		log.Fatal(err)
	}

	yyyy, err := strconv.Atoi(date[6:])
	if err != nil {
		log.Fatal(err)
	}
	return dd, mm, yyyy
}

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

// returnMenu : определяет будет ли возврат в основное меню или нет
func returnMenu(ret bool) {
	if ret == true {
		mainMenu()
	}
}

func mainMenu() {
	T := time.Now()
	fmt.Println()
	fmt.Printf(T.Format("_2.1.2006"))
	fmt.Print(" // ", cursOfToday.Valute[a-1].CharCode, " -- ", cursOfToday.Valute[a-1].Value)
	fmt.Printf(" // Меню:\n1 -- Баланс\n2 -- Купить\n3 -- Продать\n4 -- Показать историю операций\n5 -- \n6 -- \n7 -- \n8 -- Техническое меню\n9 -- func main() \n0 -- Выход из программы\n")
	fmt.Scanln(&b)

	switch b {
	case 1:
		DelFromStruct(false, true, false)
		Balans(true)
	case 2:
		rateCalculation(true, true)
	case 3:
		rateCalculation(false, true)
	case 4:
		DelFromStruct(false, true, true)
	case 5:

	case 6:

	case 7:

	case 8:
		techMenu()
	case 9:
		fmt.Println("Выход из меню")
	case 0:
		WriteTheFile(op, false)
		fmt.Println("Выход")
		os.Exit(0)
	case 11:
		DelFromStruct(true, false, true)

	default:
		fmt.Println("Введено неверное значение")
		mainMenu()
	}

}

func techMenu() {
	fmt.Printf(" // Техническое меню:\n1 -- Вычитать данные из xml, записать в файл ValCurs.bin и записать данные в структуру cursOfToday\n2 -- Прочитать информацию из файла и записать в структуру cursOfToday\n3 -- Фильтр по текущей валюте \n4 -- Вывести список доступных валют \n5 -- Сменить валюту\n6 -- Прочитать из файла историю операций\n7 -- Записать историю операций в файл\n8 -- Возврвт в основное меню - mainMenu\n9 -- Выход в func main()\n11 -- Запрос архивного курса\n12 -- Печать archiveCurses\n13 -- Преобразование даты\n0 -- Выход из программы\n")
	fmt.Scanln(&b)

	switch b {
	case 1:
		httpGet2()
		ValCursToCurs(true)
	case 2:
		readTheFile()
		ValCursToCurs2(true, true)
	case 3:
		FilterOp(true)
	case 4:
		currencySelection4(true)
	case 5:
		ratePrint2(true)
	case 6:
		readTheFile2(true)
	case 7:
		WriteTheFile(op, true)
	case 8:
		mainMenu()
	case 9:
		fmt.Println("Выход из меню")
	case 0:
		WriteTheFile(op, false)
		fmt.Println("Выход")
		os.Exit(0)
	case 11:
		CursArchive(false)
		ValCursToCurs3(true)
	case 12:
		PrintArchiveCurses(true)
		//fmt.Println(archiveCurses)
		//returnMenu(true)
	case 13:
		stringDateToInt(rate.Date)
		returnMenu(true)
	default:
		fmt.Println("Введено неверное значение")
		techMenu()
	}
}

func rateCalculation(buy, ret bool) { // расчет по выбранной валюте
	/*
		Добавить учет разряда валют, например, если курс установлен за 10 крон...
	*/

	fmt.Println("Валюта для расчета:", cursOfToday.Valute[a-1].CharCode, "  ", cursOfToday.Valute[a-1].Name)

	fmt.Println("Введите дату операции в формате ДД.ММ.ГГГГ:")
	fmt.Scanln(&dateOfPurchase)

	Dd, Mm, YYyy := stringDateToInt(dateOfPurchase)

	rateOfToday = cursOfToday.Valute[a-1].Value

	if Dd != cursOfToday.DD && Mm != cursOfToday.MM && YYyy != cursOfToday.YYYY {
		fmt.Println("Введите курс валюты (формат $$.$$$$):")
		fmt.Scanln(&rateOfPurchase)
	} else {
		rateOfPurchase = rateOfToday
	}

	fmt.Println("Введите количество валюты:")
	fmt.Scanln(&amountOfСurrency)

	sumOfPurchase = rateOfPurchase * amountOfСurrency // Сумма покупки
	fmt.Println("Сумма операции:", sumOfPurchase)

	todayCurrency = amountOfСurrency * rateOfToday // Стоимость по текущему курсу
	fmt.Println("Стоимость по текущему курсу:", todayCurrency)

	percentOfRate = ((rateOfToday / rateOfPurchase) - 1) * 10000
	percentOfRate = math.Round(percentOfRate) * 0.01
	fmt.Println("Текущий результат: \n", percentOfRate, "%\n", todayCurrency-sumOfPurchase, "руб.\n", offlineRate.Valute[a-1].Name, amountOfСurrency, "шт.")
	fmt.Println("----------------------------------------------------------------")

	SafeOperation(cursOfToday.Valute[a-1].CharCode, dateOfPurchase, buy, true, rateOfPurchase, amountOfСurrency, ret)

}

// SafeOperation : сохранение операции по покупке валюты / запись в структуру Order
func SafeOperation(ChC, D string, Opp, Fl bool, Pr, Q float64, ret bool) {

	/*
		Если произвести несколько операций с одной валютой, а потом заменить валюту и произвести еще одну операцию, то все опвалюта всех операций изменится
		Варианты:
		-- либо копать в сторону карт, чтобы для каждой валюты операции записывались отдельно
		-- (пока выбран этот вариант) либо внести обозначение валюты внутрь структуры, чтобы каждую операцию можно было идентифицировать по валюте
	*/
	T := time.Now()

	op.Fresh = T.Format("_2.1.2006")

	temp.CharCode = ChC
	temp.Date = D
	temp.Operation = Opp
	temp.Price = Pr
	temp.Quantity = Q
	temp.Flag = Fl
	op.Transaction = append(op.Transaction, temp)

	WriteTheFile(op, false)
	returnMenu(ret)
}

// readTheFile2 : Чтение из файла ValCurs.bin
func readTheFile2(ret bool) {
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

	returnMenu(ret)
}

// readTheFile3 : Чтение из файла ValCursArchive.bin
func readTheFile3(ret bool) {
	file, err := os.Open("D:/_development/_projects/DollarBill/ValCursArchive.json")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	data, err := ioutil.ReadFile("D:/_development/_projects/DollarBill/ValCursArchive.json")
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(data, &archiveCurses)
	if err != nil {
		log.Fatal(err)
	}

	returnMenu(ret)
}

// WriteTheFile : данная функция записывает в файл данные из структуры с записями всех операций
func WriteTheFile(op Order, ret bool) {

	byteValue, err := json.Marshal(op)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create("D:/_development/_projects/DollarBill/OperationDamp.json") // создание файла
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	file.Write(byteValue)

	fmt.Println("Данные записаны в файл", file.Name())

	returnMenu(ret)
}

// WriteFileValCursArchive : записывает в файл ValCursArchive.json данные из среза ArchiveCurses
func WriteFileValCursArchive(ac []Curs2, ret bool) {

	byteValue, err := json.Marshal(ac)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create("D:/_development/_projects/DollarBill/ValCursArchive.json") // создание файла
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	file.Write(byteValue)

	fmt.Println("Данные записаны в файл", file.Name())

	returnMenu(ret)
}

// DelFromStruct : вывод списка операций или удаления записей об операциях из структуры
func DelFromStruct(del, usd, ret bool) {
	if usd == true && op.Transaction[i].CharCode == cursOfToday.Valute[a-1].CharCode {
		fmt.Printf("\nСписок операций по текущей валюте:\n")
	} else {
		fmt.Printf("\nСписок операций:\n")
	}

	for i := range op.Transaction {
		if usd == true && op.Transaction[i].CharCode == cursOfToday.Valute[a-1].CharCode {
			if op.Transaction[i].Operation == true {
				fmt.Print(i+1, ")  ", op.Transaction[i].Date, " - Покупка ")
			} else {
				fmt.Print(i+1, ")  ", op.Transaction[i].Date, " - Продажа ")
			}
			fmt.Println(op.Transaction[i].Quantity, "", op.Transaction[i].CharCode, " по курсу ", op.Transaction[i].Price)
		} else if usd == false {
			if op.Transaction[i].Operation == true {
				fmt.Print(i+1, ")  ", op.Transaction[i].Date, " - Покупка ")
			} else {
				fmt.Print(i+1, ")  ", op.Transaction[i].Date, " - Продажа ")
			}
			fmt.Println(op.Transaction[i].Quantity, "", op.Transaction[i].CharCode, " по курсу ", op.Transaction[i].Price)
		}
	}

	if del == true {
		j := 100
		fmt.Println("Выбрать номер удаляемой транзакции")
		fmt.Scanln(&j)
		j--

		op.Transaction = append(op.Transaction[:j], op.Transaction[j+1:]...)

		for i := range op.Transaction {
			fmt.Println(i+1, op.Transaction[i])
		}
	}

	returnMenu(ret)
}

// FilterOp : фильтр для валют
func FilterOp(ret bool) {
	fmt.Println()
	fmt.Println("Показаны только операции с текущей валютой: ", cursOfToday.Valute[a-1].CharCode, "--", cursOfToday.Valute[a-1].Name)
	fmt.Println()

	for i, _ := range op.Transaction {
		if op.Transaction[i].CharCode == cursOfToday.Valute[a-1].CharCode {
			fmt.Println(i+1, op.Transaction[i])
		}
	}
	returnMenu(ret)
}

// Balans : суммирует все операции по конкретной валюте
func Balans(ret bool) {
	sum := 0.0
	amount := 0.0
	//cursAverage := 0.0
	k := 1.0
	opA := 0
	opS := 0
	fmt.Println()
	fmt.Println("Баланс для текущей валюты: ", cursOfToday.Valute[a-1].CharCode, "--", cursOfToday.Valute[a-1].Name)
	fmt.Println()

	for i := range op.Transaction {
		if op.Transaction[i].CharCode == cursOfToday.Valute[a-1].CharCode {
			opA++
			if op.Transaction[i].Operation != true {
				k = -1.0
				opS++
			}

			sum = sum + (op.Transaction[i].Quantity*op.Transaction[i].Price)*k
			amount = amount + op.Transaction[i].Quantity*k
		}
	}

	fmt.Println("Итого:")
	fmt.Println("Всего произведено", opA, "операций", opA-opS, "покупок валюты и", opS, "продаж")
	fmt.Println("На балансе", amount, cursOfToday.Valute[a-1].CharCode, "на сумму ", sum)
	fmt.Println("Средний курс: ", sum/amount)

	returnMenu(ret)
}

// CursArchive : запрос курса за прошедщие дни
func CursArchive(ret bool) { //добавить функцию конвертации string в нужный формат числа
	/*
		Эта функция берет данные из банковского xml-файла формирует переменную rateOld типа ValCurs
		Записывает эти данные в файл ValCurs.bin
		Но я не имею ни малейшего понятия, как это работает
	*/
	var (
		ab, bc, cd string
	)
	ab = "https://www.cbr-xml-daily.ru/daily_utf8.xml"
	bc = "?date_req="

	fmt.Println("Введите дату в формате ДД/ММ/ГГГГ:")
	fmt.Scanln(&cd)
	fmt.Println(ab + bc + cd)

	fmt.Println("Запрос...", ab+bc+cd)
	responce, err := http.Get(ab + bc + cd)
	if err != nil {
		log.Fatal(err)
	}
	defer responce.Body.Close()

	byteValue, err := ioutil.ReadAll(responce.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = xml.Unmarshal(byteValue, &rateOld)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Данные получены и записаны в rateOld")
	fmt.Println(rateOld)

	returnMenu(ret)
}

// PrintArchiveCurses : Выводит в удобном виде срез archiveCurses
func PrintArchiveCurses(ret bool) {

	for i := range archiveCurses {
		fmt.Print(archiveCurses[i].DD, ".", archiveCurses[i].MM, ".", archiveCurses[i].YYYY, "\n")
		fmt.Println(archiveCurses[i].Valute[:4])
		fmt.Println(archiveCurses[i].Valute[4:9])
		fmt.Println(archiveCurses[i].Valute[9:15])
		fmt.Println(archiveCurses[i].Valute[15:21])
		fmt.Println(archiveCurses[i].Valute[21:26])
		fmt.Println(archiveCurses[i].Valute[26:31])
		fmt.Println(archiveCurses[i].Valute[31:34])
		fmt.Println()
	}

	returnMenu(ret)
}

func main() {
	readTheFile()
	ValCursToCurs2(false, false)
	readTheFile2(false)
	readTheFile3(true)
	defer mainMenu()

	fmt.Println("func main")
	T := time.Now()
	fmt.Println(T.Format("_2.1.2006"))

}
