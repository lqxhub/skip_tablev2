package skiptablev2

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type S1[K comparable] struct {
	key K
	f   float64
}

func (s *S1[K]) Key() K {
	return s.key
}

func (s *S1[K]) Score() float64 {
	return s.f
}

type SortS1[K comparable] []*S1[K]

func (s SortS1[K]) Len() int {
	return len(s)
}

func (s SortS1[K]) Less(i, j int) bool {
	return s[i].f < s[j].f
}

func (s SortS1[K]) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

//func commonInsert[K comparable, V SkipListItem[K]](st *SkipList[K, V], num int) []*S1[K] {
//	arrS := make(SortS1[K], 0, num)
//	for i := 0; i < num; i++ {
//		f := rand.Float64()
//		s := &S1[K]{f: f}
//		st.InsertByScore(f, s)
//		arrS = append(arrS, s)
//	}
//	sort.Sort(arrS)
//	return arrS
//}

func compare[K comparable, V SkipListItem[K]](result []V, arr []*S1[K]) error {
	for i, v := range result {
		if arr[i].f != v.Score() {
			return errors.New(fmt.Sprintf("i:%d src:%v result:%v", i, arr[i], v.Score()))
		}
	}
	return nil
}

func TestNewSkipTable(t *testing.T) {
	st, err := NewDefaultSkipTable[string, *S1[string]](func(v1, v2 *S1[string]) int {
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

	N := 10000
	arrS := make(SortS1[string], 0, N)
	for i := 0; i < N; i++ {
		f := rand.Float64()
		s := &S1[string]{f: f}
		st.InsertByScore(f, s)
		arrS = append(arrS, s)
	}
	sort.Sort(arrS)

	size := st.Size()
	fmt.Println("st size:", size)

	fmt.Println()
	list := st.GetValuesByRank(1, int64(N))
	for i, node := range list {
		if arrS[i].f != node.f {
			fmt.Printf("i：%d  key:%s f:%f error", i, node.key, node.f)
		}
	}
	fmt.Println()
}

func TestSkipList_GetNodeByRank(t *testing.T) {
	st, err := NewDefaultSkipTable[string, *S1[string]](func(v1, v2 *S1[string]) int {
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

	N := 10000
	arrS := make(SortS1[string], 0, N)
	for i := 0; i < N; i++ {
		f := rand.Float64()
		s := &S1[string]{f: f}
		st.InsertByScore(f, s)
		arrS = append(arrS, s)
	}
	sort.Sort(arrS)

	size := st.Size()
	fmt.Println(size)

	l := rand.Int63n(int64(N)/2) + 1
	r := l + rand.Int63n(int64(N)/2)

	fmt.Printf("l:%d r:%d\n", l, r)

	result := st.GetValuesByRank(l, r)

	fmt.Println("result size:", len(result))
	err = compare(result, arrS[int(l)-1:])
	if err != nil {
		panic(err)
	}
	fmt.Println()
}

func TestSkipList_GetNodeByScore(t *testing.T) {
	st, err := NewDefaultSkipTable[string, *S1[string]](func(v1, v2 *S1[string]) int {
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
	rand.Seed(time.Now().UnixNano())

	const N = 10000
	arrS := make(SortS1[string], 0, N)
	for i := 0; i < N; i++ {
		f := rand.Float64()
		s := &S1[string]{f: f}
		st.InsertByScore(f, s)
		arrS = append(arrS, s)
	}
	sort.Sort(arrS)
	size := st.Size()
	fmt.Println(size)

	//测试随机范围
	n := rand.Int31n(N / 3)
	x := rand.Int31n(int32(N) - n)
	fmt.Println(n, x)

	start := arrS[x].f
	end := arrS[n+x].f

	fmt.Println("获取范围:", start, end)

	r := &SkipListFindRange{
		Min:    start,
		Max:    end,
		MinInf: true,
		MaxInf: false,
	}
	result := st.GetValuesByScore(r)

	fmt.Println("result size:", len(result))
	err = compare(result, arrS)
	if err != nil {
		panic(err)
	}
	fmt.Println("---------------------------------------------------------")

	//测试获取第一个
	fmt.Println("测试获取第一个")
	r.Min = arrS[0].Score()
	r.Max = arrS[0].Score()
	result = st.GetValuesByScore(r)
	fmt.Println("result size:", len(result))
	err = compare(result, arrS)
	if err != nil {
		panic(err)
	}
	fmt.Println("---------------------------------------------------------")

	//测试[-∞,+∞]
	fmt.Printf("测试 [-∞,+∞] [%f,%f]\n", arrS[0].Score(), arrS[N-1].Score())
	r.MinInf = true
	r.MaxInf = true
	result = st.GetValuesByScore(r)
	fmt.Println("result size:", len(result))
	err = compare(result, arrS)
	if err != nil {
		panic(err)
	}
	fmt.Println("正常", N, "实际", len(result), "个")
	fmt.Println("---------------------------------------------------------")

	//测试[-∞,n]
	n = rand.Int31n(N)
	fmt.Printf("测试 [-∞,%d] [%f,%f]\n", n, arrS[0].Score(), arrS[n-1].Score())
	r.MaxInf = false
	r.MinInf = true
	r.Max = arrS[n-1].Score()
	r.Min = 0
	result = st.GetValuesByScore(r)
	fmt.Println("result size:", len(result))
	err = compare(result, arrS)
	if err != nil {
		panic(err)
	}
	fmt.Println("正常", n, "实际", len(result), "个")
	fmt.Println("---------------------------------------------------------")

	//测试 [n,+∞]
	n = rand.Int31n(N)
	fmt.Printf("测试 [%d,+∞] [%f,%f]\n", n, arrS[n].Score(), arrS[N-1].Score())
	r.MaxInf = true
	r.MinInf = false
	r.Max = 0
	r.Min = arrS[n-1].Score()

	result = st.GetValuesByScore(r)
	fmt.Println("result size:", len(result))
	err = compare(result, arrS[n-1:])
	if err != nil {
		panic(err)
	}
	fmt.Println("正常", N-n+1, "实际", len(result), "个")
}

func TestSkipList_UpdateScore(t *testing.T) {
	st, err := NewDefaultSkipTable[string, *S1[string]](func(v1, v2 *S1[string]) int {
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

	const N = 10000
	arrS := make(SortS1[string], 0, N)
	for i := 0; i < N; i++ {
		f := rand.Float64()
		s := &S1[string]{f: f}
		st.InsertByScore(f, s)
		arrS = append(arrS, s)
	}
	sort.Sort(arrS)

	result := st.GetValuesByRank(1, N)

	err = compare(result, arrS)
	if err != nil {
		panic(err)
	}

	//更新的位置
	n := rand.Int63n(N)
	nodes := st.GetNodesByRank(n, n)
	if len(nodes) == 0 {
		panic("sl GetNodeByRank 1 not find")
	}
	newV := rand.Float64()
	arrS[n-1].f = newV
	sort.Sort(arrS)
	st.UpdateScore(nodes[0], newV)

	result = st.GetValuesByRank(1, N)
	err = compare(result, arrS)
	if err != nil {
		panic(err)
	}
}

func TestSkipList_GetNodeRank(t *testing.T) {
	st, err := NewDefaultSkipTable[string, *S1[string]](func(v1, v2 *S1[string]) int {
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

	const N = 10000
	arrS := make(SortS1[string], 0, N)
	for i := 0; i < N; i++ {
		f := rand.Float64()
		s := &S1[string]{f: f}
		st.InsertByScore(f, s)
		arrS = append(arrS, s)
	}
	sort.Sort(arrS)
	nodes := st.GetNodesByRank(1, N)

	fmt.Println("nodes size:", len(nodes))

	n := rand.Int63n(N)

	fmt.Println("n:", n+1)

	fmt.Println("rank:", st.GetNodeRank(nodes[n]))
}

func TestSkipList_Delete(t *testing.T) {
	st, err := NewDefaultSkipTable[string, *S1[string]](func(v1, v2 *S1[string]) int {
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

	const N = 10000
	for i := 0; i < N; i++ {
		f := rand.Float64()
		s := &S1[string]{f: f}
		st.InsertByScore(f, s)
	}
	fmt.Println("st size ", st.Size())

	ranks := st.GetNodesByRank(1, st.Size())
	fmt.Println("ranks len ", len(ranks))

	rand.Shuffle(int(st.Size()), func(i, j int) {
		ranks[i], ranks[j] = ranks[j], ranks[i]
	})
	for _, node := range ranks {
		st.Delete(node, st.GetUpdateList(node))
	}

	fmt.Println(st.Size())
}
