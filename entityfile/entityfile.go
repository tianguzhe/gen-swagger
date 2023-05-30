package entityfile

import (
	"os"
	"strings"

	"gen-swagger/model"
	"gen-swagger/utils"

	"github.com/fatih/color"
	"golang.org/x/exp/slices"
)

func GeneraterEntityFile(res map[string]model.DefinitionWrap) {

	// var list []string
	var build strings.Builder
	// build.WriteString(s1)
	// build.WriteString(s2)
	// s3 := build.String()
	var filterMap []string

	for k, v := range res {
		if len(v.Properties) == 0 {
			continue
		}

		// aa :=["KeyValueDTO","Response","PageList","PageRequest"]

		// if strings.Contains(k, "«") {
		color.Yellow("k1 %s", k)
		// }

		// if strings.Contains(k, "Response«") {
		// 	k = strings.Replace(strings.ReplaceAll(k, "Response«", ""), "»", "", 1)

		// 	if k == "string" || k == "int" {
		// 		continue
		// 	}
		// }

		// var re = regexp.MustCompile(`(?m)«(.*)»`)
		// if group := re.FindStringSubmatch(k); len(group) > 1 {
		// 	result := group[1]
		// 	if strings.Contains(result, ",") || result == "string" || result == "int" {
		// 		return
		// 	}
		// }

		color.Blue("k2 %s", k)

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

		build.WriteString("data class " + k + "(\n")
		for k1, v1 := range v.Properties {
			build.WriteString("/** " + v1.Description + " **/\n")
			build.WriteString("var " + k1 + ":" + utils.FormatType(k, v1) + "? = null,\n")
		}
		build.WriteString(")\n")
	}

	// aaa := strings.Join(list, "\n")

	// fmt.Println("======", aaa)

	// f, err := os.Open("/Users/laj/vsj/goPrj/gen-swagger/sssss.kt")
	// if err != nil {
	// 	fmt.Printf("%s", err.Error())
	// }
	os.WriteFile("/Users/laj/vsj/goPrj/gen-swagger/sssss.kt", []byte(build.String()), 0666)

}
