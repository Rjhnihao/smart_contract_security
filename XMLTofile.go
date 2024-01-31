package main

import (
	"encoding/xml"
	"github.com/antlr4-go/antlr/v4"
	"io/ioutil"
	"log"
	"smart_contract_security/parser"
	"strings"
)

type Element struct {
	XMLName xml.Name
	Content []byte    `xml:",innerxml"`
	Nodes   []Element `xml:",any"`
}

func ConvertToXMLTofile(data []byte) {

	input := antlr.NewInputStream(string(data))
	lexer := parser.NewSolidityLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := parser.NewSolidityParser(stream)
	tree := p.SourceUnit()

	xmlElement := ConvertToXML(tree, p)
	output, err := xml.MarshalIndent(xmlElement, "", "  ")
	if err != nil {
		log.Fatalf("在进行XML编组时出错: %s", err)
		return
	}

	xmlOutput := string(output)

	err = ioutil.WriteFile("output.xml", []byte(xmlOutput), 0644)
	if err != nil {
		log.Fatalf("写入XML文件时出错: %s", err)
		return
	}
}

func ConvertToXML(node antlr.Tree, parser *parser.SolidityParser) Element {
	switch n := node.(type) {
	case antlr.RuleNode:
		ruleContext := n.GetRuleContext()
		ruleIndex := ruleContext.GetRuleIndex()
		nodeName := parser.RuleNames[ruleIndex]
		children := n.GetChildren()
		childElements := make([]Element, 0, len(children))

		for _, child := range children {
			childElements = append(childElements, ConvertToXML(child, parser))
		}

		return Element{
			XMLName: xml.Name{Local: nodeName},
			Nodes:   childElements,
		}
	case antlr.TerminalNode:
		token := n.GetSymbol()
		tokenText := escapeXMLToken(token.GetText())
		return Element{
			XMLName: xml.Name{Local: "token"},
			Content: []byte(tokenText),
		}
	default:
		return Element{}
	}
}

func escapeXMLToken(token string) string {
	token = strings.ReplaceAll(token, "&", "&amp;")
	token = strings.ReplaceAll(token, "<", "&lt;")
	token = strings.ReplaceAll(token, ">", "&gt;")
	token = strings.ReplaceAll(token, "\"", "&quot;")
	token = strings.ReplaceAll(token, "'", "&apos;")
	return token
}
