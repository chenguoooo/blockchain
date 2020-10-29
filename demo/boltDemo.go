package main

import (
	"bolt"
	"fmt"
	"log"
)

func main() {
	db, err := bolt.Open("test.db", 0600, nil)
	//向数据库中写入数据
	//从数据库中读取数据

	if err != nil {
		log.Panic(err)
	}

	defer db.Close()
	db.Update(func(tx *bolt.Tx) error {
		//所有操作在这里
		b1 := tx.Bucket([]byte("bucketname1"))

		if b1 == nil {
			//如果b1为空，说明该桶不存在，需要创建

			b1, err = tx.CreateBucket([]byte("bucketname1"))
			if err != nil {
				log.Panic(err)
			}
		}
		//bucket已经创建完成，准备写入数据
		//写数据使用Put，读数据使用Get
		b1.Put([]byte("name1"), []byte("Lily"))
		if err != nil {
			fmt.Printf("写入数据失败name1：Lily\n")
		}
		b1.Put([]byte("name2"), []byte("Jim"))
		if err != nil {
			fmt.Printf("写入数据失败name2：Jim\n")
		}

		//读取数据
		name1 := b1.Get([]byte("name1"))
		name2 := b1.Get([]byte("name2"))
		name3 := b1.Get([]byte("name3"))

		fmt.Printf("name1:%s\n", name1)
		fmt.Printf("name2:%s\n", name2)
		fmt.Printf("name3:%s\n", name3)

		return nil

	})

}
