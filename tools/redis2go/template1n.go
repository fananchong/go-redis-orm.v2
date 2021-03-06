package main

const template1n string = `/// -------------------------------------------------------------------------------
/// THIS FILE IS ORIGINALLY GENERATED BY redis2go.exe.
/// PLEASE DO NOT MODIFY THIS FILE.
/// -------------------------------------------------------------------------------
package {{packagename}}

import (
	"errors"
	{{fmt}}
	{{strconv}}

	go_redis_orm "github.com/fananchong/go-redis-orm.v2"
	"github.com/gomodule/redigo/redis"
)

// {{classname}} : 代表 1 个 redis 对象
type {{classname}} struct {
	Key         {{key_type}}
	values      map[{{sub_key_type}}]*{{classname}}Item

	dirtyDataIn{{classname}} map[{{sub_key_type}}]int
	isLoadIn{{classname}}    bool
	dbKeyIn{{classname}}     string
	dbNameIn{{classname}}    string
	expireIn{{classname}}    uint
}

// New{{classname}} : New{{classname}} 的构造函数
func New{{classname}}(dbName string, key {{key_type}}) *{{classname}} {
	return &{{classname}}{
		Key:         key,
		values:      make(map[{{sub_key_type}}]*{{classname}}Item),
		dbNameIn{{classname}}:    dbName,
		dbKeyIn{{classname}}:     {{func_dbkey}},
		dirtyDataIn{{classname}}: make(map[{{sub_key_type}}]int),
	}
}

// HasKey : 是否存在 KEY
//          返回值，若访问数据库失败返回-1；若 key 存在返回 1 ，否则返回 0 。
func (obj{{classname}} *{{classname}}) HasKey() (int, error) {
	db := go_redis_orm.GetDB(obj{{classname}}.dbNameIn{{classname}})
	val, err := redis.Int(db.Do("EXISTS", obj{{classname}}.dbKeyIn{{classname}}))
	if err != nil {
		return -1, err
	}
	return val, nil
}

// Load : 从 redis 加载数据
func (obj{{classname}} *{{classname}}) Load() error {
	if obj{{classname}}.isLoadIn{{classname}} == true {
		return errors.New("alreay load")
	}
	db := go_redis_orm.GetDB(obj{{classname}}.dbNameIn{{classname}})
	val, err := redis.Values(db.Do("HGETALL", obj{{classname}}.dbKeyIn{{classname}}))
	if err != nil {
		return err
	}
	if len(val) == 0 {
		return go_redis_orm.ERR_ISNOT_EXIST_KEY
	}
	for i := 0; i < len(val); i += 2 {
		temp := string(val[i].([]byte))
		{{conv_subkey}}
		if err != nil {
			return err
		}
		item := New{{classname}}Item(subKey, obj{{classname}})
		err = item.Unmarshal(val[i+1].([]byte))
		if err != nil {
			return err
		}
		obj{{classname}}.values[subKey] = item
	}
	obj{{classname}}.isLoadIn{{classname}} = true
	return nil
}

// Save : 保存数据到 redis
func (obj{{classname}} *{{classname}}) Save() error {
	if len(obj{{classname}}.dirtyDataIn{{classname}}) == 0 {
		return nil
	}
	tempData := make(map[{{sub_key_type}}][]byte)
	for k := range obj{{classname}}.dirtyDataIn{{classname}} {
		if item, ok := obj{{classname}}.values[k]; ok {
			var err error
			tempData[k], err = item.Marshal()
			if err != nil {
				return err
			}
		}
	}
	db := go_redis_orm.GetDB(obj{{classname}}.dbNameIn{{classname}})
	if _, err := db.Do("HMSET", redis.Args{}.Add(obj{{classname}}.dbKeyIn{{classname}}).AddFlat(tempData)...); err != nil {
		return err
	}
	if obj{{classname}}.expireIn{{classname}} != 0 {
		if _, err := db.Do("EXPIRE", obj{{classname}}.dbKeyIn{{classname}}, obj{{classname}}.expireIn{{classname}}); err != nil {
			return err
		}
	}
	obj{{classname}}.dirtyDataIn{{classname}} = make(map[{{sub_key_type}}]int)
	return nil
}

// Delete : 从 redis 删除数据
func (obj{{classname}} *{{classname}}) Delete() error {
	db := go_redis_orm.GetDB(obj{{classname}}.dbNameIn{{classname}})
	_, err := db.Do("DEL", obj{{classname}}.dbKeyIn{{classname}})
	if err == nil {
		obj{{classname}}.isLoadIn{{classname}} = false
		obj{{classname}}.dirtyDataIn{{classname}} = make(map[{{sub_key_type}}]int)
	}
	return err
}

// NewItem : 新建 1 个子对象
func (obj{{classname}} *{{classname}}) NewItem(subKey {{sub_key_type}}) *{{classname}}Item {
	item := New{{classname}}Item(subKey, obj{{classname}})
	obj{{classname}}.values[subKey] = item
	obj{{classname}}.dirtyDataIn{{classname}}[subKey] = 1
	return item
}

// DeleteItem : 删除 1 个子对象
func (obj{{classname}} *{{classname}}) DeleteItem(subKey {{sub_key_type}}) error {
	if _, ok := obj{{classname}}.values[subKey]; ok {
		db := go_redis_orm.GetDB(obj{{classname}}.dbNameIn{{classname}})
		_, err := db.Do("HDEL", obj{{classname}}.dbKeyIn{{classname}}, subKey)
		if err != nil {
			return err
		}
		delete(obj{{classname}}.values, subKey)
		if _, ok := obj{{classname}}.dirtyDataIn{{classname}}[subKey]; ok {
			delete(obj{{classname}}.dirtyDataIn{{classname}}, subKey)
		}
	}
	return nil
}

// GetItem : 获取某个子对象
func (obj{{classname}} *{{classname}}) GetItem(subKey {{sub_key_type}}) *{{classname}}Item {
	if item, ok := obj{{classname}}.values[subKey]; ok {
		return item
	}
	return nil
}

// GetItems : 获取所有子对象
func (obj{{classname}} *{{classname}}) GetItems() []*{{classname}}Item {
	var ret []*{{classname}}Item
	for _, v := range obj{{classname}}.values {
		ret = append(ret, v)
	}
	return ret
}

// DirtyData : 获取该对象目前已脏的数据
func (obj{{classname}} *{{classname}}) DirtyData() (map[{{sub_key_type}}][]byte, error) {
	data := make(map[{{sub_key_type}}][]byte)
	for k := range obj{{classname}}.dirtyDataIn{{classname}} {
		if item, ok := obj{{classname}}.values[k]; ok {
			var err error
			data[k], err = item.Marshal()
			if err != nil {
				return nil, err
			}
		}
	}
	obj{{classname}}.dirtyDataIn{{classname}} = make(map[{{sub_key_type}}]int)
	return data, nil
}

// Save2 : 保存数据到 redis 的第 2 种方法
func (obj{{classname}} *{{classname}}) Save2(dirtyData map[{{sub_key_type}}][]byte) error {
	if len(dirtyData) == 0 {
		return nil
	}
	db := go_redis_orm.GetDB(obj{{classname}}.dbNameIn{{classname}})
	if _, err := db.Do("HMSET", redis.Args{}.Add(obj{{classname}}.dbKeyIn{{classname}}).AddFlat(dirtyData)...); err != nil {
		return err
	}
	if obj{{classname}}.expireIn{{classname}} != 0 {
		if _, err := db.Do("EXPIRE", obj{{classname}}.dbKeyIn{{classname}}, obj{{classname}}.expireIn{{classname}}); err != nil {
			return err
		}
	}
	return nil
}

// IsLoad : 是否已经从 redis 导入数据
func (obj{{classname}} *{{classname}}) IsLoad() bool {
	return obj{{classname}}.isLoadIn{{classname}}
}

// Expire : 向 redis 设置该对象过期时间
func (obj{{classname}} *{{classname}}) Expire(v uint) {
	obj{{classname}}.expireIn{{classname}} = v
}
`

const convSubKeyFuncStringInt = `tempUint64, err := strconv.ParseUint(temp, 10, 64)
subKey := {{sub_key_type}}(tempUint64)`

const convSubKeyFuncStringStr = `subKey := temp`
