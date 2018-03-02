package main

import (
	"fmt"

	go_redis_orm "github.com/fananchong/go-redis-orm.v2"
)

func main() {

	dbName := "db1"

	go_redis_orm.SetNewRedisHandler(go_redis_orm.NewDefaultRedisClient)
	go_redis_orm.CreateDB(dbName, []string{"192.168.1.12:16379"}, "", 0)

	// key值为1的 TestStruct2 数据
	data1 := NewRD_TestStruct1(1)
	data1.SetMyb(true)
	data1.SetMyf1(1.5)
	data1.SetMyi5(100)
	data1.SetMys1("hello")
	err := data1.Save(dbName)
	if err != nil {
		panic(err)
	}

	data2 := NewRD_TestStruct1(1)
	err = data2.Load(dbName)

	if err == nil {
		if data2.GetMyb() != true ||
			data2.GetMyf1() != 1.5 ||
			data2.GetMyi5() != 100 ||
			data2.GetMys1() != "hello" {
			panic("#2")
		}
	} else {
		panic(err)
	}

	err = data2.Delete(dbName)
	if err != nil {
		panic(err)
	}

	err = data2.Load(dbName)
	if data2.IsLoad() {
		panic("#6")
	}

	fmt.Println("OK")
}
