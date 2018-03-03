# go-redis-orm.v2

本库通过定义json文件，使用工具生成redis orm类文件

可以处理1对1类型的数据、以及1对N类型的数据


## 1对1类型数据例子

```go
func test11() {
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
```


## 1对N类型的数据例子

```go
func test1n() {
	dbName := "db2"

	go_redis_orm.SetNewRedisHandler(go_redis_orm.NewDefaultRedisClient)
	go_redis_orm.CreateDB(dbName, []string{"192.168.1.12:16379"}, "", 0)

	data1 := NewTestStruct2(dbName, 8)
	item1 := data1.NewItem(1)
	item1.SetMyf2(99.9)
	item2 := data1.NewItem(2)
	item2.SetMys1("hello")
	item2.SetMys2([]byte("world"))
	err := data1.Save()
	if err != nil {
		panic(err)
	}

	data2 := NewTestStruct2(dbName, 8)
	err = data2.Load()
	if err != nil {
		panic(err)
	}
	fmt.Printf("2: %+v\n", data2.GetItem(1))
	fmt.Printf("2: %+v\n", data2.GetItem(2))
	data2.DeleteItem(1)
	data2.Save()

	data3 := NewTestStruct2(dbName, 8)
	data3.Load()
	for _, v := range data3.GetItems() {
		fmt.Printf("3: %+v\n", v)
	}
	data3.Delete()
	data3.Save()

	data4 := NewTestStruct2(dbName, 8)
	data4.Load()
	fmt.Printf("4: item count = %d\n", len(data4.GetItems()))

	fmt.Println("OK")
}
```


## 使用方法

  1. 定义json文件

    格式参考：example/redis_def/*.json

  1. 使用release/redis2go.exe生成go文件


## 编译

执行下列语句：

```dos
git.exe submodule update --init -- "tools/build"
build.bat
```
