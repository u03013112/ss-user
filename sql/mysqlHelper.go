package sql

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	mysqlPasswd = "!@#sspaas@U0"
	mysqlDBName = "ss_main"
)

// BaseModel : 为了加json tag
type BaseModel struct {
	ID        uint       `gorm:"primary_key" json:"id,omitempty"`
	CreatedAt time.Time  `sql:"index" json:"createdTime,omitempty"`
	UpdatedAt time.Time  `sql:"index" json:"updatedTime,omitempty"`
	DeletedAt *time.Time `sql:"index" json:"deletedTime,omitempty"`
}

// InitDB ：初始化数据库，如果失败直接panic，成功返回db
func InitDB() {
	mysqlAddr := os.Getenv("MYSQL")

	if mysqlAddr == "" {
		mysqlAddr = "mysql"
	}
	db, err := gorm.Open("mysql", fmt.Sprintf("root:%s@tcp(%s:3306)/%s?charset=utf8&loc=%s&parseTime=True", mysqlPasswd, mysqlAddr, mysqlDBName, "Asia%2FShanghai"))
	if err != nil {
		panic(err)
	}
	// defer db.Close()
	instance = db
}

var instance *gorm.DB

// GetInstance : 获得单例
func GetInstance() *gorm.DB {
	if instance == nil {
		// 这里应该用不上，直接在main开始位置初始化，这里应该走不到
		InitDB()
	}
	return instance.Debug()
	// return instance
}

// ListFilter : 列表通用过滤器
type ListFilter struct {
	PageNumber int64  `json:"pageNumber,omitempty"`
	PageSize   int64  `json:"pageSize,omitempty"`
	OrderBy    string `json:"orderBy,omitempty"`
	Order      string `json:"order,omitempty"`
	Search     string `json:"search,omitempty"`
	StartTime  int64  `json:"startTime,omitempty"`
	EndTime    int64  `json:"endTime,omitempty"`
}

// getListFilter ：从in里获取ListFilter，如果没有就返回nil，过滤器字段必须匹配，这个项目的过滤器都要是这些字段
func getListFilter(in interface{}) *ListFilter {
	if in != nil {
		var filter ListFilter
		jsonStr, _ := json.Marshal(in)
		json.Unmarshal(jsonStr, &filter)
		return &filter
	}
	return nil
}

// PassListFilter :通用列表过滤器,要求输入db要有Model，否则不能计数
func PassListFilter(db *gorm.DB, in interface{}) (int64, *gorm.DB) {
	var totalCount int64
	filter := getListFilter(in)
	if filter == nil { //没有过滤器
		db.Count(&totalCount)
		return totalCount, db
	}

	if filter.StartTime != 0 || filter.EndTime != 0 {
		s := time.Unix(filter.StartTime, 0)
		e := time.Unix(filter.EndTime, 0)
		db = db.Where("updated_at > ?", s).Where("updated_at < ?", e)
	}

	orderBy := "updated_at" //数据库字段就是这个
	if filter.OrderBy != "" {
		orderBy = filter.OrderBy
		// 字段转换
		changeList := []([]string){
			{"createdTime", "created_at"},
			{"updatedTime", "updated_at"},
		}
		for _, c := range changeList {
			if filter.OrderBy == c[0] {
				orderBy = c[1]
			}
		}
	}

	order := "asc"
	if strings.ToLower(filter.Order) == "desc" {
		order = "desc"
	}

	db = db.Order(orderBy + " " + order)

	db.Count(&totalCount)

	pageSiz := filter.PageSize
	if filter.PageSize > 1000 {
		pageSiz = 1000
	}
	if filter.PageSize <= 0 {
		pageSiz = 1
	}

	pageNum := filter.PageNumber
	pageNumMax := totalCount / pageSiz
	// TODO： 这里有bug，最后一页零头可能取不到

	if pageNum > pageNumMax {
		pageNum = pageNumMax
	}

	db = db.Offset(pageSiz * (pageNum - 1)).Limit(pageSiz)

	return totalCount, db
}
