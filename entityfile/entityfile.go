package entityfile

import (
	"os"
	"strings"

	"gen-swagger/model"
	"gen-swagger/utils"

	"golang.org/x/exp/slices"
)

func GeneraterEntityFile(res map[string]model.DefinitionWrap, serviceName string) {

	var strBuild strings.Builder

	var filterMap []string

	for k, v := range res {
		if len(v.Properties) == 0 {
			continue
		}

		if strings.HasPrefix(k, "Response") {
			k = "Response<T>"
		} else if strings.HasPrefix(k, "PageList") {
			k = "PageList<T>"
		} else if strings.HasPrefix(k, "PageRequest") {
			k = "PageRequest<T>"
		} else if strings.HasPrefix(k, "KidRequest") {
			k = "KidRequest"
		} else if strings.HasPrefix(k, "KeyValueDTO") {
			k = "KeyValueDTO"
		}

		if !slices.Contains(filterMap, k) {
			filterMap = append(filterMap, k)
		} else {
			continue
		}

		strBuild.WriteString("data class " + k + "(\n")
		for k1, v1 := range v.Properties {
			strBuild.WriteString("/** " + v1.Description + " **/\n")
			strBuild.WriteString("var " + k1 + ":" + utils.FormatType(k, v1) + "? = null,\n")
		}
		strBuild.WriteString(")\n")
	}

	os.MkdirAll("./entity/", 0755)

	os.WriteFile("./entity/"+utils.UpperCamelCase(serviceName)+"Entity.kt", []byte(strBuild.String()), 0666)
}
