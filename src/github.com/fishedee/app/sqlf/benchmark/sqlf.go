package main

import (
	"fmt"
	. "github.com/fishedee/app/log"
	. "github.com/fishedee/app/sqlf"
)

type SqlfDb struct {
	db SqlfDB
}

func (this *SqlfDb) Init() {
	var err error
	log, err := NewLog(LogConfig{
		Driver: "console",
	})
	if err != nil {
		panic(err)
	}
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s?parseTime=true", USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)
	this.db, err = NewSqlfDB(log, SqlfDBConfig{
		Driver:     "mysql",
		SourceName: dsn,
		Debug:      false,
	})
	if err != nil {
		panic(err)
	}
}

func (this *SqlfDb) GetAllMaterial() []Material {
	materials := []Material{}
	this.db.MustQuery(&materials, "select materialId,productId from t_material")

	return materials
}

func (this *SqlfDb) GetProduct(productIds []int) []Product {
	products := []Product{}

	this.db.MustQuery(&products, "select ?.column from t_product where productId in (?)", products, productIds)
	return products
}
