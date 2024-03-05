package frame

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"
	"text/template"
)

const (
	metadataTpl = `package api

import (
    {{- range $i, $r := .packages }}
    {{ $r }}
    {{- end }}

	"github.com/xm-chentl/gocore/frame/handles"
)

func init() {
	handles.Register(handles.Handlers{
		{{- range $i, $r := .entrys }}
		{{ $r.Register }}{{ end }}
	})
}`
)

var (
	// 0 项目名 1`
	formatQuote = `"%s/%s%s"`
	regexAPI    = regexp.MustCompile("[A-Za-z0-9]+API")
)

type apiInfo struct {
	Import   string
	Register string

	endpoint string // 终端名
	name     string //
	path     string // api文件路径
	quote    string // 引用名
	target   string // 目标
}

/*
	1. 读取目录
	2. 分析结构
	3. 读取API结构体文件
	4. 生成源文件
*/

func GenerateAPI(dir string) (err error) {
	virtualPath := dir
	apiDir := dir
	rootDir, _ := os.Getwd()
	if virtualPath != "" && virtualPath[len(virtualPath)-1:] == "/" {
		virtualPath = virtualPath[0 : len(virtualPath)-1]
	}
	// fmt.Println("root > ", rootDir, " | ", apiDir, " | ", virtualPath)
	if apiDir == "" {
		apiDir = path.Join(rootDir, "internal", "api")
	} else {
		apiDir = path.Join(rootDir, apiDir)
	}

	apiDirFileArray, err := os.ReadDir(apiDir)
	if err != nil {
		return
	}

	apiInfoArray := make([]apiInfo, 0)
	for _, dirEntry := range apiDirFileArray {
		if dirEntry.IsDir() {
			getApiInfo(apiDir, dirEntry, &apiInfoArray)
		}
	}

	var wg sync.WaitGroup
	for index, apiInfo := range apiInfoArray {
		wg.Add(1)
		go getApiName(apiInfo.path, &apiInfoArray[index], &wg)
	}
	wg.Wait()

	// 模块地址
	modelName := getModelName()
	// 整合数值
	packageMap := make(map[string]string)
	packageExistMap := make(map[string]string)
	deleteIndexArray := make([]int, 0)
	for index, info := range apiInfoArray {
		if info.name == "" {
			deleteIndexArray = append(deleteIndexArray, index)
			continue
		}

		apiInfoArray[index].target = strings.ReplaceAll(path.Dir(info.target), apiDir, "")
		importQuote := fmt.Sprintf(formatQuote, modelName, virtualPath, apiInfoArray[index].target)
		// fmt.Println(">>>> ", apiInfoArray[index].target)
		_, ok := packageExistMap[info.quote]
		if ok {
			if _, ok = packageMap[importQuote]; !ok {
				apiInfoArray[index].quote += "1"
				importQuote = apiInfoArray[index].quote + " " + importQuote
				packageMap[importQuote] = apiInfoArray[index].quote
				packageExistMap[apiInfoArray[index].quote] = importQuote
			}
		} else {
			packageMap[importQuote] = info.quote
			packageExistMap[info.quote] = importQuote
		}

		apiInfoArray[index].Import = importQuote
		content := fmt.Sprintf(`"%s/%s": &%s.%s{},`,
			apiInfoArray[index].target,
			apiInfoArray[index].endpoint,
			strings.ReplaceAll(apiInfoArray[index].quote, "-", ""), //  todo: 过滤特殊处理(package 不含特殊字符，如包名为login-log， package loginlog)
			info.name,
		)
		apiInfoArray[index].Register = content
	}
	fmt.Println(">>>> ", apiInfoArray[0].target)
	for _, index := range deleteIndexArray {
		apiInfoArray = append(apiInfoArray[:index], apiInfoArray[index+1:]...)
	}

	var packageArray []string
	for key := range packageMap {
		packageArray = append(packageArray, key)
	}

	metadataTemple, err := template.New("").Parse(metadataTpl)
	if err != nil {
		return
	}

	filePath := path.Join(apiDir, "register.go")
	if _, err = os.Stat(filePath); err != nil {
		// 不存在
		_, _ = os.Create(filePath)
	}

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return
	}
	// 清空内容
	_ = file.Truncate(0)
	_, _ = file.Seek(0, 0)
	defer file.Close()

	writer := bufio.NewWriter(io.MultiWriter(file))
	err = metadataTemple.Execute(writer, map[string]interface{}{
		"entrys":   apiInfoArray,
		"packages": packageArray,
	})
	if err != nil {
		return
	}
	err = writer.Flush()

	return
}

func getApiInfo(dirPath string, entry os.DirEntry, apiInfoArray *[]apiInfo) {
	if entry.IsDir() {
		dirPath = path.Join(dirPath, entry.Name())
		dirs, _ := os.ReadDir(dirPath)
		for _, dirEntry := range dirs {
			getApiInfo(dirPath, dirEntry, apiInfoArray)
		}
	} else {
		dirNameArray := strings.Split(dirPath, "/")
		*apiInfoArray = append(*apiInfoArray, apiInfo{
			endpoint: strings.Split(entry.Name(), ".")[0],
			path:     path.Join(dirPath, entry.Name()),
			quote:    dirNameArray[len(dirNameArray)-1],
			target:   path.Join(dirPath, entry.Name()),
		})
	}
}

func getApiName(apiPath string, apiInfo *apiInfo, wg *sync.WaitGroup) {
	var err error
	var goFile *os.File
	var name string
	defer func() {
		if goFile != nil {
			goFile.Close()
		}
		if err != nil {
			log.Fatal(err)
		}
		apiInfo.name = name
		wg.Done()
	}()

	goFile, err = os.Open(apiPath)
	if err != nil {
		return
	}

	content, err := io.ReadAll(goFile)
	if err != nil {
		return
	}

	nameArray := regexAPI.FindAllString(string(content), 1)
	if len(nameArray) == 0 {
		return
	}
	name = nameArray[0]
}

func getModelName() (res string) {
	root, _ := os.Getwd()
	modFilePath := path.Join(root, "go.mod")
	f, err := os.Open(modFilePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		res = scanner.Text()
		break
	}
	if res != "" {
		res = strings.ReplaceAll(res, "module ", "")
	}

	return
}
