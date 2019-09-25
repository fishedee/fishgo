package main

import (
	"fmt"
	"github.com/gocraft/dbr"
)

type DbrDb struct {
	db *dbr.Connection
}

func (this *DbrDb) Init() {
	var err error
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s?parseTime=true", USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)
	this.db, err = dbr.Open("mysql", dsn, nil)
	if err != nil {
		panic(err)
	}
}

func (this *DbrDb) GetAllMaterial() []Material {
	materials := []Material{}
	this.db.NewSession(nil).Select("materialId,productId").From("t_material").Load(&materials)

	return materials
}

func (this *DbrDb) GetProduct(productIds []int) []Product {
	products := []Product{}

	this.db.NewSession(nil).Select("*").From("t_material").Where("productId in ?", productIds).Load(&products)
	return products
}
