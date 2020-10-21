package sqlDB

import (
	"fmt"
	"github.com/tsingeye/FreeEhome/tools/logs"
)

//初始化创建设备信息表-DeviceList
func InitDeviceList() {
	if !gormDB.HasTable(&DeviceList{}) {
		//不存在，创建表格
		CreateTable(&DeviceList{})
	} else {
		//表格存在，则更新Status=OFF
		deviceList := map[string]interface{}{
			"Status": "OFF",
		}
		Updates(&DeviceList{}, deviceList)
	}
}

//初始化创建通道列表-ChannelList
func InitChannelList() {
	if !gormDB.HasTable(&ChannelList{}) {
		//不存在，创建表格
		CreateTable(&ChannelList{})
	} else {
		//表格存在，则更新Status=OFF
		channelList := map[string]interface{}{
			"Status": "OFF",
		}
		Updates(&ChannelList{}, channelList)
	}
}

//创建表格，如果表格不存在
func CreateTable(tab interface{}) {
	tableName := GetTableName(tab)
	switch dbType {
	case "mysql":
		if err := gormDB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(tab).Error; err != nil {
			logs.PanicLogger.Panicln(fmt.Sprintf("create %s table failed: %s", tableName, err))
		}
	case "sqlite3":
		if err := gormDB.CreateTable(tab).Error; err != nil {
			logs.PanicLogger.Panicln(fmt.Sprintf("create %s table failed: %s", tableName, err))
		}
	}
	logs.BeeLogger.Info("create %s table success", tableName)
}

//解析获得用户要操作的表名
func GetTableName(value interface{}) (tableName string) {
	scope := gormDB.NewScope(value)
	if name, ok := value.(string); ok {
		tableName = name
	} else {
		tableName = scope.TableName()
	}

	return
}

//查询某个表符合条件的总数，query参数为空字符串则查询整个表的总数，否则为按条件查询
func Count(tab interface{}, query string, where ...interface{}) (count int64) {
	tableName := GetTableName(tab)

	if query == "" {
		if err := gormDB.Table(tableName).Count(&count); err != nil {
			logs.BeeLogger.Error("query total count from %s's table error: %s", tableName, err)
		}
	} else {
		if err := gormDB.Model(tab).Where(query, where...).Count(&count); err != nil {
			logs.BeeLogger.Error("condition query count from %s's table error: %s", tableName, err)
		}
	}

	return
}

//添加记录到表格中，存在则更新，不存在则插入
func Save(tab interface{}) bool {
	tableName := GetTableName(tab)
	if err := gormDB.Save(tab).Error; err != nil {
		logs.BeeLogger.Error("%s's table save record error: %s", tableName, err)
		return false
	}
	return true
}

//分页查询
func Limit(out interface{}, tableName string, limit, offset int, where ...interface{}) {
	if err := gormDB.Limit(limit).Offset(offset).Find(out, where...).Error; err != nil {
		logs.BeeLogger.Error("paging query error in the %s table: %s", tableName, err)
	}
}

//执行更新操作，更新更改字段
func Updates(tbl interface{}, data interface{}) {
	tableName := GetTableName(tbl)

	if err := gormDB.Model(tbl).Updates(data).Error; err != nil {
		logs.BeeLogger.Error("%s table run db.Updates() failed: %s", tableName, err)
		return
	}
}

//查询符合条件的第一条记录，有数据返回0，无数据返回1，其他查询出错返回-1
func First(tab interface{}, whereMap interface{}) int {
	tableName := GetTableName(tab)
	if err := gormDB.First(tab, whereMap).Error; err != nil {
		logs.BeeLogger.Error("%s table run db.First() failed: %s", tableName, err)
		//gormDB.Where(data).First(tbl).RecordNotFound()
		if err.Error() == "record not found" {
			//查询失败，无符合条件的数据
			return 1
		}
		return -1
	}

	return 0
}

//查询所有记录，参数一为结构体数组指针，参数二为表名，后面参数是查询套件
func Find(out interface{}, tableName string, where ...interface{}) {
	if err := gormDB.Find(out, where...).Error; err != nil {
		logs.BeeLogger.Error("%s table run db.Find() failed: %s", tableName, err)
	}
}

//创建记录
func Create(tab interface{}) bool {
	tableName := GetTableName(tab)
	if err := gormDB.Create(tab).Error; err != nil {
		logs.BeeLogger.Error("error inserting records into the %s table: %s", tableName, err)
	} else {
		return true
	}
	return false
}

//更新数据库中表格的状态
func UpdateTable(sql string) {
	if err := gormDB.Exec(sql).Error; err != nil {
		logs.BeeLogger.Error("update table error: %s", err)
	}
}
