package utils

import (
	"fmt"
	"regexp"
	"strings"

	"gen-swagger/model"

	"github.com/fatih/color"
)

func FormatType(className string, prop model.Properties) string {

	if prop.Ref != "" {
		if strings.Contains(className, "<T>") {
			return "T"
		} else {
			return ReadRefObject(prop.Ref)
		}
	}

	typeName := prop.TypeGo
	typeFormat := prop.Format

	arrayItemType := prop.Items.TypeGo
	arrayItemRef := prop.Items.Ref

	objectType := prop.AdditionalProperties.TypeGo
	objectAdditionRef := prop.AdditionalProperties.Ref

	switch typeName {
	case "array":
		if arrayItemType != "" {
			return fmt.Sprintf("List<%s>", _nativeType(arrayItemType, "not format"))
		}

		// 泛型处理
		if strings.Contains(className, "<T>") {
			return "List<T>"
		}

		arrayRef := ReadRefObject(arrayItemRef)

		// 历史遗留问题 兼容
		if strings.Contains(arrayRef, "KeyValueDTO<") {
			return "List<KeyValueDTO>"
		}

		return arrayRef
	case "object":
		if objectType != "" {
			return fmt.Sprintf("List<%s>", _nativeType(objectType, "not format"))
		}

		objectRef := ReadRefObject(objectAdditionRef)

		return objectRef

	default:
		return _nativeType(typeName, typeFormat)
	}
}

func ReadRefObject(ref string) string {
	if strings.HasPrefix(ref, "#/definitions") {
		refSplit := strings.Split(ref, "/")

		tmpObject := refSplit[len(refSplit)-1]

		// 处理异常数据
		if strings.Contains(tmpObject, "«") {

			replacer := strings.NewReplacer("«", "<", "»", ">")
			tmpObject = replacer.Replace(tmpObject)

			if strings.Contains(tmpObject, "KidRequest") {
				return "KidRequest"
			}

			re := regexp.MustCompile("<(.*)>")
			result := re.FindStringSubmatch(tmpObject)
			formatType := _nativeType(result[1], "not format")

			tmpObject = strings.ReplaceAll(tmpObject, result[1], formatType)
			color.Red("泛型异常 ==== %s", tmpObject)

			return tmpObject

		}

		return tmpObject
	}

	color.Red("ref 数据结构异常 === %s" + ref)
	return ref
}

func _nativeType(typeName string, format string) string {
	switch typeName {
	case "int", "int32", "integer", "number":
		if format == "int64" {
			return "Long"
		}
		if format == "double" {
			return "Double"
		}
		if format == "float" {
			return "Float"
		}
		if format == "int" || format == "int32" {
			return "Int"
		}
		return "Double"
	case "long", "int64":
		return "Long"
	case "bigdecimal", "double":
		return "Double"
	case "float":
		return "Float"
	case "string":
		return "String"
	case "boolean":
		return "Boolean"
	default:
		color.Red("Error type === %s", typeName)
		return typeName
	}
}
