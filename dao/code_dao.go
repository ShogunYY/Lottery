package dao

import (
	"github.com/go-xorm/xorm"
	"log"
	models "lottery/model"
)

type CodeDao struct {
	engine *xorm.Engine
}

func newCodeDao(engine *xorm.Engine) *CodeDao {
	return &CodeDao{
		engine: engine,
	}
}

func (d *CodeDao) Get(id int) *models.LtGift {

	data := &models.LtGift{Id: id}
	ok, err := d.engine.Get(data)
	if ok && err == nil {
		return data
	}

	data.Id = 0
	return data
}

func (d *CodeDao) GetAll() []models.LtGift {

	dataList := make([]models.LtGift, 0)
	err := d.engine.Desc("id").Find(&dataList)
	if err != nil {
		log.Println("gift_dao GetAll() error: ", err)
		return dataList
	}

	return dataList
}

func (d *CodeDao) CountAll() int64 {

	count, err := d.engine.Count(&models.LtGift{})
	if err != nil {
		log.Println("gift_dao CountAll() error: ", err)
		return 0
	}

	return count
}

func (d *CodeDao) Delete(id int) error {

	data := &models.LtGift{Id: id, SysStatus: 1}
	_, err := d.engine.Id(data.Id).Update(data)
	return err
}

func (d *CodeDao) Update(data *models.LtGift, columns []string) error {

	_, err := d.engine.Id(data.Id).MustCols(columns...).Update(data)
	return err
}

func (d *CodeDao) Create(data *models.LtGift) error {

	_, err := d.engine.Insert(data)
	return err
}
