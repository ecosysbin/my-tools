package model

import (
	"fmt"
	"regexp"
	"service/util"
	"time"
)

const (
	URL_QUEUE    = "url_queue"
	URL_VIST_SET = "url_vist_set"
)

func ProductNewsUrl(htmlurl string) {
	fmt.Println("begin product news url!")
	html := util.GetResStr(htmlurl)
	// 如何判断地址是否有效，后面考虑
	htmlstr, err := html.String()
	if err != nil {
		fmt.Println("Get html err, ", err.Error())
	}

	rule := `<a href="(http://sports.sina.com.cn/basketball/nba/.*.shtml)" target="_blank">`
	reg := regexp.MustCompile(rule)
	result := reg.FindAllStringSubmatch(htmlstr, -1)

	for i := 0; i < len(result); i++ {
		if !ISVist(URL_VIST_SET, result[i][1]) {
			PutinQueue(URL_QUEUE, result[i][1])
			AddToSet(URL_VIST_SET, result[i][1])
		}
	}
	fmt.Printf("Product url with %s finish", htmlurl)
}

func ConsumeNewsUrl() {
	// 是否需要使用信号量channel来控制消费，后续思考
	for {
		newsUrl := PopFromQueue(URL_QUEUE)
		newsstr, err := util.GetResStr(newsUrl).String()
		if err != nil {
			fmt.Printf("avidid news url, %s", newsUrl)
			continue
		}
		// 可能是不同的mysql版本要求，这个Id不传就会插入失败(最新：数据库id为主键递增就可以不传)
		title := getTitle(newsstr)
		content := getContent(newsstr)
		publishdate := getPublishDate(newsstr)
		autor := getAutor(newsstr)
		news := News{Title: title, Content: content, Publishdate: publishdate, Autor: autor}
		AddNews(&news)
		time.Sleep(time.Second)
	}
}

func getTitle(news string) string {
	fmt.Println("get title")
	rule := `<title>(.*)</title>`
	return getStringWithRule(news, rule)
}

func getContent(news string) string {
	fmt.Println("get content")
	rule := `<meta name="description" content="(.*)" />`
	return getStringWithRule(news, rule)
}

func getAutor(news string) string {
	rule := `<meta property="article:author" content="(.*)" />`
	return getStringWithRule(news, rule)
}

func getPublishDate(news string) string {
	rule := `<meta property="article:published_time" content="(.*)" />`
	timestap := getStringWithRule(news, rule)
	return timestap[0:9]
}

func getStringWithRule(str, rule string) string {
	reg := regexp.MustCompile(rule)
	result := reg.FindAllStringSubmatch(str, -1)
	return result[0][1]
}
