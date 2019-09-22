package main

import (
	"fmt"
	. "github.com/fishedee/app/sqlf"
)

type SqlfDb struct {
	db SqlDB
}

func (this *SqlfDb) Init() {
	var err error
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s?parseTime=true", USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)
	this.db, err = NewSqlDB(SqlDBConfig{
		Driver:     "mysql",
		SourceName: dsn,
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
