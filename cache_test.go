package goswift

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	// "github.com/leoantony72/goswift"
	"github.com/leoantony72/goswift/expiry"
)

// const (
// 	ErrKeyNotFound  = "key does not Exists"
// 	ErrNotHashvalue = "not a Hash value/table"
// )

func TestSet(t *testing.T) {
	cache := NewCache()

	key := "name"
	val := "leoantony"
	cache.Set(key, val, 0)

	getValue, err := cache.Get(key)
	if err != nil {
		if err.Error() == ErrKeyNotFound {
			t.Errorf("key `%s`: %s", key, ErrKeyNotFound)
			return
		}
		return
	}
	if getValue.(string) != val {
		t.Errorf("val not the same")
		return
	}
}

func TestGet(t *testing.T) {
	c := NewCache()
	c.Set("age", 12, 0)

	val, err := c.Get("age")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	if val.(int) != 12 {
		t.Errorf("Expected Value: 12(int) ,Gotten: %d", val)
		return
	}

	// Key does not exists
	_, err = c.Get("name")
	if err == nil {
		t.Errorf("Expected Error: %s", ErrKeyNotFound)
		return
	}

	//expiry provided- expiry>>time.Now()
	c.Set("place", "Kerala", 150000)
	val, err = c.Get("place")
	if err != nil {
		t.Errorf(err.Error())
		t.Errorf("key %s Might be expired", "place")
		return
	}
	if val.(string) != "Kerala" {
		t.Errorf("Expected Value: %s, Gotten: %s", "Kerala", val.(string))
		return
	}

	c.Set("country", "India", 100)
	_, err = c.Get("country")
	if err != nil {
		if err.Error() != ErrKeyNotFound {
			t.Errorf("key %s Should be Expired", "country")
			return
		}
	}

}

func TestUpdate(t *testing.T) {
	c := NewCache()

	key := "users:bob"
	value := "Cool shirt"
	c.Set(key, value, 0)

	data, err := c.Get(key)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	if data.(string) != value {
		t.Errorf("Expected Value: %s ,Gotten: %s", value, data)
		return
	}

	newValue := "Chemistry sucks"
	c.Update(key, newValue)

	data, err = c.Get(key)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	if data.(string) != newValue {
		t.Errorf("Expected Value: %s ,Gotten: %s", newValue, data)
		return
	}

	//key does not exist
	key = "water"
	err = c.Update(key, "H2O")
	if err == nil {
		t.Errorf("Expected Err: %s, Gotten: ERR NIL", ErrKeyNotFound)
		return
	}

	if err.Error() != ErrKeyNotFound {
		t.Errorf("Expected Err: %s, Gotten: %s", ErrKeyNotFound, err.Error())
		return
	}
}

func TestDel(t *testing.T) {
	c := NewCache()
	key := "users:bob"
	value := "Cool shirt"
	c.Set(key, value, 0)

	ok := c.Exists(key)
	if !ok {
		t.Errorf("Expected Value: %v, Gotten: %v", true, ok)
		return
	}

	c.Del(key)
	ok = c.Exists(key)
	if ok {
		t.Errorf("Expected Value: %v, Gotten: %v", false, ok)
		return
	}

	//Key does not exist
	key = "users:varun"
	c.Del(key)

	// Key with Expiry
	key = "users:kingbob"
	c.Set(key, "bobbb!", 10000)
	c.Del(key)
}
func TestHset(t *testing.T) {
	c := NewCache()
	key := "users:John:metadata"
	c.Hset(key, "name", "John")
	c.Hset(key, "age", 20)
	c.Hset(key, "place", "Thrissur")
	c.Hset(key, "people", []string{"bob", "tony", "henry"})

	data, err := c.HGetAll(key)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	name := data["name"].(string)
	age := data["age"].(int)
	place := data["place"].(string)
	people := data["people"].([]string)

	expectedArrayValues := []string{"bob", "tony", "henry"}

	if name != "John" {
		t.Errorf("Expected Value: %s, Gotten: %s", "John", name)
		return
	}

	if age != 20 {
		t.Errorf("Expected Value: %d, Gotten: %d", 20, age)
		return
	}

	if place != "Thrissur" {
		t.Errorf("Expected Value: %s, Gotten: %s", "Thrissur", place)
		return
	}

	i := 0
	t.Run("Hash :Array Data Type", func(t *testing.T) {

		for _, val := range expectedArrayValues {
			if val != people[i] {
				t.Errorf("Expected Value: %s, Gotten: %s", val, people[i])
				return
			}
			i++
		}
	})

}

func TestHGet(t *testing.T) {
	c := NewCache()

	key := "users:Jhon:data"
	field := "age"
	value := 20
	c.Hset(key, field, value)

	data, err := c.HGet(key, field)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	if data.(int) != value {
		t.Errorf("Expected Value: %d, Gotten: %d", value, data)
		return
	}

	// key does not exists
	key = "fruits"
	_, err = c.HGet(key, "sweet")
	if err == nil {
		t.Errorf("Expected Err: %s, Gotten: ERR NIL", ErrKeyNotFound)
		return
	}

	if err.Error() != ErrKeyNotFound {
		t.Errorf("Expected Err: %s, Gotten: %s", ErrKeyNotFound, err.Error())
		return
	}

	// field does not exist
	key = "fruits"
	field = "bitter"
	c.Hset(key, field, "lemons")
	_, err = c.HGet(key, "sweet")
	if err == nil {
		t.Errorf("Expected Err: %s, Gotten: ERR NIL", ErrFieldNotFound)
		return
	}

	if err.Error() != ErrFieldNotFound {
		t.Errorf("Expected Err: %s, Gotten: %s", ErrFieldNotFound, err.Error())
		return
	}

	// Not an Hash Value
	key = "fruits"
	v := "orange"
	c.Set(key, v, 0)
	_, err = c.HGet(key, "sweet")
	if err == nil {
		t.Errorf("Expected Err: %s, Gotten: ERR NIL", ErrNotHashvalue)
		return
	}

	if err.Error() != ErrNotHashvalue {
		t.Errorf("Expected Err: %s, Gotten: %s", ErrNotHashvalue, err.Error())
		return
	}
}

func TestHgetAll(t *testing.T) {
	c := NewCache()

	//key does not exists
	key := "users:bob"
	_, err := c.HGetAll(key)
	if err == nil {
		t.Errorf("Expected Err: %s, Gotten: ERR NIL", ErrKeyNotFound)
		return
	}

	if err.Error() != ErrKeyNotFound {
		t.Errorf("Expected Err: %s, Gotten: %s", ErrKeyNotFound, err.Error())
		return
	}

	//value not hash
	key = "fruits"
	v := "orange"
	c.Set(key, v, 0)
	_, err = c.HGetAll(key)
	if err == nil {
		t.Errorf("Expected Err: %s, Gotten: ERR NIL", ErrNotHashvalue)
		return
	}

	if err.Error() != ErrNotHashvalue {
		t.Errorf("Expected Err: %s, Gotten: %s", ErrNotHashvalue, err.Error())
		return
	}

}

func TestExist(t *testing.T) {
	c := NewCache()

	key := "users:bob"
	c.Set(key, "mexican alien", 4000)

	ok := c.Exists(key)
	if !ok {
		t.Errorf("Expected Value: %v, Gotten: %v", true, ok)
		return
	}

	key = "users:john"
	t.Run("Key does not exist", func(t *testing.T) {
		ok := c.Exists(key)
		if ok {
			t.Errorf("Expected Value: %v, Gotten: %v", false, ok)
			return
		}
	})
}

func TestGetAllData(t *testing.T) {
	c := NewCache()

	keys := []string{"name", "age", "idk"}
	c.Set(keys[0], "bob", 0)
	c.Set(keys[1], 22, 0)
	c.Set(keys[2], "idk", 0)
	data, _ := c.AllData()

	for i := 0; i < len(keys); i++ {
		if _, ok := data[keys[i]]; !ok {
			t.Errorf("Key:%s does't Exist", keys[i])
		}
	}

}

func TestDeleteExpiredKeys(t *testing.T) {
	c := NewCache().(*Cache)

	c.Set("key1", "t1", 1000)
	c.Set("key2", "t1", 2000)
	c.Set("key3", "t1", 10000)

	time.Sleep(time.Second * 3)
	testDeleteExpiredKeys(c)

	if c.Exists("key1") || c.Exists("key2") {
		t.Errorf("key1 & key2 has not been removed")
		return
	}
	if !c.Exists("key3") {
		t.Errorf("key3 should exists")
	}
	c.Del("key3")
	testDeleteExpiredKeys(c)

	c.Set("key4", "t4", 0)
	testDeleteExpiredKeys(c)

}

// func TestCache(t *testing.T) {
// 	c := goswift.NewCache()

// 	fmt.Println(time.Now().Unix())
// 	c.Set("leo", 23000, "kinglol")
// 	c.Set("name", 9000, "leoantony")
// 	c.Set("jsondata", 6000, "THIS IS A TEST ")
// 	exp := 3000
// 	var wg sync.WaitGroup
// 	for i := 0; i < 1000; i++ {
// 		wg.Add(3)
// 		go AddNode(c, exp, &wg)
// 		go AddNode(c, exp, &wg)
// 		go AddNode(c, exp, &wg)
// 	}
// 	c.Set("idk", 2000, "THIS IS A TEST ")
// 	c.Set("boiz", 7000, "THIS IS A TEST ")
// 	c.Set("no name", 10000, "THIS IS A TEST ")

// 	wg.Wait()

// 	PrintALL(c)

// 	c.Del("no name")
// 	interval := 1 * time.Second
// 	// fmt.Println(interval)
// 	ticker := time.NewTicker(interval)
// 	defer ticker.Stop()
// 	for {
// 		select {
// 		case <-ticker.C:
// 			PrintALL(c)
// 			PrintALLH(c)
// 		}
// 	}
// }

// func PrintALLH(c CacheFunction) {
// 	d := c.AllDataHeap()
// 	// fmt.Println(d)
// 	counTer := 0
// 	for s, v := range d {
// 		fmt.Println(s, v)
// 		counTer += 1
// 	}
// 	fmt.Println("total Heap Data: ", counTer)
// 	fmt.Println("----------------------")
// }

func AddNode(c CacheFunction, exp int, wg *sync.WaitGroup) {
	defer wg.Done()
	key := uuid.New()
	v := uuid.New()
	c.Set(key.String(), v.String(), exp)
}

func Print(h *expiry.Heap) {
	for _, b := range h.Data {
		fmt.Println(b)
	}
}
