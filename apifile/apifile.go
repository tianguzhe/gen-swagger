package apifile

import (
	"fmt"
	"gen-swagger/model"
	"gen-swagger/utils"
	"os"
	"strings"

	"github.com/fatih/color"
)

type ApiGenModel struct {
	ApiMethod string
	FuncName  string
	UrlPath   string
	Params    []model.Parameter
	Summary   string
	Response  string
}

func GeneraterApiFile(paths map[string]map[string]model.Path) {

	urlPathMap := make(map[string][]ApiGenModel, 20)

	var build strings.Builder

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
				if v.In != "header" {
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

	for k, v := range urlPathMap {

		build.WriteString(fmt.Sprintf("interface %s {", k))

		for _, v2 := range v {

			// @GET("jinyeyahe-order/[api_version]/pt/product-orders/action/app-cancel")
			// fun appCanCel(@QueryMap params: MutableMap<String, Any>): Observable<BaseModel<Boolean>>

			color.Yellow("v2.Response %s", v2.Response)

			if v2.ApiMethod == "get" {
				build.WriteString(fmt.Sprintf("@GET(\"%s\")\n", v2.UrlPath))
				build.WriteString(fmt.Sprintf("fun %s(@QueryMap params: MutableMap<String, Any>): %s\n\n", v2.FuncName, strings.ReplaceAll(v2.Response, "Response", "BaseModel")))
			} else if v2.ApiMethod == "post" {
				// @POST("jinyeyahe-order/[api_version]/pt/product-orders/action/storeMeituanDeliver")
				// fun storeMeituanDeliver(@Body dto: ManagerOperaDTO): Observable<BaseModel<Boolean>>

				build.WriteString(fmt.Sprintf("@POST(\"%s\")\n", v2.UrlPath))
				fmt.Println(v2.UrlPath)
				build.WriteString(fmt.Sprintf("fun %s(@Body %s: %s): %s\n\n", v2.FuncName, v2.Params[0].Name, utils.ReadRefObject(v2.Params[0].Schema.Ref), strings.ReplaceAll(v2.Response, "Response", "BaseModel")))

			}

		}

		build.WriteString("}\n\n")

	}

	os.WriteFile("/Users/laj/vsj/goPrj/gen-swagger/api.kt", []byte(build.String()), 0666)

}
