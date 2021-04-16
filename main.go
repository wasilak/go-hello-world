package main

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// Default middleware config
var store *session.Store

type Header struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Response struct {
	Counter  int `json:"counter"`
	Host  string `json:"host"`
	Hostnames  map[string]int `json:"hostnames"`
	Headers []Header `json:"headers"`
}

func keyGen() string {
	return "32543gv45b45"
}

func main() {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	appLogger := logger.New(logger.Config{
		Format: `{"pid": "${pid}", "requestid": "${locals:requestid}", "status": "${status}", "method": "${method}", "path": "${path}"}​` + "\n",
	})

	app.Use(appLogger)

	var ConfigDefault = session.Config{
	    CookieName:   "go_hello_world_session",
	    KeyGenerator: keyGen,
	}

	store := session.New(ConfigDefault)

	app.Get("/", func(c *fiber.Ctx) error {

		// get session from storage
    sess, err := store.Get(c)
    if err != nil {
        panic(err)
    }

    // save session
    defer sess.Save()

		var response Response

		response.Host = string(c.Request().Host())

		if val, ok := sess.Get("counter").(int); ok {
			response.Counter = val
		}
		response.Counter = response.Counter + 1
		sess.Set("counter", response.Counter)
		// fmt.Printf("%+v\n", response.Counter)

		if nil == sess.Get("hostnames") {
			response.Hostnames = make(map[string]int)
			response.Hostnames[response.Host]++
			// sess.Set("hostnames", make(map[string]int))

			bytes, err := json.Marshal(response.Hostnames)
			if err != nil {
				panic(err)
			}
			sess.Set("hostnames", string(bytes))
		}
		// response.Hostnames = sess.Get("hostnames").(map[string]int)
		// response.Hostnames[response.Host]++

		json.Unmarshal([]byte(sess.Get("hostnames").(string)), &response.Hostnames)

		response.Hostnames[response.Host]++

		bytes, err := json.Marshal(response.Hostnames)
		if err != nil {
			panic(err)
		}
		sess.Set("hostnames", string(bytes))

		fmt.Printf("simple get %+v\n", response.Hostnames)
		// fmt.Printf("get with conversion %+v\n", sess.Get("hostnames").(map[string]int))

		// sess.Set("hostnames", response.Hostnames)


		// if _, ok := sess.Get("hostnames").(map[string]int); !ok {
		// if sess.Get("hostnames") == nil {
		// 	m := make(map[string]int)
		// 	// m[response.Host] = 1
		// 	// // val[response.Host] = 1
		// 	m[response.Host] = 1
		// 	response.Hostnames = m

		// 	fmt.Println("doesn't exists")
		// } else {

		// 	response.Hostnames[response.Host] = response.Hostnames[response.Host] + 1
		// 	// sess.Set("hostnames", response.Hostnames)
		// 	fmt.Println("does exist")
		// }

		// sess.Set("hostnames", response.Hostnames)

		// if val, ok := response.Hostnames[response.Host]; ok {
		// 	response.Hostnames[response.Host] = val + 1
		// }
		// response.Counter = response.Counter + 1
		// sess.Set("hostnames", response.Hostnames)
		// fmt.Printf("%+v\n", response.Hostnames)

    // if counter == nil {
    // 	sess.Set("counter", 1)
    // } else {
    // 	counterNewValue := counter.(int) + 1
    // 	response.Counter = counterNewValue
    // 	sess.Set("counter", counterNewValue)
    // }

    // hosts := sess.Get("hosts")

    // if hosts == nil {
    // 	var hostsSlice map[string]int
    // 	sess.Set("hosts", hostsSlice)
    // } else {
    // 	hostsSliceFromSession := hosts.(map[string]int)
    // 	if val, ok := hostsSliceFromSession[response.Host]; ok {
    // 		hostsSliceFromSession[response.Host] = val + 1
    // 	}
    // 	response.Hostnames = hostsSliceFromSession
    // 	sess.Set("hosts", hostsSliceFromSession)
    // }

		c.Request().Header.VisitAll(func(key, value []byte) {
			tmp := Header{
				Name:  string(key),
				Value: string(value),
			}
			response.Headers = append(response.Headers, tmp)
		})

		// bytes, err := json.Marshal(response)
		// if err != nil {
		// 	panic(err)
		// }

		// fmt.Println(string(bytes))

		return c.JSON(response)
		// return c.Request().Header
	})

	app.Listen(":5000")
}
