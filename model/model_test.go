package model

import (
	"fmt"
	"log"
	"testing"
)

type OrderInfo struct {
	Id       int
	Oid      string
	Username string
}

func (o *OrderInfo) TableName() string {
	return "order_info"
}

func ExampleReader_GetOne() {
	Init()
	RegisterModel(new(OrderInfo))

	var order *OrderInfo

	err := Read(new(OrderInfo)).Filter("id", 87).GetOne(&order)

	if err != nil {
		log.Panic(err.Error())
	}

	fmt.Println(order.Id)
}

func ExampleAdd() {
	Init()
	RegisterModel(new(OrderInfo))

	var order OrderInfo
	order.Oid = "3423328"
	order.Username = "JYGO"

	lastInsertId := Add(order)

	fmt.Println(lastInsertId)
}

func ExampleReader_GetAll() {
	Init()
	RegisterModel(new(OrderInfo))

	var orders []*OrderInfo

	num, err := Read(new(OrderInfo)).Filter("id", 87).GetAll(&orders)

	if err != nil {
		log.Panic(err.Error())
	}

	fmt.Println(num)
}

func TestInit(t *testing.T) {
	Init()
	RegisterModel(new(OrderInfo))
}

func TestReaderGetOne(t *testing.T) {
	t.Parallel()
	var order *OrderInfo

	err := Read(new(OrderInfo)).Filter("id", 87).GetOne(&order)

	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println(order.Id)
}

func TestAdd(t *testing.T) {
	t.Parallel()

	var order OrderInfo

	order.Oid = "3423328"
	order.Username = "jiyic1"

	lastInsertId := Add(order)

	fmt.Println(lastInsertId)
}
