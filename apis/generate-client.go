package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

// APIDefinition 表示整个.api文件的定义
type APIDefinition struct {
	Types    []TypeDefinition
	Services []ServiceDefinition
}

// TypeDefinition 表示.api文件中的类型定义
type TypeDefinition struct {
	Name   string
	Fields []FieldDefinition
}

// FieldDefinition 表示类型中的字段定义
type FieldDefinition struct {
	Name    string
	Type    string
	Tag     string
	Comment string
}

// ServiceDefinition 表示.api文件中的服务定义
type ServiceDefinition struct {
	Name    string
	Methods []MethodDefinition
}

// MethodDefinition 表示服务中的方法定义
type MethodDefinition struct {
	Name       string
	HTTPMethod string
	Path       string
	Handler    string
	Request    string
	Response   string
}

// 解析器实现
func parseAPIDefinition(content string) (*APIDefinition, error) {
	apiDef := &APIDefinition{}

	// 解析类型定义
	fmt.Println("parse struct code")
	typeRegex := regexp.MustCompile(`type\s*\(\s*([\s\S]*?)\s*\)`)
	typeMatches := typeRegex.FindStringSubmatch(content)
	if len(typeMatches) > 1 {
		typeContent := typeMatches[1]
		apiDef.Types = parseTypes(typeContent)
	}

	// 解析服务定义
	fmt.Println("parse service code")
	serviceRegex := regexp.MustCompile(`service\s+([\w-]+)\s*\{\s*([\s\S]*?)\s*\}`)
	serviceMatches := serviceRegex.FindAllStringSubmatch(content, -1)

	for _, match := range serviceMatches {
		if len(match) < 3 {
			continue
		}

		serviceName := convertToCamelCase(match[1])

		serviceContent := match[2]

		service := ServiceDefinition{
			Name:    serviceName,
			Methods: parseMethods(serviceContent),
		}

		apiDef.Services = append(apiDef.Services, service)
	}

	return apiDef, nil
}

// ConvertToCamelCase 将中划线分隔的字符串转换为驼峰命名（首字母大写）
func convertToCamelCase(s string) string {
	if s == "" {
		return ""
	}

	// 按中划线分割字符串
	parts := strings.Split(s, "-")
	var result strings.Builder

	for _, part := range parts {
		if len(part) == 0 {
			continue
		}

		// 将每个部分的首字母大写，其余字符保持原样
		runes := []rune(part)
		runes[0] = unicode.ToUpper(runes[0])
		result.WriteString(string(runes))
	}

	return result.String()
}

// 解析类型定义
func parseTypes(content string) []TypeDefinition {
	var types []TypeDefinition

	// 匹配每个类型定义
	typeRegex := regexp.MustCompile(`(\w+)\s*\{\s*([\s\S]*?)\s*\}`)
	matches := typeRegex.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) < 3 {
			continue
		}

		typeName := match[1]
		fieldsContent := match[2]

		fields := parseFields(fieldsContent)

		types = append(types, TypeDefinition{
			Name:   typeName,
			Fields: fields,
		})
	}

	return types
}

// 解析字段定义
func parseFields(content string) []FieldDefinition {
	var fields []FieldDefinition

	// 匹配每个字段
	fieldRegex := regexp.MustCompile(`(\w+)\s+([\w.]+)\s*(\` + "`" + `[^` + "`" + `]*` + "`" + `)?\s*(//.*)?`)
	matches := fieldRegex.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		field := FieldDefinition{
			Name:    match[1],
			Type:    match[2],
			Comment: strings.TrimPrefix(match[4], "// "),
		}

		if len(match) > 3 {
			field.Tag = strings.Trim(match[3], "`")
		}

		fields = append(fields, field)
	}

	return fields
}

// 解析方法定义
func parseMethods(content string) []MethodDefinition {
	var methods []MethodDefinition

	// 匹配每个方法
	methodRegex := regexp.MustCompile(`@handler\s+(\w+)\s+(\w+)\s+([^\s]+)\s+\((\w+)\)\s+returns\s+\((\w+)\)`)
	matches := methodRegex.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) < 6 {
			continue
		}

		methods = append(methods, MethodDefinition{
			Name:       match[1],
			HTTPMethod: match[2],
			Path:       match[3],
			Handler:    match[1],
			Request:    match[4],
			Response:   match[5],
		})
	}

	return methods
}

// 代码生成器
func generateClientCode(apiDef *APIDefinition, packageName string) ([]byte, error) {
	var buf bytes.Buffer

	// 生成包声明
	fmt.Println("generate package code")
	buf.WriteString(fmt.Sprintf("package %s\n\n", packageName))

	// 导入必要的包
	buf.WriteString(`import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

`)

	buf.WriteString(`// _forceBytesImport 强制导入 bytes 包，避免条件编译错误
func _forceBytesImport() {
	_ = bytes.Buffer{}
}

`)

	// 生成类型定义
	fmt.Println("generate struct code")
	fmt.Printf("len apiDef.Types %d", len(apiDef.Types))
	for _, typeDef := range apiDef.Types {
		buf.WriteString(fmt.Sprintf("// %s 对应.api文件中的%s类型\n", typeDef.Name, typeDef.Name))
		buf.WriteString(fmt.Sprintf("type %s struct {\n", typeDef.Name))

		for _, field := range typeDef.Fields {
			if field.Comment != "" {
				buf.WriteString(fmt.Sprintf("	// %s\n", field.Comment))
			}

			if field.Tag != "" {
				buf.WriteString(fmt.Sprintf("	%s %s `%s`\n", field.Name, field.Type, field.Tag))
			} else {
				buf.WriteString(fmt.Sprintf("	%s %s\n", field.Name, field.Type))
			}
		}

		buf.WriteString("}\n\n")
	}

	// 生成客户端结构体
	fmt.Println("generate struct client code")
	buf.WriteString(fmt.Sprintf("// %sClient 是访问%s服务的客户端\n", apiDef.Services[0].Name, apiDef.Services[0].Name))
	buf.WriteString(fmt.Sprintf("type %sClient struct {\n", apiDef.Services[0].Name))
	buf.WriteString("	domain string\n")
	buf.WriteString("	client *http.Client\n")
	buf.WriteString("}\n\n")

	// 生成构造函数
	buf.WriteString(fmt.Sprintf("// New%s 创建一个新的%sClient实例\n", apiDef.Services[0].Name, apiDef.Services[0].Name))
	buf.WriteString(fmt.Sprintf("func New%s(domain string) *%sClient {\n", apiDef.Services[0].Name, apiDef.Services[0].Name))
	buf.WriteString(fmt.Sprintf("	return &%sClient{\n", apiDef.Services[0].Name))
	buf.WriteString("		domain: domain,\n")
	buf.WriteString("		client: &http.Client{},\n")
	buf.WriteString("	}\n")
	buf.WriteString("}\n\n")

	// 生成方法
	for _, service := range apiDef.Services {
		for _, method := range service.Methods {
			// 1. 解析当前请求结构体的参数类型（path/form/header/json）
			var (
				requestType  *TypeDefinition
				pathParams   []FieldDefinition // path:"name"
				formParams   []FieldDefinition // form:"delete,optional"
				headerParams []FieldDefinition // header:"authorization"
				jsonParams   []FieldDefinition // json:"name"
			)

			// 查找请求结构体定义
			for _, typ := range apiDef.Types {
				if typ.Name == method.Request {
					requestType = &typ
					break
				}
			}

			// 分类参数
			if requestType != nil {
				for _, field := range requestType.Fields {
					switch {
					case strings.Contains(field.Tag, `path:"`):
						pathParams = append(pathParams, field)
					case strings.Contains(field.Tag, `form:"`):
						formParams = append(formParams, field)
					case strings.Contains(field.Tag, `header:"`):
						headerParams = append(headerParams, field)
					case strings.Contains(field.Tag, `json:"`):
						jsonParams = append(jsonParams, field)
					}
				}
			}

			// 2. 生成方法签名
			buf.WriteString(fmt.Sprintf("// %s 对应.api文件中的%s接口\n", method.Handler, method.Handler))
			buf.WriteString(fmt.Sprintf("func (c *%sClient) %s(ctx context.Context, req %s) (*%s, error) {\n",
				service.Name, method.Handler, method.Request, method.Response))

			// 3. 构建URL（拼接path参数）
			buf.WriteString("	fullURL := fmt.Sprintf(\"%s%s\", c.domain, \"")
			pathTemplate := method.Path
			for _, field := range pathParams {
				// 提取path标签的参数名（如 `path:"name"` 中的 "name"）
				paramKey := extractTagValue(field.Tag, "path")
				if paramKey == "" {
					paramKey = field.Name // 默认使用字段名
				}
				// 替换路径中的占位符（如 /v1/user/:name -> /v1/user/"+req.Name）
				pathTemplate = strings.ReplaceAll(pathTemplate, ":"+paramKey, "\" + req."+field.Name+" + \"")
			}
			buf.WriteString(pathTemplate + "\")\n")

			// 4. 处理form参数（拼接为查询字符串）
			if len(formParams) > 0 {
				buf.WriteString("\n	// 处理form参数（查询字符串）\n")
				buf.WriteString("	query := url.Values{}\n")
				for _, field := range formParams {
					paramKey := extractTagValue(field.Tag, "form")
					if paramKey == "" {
						paramKey = field.Name
					}
					// 处理optional标记（可选参数）
					if strings.Contains(field.Tag, "optional") {
						buf.WriteString(fmt.Sprintf("	if req.%s != %v {\n", field.Name, getZeroValue(field.Type)))
						buf.WriteString(fmt.Sprintf("		query.Add(\"%s\", fmt.Sprintf(\"%%v\", req.%s))\n", paramKey, field.Name))
						buf.WriteString("	}\n")
					} else {
						buf.WriteString(fmt.Sprintf("	query.Add(\"%s\", fmt.Sprintf(\"%%v\", req.%s))\n", paramKey, field.Name))
					}
				}
				buf.WriteString("	if len(query) > 0 {\n")
				buf.WriteString("		fullURL += \"?\" + query.Encode()\n")
				buf.WriteString("	}\n")
			}

			// 5. 创建HTTP请求（根据方法处理请求体）
			httpMethod := strings.ToUpper(method.HTTPMethod)
			buf.WriteString(fmt.Sprintf("\n	// 创建HTTP %s请求\n", httpMethod))

			switch httpMethod {
			case "GET", "DELETE":
				// GET/DELETE 无请求体
				buf.WriteString(fmt.Sprintf("	reqObj, err := http.NewRequestWithContext(ctx, \"%s\", fullURL, nil)\n", httpMethod))
			case "POST", "PUT":
				// POST/PUT 用json参数作为请求体
				buf.WriteString("	reqBody, err := json.Marshal(req)\n")
				buf.WriteString("	if err != nil {\n")
				buf.WriteString("		return nil, err\n")
				buf.WriteString("	}\n")
				buf.WriteString(fmt.Sprintf("	reqObj, err := http.NewRequestWithContext(ctx, \"%s\", fullURL, bytes.NewBuffer(reqBody))\n", httpMethod))
			}

			// 检查请求创建错误
			buf.WriteString("	if err != nil {\n")
			buf.WriteString(fmt.Sprintf("		return nil, fmt.Errorf(\"failed to create %s request to %%s: %%w\", fullURL, err)\n", httpMethod))
			buf.WriteString("	}\n")

			// 6. 设置header参数
			if len(headerParams) > 0 {
				buf.WriteString("\n	// 设置header参数\n")
				for _, field := range headerParams {
					paramKey := extractTagValue(field.Tag, "header")
					if paramKey == "" {
						paramKey = field.Name
					}
					buf.WriteString(fmt.Sprintf("	reqObj.Header.Set(\"%s\", req.%s)\n", paramKey, field.Name))
				}
			}

			// 7. 设置Content-Type（POST/PUT默认json，GET/DELETE可选）
			if httpMethod == "POST" || httpMethod == "PUT" {
				buf.WriteString("\n	// 设置请求体类型\n")
				buf.WriteString("	reqObj.Header.Set(\"Content-Type\", \"application/json\")\n")
			}

			// 8. 发送请求并处理响应（与之前逻辑一致）
			buf.WriteString("\n	// 发送请求\n")
			buf.WriteString("	resp, err := c.client.Do(reqObj)\n")
			buf.WriteString("	if err != nil {\n")
			buf.WriteString("		return nil, err\n")
			buf.WriteString("	}\n")
			buf.WriteString("	defer resp.Body.Close()\n")

			// 读取响应内容
			buf.WriteString("\n	// 读取响应内容\n")
			buf.WriteString("	body, err := io.ReadAll(resp.Body)\n")
			buf.WriteString("	if err != nil {\n")
			buf.WriteString("		return nil, err\n")
			buf.WriteString("	}\n")

			// 检查HTTP状态码
			buf.WriteString("\n	// 检查HTTP状态码\n")
			buf.WriteString("	if resp.StatusCode != http.StatusOK {\n")
			buf.WriteString("		return nil, fmt.Errorf(\"request failed with status code: %d, body: %s\", resp.StatusCode, string(body))\n")
			buf.WriteString("	}\n")

			// 解析JSON响应
			buf.WriteString("\n	// 解析JSON响应\n")
			buf.WriteString(fmt.Sprintf("	var response %s\n", method.Response))
			buf.WriteString("	if err := json.Unmarshal(body, &response); err != nil {\n")
			buf.WriteString("		return nil, err\n")
			buf.WriteString("	}\n")

			buf.WriteString("\n	return &response, nil\n")
			buf.WriteString("}\n\n")
		}
	}

	// 格式化生成的代码
	formattedCode, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("formatting error: %v\nOriginal code:\n%s", err, buf.String())
	}

	return formattedCode, nil
}

// 提取路径中的参数
func extractPathParams(path string) []string {
	var params []string
	re := regexp.MustCompile(`:[^/?]+`)
	params = re.FindAllString(path, -1)
	return params
}

// 检查路径是否包含查询参数
func hasQueryParams(path string) bool {
	return strings.Contains(path, "?")
}

// 提取查询参数
func extractQueryParams(path string) []string {
	if !hasQueryParams(path) {
		return nil
	}

	parts := strings.Split(path, "?")
	if len(parts) < 2 {
		return nil
	}

	queryStr := parts[1]
	params := strings.Split(queryStr, "&")
	return params
}

func main() {
	// 定义命令行参数
	apiFile := flag.String("api", "", "Path to .api file")
	output := flag.String("output", "", "Output file name")
	pkgName := flag.String("package", "client", "Package name for generated code")
	flag.Parse()

	// 检查参数
	if *apiFile == "" {
		fmt.Println("Error: .api file path is required")
		flag.Usage()
		return
	}
	fmt.Println("begin read api file")
	// 读取.api文件
	content, err := os.ReadFile(*apiFile)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	// 解析.api文件
	apiDef, err := parseAPIDefinition(string(content))
	if err != nil {
		fmt.Printf("Error parsing API definition: %v\n", err)
		return
	}

	// 生成客户端代码
	clientCode, err := generateClientCode(apiDef, *pkgName)
	if err != nil {
		fmt.Printf("Error generating client code: %v\n", err)
		return
	}

	// 写入输出文件
	outputFile := *output
	if outputFile == "" {
		// 默认使用输入文件名的基本名称加上 "_client.go" 后缀
		baseName := filepath.Base(*apiFile)
		ext := filepath.Ext(baseName)
		outputFile = strings.TrimSuffix(baseName, ext) + "_client.go"
	}
	outputDir := filepath.Dir(outputFile)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}
	err = os.WriteFile(outputFile, clientCode, 0644)
	if err != nil {
		fmt.Printf("Error writing output file: %v\n", err)
		return
	}

	fmt.Printf("Client code generated successfully: %s\n", outputFile)
}

// 辅助函数：获取类型的零值（用于optional参数判断）
func getZeroValue(t string) string {
	switch t {
	case "string":
		return "\"\""
	case "bool":
		return "false"
	case "int", "int64", "float64":
		return "0"
	default:
		return "\"\""
	}
}

// 辅助函数：提取标签中的值（如 `header:"authorization"` 提取 "authorization"）
func extractTagValue(tag, key string) string {
	re := regexp.MustCompile(fmt.Sprintf(`%s:"([^",]+)`, key))
	matches := re.FindStringSubmatch(tag)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}
