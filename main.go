package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

func round(x, unit float64) float64 {
	return math.Round(x/unit) * unit
}

func readInput() string {
	textIn, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	txt := strings.TrimSpace(textIn)
	return txt
}

func fillMap(intSlice []int) map[int]int {
	mp := map[int]int{}
	for _, elem := range intSlice {
		mp[elem] += 1
	}
	return mp
}

func getSlice(txt string) ([]int, error) {
	txtSlice := strings.Split(txt, ", ")
	intSlice := make([]int, 0, len(txtSlice))
	sum := 0
	for _, elem := range txtSlice {
		strInt, err := strconv.Atoi(elem)
		if err != nil {
			return []int{}, fmt.Errorf("int's conversion error")
		}
		if strInt < 0 {
			return []int{}, fmt.Errorf("contains negative numbers")
		}
		sum += strInt
		intSlice = append(intSlice, strInt)
	}
	if sum == 0 {
		return []int{}, fmt.Errorf("series contains only zero")
	}
	if len(intSlice) < 4 {
		return []int{}, fmt.Errorf("the size of series isn't enough")
	}
	return intSlice, nil
}

func getUniqueSlice(mp map[int]int) []int {
	sl := make([]int, 0, len(mp))
	for key := range mp {
		sl = append(sl, key)
	}
	return sl
}

type StrCriteria struct {
	max     int
	preMax  int
	postMin int
	min     int
}

func newStore() *StrCriteria {
	return &StrCriteria{}
}

func (s *StrCriteria) fillCriteria(sl []int) {
	lenSl := len(sl)
	s.max = sl[lenSl-1]
	s.preMax = sl[lenSl-2]
	s.min = sl[0]
	s.postMin = sl[1]
}

func (s *StrCriteria) getCriteria() (float64, float64) {
	res1 := float64(s.max-s.preMax) / float64(s.max-s.postMin)
	res2 := float64(s.postMin-s.min) / float64(s.preMax-s.min)
	return round(res1, 0.0005), round(res2, 0.0005)
}

func checkCrt(crt float64, serSize int) bool {
	var b bool
	if crt > checkCriteria[serSize] {
		b = true
	}
	return b
}

func getRelValFreqSl(mp map[int]int, uniqueSlice []int, serSize int) []float64 {
	res := make([]float64, 0, len(uniqueSlice))
	for _, elem := range uniqueSlice {
		ans := mp[elem]
		res = append(res, round((float64(ans)/float64(serSize))*100, 0.005))
	}
	return res
}

/*
	образцы вариационного ряда -

11, 8, 9, 10, 8, 6, 7, 7, 9, 11, 10, 6, 5, 11, 10, 7, 9, 11, 10
1, 2, 1, 2, 2, 4, 3, 3, 25, 2, 1, 25
5, 8, 8, 9, 9, 5, 8, 5, 2, 7, 9, 5, 2, 10, 8, 7, 5, 7, 3, 9, 5, 3, 9, 7
5, 8, 8, 9, 9, 5, 8, 5, 2, 7, 9, 5, 2, 10, 8, 5, 8, 8, 9, 9, 5, 8, 5, 2, 7, 9, 5, 2, 10, 8, 9
*/
func main() {
	fmt.Print("Введите вариационный ряд из положительных чисел.\n")
	fmt.Println("Обратите внимание, как должны отделяться числа друг от друга, например\n1, 2, 3")
	txt := readInput()  // считываем вариационный ряд как строку
	start := time.Now() // начало отсчета времени выполнения
	chUniqueSlice := make(chan []int)
	chRelValFreqSlice := make(chan []float64)
	intSlice, err := getSlice(txt) // преобразуем в интовый слайс
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sort.Ints(intSlice)
	serSize := len(intSlice)          // serSize - это объем выборки
	mp := fillMap(intSlice)           // заполняем мапу, ключ: вариант, значение: частота его повторения
	uniqueSlice := getUniqueSlice(mp) // получаем массив уникальных вариантов (которые не повторяются)
	sort.Ints(uniqueSlice)
	relValFreqSlice := getRelValFreqSl(mp, uniqueSlice, serSize)
	go func() {
		for i := 0; i < 3; i++ {
			chUniqueSlice <- uniqueSlice
			chRelValFreqSlice <- relValFreqSlice
		}
	}()
	// подача данных в рисовальщик гистограммы
	var wg1 sync.WaitGroup
	wg1.Add(3)
	go func(wg1 *sync.WaitGroup) {
		defer wg1.Done()
		a1 := <-chUniqueSlice
		a2 := <-chRelValFreqSlice
		errDrawBarChart := drawBarChart(a1, a2)
		if errDrawBarChart != nil {
			fmt.Println("This is DrawBarChart error", errDrawBarChart)
			os.Exit(1)
		}
	}(&wg1)
	// подача данных в эксель
	go func(wg1 *sync.WaitGroup) {
		defer wg1.Done()
		a1 := <-chUniqueSlice
		a2 := <-chRelValFreqSlice
		errExcel := createExcelTable(a1, mp, a2)
		if errExcel != nil {
			fmt.Println("This is Excel error", errExcel)
			os.Exit(1)
		}
	}(&wg1)
	// подача далее в рисовальщик полигона распределения
	go func(wg1 *sync.WaitGroup) {
		defer wg1.Done()
		a1 := <-chUniqueSlice
		a2 := <-chRelValFreqSlice
		line := drawLine(a1, a2)
		errRenderLine := renderLine(line)
		if errRenderLine != nil {
			fmt.Println("This is RenderLine error")
			os.Exit(1)
		}
	}(&wg1)
	wg1.Wait()
	// подача закончена
	fmt.Println("ВНИМАНИЕ, ОТВЕТ!")
	fmt.Printf("Объем выборки составил: %d.", serSize)
	fmt.Println("\nТаблица `вариант:частота:относительная частота варианта в процентах` сконвертирована программой в Excel.")
	fmt.Println("Гистограмма данных и полигон распределения готовы!")
	if serSize > 30 {
		fmt.Println("Так как объем выборки больше 30, сравнивать рассчитанные критерии с табличными некорректно!")
		finish := time.Since(start)
		fmt.Println("Программа выполнена! Время выполнения составило", finish)
		os.Exit(1)
	}
	var wg sync.WaitGroup
	wg.Add(1)
	//проверка критериев
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		structure3 := newStore() // конструктор структуры типа *StrCriteria
		structure3.fillCriteria(uniqueSlice)
		crt1, crt2 := structure3.getCriteria() // получаем высчитанные ранее критерии
		b1 := checkCrt(crt1, serSize)          // проверяем первый критерий
		if b1 {
			fmt.Printf("Критерий К1 = %.4f > табличного критерия Кт = %.4f для объема выборки %v,\nпоэтому следует исключить вариант %v из вариационного ряда и прогнать обновленный ряд через программу.", crt1, checkCriteria[serSize], serSize, structure3.max)
		}
		b2 := checkCrt(crt2, serSize) // проверяем второй критерий
		if b2 {
			fmt.Printf("Критерий К2 = %.4f > табличного критерия Кт = %.4f для объема выборки %v,\nпоэтому следует исключить вариант %v из вариационного ряда и прогнать обновленный ряд через программу.", crt2, checkCriteria[serSize], serSize, structure3.min)
		}
		if !b1 && !b2 {
			fmt.Print("Исключать варианты из вариационного ряда не нужно, ")
			fmt.Printf("т.к. твой критерий К1 = %.4f и К2 = %.4f меньше табличого %.4f.", crt1, crt2, checkCriteria[serSize])
		}
		finish := time.Since(start)
		fmt.Println("\nВремя выполнения программы составило", finish)
	}(&wg)
	wg.Wait()
}
