package main

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"testing"
)

// проверка на возможные невалидные данные
func TestGetSlice(t *testing.T) {
	err1 := errors.New("int's conversion error")
	err2 := errors.New("the size of series isn't enough")
	err3 := errors.New("series contains only zero")
	err4 := errors.New("contains negative numbers")
	tests := []struct {
		got string
		sl  []int
		err error
	}{
		{"1, 2, 3,", []int{}, err1},
		{"1, 2, 3", []int{}, err2},
		{"1, -2, 3, 4, 5, 6", []int{}, err4},
		{"google, 2, 3, 4, 5", []int{}, err1},
		{"1, 2, 4, !2", []int{}, err1},
		{"1; 2; 3", []int{}, err1},
		{"2, 2", []int{}, err2},
		{"'page, 0g, jb'", []int{}, err1},
		{"0, 0, 0, 0", []int{}, err3},
		{"0, 0", []int{}, err3},
		{"1, 2, 6, 5, 4, -4, 4, 1, 2, 6, j", []int{}, err4},
		{"1, 2, 2, 4, 5, 6", []int{1, 2, 2, 4, 5, 6}, nil},
		{"1, 10, 10, 2, 2, 2, 2", []int{1, 10, 10, 2, 2, 2, 2}, nil},
	}
	for idx, test := range tests {
		name := fmt.Sprintf("CASE %d, [%s], want %d, error '%s'", idx+1, test.got, test.sl, test.err)
		t.Run(name, func(t *testing.T) {
			got, err := getSlice(test.got)
			switch err {
			case nil:
				if !reflect.DeepEqual(got, test.sl) {
					t.Errorf("in CASE %d got %d, want %d", idx, got, test.sl)
				}
			default:
				if err.Error() != test.err.Error() {
					t.Errorf("in CASE %d got %d, %s, want %d, %s", idx, got, err.Error(), test.sl, test.err.Error())
				}
			}
		})
	}
}

type Str struct {
	series     int
	freq       int
	relVarFreq float64
}

type SetFull struct {
	set []Str
}

// ф-ция для заполнения общей структуры значений
func runFillGetStr(mp map[int]int, uniqueSlice []int, relValFreqSl []float64) *SetFull {
	structure2 := SetFull{}
	for i := range uniqueSlice {
		ans := mp[uniqueSlice[i]]
		structure1 := Str{series: uniqueSlice[i], freq: ans, relVarFreq: relValFreqSl[i]}
		structure2.set = append(structure2.set, structure1)
	}
	return &structure2
}

// функция, аналогичная getRelValFreq в ф-ции main, но без округления результата - для корректного сравнения
func getRelValFreqSlAnalog(mp map[int]int, uniqueSlice []int, serSize int) []float64 {
	res := make([]float64, 0, len(uniqueSlice))
	for _, elem := range uniqueSlice {
		ans := mp[elem]
		res = append(res, 100*(float64(ans)/float64(serSize)))
	}
	return res
}

func TestRunFillGetStr(t *testing.T) {

	tests := []struct {
		sl []int
	}{
		{[]int{11, 8, 9, 10, 8, 6, 7, 7, 9, 11, 10, 6, 5, 11, 10, 7, 9, 11, 10}},
		{[]int{1, 2, 2, 0, 0, 4, 17, 17, 23, 44, 12, 22, 24}},
		{[]int{5, 8, 8, 9, 9, 5, 8, 5, 2, 7, 9, 5, 2, 10, 8, 7, 5, 7, 3, 9, 5, 3, 9, 7}},
		{[]int{5, 8, 8, 9, 9, 5, 8, 5, 2, 7, 9, 5, 2, 10, 8, 5, 8, 8, 9, 9, 5, 8, 5, 2, 7, 9, 5, 2, 10, 8, 9}},
		{[]int{5, 8, 8, 9, 9, 5, 8, 5, 2, 7, 9, 5, 2, 10, 8, 5, 8, 8, 9, 9, 5, 8, 5, 2, 7, 9, 5, 2, 10, 8, 9, 5, 8, 8, 9, 9, 5, 8, 5, 2, 7, 9, 5, 2, 10, 8, 5, 8, 8, 9, 9, 5, 8, 5, 2, 7, 9, 5, 2, 10, 8, 9, 5, 8, 8, 9, 9, 5, 8, 5, 2, 7, 9, 5, 2, 10, 8, 5, 8, 8, 9, 9, 5, 8, 5, 2, 7, 9, 5, 2, 10, 8, 9, 23, 22, 12, 2, 3, 1, 11}},
	}

	for idx, test := range tests {
		name := fmt.Sprintf("CASE %d, serSize %d, got slice %d", idx+1, len(test.sl), test.sl)
		t.Run(name, func(t *testing.T) {
			sumRelValFreq := 0.00
			sumSerSize := 0
			sumElemSlice := 0
			f := func(sl []int) int {
				sum := 0
				for i := range sl {
					sum += sl[i]
				}
				return sum
			}
			sumElemSliceFirst := f(test.sl)
			mp := fillMap(test.sl)
			uniqueSl := getUniqueSlice(mp)
			serSize := len(test.sl)
			sort.Ints(uniqueSl)
			relValFreq := getRelValFreqSlAnalog(mp, uniqueSl, serSize)
			got := runFillGetStr(mp, uniqueSl, relValFreq)
			for i, elem := range got.set {
				sumRelValFreq += elem.relVarFreq
				sumSerSize += 1 * elem.freq
				sumElemSlice += elem.series * elem.freq
				newRes := float64(elem.freq) / float64(serSize)
				if elem.relVarFreq != 100*newRes {
					t.Errorf("in case %d got relValFreq %f want %f", i, elem.relVarFreq, 100*newRes)
				}
			}
			switch {
			case int(sumRelValFreq) != 100:
				t.Errorf("in case %d got sumRelValFreq %.2f, want %.2f", idx, sumRelValFreq, 100.00)
			case sumSerSize != serSize:
				t.Errorf("in case %d got sumSerSize %d, want %d", idx, sumSerSize, serSize)
			case sumElemSlice != sumElemSliceFirst:
				t.Errorf("in case %d got sumElemSlice %d, want %d", idx, sumElemSlice, sumElemSliceFirst)
			}
		})
	}

}
