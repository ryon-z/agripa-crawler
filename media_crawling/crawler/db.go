package crawler

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"media_crawling/alarm"
	"media_crawling/config"
	"media_crawling/models"
	"media_crawling/util"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func getHostAndDbName(dbUsage string) (string, string) {
	var host string
	var dbName string
	secret := config.Secret
	checkDbUsageCorrect(dbUsage)
	if dbUsage == "operation" {
		host = secret["db:operation_host"]
		dbName = secret["db:operation_database"]
	} else {
		host = secret["db:collection_host"]
		dbName = secret["db:collection_database"]
	}

	return host, dbName
}

func checkDbUsageCorrect(dbUsage string) {
	if !util.InArray(dbUsage, []string{"operation", "collection"}) {
		errorMessage := fmt.Sprintf("유효하지 않은 dbUsage. 입력 받은 dbUsage: %s", dbUsage)
		alarm.PostMessage("default", errorMessage)
		panic(errors.New(errorMessage))
	}
}

// GetDB : db 커넥션 생성
func GetDB(dbUsage string) *gorm.DB {
	secret := config.Secret
	host, dbName := getHostAndDbName(dbUsage)

	source := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", secret["db:user"], secret["db:password"], host, secret["db:port"], dbName)
	db, err := gorm.Open(secret["db:rdbms"], source)
	if err != nil {
		alarm.PostMessage("default", err.Error())
		panic(err)
	}

	return db
}

// CloseDB : db 커넥션 끊음
func CloseDB(db *gorm.DB) {
	db.Close()
}

// DatetimeFormat : 현재 사용하는 DB의 datetime 데이터 유형의 예시, Golang에서 time의 형변환 시 사용
const DatetimeFormat = "2006-01-02 15:04:05"

// ImportDataToDB : data를 받아서 models.Model Table로 import
func ImportDataToDB(db *gorm.DB, modelName string, tableName string, data []interface{}) {
	model := models.GetModel(modelName)
	if tableName == "" {
		tableName = model.TableName()
	}
	ModelColumns := model.Columns()
	var columns []string

	// 데이터가 없으면 함수 종료
	if len(data) <= 0 {
		errorMessage := fmt.Sprintf("DB에 insert 할 데이터가 없습니다. modelName: %s", modelName)
		alarm.PostMessage("default", errorMessage)
		panic(errors.New(errorMessage))
	}

	// columns 문자열 배열 획득
	row := data[0]
	v := reflect.ValueOf(row)
	typeOfRow := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fieldName := typeOfRow.Field(i).Name
		if util.InArray(fieldName, ModelColumns) {
			columns = append(columns, typeOfRow.Field(i).Name)
		}
	}

	sql := fmt.Sprintf("INSERT IGNORE INTO %s (%s) VALUES ", tableName, strings.Join(columns, ","))

	for index0, row := range data {
		var sqlElems []string
		v := reflect.ValueOf(row)

		for i := 0; i < v.NumField(); i++ {
			fieldName := typeOfRow.Field(i).Name
			originFieldValue := v.Field(i).Interface()
			var fieldValue string
			switch fieldValueType := (v.Field(i).Interface()).(type) {
			case string:
				fieldValue = originFieldValue.(string)
			case int:
				fieldValue = strconv.Itoa(originFieldValue.(int))
			default:
				errorMessage := fmt.Sprintf("구조체 fieldValue type error. type: %s", fieldValueType)
				alarm.PostMessage("default", errorMessage)
				panic(errors.New(errorMessage))
			}

			if util.InArray(fieldName, ModelColumns) {
				fieldValue = strings.ReplaceAll(fieldValue, "\"", "")
				elem := "\"" + fieldValue + "\""
				sqlElems = append(sqlElems, elem)
			}
		}

		rowSQL := strings.Join(sqlElems, ",")
		rowSQL = "(" + rowSQL
		if util.IsLastElement(index0, len(data)) {
			rowSQL = rowSQL + ");"
		} else {
			rowSQL = rowSQL + "), "
		}

		sql += rowSQL
	}
	// fmt.Println(sql)

	// 데이터 삽입
	result := db.Exec(sql)
	if result.Error != nil {
		fmt.Println("import data to db 에러 발생")

		alarm.PostMessage("default", result.Error.Error())
		panic(result.Error)
	}

	fmt.Println("import data to db 정상 종료")
}

// RunQuery : 단순 쿼리 실행
func RunQuery(db *gorm.DB, sqlQuery string) {
	fmt.Printf("RunQuery :: sqlQuery: %s\n", sqlQuery)

	result := db.Exec(sqlQuery)
	if result.Error != nil {
		alarm.PostMessage("default", result.Error.Error())
		panic(result.Error)
	}
}

// Dump : dump 디렉토리에 특정 테이블을 dump 함
func Dump(dbUsage string, tableName string) {
	host, dbName := getHostAndDbName(dbUsage)
	secret := config.Secret
	dumpDirPath := config.Conf.DumpDirPath
	functionName := "Dump"

	util.MakeDirIfNotExists(dumpDirPath)

	hostPhrase := fmt.Sprintf("-h%s", host)
	userPhrase := fmt.Sprintf("-u%s", secret["db:user"])
	passwordPhrase := fmt.Sprintf("-p%s", secret["db:password"])
	portPhrase := fmt.Sprintf("-P%s", secret["db:port"])
	dbNamePhrase := fmt.Sprintf("%s", dbName)

	cmd := exec.Command("mysqldump", hostPhrase, userPhrase, passwordPhrase, portPhrase, dbNamePhrase, tableName)
	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()
	stdout := outbuf.String()
	stderr := errbuf.String()
	util.CheckCondition(stderr != "", functionName, stderr)
	util.CheckError(err, functionName)

	dumpFilePath := fmt.Sprintf("%s/%s.sql", dumpDirPath, tableName)
	err = ioutil.WriteFile(dumpFilePath, []byte(stdout), 0644)
	util.CheckError(err, functionName)
}

// ImportDumpFile : 특정 dump 파일을 import
func ImportDumpFile(dbUsage string, tableName string) {
	host, dbName := getHostAndDbName(dbUsage)
	secret := config.Secret
	dumpDirPath := config.Conf.DumpDirPath
	functionName := "ImportDumpFile"

	dumpFilePath := fmt.Sprintf("%s/%s.sql", dumpDirPath, tableName)
	hostPhrase := fmt.Sprintf("-h%s", host)
	userPhrase := fmt.Sprintf("-u%s", secret["db:user"])
	passwordPhrase := fmt.Sprintf("-p%s", secret["db:password"])
	portPhrase := fmt.Sprintf("-P%s", secret["db:port"])
	dbNamePhrase := fmt.Sprintf("%s", dbName)

	cmd := exec.Command("mysql", hostPhrase, userPhrase, passwordPhrase, portPhrase, dbNamePhrase)
	dump, err := os.Open(dumpFilePath)
	util.CheckError(err, functionName)
	cmd.Stdin = dump

	var errbuf bytes.Buffer
	cmd.Stderr = &errbuf

	err = cmd.Run()
	stderr := errbuf.String()
	util.CheckCondition(stderr != "", functionName, stderr)
	util.CheckError(err, functionName)
}

// MoveDataUsingDump : mariaDB shell dump 명령을 통해 fromDB에서 toDB로 테이블을 복사
func MoveDataUsingDump(fromDbUsage string, toDbUsage string, tableName string, copyOrigin bool) {
	fromDB := GetDB(fromDbUsage)
	defer CloseDB(fromDB)
	toDB := GetDB(toDbUsage)
	defer CloseDB(toDB)
	newTableName := tableName + "_NEW"

	if copyOrigin {
		RunQuery(fromDB, fmt.Sprintf("CREATE TABLE %s LIKE %s", newTableName, tableName))
		RunQuery(fromDB, fmt.Sprintf("INSERT INTO %s SELECT * FROM %s", newTableName, tableName))
	} else {
		RunQuery(fromDB, fmt.Sprintf("RENAME TABLE %s TO %s", tableName, newTableName))
	}

	Dump(fromDbUsage, newTableName)

	if copyOrigin {
		RunQuery(fromDB, fmt.Sprintf("DROP TABLE %s", newTableName))
	} else {
		RunQuery(fromDB, fmt.Sprintf("RENAME TABLE %s TO %s", newTableName, tableName))
	}

	ImportDumpFile(toDbUsage, newTableName)

	tmpTableName := tableName + "_TMP"
	_, toDbName := getHostAndDbName(toDbUsage)
	if IsTableExists(toDB, toDbName, tableName) {
		RunQuery(toDB, fmt.Sprintf("RENAME TABLE %s to %s", tableName, tmpTableName))
		RunQuery(toDB, fmt.Sprintf("RENAME TABLE %s to %s", newTableName, tableName))
		RunQuery(toDB, fmt.Sprintf("DROP TABLE %s", tmpTableName))
	} else {
		RunQuery(toDB, fmt.Sprintf("RENAME TABLE %s to %s", newTableName, tableName))
	}
}

// MoveDataUsingGorm : Gorm 기능을 이용하여 fromDB에서 toDB로 테이블을 복사
func MoveDataUsingGorm(fromDbUsage string, toDbUsage string, modelName string, sqlQuery string) {
	fromDB := GetDB(fromDbUsage)
	defer CloseDB(fromDB)
	toDB := GetDB(toDbUsage)
	defer CloseDB(toDB)

	model := models.GetModel(modelName)
	tableName := model.TableName()

	CreateRemoteTableIfNotExists(fromDB, toDB, tableName, tableName)

	modelArrAddr := models.GetModelArrAddr(modelName)
	fmt.Printf("Run Raw Query :: sqlQuery: %s", sqlQuery)
	fromDB.Raw(sqlQuery).Scan(modelArrAddr)

	data := models.GetInterfaceArr(modelArrAddr)
	ImportDataToDB(toDB, modelName, "", data)
}

// CreateTableIfNotExists : 기존 테이블을 참고하여 테이블 생성
func CreateTableIfNotExists(db *gorm.DB, newTableName string, referenceTableName string) {
	sqlQuery := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s LIKE %s", newTableName, referenceTableName)

	// 데이터 삽입
	result := db.Exec(sqlQuery)
	util.CheckError(result.Error, "CreateTableIfNotExists")
}

type showCreateTableResponse struct {
	Table string `gorm:"column:Table"`
	Query string `gorm:"column:Create Table"`
}

// CreateRemoteTableIfNotExists : 기존 테이블을 참고하여 다른 DB에 테이블 생성
func CreateRemoteTableIfNotExists(fromDB *gorm.DB, toDB *gorm.DB, newTableName string, referenceTableName string) {
	var response showCreateTableResponse
	sqlQuery := fmt.Sprintf("SHOW CREATE TABLE %s", referenceTableName)
	result := fromDB.Raw(sqlQuery).Scan(&response)
	util.CheckError(result.Error, "CreateRemoteTableIfNotExists")

	sqlQuery = strings.ReplaceAll(response.Query, "CREATE TABLE", "CREATE TABLE IF NOT EXISTS")
	RunQuery(toDB, sqlQuery)
}

// OneColumnQueryResponse : 하나의 컬럼의 결과값을 갖는 쿼리 결과값 구조체
type OneColumnQueryResponse struct {
	Result string `gorm:"column:result"`
}

// IsTableExists : 입력 받은 테이블이 입력 받은 DB에 존재하는지 확인
func IsTableExists(db *gorm.DB, dbName string, tableName string) bool {
	var response OneColumnQueryResponse
	sqlQuery := fmt.Sprintf(`
		SELECT table_name AS result 
		FROM information_schema.tables
		WHERE table_schema = "%s"
		AND table_name = "%s"
	;`, dbName, tableName)
	fmt.Println("sqlQuery", sqlQuery)
	db.Raw(sqlQuery).Scan(&response)

	return response.Result != ""
}
