package dao

import (
	"github.com/go-xorm/xorm"
	"log"
	models "lottery/model"
)

type BlackipDao struct {
	engine *xorm.Engine
}

func newBlackipDao(engine *xorm.Engine) *BlackipDao {
	return &BlackipDao{
		engine: engine,
	}
}

func (d *BlackipDao) Get(id int) *models.LtGift {

	data := &models.LtGift{Id: id}
	ok, err := d.engine.Get(data)
	if ok && err == nil {
		return data
	}

	data.Id = 0
	return data
}

func (d *BlackipDao) GetAll() []models.LtGift {

	dataList := make([]models.LtGift, 0)
	err := d.engine.Asc("sys_status").Asc("displayorder").Find(&dataList)
	if err != nil {
		log.Println("gift_dao GetAll() error: ", err)
		return dataList
	}

	return dataList
}

func (d *BlackipDao) CountAll() int64 {

	count, err := d.engine.Count(&models.LtGift{})
	if err != nil {
		log.Println("gift_dao CountAll() error: ", err)
		return 0
	}

	return count
}

func (d *BlackipDao) Delete(id int) error {

	data := &models.LtGift{Id: id, SysStatus: 1}
	_, err := d.engine.Id(data.Id).Update(data)
	return err
}

func (d *BlackipDao) Update(data *models.LtGift, columns []string) error {

	_, err := d.engine.Id(data.Id).MustCols(columns...).Update(data)
	return err
}

func (d *BlackipDao) Create(data *models.LtGift) error {

	_, err := d.engine.Insert(data)
	return err
}

func (d *BlackipDao) GetByIp(ip string) *models.LtBlackip {

	dataList := make([]models.LtBlackip, 0)
	err := d.engine.Where("ip=?", ip).Desc("id").Limit(1).Find(&dataList)

	if err != nil || len(dataList) < 1 {
		return nil
	}

	return &dataList[0]
}
