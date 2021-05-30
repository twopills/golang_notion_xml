package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	s "go-notion-issue/secret"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Feed struct {
	XMLName xml.Name `xml:"feed"`
	Entry   []Entry  `xml:"entry"`
}

type Entry struct {
	XMLName     xml.Name    `xml:"entry"`
	Link        string      `xml:"id"`
	Author      []Author    `xml:"author"`
	Assignees   []Assignees `xml:"assignees"`
	Title       string      `xml:"title"`
	Description string      `xml:"description"`
	Milestone   string      `xml:"milestone"`
}

type Task struct {
	Author      string
	Title       string
	Description string
	Milestone   string
}

type Assignees struct {
	XMLName  xml.Name  `xml:"assignees"`
	Assignee []Assigne `xml:"assignee"`
}
type Assigne struct {
	XMLName xml.Name `xml:"assignee"`
	Name    string   `xml:"name"`
}

type Author struct {
	XMLName xml.Name `xml:"author"`
	Name    string   `xml:"name"`
}

func main() {
	openAndReadXml()

	//getNotionPage()
	//getBlockOnPage()
	//fmt.Println("Suca")
}

func openAndReadXml() {
	xmlFile, err := os.Open("./xml_documents/example-gitlab-issue.xml")

	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println("Successfully Opened users.xml")

	defer xmlFile.Close()

	byteValue, _ := ioutil.ReadAll(xmlFile)

	var feed Feed

	error := xml.Unmarshal(byteValue, &feed)
	if error != nil {
		log.Panic(error)
	}
	var arr []Task

	for i := 0; i < len(feed.Entry); i++ {
		for j := 0; j < len(feed.Entry[i].Author); j++ {
			//fmt.Println(feed.Entry[i].Author[j].Name)
			var obj Task = Task{
				Author:      feed.Entry[i].Author[j].Name,
				Title:       feed.Entry[i].Title,
				Description: feed.Entry[i].Description,
				Milestone:   feed.Entry[i].Milestone,
			}
			arr = append(arr, obj)
		}
	}

	// for _, value := range feed.Entry {
	// 	for _, assignees := range value.Assignees {
	// 		for _, assignee := range assignees.Assignee {
	// 			fmt.Println(assignee.Name)

	// 		}
	// 	}
	// }

	fmt.Println(arr[0])

	addToDoFromXmlToNotionPage(arr)
}

func addToDoFromXmlToNotionPage(arr []Task) {
	url := "https://api.notion.com/v1/blocks/" + s.BlockIdPatch
	method := "PATCH"
	author := arr[0].Author
	fmt.Print(author)

	var block []string
	for _, value := range arr {

		author := value.Author
		title := value.Title
		description := value.Description
		milestone := value.Milestone

		first_block := `{
			"children": [
				{
					"object": "block",
					"type": "heading_2",
					"heading_2": {
						"text": [
							{
								"type": "text",
								"text": {
									"content": "To do:"
								},
								"annotations": {
									"bold": true,
									"italic": false,
									"strikethrough": false,
									"underline": false,
									"code": false,
									"color": "orange_background"
								}
							}
						]
					}
				},`

		block = append(block, `{
			"object": "block",
			"type": "bulleted_list_item",
			"bulleted_list_item": {
				"text": [
					{
						"type": "text",
						"text": {
							"content": "`+title+`",
							"link": null
						}
					}
				]
			}
},
{
			
			"object": "block",
			"type": "bulleted_list_item",
			"bulleted_list_item": {
				"text": [
					{
						"type": "text",
						"text": {
							"content": "`+author+`",
							"link": null
						}
					}
				]
			}
},{
			"object": "block",
			"type": "bulleted_list_item",
			"bulleted_list_item": {
				"text": [
					{
						"type": "text",
						"text": {
							"content": "`+description+`",
							"link": null
						}
					}
				]
			}
},{
			"object": "block",
			"type": "bulleted_list_item",
			"bulleted_list_item": {
				"text": [
					{
						"type": "text",
						"text": {
							"content": "`+milestone+`",
							"link": null
						}
					}
				]
			}
}]}`)

		fmt.Println(first_block + block[0])
		payload := strings.NewReader(first_block + block[0])

		client := &http.Client{}
		req, err := http.NewRequest(method, url, payload)

		if err != nil {
			fmt.Println(err)
			return
		}
		req.Header.Add("Authorization", "Bearer "+s.Token)
		req.Header.Add("Notion-Version", "2021-05-13")
		req.Header.Add("Content-Type", "application/json")

		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatalln(err)
			fmt.Println(body)
			return
		}

		fmt.Println("Tutto ok")
	}
}

func getNotionPage() {
	url := "https://api.notion.com/v1/pages/"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", "Bearer "+s.Token)
	req.Header.Add("Notion-Version", "2021-05-13")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(jsonPrettyPrint(string(body)))
}

func getBlockOnPage() {
	pagesize := "10"
	url := "https://api.notion.com/v1/blocks/" + s.BlockId + pagesize
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", "Bearer "+s.Token)
	req.Header.Add("Notion-Version", "2021-05-13")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(jsonPrettyPrint(string(body)))
}

func jsonPrettyPrint(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "\t")
	if err != nil {
		return in
	}
	return out.String()
}

func connectToNotion() {
	url := "https://api.notion.com/v1/databases/" + s.DatabaseId
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", "Bearer "+s.Token)
	req.Header.Add("Notion-Version", "2021-05-13")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(jsonPrettyPrint(string(body)))
}
