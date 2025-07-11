package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
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
	serviceRegex := regexp.MustCompile(`service\s+(\w+)\s*\{\s*([\s\S]*?)\s*\}`)
	serviceMatches := serviceRegex.FindAllStringSubmatch(content, -1)

	for _, match := range serviceMatches {
		if len(match) < 3 {
			continue
		}

		serviceName := match[1]
		serviceContent := match[2]

		service := ServiceDefinition{
			Name:    serviceName,
			Methods: parseMethods(serviceContent),
		}

		apiDef.Services = append(apiDef.Services, service)
	}

	return apiDef, nil
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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

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
	fmt.Println("generate function code")
	for _, service := range apiDef.Services {
		for _, method := range service.Methods {
			// 解析路径参数
			pathParams := extractPathParams(method.Path)

			// 生成方法签名
			buf.WriteString(fmt.Sprintf("// %s 对应.api文件中的%s接口\n", method.Handler, method.Handler))
			buf.WriteString(fmt.Sprintf("func (c *%sClient) %s(", apiDef.Services[0].Name, method.Handler))

			// 添加请求参数
			if method.Request != "" {
				buf.WriteString(fmt.Sprintf("req %s", method.Request))
			}

			buf.WriteString(fmt.Sprintf(") (*%s, error) {\n", method.Response))

			// 构建URL
			buf.WriteString("	url := fmt.Sprintf(\"%s%s\", c.domain, \"")

			// 处理路径参数
			pathTemplate := method.Path
			for _, param := range pathParams {
				paramName := strings.Trim(param, ":/")
				pathTemplate = strings.Replace(pathTemplate, paramName, "%s", 1)
			}

			buf.WriteString(pathTemplate)

			// 添加路径参数值
			if len(pathParams) > 0 {
				buf.WriteString("\",")
				for i, param := range pathParams {
					paramName := strings.Trim(param, ":/")
					buf.WriteString(fmt.Sprintf("req.%s", paramName))
					if i < len(pathParams)-1 {
						buf.WriteString(", ")
					}
				}
			} else {
				buf.WriteString("\")")
			}

			// 处理查询参数
			if hasQueryParams(method.Path) {
				buf.WriteString("\n\n	// 处理查询参数")
				buf.WriteString("\n	query := url.Values{}")

				// 提取查询参数
				queryParams := extractQueryParams(method.Path)
				for _, param := range queryParams {
					paramName := strings.Split(param, "=")[0]
					buf.WriteString(fmt.Sprintf("\n	if req.%s != \"\" {", paramName))
					buf.WriteString(fmt.Sprintf("\n		query.Add(\"%s\", req.%s)", paramName, paramName))
					buf.WriteString("\n	}")
				}

				buf.WriteString("\n	if len(query) > 0 {")
				buf.WriteString("\n		url += \"?\" + query.Encode()")
				buf.WriteString("\n	}")
			}

			// 创建HTTP请求
			buf.WriteString(fmt.Sprintf("\n\n	// 创建HTTP %s请求\n", method.HTTPMethod))

			if method.HTTPMethod == "GET" || method.HTTPMethod == "DELETE" {
				buf.WriteString(fmt.Sprintf("	reqObj, err := http.NewRequest(\"%s\", url, nil)\n", method.HTTPMethod))
			} else {
				buf.WriteString("	var reqBody []byte\n")
				buf.WriteString("	if req != nil {\n")
				buf.WriteString("		reqBody, err = json.Marshal(req)\n")
				buf.WriteString("		if err != nil {\n")
				buf.WriteString("			return nil, err\n")
				buf.WriteString("		}\n")
				buf.WriteString("	}\n")
				buf.WriteString(fmt.Sprintf("	reqObj, err := http.NewRequest(\"%s\", url, bytes.NewBuffer(reqBody))\n", method.HTTPMethod))
			}

			buf.WriteString("	if err != nil {\n")
			buf.WriteString("		return nil, err\n")
			buf.WriteString("	}\n")

			// 设置请求头
			buf.WriteString("\n	// 设置请求头\n")
			buf.WriteString("	reqObj.Header.Set(\"Content-Type\", \"application/json\")\n")

			// 发送请求
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
	content, err := ioutil.ReadFile(*apiFile)
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

	err = ioutil.WriteFile(outputFile, clientCode, 0644)
	if err != nil {
		fmt.Printf("Error writing output file: %v\n", err)
		return
	}

	fmt.Printf("Client code generated successfully: %s\n", outputFile)
}
