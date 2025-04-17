package utils

import (
	"net/url"
	"strings"
)

func ExtractDomain(inputURL string) string {
	// 解析输入的URL
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return inputURL
	}

	// 提取URL的Scheme（协议）和Host（域名和端口）
	domain := parsedURL.Scheme + "://" + parsedURL.Host
	return domain
}

func ExtractURLPath(inputURL string) string {
	// 解析输入的URL
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return inputURL
	}

	// 提取URL的Scheme（协议）和Host（域名和端口）
	return parsedURL.Path
}

const (
	MaskedMiddle  = "******"
	VisiblePrefix = 6
	VisibleSuffix = 6
)

func MaskToken(token string) string {
	token = strings.TrimSpace(token)

	length := len(token)

	if length <= VisiblePrefix+VisibleSuffix {
		return token
	}

	var builder strings.Builder

	builder.WriteString(token[:VisiblePrefix])

	builder.WriteString(MaskedMiddle)

	builder.WriteString(token[length-VisibleSuffix:])

	return builder.String()
}
