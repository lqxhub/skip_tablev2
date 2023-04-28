package skiptablev2

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"testing"
	"time"
)

const N = 100

func init() {
	rand.Seed(time.Now().UnixNano())
}

type StItem[K comparable] struct {
	f float64
	k string
}

func (s *StItem[K]) Key() string {
	return s.k
}

func (s *StItem[K]) Score() float64 {
	return s.f
}

type SortStItem[K comparable] []*StItem[K]

func (s SortStItem[K]) Len() int {
	return len(s)
}

func (s SortStItem[K]) Less(i, j int) bool {
	return s[i].f < s[j].f
}

func (s SortStItem[K]) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func compareItem(item1, item2 *StItem[string]) bool {
	return item1.k == item2.k
}

func CreateStItem() *StItem[string] {
	f := rand.Float64()
	return &StItem[string]{
		f: f,
		k: strconv.FormatFloat(f, 'f', -1, 64),
	}
}

func NewTestSortSet() *SortSet[string, *StItem[string]] {
	sortSet, err := NewDefaultSortSet[string, *StItem[string]](func(v1, v2 *StItem[string]) int {
		if v1.f == v2.f {
			return 0
		} else if v1.f < v2.f {
			return -1
		} else {
			return 1
		}
	})
	if err != nil {
		panic(err)
	}
	return sortSet
}

// ////////////////////////////////////////////////////////////////
func TestNewSortSet(t *testing.T) {
	sortSet := NewTestSortSet()
	fmt.Println(sortSet.Count())
}

func TestSortSet_Add(t *testing.T) {
	sortSet := NewTestSortSet()
	fmt.Println(sortSet.Count())

	for i := 0; i < N; i++ {
		sortSet.Add(CreateStItem())
	}

	fmt.Println(sortSet.Count())
}

func TestSortSet_Range(t *testing.T) {
	sortSet := NewTestSortSet()
	arr := make(SortStItem[string], 0, N)

	for i := 0; i < N; i++ {
		item := CreateStItem()
		sortSet.Add(item)
		arr = append(arr, item)
	}
	sort.Sort(arr)
	result := sortSet.Range(0, -1)
	fmt.Printf("l:%d,r:%d result len:%d \n", 0, -1, len(result))
	for i, r := range result {
		if !compareItem(arr[i], r) {
			fmt.Printf("item:%v rank:%d arrItem:%v error\n", r, i, arr[i])
		}
	}

	for i := 0; i < 10; i++ {
		perm := rand.Perm(N)
		l := int64(perm[0])
		r := int64(perm[1])
		if l > r {
			l, r = r, l
		}
		result := sortSet.Range(l, r)
		fmt.Printf("l:%d,r:%d result len:%d \n", l, r, len(result))
		for i, r := range arr[l:r] {
			if !compareItem(result[i], r) {
				fmt.Printf("item:%v rank:%d arrItem:%v error\n", r, i, arr[i])
			}
		}
	}
}

func TestSortSet_RevRange(t *testing.T) {
	sortSet := NewTestSortSet()
	arr := make(SortStItem[string], 0, N)

	for i := 0; i < N; i++ {
		item := CreateStItem()
		sortSet.Add(item)
		arr = append(arr, item)
	}
	sort.Sort(sort.Reverse(arr))
	result := sortSet.RevRange(0, -1)
	fmt.Printf("l:%d,r:%d result len:%d\n", 0, -1, len(result))
	for i, r := range result {
		if !compareItem(arr[i], r) {
			fmt.Printf("item:%v rank:%d arrItem:%v error\n", r, i, arr[i])
		}
	}

	for i := 0; i < 10; i++ {
		perm := rand.Perm(N)
		l := int64(perm[0])
		r := int64(perm[1])
		if l > r {
			l, r = r, l
		}
		result := sortSet.RevRange(l, r)
		fmt.Printf("l:%d,r:%d result len:%d \n", l, r, len(result))
		for i, r := range arr[l:r] {
			if !compareItem(result[i], r) {
				fmt.Printf("item:%v rank:%d arrItem:%v error\n", r, i, arr[i])
			}
		}
	}
}

func TestSortSet_Rank(t *testing.T) {
	sortSet := NewTestSortSet()

	arr := make(SortStItem[string], 0, N)
	for i := 0; i < N; i++ {
		item := CreateStItem()
		arr = append(arr, item)
		sortSet.Add(item)
	}

	sort.Sort(arr)

	fmt.Println("sortSet len", sortSet.Count())
	perm := rand.Perm(arr.Len())

	for _, i := range perm {
		item := arr[i]
		rank := sortSet.Rank(item.k)
		if int64(i) != rank {
			fmt.Printf("item:%v rank:%d arrItem:%v error\n", item, i, arr[i])
		}
	}
}

func TestSortSet_RangeByScore(t *testing.T) {
	sortSet := NewTestSortSet()

	arr := make(SortStItem[string], 0, N)
	for i := 0; i < N; i++ {
		item := CreateStItem()
		arr = append(arr, item)
		sortSet.Add(item)
	}

	sort.Sort(arr)

	findRange := &SkipListFindRange{
		Min:    0,
		Max:    0,
		MinInf: true,
		MaxInf: true,
	}

	result := sortSet.RangeByScore(findRange)
	fmt.Printf("测试随机取范围 [-∞,+∞] result len:%d \n", len(result))
	for i, r := range result {
		if !compareItem(r, arr[i]) {
			fmt.Printf("item:%v rank:%d arrItem:%v error\n", r, i, arr[i])
		}
	}

	//测试范围
	n := rand.Int31n(N)
	findRange.MinInf = true
	findRange.Min = 0
	findRange.MaxInf = false
	findRange.Max = arr[n].Score()
	result = sortSet.RangeByScore(findRange)
	fmt.Printf("测试随机取范围 [-∞,%f] result len:%d \n", arr[n].Score(), len(result))

	for i, r := range arr[:n] {
		if !compareItem(r, result[i]) {
			fmt.Printf("item:%v rank:%d arrItem:%v error\n", r, i, arr[i])
		}
	}

	//测试范围
	n = rand.Int31n(N)
	findRange.MinInf = false
	findRange.Min = arr[n].Score()
	findRange.MaxInf = true
	findRange.Max = 0
	result = sortSet.RangeByScore(findRange)
	fmt.Printf("测试随机取范围 [%f,+∞] result len:%d \n", arr[n].Score(), len(result))
	for i, r := range arr[n:] {
		if !compareItem(r, result[i]) {
			fmt.Printf("item:%v rank:%d arrItem:%v error\n", r, i, arr[i])
		}
	}

	for i := 0; i < 10; i++ {
		perm := rand.Perm(N)
		l := perm[0]
		r := perm[1]
		if l > r {
			l, r = r, l
		}

		findRange.MinInf = false
		findRange.Min = arr[l].Score()
		findRange.MaxInf = false
		findRange.Max = arr[r].Score()
		result = sortSet.RangeByScore(findRange)
		fmt.Printf("测试随机取范围 [%f,%f] result len:%d\n", arr[l].Score(), arr[r].Score(), len(result))
		for i, r := range arr[l:r] {
			if !compareItem(r, result[i]) {
				fmt.Printf("item:%v rank:%d arrItem:%v error\n", r, i, arr[i])
			}
		}
	}
}

func TestSortSet_RevRangeByScore(t *testing.T) {
	sortSet := NewTestSortSet()

	arr := make(SortStItem[string], 0, N)
	for i := 0; i < N; i++ {
		item := CreateStItem()
		arr = append(arr, item)
		sortSet.Add(item)
	}

	sort.Sort(arr)

	findRange := &SkipListFindRange{
		Min:    0,
		Max:    0,
		MinInf: true,
		MaxInf: true,
	}

	sort.Sort(sort.Reverse(arr))
	result := sortSet.RevRangeByScore(findRange)

	fmt.Printf("测试随机取范围 [-∞,+∞] resutl len:%d\n", len(result))
	for i, r := range result {
		if !compareItem(r, arr[i]) {
			fmt.Printf("item:%v rank:%d arrItem:%v error\n", r, i, arr[i])
		}
	}

	//测试范围
	n := rand.Int31n(N)
	findRange.MinInf = true
	findRange.Min = 0
	findRange.MaxInf = false
	findRange.Max = arr[n].Score()
	result = sortSet.RevRangeByScore(findRange)
	fmt.Printf("测试随机取范围 [-∞,%f] resutl len:%d \n", arr[n].Score(), len(result))
	for i, r := range arr[:n] {
		if !compareItem(r, result[i]) {
			fmt.Printf("item:%v rank:%d arrItem:%v error\n", r, i, arr[i])
		}
	}

	//测试范围
	n = rand.Int31n(N)
	findRange.MinInf = false
	findRange.Min = arr[n].Score()
	findRange.MaxInf = true
	findRange.Max = 0
	result = sortSet.RevRangeByScore(findRange)
	fmt.Printf("测试随机取范围 [%f,+∞] resutl len:%d \n", arr[n].Score(), len(result))
	for i, r := range arr[n:] {
		if !compareItem(r, result[i]) {
			fmt.Printf("item:%v rank:%d arrItem:%v error\n", r, i, arr[i])
		}
	}

	for i := 0; i < 10; i++ {
		perm := rand.Perm(N)
		l := perm[0]
		r := perm[1]
		if l > r {
			l, r = r, l
		}

		findRange.MinInf = false
		findRange.Min = arr[l].Score()
		findRange.MaxInf = false
		findRange.Max = arr[r].Score()
		result = sortSet.RevRangeByScore(findRange)
		fmt.Printf("测试随机取范围 [%f,%f] resutl len:%d \n", arr[l].Score(), arr[r].Score(), len(result))
		for i, r := range arr[l:r] {
			if !compareItem(r, result[i]) {
				fmt.Printf("item:%v rank:%d arrItem:%v error\n", r, i, arr[i])
			}
		}
	}
}

func TestSortSet_RemoveRangeByRank(t *testing.T) {
	sortSet := NewTestSortSet()
	arr := make(SortStItem[string], 0, N)
	for i := 0; i < N; i++ {
		item := CreateStItem()
		arr = append(arr, item)
		sortSet.Add(item)
	}

	sort.Sort(arr)

	var deleteArr = func(i int) {
		arr = append(arr[:i], arr[i+1:]...)
	}

	for len(arr) > 0 {
		n := rand.Int63n(int64(len(arr)))
		sortSet.RemoveRangeByRank(n, n)
		deleteArr(int(n))
		result := sortSet.Range(0, -1)
		fmt.Printf("delete rank:%d sortSet len:%d \n", n, sortSet.Count())
		for i, r := range result {
			if !compareItem(arr[i], r) {
				fmt.Printf("item:%v rank:%d arrItem:%v error\n", r, i, arr[i])
			}
		}
	}
}

func TestSortSet_RemoveRangeByScore(t *testing.T) {
	sortSet := NewTestSortSet()
	arr := make(SortStItem[string], 0, N)
	for i := 0; i < N; i++ {
		item := CreateStItem()
		arr = append(arr, item)
		sortSet.Add(item)
	}

	sort.Sort(arr)

	var deleteArr = func(i int) {
		arr = append(arr[:i], arr[i+1:]...)
	}

	for len(arr) > 0 {
		n := rand.Int63n(int64(len(arr)))
		score := arr[n].Score()
		sortSet.RemoveRangeByScore(score, score)
		deleteArr(int(n))
		result := sortSet.Range(0, -1)
		fmt.Printf("delete rank:%d socre:%f sortSet len:%d \n", n, score, sortSet.Count())
		for i, r := range result {
			if !compareItem(arr[i], r) {
				fmt.Printf("item:%v rank:%d arrItem:%v error\n", r, i, arr[i])
			}
		}
	}
}

func BenchmarkSortSet_RankBench(b *testing.B) {
	const SIZE = 100000
	arr := make(SortStItem[string], 0, SIZE)
	sortSet := NewTestSortSet()
	for i := 0; i < SIZE; i++ {
		item := CreateStItem()
		sortSet.Add(item)
		arr = append(arr, item)
	}
	for i := 0; i < b.N; i++ {
		n := rand.Int63n(SIZE)
		sortSet.Rank(arr[n].k)
	}
}

func BenchmarkSortSet_RankBenchArr(b *testing.B) {
	const SIZE = 100000
	arr := make(SortStItem[string], 0, SIZE)
	for i := 0; i < SIZE; i++ {
		item := CreateStItem()
		arr = append(arr, item)
	}
	for i := 0; i < b.N; i++ {
		n := rand.Int63n(SIZE)
		sort.Search(SIZE, func(i int) bool {
			return arr[i].f <= arr[n].f
		})
	}
}
