package biligo

import (
	"embed"
	"fmt"
	"reflect"
	"strings"
	"text/template"
)

const (
	typeTemplatesDir = "template/type"
	typeTemplatesExt = ".tmpl"
)

//go:embed template/type/*.tmpl
var typeTemplates embed.FS

var typeTemplate = template.New("biligo")

var loadedTemplates = make(map[string]struct{}) // type name without path

func init() {
	typeTemplate.Funcs(template.FuncMap{
		"TypeOf":     reflect.TypeOf,
		"TrimPrefix": strings.TrimPrefix,
		"TrimSpace":  strings.TrimSpace,
		"Join":       strings.Join,
		"FmtNum":     formatNumber[int],
		"FmtDur":     formatDuration[int],
		"FmtItv":     formatInterval[int],
		"FmtTime":    formatTime[int],
		"Percent":    formatPercent[int],
		"LST":        liveStatusText[int],
	})
	const MAIN_TEMPLATE_PLACEHOLDER = // 不使用主模板
	`THIS IS A PLACEHOLDER. IF YOU SEE THIS, THEN YOU HAVE MISTAKENLY CALLED [*template.Template.Execute] FOR {{TypeOf .}}`
	_, err := typeTemplate.Parse(MAIN_TEMPLATE_PLACEHOLDER)
	if err != nil {
		panic(err)
	}
}

// getTypeName 返回类型名称, 无包名
func getTypeName[T any]() (typ string, isInside bool) {
	t := reflect.TypeFor[T]()
	if t == nil {
		panic(fmt.Errorf("type %T is nil", t))
	}
	parts := strings.Split(t.String(), ".")
	isInside = len(parts) > 1 && strings.Contains(parts[0], "biligo")
	if len(parts) == 0 {
		typ = t.String()
	} else {
		typ = parts[len(parts)-1]
	}
	return
}

func assertType[T any]() string {
	typ, isInside := getTypeName[T]()
	if !isInside {
		panic(fmt.Errorf("type %T is not a biligo type", *(new(T))))
	}
	return typ
}

// getEmbeddedTemplate 返回对应类型的内置模板,
// 找不到时会直接 panic
func getEmbeddedTemplate[T any]() string {
	typ := assertType[T]() // 只能是 biligo 的类型

	files, err := typeTemplates.ReadDir(typeTemplatesDir) // path末尾不能有'/'
	if err != nil {
		panic(err)
	}

	filename := typ + typeTemplatesExt
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if file.Name() != filename {
			continue
		}

		data, err := typeTemplates.ReadFile(typeTemplatesDir + "/" + file.Name())
		if err != nil {
			panic(fmt.Errorf("failed to read template file %s: %v", file.Name(), err))
		}
		if len(data) == 0 {
			panic(fmt.Errorf("template file %s is empty", file.Name()))
		}

		return toString(data)
	}

	panic(fmt.Errorf("template file not found for type %s", typ))
}

type TemplateConfig struct {
	Template string
	Funcs    template.FuncMap
}

// IsTemplateOk 仅检查对应类型是否调用过 [SetTemplateFor] 且成功执行
func IsTemplateOk[T any]() bool {
	name, _ := getTypeName[T]()
	_, ok := loadedTemplates[name]
	return ok
}

// SetTemplateFor 为指定类型设置模板,
// 内置类型在不指定模板时使用内置模板,
// 外部类型在不指定模板时会 panic
func SetTemplateFor[T any](conf ...TemplateConfig) (err error) {
	typ, isInside := getTypeName[T]()

	var tmpl string
	var funcs template.FuncMap
	if len(conf) > 0 {
		tmpl = conf[0].Template
		funcs = conf[0].Funcs
	}

	if funcs != nil {
		typeTemplate.Funcs(funcs)
	}
	if tmpl == "" {
		if !isInside {
			panic(fmt.Errorf("template for type %s not found", typ))
		} else {
			_, err = typeTemplate.Parse(getEmbeddedTemplate[T]())
			if err != nil {
				panic(fmt.Errorf("failed to parse template for type %s: %v", typ, err))
			}
		}
	} else {
		_, err = typeTemplate.Parse(tmpl)
		if err != nil {
			return
		}
	}

	loadedTemplates[typ] = struct{}{}
	return nil
}

// Templatable 用于明确声明类型支持模板格式化
type Templatable interface {
	DoTemplate() string
}

// DoTemplate 填充时会根据类型选择模板,
// 外部类型在使用前需要通过 [SetTemplateFor] 函数自定义模板
func DoTemplate[T any](v *T) string {
	if v == nil {
		return fmt.Sprint(v)
	}
	typ, _ := getTypeName[T]()
	_, ok := loadedTemplates[typ]
	if !ok {
		SetTemplateFor[T]()
	}
	sb := strings.Builder{}
	err := typeTemplate.ExecuteTemplate(&sb, typ, v)
	if err != nil {
		return err.Error()
	}
	return sb.String()
}
