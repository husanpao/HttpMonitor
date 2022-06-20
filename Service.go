package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Question struct {
	Options      []string `json:"options"`
	QuesTypeStr  string   `json:"quesTypeStr"`
	Content      string   `json:"content"`
	RightOptions []string `json:"rightOptions"`
	QuesID       string   `json:"quesId"`
}

var questions []Question
var Answers map[string]Question

type Data struct {
	Result struct {
		Msg string `json:"msg,omitempty"`
	} `json:"result,omitempty"`
	Data struct {
		IsRight         bool     `json:"isRight,omitempty"`
		AnsweredOptions []string `json:"answeredOptions,omitempty"`
		Ques            struct {
			QuesNo      int      `json:"quesNo,omitempty"`
			Options     []string `json:"options,omitempty"`
			QuesTypeStr string   `json:"quesTypeStr,omitempty"`
			QuesID      string   `json:"quesId,omitempty"`
			Content     string   `json:"content,omitempty"`
			QuesType    int      `json:"quesType,omitempty"`
		} `json:"ques,omitempty"`
		RightOptions []string `json:"rightOptions,omitempty"`
	} `json:"data,omitempty"`
}

func initJson() {
	jsonFile, err := os.Open("answer.json")
	if err != nil {
		log.Printf("load json error:%s", err.Error())
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &questions)
	if err != nil {
		log.Printf("%s", err.Error())
		return
	}
	Answers = make(map[string]Question)
	for _, value := range questions {
		Answers[value.Content] = value
	}
	log.Printf("题库数量：%d", len(Answers))
}
func answer(req Req) string {
	if strings.HasSuffix(req.Url, "regist/activity") || strings.HasSuffix(req.Url, "regist/competition") {
		req.Data = strings.ReplaceAll(req.Data, "memberName\":\"", "memberName\":\"[小助手]")
	}
	if strings.HasSuffix(req.Url, "ques/startCompetition") || strings.HasSuffix(req.Url, "ques/answerQues") {
		var data Data
		err := json.Unmarshal([]byte(req.Data), &data)
		if err != nil {
			log.Printf("%s", err.Error())
			return ""
		}
		ques := data.Data.Ques
		content := ques.Content
		obj, ok := Answers[content]
		if ok {
			newopt := make([]string, len(obj.Options))
			for i, option := range ques.Options {
				for _, right := range obj.RightOptions {
					if option == right {
						option = option + "(正确)"
						break
					}
				}
				newopt[i] = option
			}
			data.Data.Ques.Options = newopt
			str, err := json.Marshal(data)
			if err != nil {
				return req.Data
			}
			return string(str)
		}
	}
	return req.Data
}
