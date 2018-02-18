package container

import (
	. "github.com/fishedee/assert"
	"math/rand"
	"testing"
)

func TestTrieAllMatch(t *testing.T) {
	trieTree := NewTrieTree()

	testFoundData := []string{"a", "ab", "ba", "1", "2", "3c_4", "你", "a你", "你好", "gm", "g"}
	testNotFoundData := []string{"c", "abc", "ac", "b", "bac", "_4", "c_", "a ", " a", "你a", "好"}

	for _, data := range testFoundData {
		trieTree.Set(data, data+"_value")
	}

	trie := trieTree.ToTrieArray()
	for _, data := range testFoundData {
		result := trie.ExactMatch(data)
		AssertEqual(t, result, data+"_value")
		result2 := trieTree.Get(data)
		AssertEqual(t, result2, data+"_value")
	}
	for _, data := range testNotFoundData {
		result := trie.ExactMatch(data)
		AssertEqual(t, result, nil)
		result2 := trieTree.Get(data)
		AssertEqual(t, result2, nil)
	}
}

func TestTriePrefixMatch(t *testing.T) {
	trieTree := NewTrieTree()

	testData := []string{"ab", "/ab", "cde", "abeg", "ac", "a你", "你", "你好", "你去", "hm", "hmgch", "hmgbd"}
	testFindData := map[string][]TrieMatch{
		"ab": []TrieMatch{
			{"ab", "ab_value"},
		},
		"abc": []TrieMatch{
			{"ab", "ab_value"},
		},
		"/abc": []TrieMatch{
			{"/ab", "/ab_value"},
		},
		"cdef/g": []TrieMatch{
			{"cde", "cde_value"},
		},
		"abegh": []TrieMatch{
			{"ab", "ab_value"},
			{"abeg", "abeg_value"},
		},
		"ac": []TrieMatch{
			{"ac", "ac_value"},
		},
		"acm": []TrieMatch{
			{"ac", "ac_value"},
		},
		"a你": []TrieMatch{
			{"a你", "a你_value"},
		},
		"你": []TrieMatch{
			{"你", "你_value"},
		},
		"你a": []TrieMatch{
			{"你", "你_value"},
		},
		"你好": []TrieMatch{
			{"你", "你_value"},
			{"你好", "你好_value"},
		},
		"你好吗": []TrieMatch{
			{"你", "你_value"},
			{"你好", "你好_value"},
		},
		"你去": []TrieMatch{
			{"你", "你_value"},
			{"你去", "你去_value"},
		},
		"a":   []TrieMatch{},
		"/a":  []TrieMatch{},
		"cck": []TrieMatch{},
		"cd":  []TrieMatch{},
		"hmgb": []TrieMatch{
			{"hm", "hm_value"},
		},
	}

	for _, data := range testData {
		trieTree.Set(data, data+"_value")
	}

	trie := trieTree.ToTrieArray()
	for key, value := range testFindData {
		result := trie.PrefixMatch(key)
		AssertEqual(t, result, value)

		resultKey2, resultValue2 := trie.LongestPrefixMatch(key)
		var result2 interface{}
		if resultValue2 != nil {
			result2 = TrieMatch{resultKey2, resultValue2}
		} else {
			result2 = nil
		}
		var value2 interface{}
		if len(value) == 0 {
			value2 = nil
		} else {
			value2 = value[len(value)-1]
		}

		AssertEqual(t, result2, value2)

	}
}

func TestTrieFull(t *testing.T) {
	trieTree := NewTrieTree()

	testData := []string{"", "abc", "bc", "bcd"}
	testFindData := map[string][]TrieMatch{
		"": []TrieMatch{
			{"", "_value"},
		},
		"abc": []TrieMatch{
			{"", "_value"},
			{"abc", "abc_value"},
		},
		"abcg": []TrieMatch{
			{"", "_value"},
			{"abc", "abc_value"},
		},
		"b": []TrieMatch{
			{"", "_value"},
		},
		"bc": []TrieMatch{
			{"", "_value"},
			{"bc", "bc_value"},
		},
		"bcg": []TrieMatch{
			{"", "_value"},
			{"bc", "bc_value"},
		},
		"bcd": []TrieMatch{
			{"", "_value"},
			{"bc", "bc_value"},
			{"bcd", "bcd_value"},
		},
		"c": []TrieMatch{
			{"", "_value"},
		},
	}

	for _, data := range testData {
		trieTree.Set(data, data+"_value")
	}

	trie := trieTree.ToTrieArray()
	for key, value := range testFindData {
		result := trie.PrefixMatch(key)
		AssertEqual(t, result, value, key)

		_, result2 := trie.LongestPrefixMatch(key)
		var value2 interface{}
		if len(value) == 0 {
			value2 = nil
		} else {
			value2 = value[len(value)-1].Value
		}
		AssertEqual(t, result2, value2, key)
	}
}

func TestTrieWalk(t *testing.T) {
	trieTree := NewTrieTree()

	invalidUtf8 := string([]byte{97, 228, 189, 160, 229})
	testData := []string{"", "abc", "bc", "bcd", "b你", "a你好", "a你去"}
	walkData := []struct {
		key         string
		value       interface{}
		parentKey   string
		parentValue interface{}
	}{
		{"", "_value", "", nil},
		{"a", nil, "", "_value"},
		{"b", nil, "", "_value"},
		{"abc", "abc_value", "a", nil},
		{invalidUtf8, nil, "a", nil},
		{"bc", "bc_value", "b", nil},
		{"b你", "b你_value", "b", nil},
		{"a你去", "a你去_value", invalidUtf8, nil},
		{"a你好", "a你好_value", invalidUtf8, nil},
		{"bcd", "bcd_value", "bc", "bc_value"},
	}

	for _, data := range testData {
		trieTree.Set(data, data+"_value")
	}

	index := 0
	trieTree.Walk(func(key string, value interface{}, parentKey string, parentValue interface{}) {
		result := walkData[index]
		index++
		AssertEqual(t, result.key, key, index)
		AssertEqual(t, result.value, value, index)
		AssertEqual(t, result.parentKey, parentKey, index)
		AssertEqual(t, result.parentValue, parentValue, index)
	})
	AssertEqual(t, index, len(walkData))
}

func getSingleData(count int) string {
	var randStr = []byte("0123456789abcdefghijklmnopqrstuvwxyz")
	result := make([]byte, count)
	rand.Read(result)
	for singleIndex, singleByte := range result {
		result[singleIndex] = randStr[int(singleByte)%len(randStr)]
	}
	return string(result)
}

func getData(prefix string, count int, size int) []string {
	result := []string{}
	for i := 0; i != count; i++ {
		result = append(result, prefix+getSingleData(size))
	}
	return result
}

func benchmarkTrieSpeed(prefix string, b *testing.B, length int) {
	insertData := getData(prefix, 1000, length)
	findData := getData(prefix, 1000, length)
	trieTree := NewTrieTree()
	for _, singleData := range insertData {
		trieTree.Set(singleData, true)
	}
	trie := trieTree.ToTrieArray()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		single := findData[i%len(findData)]
		trie.ExactMatch(single)
	}
}

func benchmarkMapSpeed(prefix string, b *testing.B, length int) {
	insertData := getData(prefix, 1000, length)
	findData := getData(prefix, 1000, length)
	mapper := map[string]bool{}
	for _, singleData := range insertData {
		mapper[singleData] = true
	}

	b.ResetTimer()
	for i := 0; i != b.N; i++ {
		single := findData[i%len(findData)]
		_, _ = mapper[single]
	}
}

func BenchmarkTrieSpeedShort(b *testing.B) {
	benchmarkTrieSpeed("", b, 20)
}

func BenchmarkMapSpeedShort(b *testing.B) {
	benchmarkMapSpeed("", b, 20)
}

func BenchmarkTrieSpeedLong(b *testing.B) {
	benchmarkTrieSpeed("", b, 100)
}

func BenchmarkMapSpeedLong(b *testing.B) {
	benchmarkMapSpeed("", b, 100)
}

func BenchmarkTrieSpeedPrefixLong(b *testing.B) {
	benchmarkTrieSpeed("/abcefgdsfa/abcefgdsfa/abcefgdsfa/abcefgdsfa/abcefgdsfa/abcefgdsfa/abcefgdsfa/abcefgdsfa/abcefgdsfa/abcefgdsfa/abcefgdsfa/", b, 100)
}

func BenchmarkMapSpeedPrefixLong(b *testing.B) {
	benchmarkMapSpeed("/abcefgdsfa/abcefgdsfa/abcefgdsfa/abcefgdsfa/abcefgdsfa/abcefgdsfa/abcefgdsfa/abcefgdsfa/abcefgdsfa/abcefgdsfa/abcefgdsfa/", b, 100)
}

func BenchmarkTrieSpeedPrefixLong2(b *testing.B) {
	benchmarkTrieSpeed("/abcefgdsfa/abcefgdsfa/abcefgdsfa/abcefgdsfa/abcefgdsfa/abcefgdsfa/abcefgdsfa/abcefgdsfa/abcefgdsfa/abcefgdsfa/abcefgdsfa/", b, 0)
}

func BenchmarkMapSpeedPrefixLong2(b *testing.B) {
	benchmarkMapSpeed("/abcefgdsfa/abcefgdsfa/abcefgdsfa/abcefgdsfa/abcefgdsfa/abcefgdsfa/abcefgdsfa/abcefgdsfa/abcefgdsfa/abcefgdsfa/abcefgdsfa/", b, 0)
}
