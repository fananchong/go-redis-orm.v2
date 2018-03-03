package main

import (
	"fmt"

	go_redis_orm "github.com/fananchong/go-redis-orm.v2"
)

func main() {

	dbName := "db1"

	go_redis_orm.SetNewRedisHandler(go_redis_orm.NewDefaultRedisClient)
	go_redis_orm.CreateDB(dbName, []string{"192.168.1.12:16379"}, "", 0)

	// key值为1的 TestStruct1 数据
	data1 := NewTestStruct1(dbName, 1)
	data1.SetMyb(true)
	data1.SetMyf1(1.5)
	data1.SetMyi5(100)
	data1.SetMys1("hello")
	data1.SetMys2([]byte("world"))
	err := data1.Save()
	if err != nil {
		panic(err)
	}

	data2 := NewTestStruct1(dbName, 1)
	err = data2.Load()

	if err == nil {
		if data2.GetMyb() != true ||
			data2.GetMyf1() != 1.5 ||
			data2.GetMyi5() != 100 ||
			data2.GetMys1() != "hello" ||
			string(data2.GetMys2()) != "world" {
			panic("#1")
		}
	} else {
		panic(err)
	}

	err = data2.Delete()
	if err != nil {
		panic(err)
	}

	var hasKey int
	hasKey, err = data2.HasKey()
	if hasKey != 0 {
		panic("#2")
	}

	fmt.Println("OK")
}
