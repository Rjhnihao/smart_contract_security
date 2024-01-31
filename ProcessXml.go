// ProcessXml.go
package main

import (
	"github.com/antchfx/xmlquery"
	"io/ioutil"
	"log"
	"strings"
)

func ProcessXml() (string, error) {

	rulesData, err := ioutil.ReadFile("solidity-rules.xml")
	if err != nil {
		return "", err
	}

	doc, err := xmlquery.Parse(strings.NewReader(string(rulesData)))
	if err != nil {
		return "", err
	}

	patternNodes := xmlquery.Find(doc, "//Pattern")
	var results []string

	for _, patternNode := range patternNodes {
		patternId := patternNode.SelectAttr("patternId")
		xpathNode := xmlquery.FindOne(patternNode, "./XPath/text()")
		xpathString := xpathNode.InnerText()
		xpathString = strings.TrimSpace(strings.ReplaceAll(xpathString, "\n", " "))

		data, err := ioutil.ReadFile("output.xml")
		if err != nil {
			log.Printf("读取output.xml时出错: %s\n", err)
			continue
		}

		dataDoc, err := xmlquery.Parse(strings.NewReader(string(data)))
		if err != nil {
			log.Printf("解析output.xml时出错: %s\n", err)
			continue
		}

		dataNodes := xmlquery.Find(dataDoc, xpathString)

		if len(dataNodes) > 0 {
			results = append(results, "匹配到漏洞："+patternId)
		}
	}

	return strings.Join(results, "\n"), nil
}
