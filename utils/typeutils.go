package utils

import (
	"fmt"
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
		if strings.Contains(arrayRef, "KeyValueDTO«") {
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

	// // 对象
	// if prop.Ref != "" {
	// 	ref := strings.Split(prop.Ref, "/")
	// 	return ref[len(ref)-1]
	// }

	// var format = prop.Format

	// switch prop.TypeGo {
	// case "int", "int32", "integer", "number":
	// 	if format == "int64" {
	// 		return "Long"
	// 	}
	// 	if format == "double" {
	// 		return "Double"
	// 	}
	// 	if format == "float" {
	// 		return "Float"
	// 	}
	// 	if format == "int" || format == "int32" {
	// 		return "Int"
	// 	}
	// 	return "Double"
	// case "long", "int64":
	// 	return "Long"
	// case "bigdecimal", "double":
	// 	return "Double"
	// case "float":
	// 	return "Float"
	// case "string":
	// 	return "String"
	// case "boolean":
	// 	return "Boolean"
	// case "array":
	// 	if prop.Items.TypeGo != "" {
	// 		return "List<" + prop.Items.TypeGo + ">"
	// 	}
	// 	ref := strings.Split(prop.Items.Ref, "/")
	// 	if strings.HasPrefix(ref[len(ref)-1], "KeyValueDTO") {
	// 		return "List<KeyValueDTO>"
	// 	}
	// 	return "List<" + ref[len(ref)-1] + ">"
	// case "object":
	// 	if prop.AdditionalProperties.TypeGo != "" {
	// 		return prop.AdditionalProperties.TypeGo
	// 	}
	// 	ref := strings.Split(prop.AdditionalProperties.Ref, "/")
	// 	return ref[len(ref)-1]
	// default:
	// 	color.Red("%#v ==== ", prop)
	// 	return ""
	// }

}

func ReadRefObject(ref string) string {
	if strings.HasPrefix(ref, "#/definitions") {
		refSplit := strings.Split(ref, "/")

		tmpObject := refSplit[len(refSplit)-1]

		// 处理异常数据
		if strings.Contains(tmpObject, "«") {
			color.Red("泛型异常 ==== %s", tmpObject)

			replacer := strings.NewReplacer("«", "<", "»", ">")
			tmpObject = replacer.Replace(tmpObject)
		}

		return tmpObject
	}
	panic("ref 数据结构异常" + ref)
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
