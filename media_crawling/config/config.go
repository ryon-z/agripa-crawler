package config

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"gopkg.in/ini.v1"
)

// Conf : 설정값
var Conf Config

// Secret : 보안 정보
var Secret map[string]string

func init() {
	Conf = GetConfig()
	Secret = GetSecret(Conf.SecretIniPath)
}

// Config : 설정
type Config struct {
	SecretIniPath string
	LogDirPath    string
	DataDirPath   string
	SubLogDirPath string
	DumpDirPath   string
	Timezone      *time.Location
	NumRetrying   int
}

func getSubLogDirPath(logDirPath string) string {
	// 한국 시간 기준
	loc, _ := time.LoadLocation("Asia/Seoul")
	t := time.Now().In(loc)
	datetime := fmt.Sprintf("%d%02d%02d", t.Year(), t.Month(), t.Day())
	path := fmt.Sprintf("%s/%s", logDirPath, datetime)

	return filepath.FromSlash(path)
}

// GetConfig : Config 구조체 획득
func GetConfig() Config {
	conf := Config{}

	pwd := getWorkingDirPath()

	conf.SecretIniPath = makePath([]string{pwd, "secret.ini"})
	conf.LogDirPath = makePath([]string{pwd, "logs"})
	conf.DataDirPath = makePath([]string{pwd, "data"})
	conf.SubLogDirPath = getSubLogDirPath(conf.LogDirPath)
	conf.DumpDirPath = makePath([]string{pwd, "dump"})
	conf.Timezone, _ = time.LoadLocation("Asia/Seoul")
	conf.NumRetrying = 5

	return conf
}

// GetSecret : ini 파일 로드
func GetSecret(path string) map[string]string {
	cfg, err := ini.Load(path)
	if err != nil {
		// Secret.ini 로드 전이라 채널 ID를 알 수 없어 메세지 전송 불가능
		panic(err)
	}

	var secret map[string]string
	secret = make(map[string]string)

	// naver API 정보 로드
	secret["naver:id"] = cfg.Section("naver").Key("id").String()
	secret["naver:secret"] = cfg.Section("naver").Key("secret").String()

	// DB 정보 로드
	secret["db:rdbms"] = cfg.Section("db").Key("rdbms").String()
	secret["db:user"] = cfg.Section("db").Key("user").String()
	secret["db:password"] = cfg.Section("db").Key("password").String()
	secret["db:operation_host"] = cfg.Section("db").Key("operation_host").String()
	secret["db:collection_host"] = cfg.Section("db").Key("collection_host").String()
	secret["db:operation_database"] = cfg.Section("db").Key("operation_database").String()
	secret["db:collection_database"] = cfg.Section("db").Key("collection_database").String()
	secret["db:port"] = cfg.Section("db").Key("port").String()

	// youtube API 정보 로드
	secret["youtube:apikey"] = cfg.Section("youtube").Key("apikey").String()

	// slack token 로드
	secret["slack:token"] = cfg.Section("slack").Key("token").String()
	secret["slack:channel_id"] = cfg.Section("slack").Key("channel_id").String()

	// data.go.kr 로드
	secret["data.go.kr:secret"] = cfg.Section("data.go.kr").Key("secret").String()

	// mafra 로드
	secret["mafra:secret"] = cfg.Section("mafra").Key("secret").String()

	return secret
}

// getWorkingDirPath : Working Directory 경로 얻기
// util.getWorkingDirPath 사용 시 alarm과 config, util이 순환 import가 되기 때문에 자체 함수 생성
func getWorkingDirPath() string {
	var (
		_, b, _, _ = runtime.Caller(0)
		basepath   = filepath.Dir(filepath.Dir(b))
	)

	return basepath
}

// makePath : 문자열 배열을 받아 조합하여 경로 생성
// util.MakePath 사용 시 alarm과 config, util이 순환 import가 되기 때문에 자체 함수 생성
func makePath(pathElems []string) string {
	path := strings.Join(pathElems, "/")

	return filepath.FromSlash(path)
}
