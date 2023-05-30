package main

import (
	"io"
	"net/http"

	"gen-swagger/apifile"
	"gen-swagger/model"

	"github.com/bytedance/sonic"
	"github.com/fatih/color"
)

func getUrl(url string) {

	var respMap model.Swagger

	resp, err := http.Get(url)
	if err != nil {
		color.Red("error %s", err.Error())
	}

	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)

	sonic.Unmarshal(bytes, &respMap)

	// 生成 数据模型
	// entityfile.GeneraterEntityFile(respMap.Definitions)

	//生成 api 数据
	apifile.GeneraterApiFile(respMap.Paths)

}

func main() {
	getUrl("https://petstore.swagger.io/v2/swagger.json")
}
