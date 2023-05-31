package apifile

import (
	"fmt"
	"gen-swagger/model"
	"gen-swagger/utils"
	"os"
	"strings"
)

type ApiGenModel struct {
	ApiMethod string
	FuncName  string
	UrlPath   string
	Params    []model.Parameter
	Summary   string
	Response  string
}

func GeneraterApiFile(paths map[string]map[string]model.Path, serviceName string) {

	urlPathMap := make(map[string][]ApiGenModel, 20)

	var strBuild strings.Builder
	var providerBuild strings.Builder

	for urlPath, v := range paths {

		tmpParams := make([]model.Parameter, 0)

		if !strings.Contains(urlPath, "action") ||
			strings.Contains(urlPath, "excel") ||
			strings.Contains(urlPath, "meituan") ||
			strings.Contains(urlPath, "/pv/") ||
			strings.Contains(urlPath, "/pm/") {
			continue
		}

		pathParams := strings.Split(urlPath, "/")

		urlService := utils.UpperCamelCase(pathParams[len(pathParams)-3])

		if urlPathMap[urlService] == nil {
			urlPathMap[urlService] = make([]ApiGenModel, 0)
		}

		for ApiMethod, pathInfo := range v {

			for _, v := range pathInfo.Parameters {
				if v.In != "header" && v.In != "path" {
					tmpParams = append(tmpParams, v)
				}
			}

			if pathInfo.Responses.HttpOk.Schema.Ref == "" {
				continue
			}

			responseRef := utils.ReadRefObject(pathInfo.Responses.HttpOk.Schema.Ref)

			urlPathMap[urlService] = append(urlPathMap[urlService], ApiGenModel{
				FuncName:  utils.LowerCamelCase(pathParams[len(pathParams)-1]),
				ApiMethod: ApiMethod,
				UrlPath:   urlPath,
				Summary:   pathInfo.Summary,
				Params:    tmpParams,
				Response:  responseRef,
			})

		}

	}

	strBuild.WriteString(fmt.Sprintf("interface %sApiServer {\n", utils.UpperCamelCase(serviceName)))
	providerBuild.WriteString(fmt.Sprintf("object Provide%sApiServer {\n", utils.UpperCamelCase(serviceName)))

	for k, v := range urlPathMap {

		strBuild.WriteString(fmt.Sprintf("interface %s {\n", k))

		providerBuild.WriteString(fmt.Sprintf("fun provide%s(): %s.%s = \n", k, utils.UpperCamelCase(serviceName), k))
		providerBuild.WriteString(fmt.Sprintf("RetrofitFactory.instance.createService(%s.%s::class.java)\n", utils.UpperCamelCase(serviceName), k))

		for _, v2 := range v {

			strBuild.WriteString("/**\n")
			strBuild.WriteString(fmt.Sprintf("@description  %s\n", v2.Summary))

			for _, param := range v2.Params {
				strBuild.WriteString(fmt.Sprintf("@%s %s : %s (require : %t)\n", param.In, param.Name, getType(param), param.Required))
			}
			strBuild.WriteString("*/\n")

			if v2.ApiMethod == "get" {
				strBuild.WriteString(fmt.Sprintf("@GET(\"%s%s\")\n", serviceName, v2.UrlPath))
				strBuild.WriteString(fmt.Sprintf("fun %s(@QueryMap params: MutableMap<String, Any>): %s\n\n", v2.FuncName, strings.ReplaceAll(v2.Response, "Response", "BaseModel")))
			} else if v2.ApiMethod == "post" {
				strBuild.WriteString(fmt.Sprintf("@POST(\"%s%s\")\n", serviceName, v2.UrlPath))
				if len(v2.Params) != 0 {
					strBuild.WriteString(fmt.Sprintf("fun %s(@Body %s: %s): %s\n\n", v2.FuncName, v2.Params[0].Name, utils.ReadRefObject(v2.Params[0].Schema.Ref), strings.ReplaceAll(v2.Response, "Response", "BaseModel")))
				} else {
					strBuild.WriteString(fmt.Sprintf("fun %s(): %s\n\n", v2.FuncName, strings.ReplaceAll(v2.Response, "Response", "BaseModel")))
				}
			} else if v2.ApiMethod == "put" {
				strBuild.WriteString(fmt.Sprintf("@PUT(\"%s%s\")\n", serviceName, v2.UrlPath))
				strBuild.WriteString(fmt.Sprintf("fun %s(@Body %s: %s): %s\n\n", v2.FuncName, v2.Params[0].Name, utils.ReadRefObject(v2.Params[0].Schema.Ref), strings.ReplaceAll(v2.Response, "Response", "BaseModel")))
			}

		}

		strBuild.WriteString("}\n")

	}

	strBuild.WriteString("}\n")
	providerBuild.WriteString("}\n")

	os.MkdirAll("./apiservice/", 0755)
	os.WriteFile("./apiservice/"+utils.UpperCamelCase(serviceName)+"ApiService.kt", []byte(strBuild.String()), 0666)

	os.MkdirAll("./provider/", 0755)
	os.WriteFile("./provider/"+"Provide"+utils.UpperCamelCase(serviceName)+"ApiService.kt", []byte(providerBuild.String()), 0666)

}

func getType(p model.Parameter) string {
	if p.TypeGo != "" {
		return p.TypeGo
	}

	return utils.ReadRefObject(p.Schema.Ref)
}
