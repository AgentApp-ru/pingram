package main

import (
	"fmt"
	"sync"
	"time"
)

func right() {
    roomA := new(sync.Map)
	roomA.Store("vasya", 1)
	roomA.Store("petya", 2)
	roomB := new(sync.Map)
	roomB.Store("sasha", 3)
	roomB.Store("masha", 4)

	people := new(sync.Map)
	people.Store("roomA", roomA)
	people.Store("roomB", roomB)

	people.Range(func(key, value interface{}) bool {
		go checkBig(value.(*sync.Map))
		// fmt.Printf("%v\n", key)
		// fmt.Printf("%v\n", value)
		return true
	})
}

func wrong() {
    people := map[string]map[string]int{
		"roomA": {
            "vasya": 1,
            "petya": 2,
        },
        "roomB": {
            "sasha": 3,
            "masha": 4,
        },
	}

    for _, v := range people {
        go checkBigWrong(v)
    }
}

func checkBig(b *sync.Map) {
	b.Range(func(key, value interface{}) bool {
		fmt.Printf("%v %v\n", key, value)
		return true
	})
	fmt.Println("======")
	b.Range(func(key, value interface{}) bool {
		valueInt := value.(int)
		valueInt = valueInt * 2
		b.Store(key, valueInt)
		return true
	})
	b.Range(func(key, value interface{}) bool {
		fmt.Printf("%v %v\n", key, value)
		return true
	})
}

// func checkSmall()


func checkBigWrong(b map[string]int) {
	for key, value := range b {
        fmt.Printf("%v %v\n", key, value)
    }
	fmt.Println("======")
	for key, value := range b {
        b[key] = value * 2
    }
	for key, value := range b {
        fmt.Printf("%v %v\n", key, value)
    }
}

func checkSmall()


func main() {
	// wg = sync.WaitGroup()
    // right()
	wrong()

	

	time.Sleep(2 * time.Second)

	// fmt.Println(people.Range())
}
