package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"gen-swagger/apifile"
	"gen-swagger/entityfile"
	"gen-swagger/model"

	"github.com/bytedance/sonic"
	"github.com/fatih/color"
)

func getUrl(url string, serviceName string) {

	var respMap model.Swagger

	resp, err := http.Get(url)
	if err != nil {
		color.Red("error %s", err.Error())
	}

	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)

	sonic.Unmarshal(bytes, &respMap)

	// 生成 数据模型
	entityfile.GeneraterEntityFile(respMap.Definitions, serviceName)

	// //生成 api 数据
	apifile.GeneraterApiFile(respMap.Paths, serviceName)

}

type ApiConfig struct {
	ServiceList []struct {
		ServiceName string `json:"serviceName"`
		FileName    string `json:"fileName"`
		ApiParams   string `json:"apiParams"`
	} `json:"serviceNames"`
	Api struct {
		BaseUrl string `json:"baseUrl"`
		Version string `json:"version"`
	} `json:"api"`
}

func main() {

	var config ApiConfig

	f, err := os.Open("./apiC.json")
	if err != nil {
		color.Red("read config json error: %s \n", err.Error())
	}
	bytes, _ := io.ReadAll(f)
	sonic.Unmarshal(bytes, &config)

	for _, v := range config.ServiceList {

		finalUrl := fmt.Sprintf("%s/%s/%s?%s", config.Api.BaseUrl, v.ServiceName, config.Api.Version, v.ApiParams)

		getUrl(finalUrl, v.ServiceName)
	}

	// getUrl("https://petstore.swagger.io/v2/swagger.json")
}
