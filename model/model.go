package model

import (
	"JYGO/config"
	"JYGO/db"
	"fmt"
	"reflect"
)

var (
	modelContainer = &models {
		container : make(map[string]*modelInfo),
		containerByFullName:make(map[string]*modelInfo),
	}
)

type modelInfo struct {
	modelName		string
	model			interface{}
	fields			[]string
	table 			string
	connection 		string

}

type models struct {
	container 				map[string]*modelInfo
	containerByFullName		map[string]*modelInfo
}

type ContractModel interface {
	GetOne()	*modelInfo
	Get()		[]*modelInfo
	Count()		int64
}

func Init() {
	_ = config.Init()
	_ = db.Db.Init()
}

func RegisterModel(models ...interface{}) {
	for _, model := range models {
		register(model)
	}
}

func register(model interface{}) {
	sv := reflect.ValueOf(model)
	st := reflect.Indirect(sv).Type()

	if sv.Kind() != reflect.Ptr {
		panic(fmt.Errorf("cannot use non-ptr model struct `%s`", getFullName(sv)))
	}

	if st.Kind() == reflect.Ptr {
		panic(fmt.Errorf("only allow ptr model struct, it looks you use two reference to the struct `%s`", st))
	}

	name := getFullName(sv)

	if _,ok := modelContainer.fetchModelByFullName(name); ok {
		panic(fmt.Errorf("model `%s` repeat register, must be unique\n", name))
	}

	table := getTableName(sv)

	if _,ok := modelContainer.fetchModelByTable(table); ok {
		panic(fmt.Errorf("table name `%s` repeat register, must be unique\n", table))
	}

	modelInfo := newModelInfo(sv)
	modelInfo.model = model
	modelInfo.table = table

	modelContainer.add(table,modelInfo)

}

func (m *models)fetchModelByFullName(name string)(*modelInfo,bool) {
	mi,ok := m.containerByFullName[name]
	return mi, ok
}

func (m *models)fetchModelByTable(table string) (*modelInfo,bool) {
	mi,ok := m.container[table]
	return mi,ok
}

func (m *models)fetchModel(model interface{},needPtr bool)(*modelInfo,reflect.Value,bool) {
	sv := reflect.ValueOf(model)
	snd := reflect.Indirect(sv)

	if needPtr && sv.Kind() != reflect.Ptr {
		panic(fmt.Errorf("cannot use non-ptr model struct `%s`", getFullName(sv)))
	}

	name := getFullName(sv)

	if mi,ok := m.fetchModelByFullName(name); ok {
		return mi,snd,ok
	}
	return nil,snd,false
}

func (m *models)add(table string,model *modelInfo) bool {
	m.container[table] = model
	m.containerByFullName[model.modelName] = model

	return true
}

// Read 获取指定model的Reader 从而可以进行数据的读取
func Read(model interface{}) ReaderContract {
	mi ,_, ok := modelContainer.fetchModel(model,true)
	if !ok {
		panic(fmt.Errorf("model `%s` has not been registered！", reflect.Indirect(reflect.ValueOf(model)).Type().Name()))
	}
	r := new(Reader)

	r.model = mi
	return r
}

func Add(model interface{}) (lastInsertId int64) {
	mi,snd,ok := modelContainer.fetchModel(model,false)

	if !ok {
		panic(fmt.Errorf("model `%s` has not been registered！", reflect.Indirect(reflect.ValueOf(model)).Type().Name()))
	}

	var insertData []interface{}

	for index,field := range mi.fields {
		fieldObj := snd.Field(index)

		insertData = append(insertData,[]interface{}{field,fieldObj.Interface()})
	}

	connect := getConnector(mi)

	result,err := connect.Table(mi.table).Add(insertData...)

	if err != nil {
		panic(fmt.Errorf("Insert `%s` Failed ",mi.table))
	}

	lastInsertId ,_ = result.LastInsertId()

	return lastInsertId
}


func getConnector(mi *modelInfo) db.BaseDbContract {
	connect := db.Db.Connector()
	if mi.connection != "" {
		connect = db.Db.GetConnection(mi.connection)
	}

	return connect
}
