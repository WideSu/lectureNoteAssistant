package app

import (
	"bytes"
	"lectureNoteAssistant/app/tool"
	"regexp"
	"strconv"
	"strings"
)

//remove stop words
func ModalWordsFilter(s string, w string) string {
	tmpText := strings.ReplaceAll(s, w, "")
	if strings.TrimSpace(tmpText) == "" || tool.CheckOnlySymbolText(strings.TrimSpace(tmpText)) {
		return ""
	} else {
		//Attempt to filter repeating particles
		compile, e := regexp.Compile(w + "{2,}")
		if e != nil {
			return s
		}
		return compile.ReplaceAllString(s, "")
	}
}

//Custom rule filtering
func DefinedWordRuleFilter(s string, rule *AppDefinedFilterRule) string {
	if rule.Way == FILTER_TYPE_STRING {
		//filter text by rules
		s = strings.ReplaceAll(s, rule.Target, rule.Replace)
	} else if rule.Way == FILTER_TYPE_REGX {
		//filter text by regex expressions
		compile, e := regexp.Compile(rule.Target)
		if e != nil {
			return s
		}
		s = compile.ReplaceAllString(s, rule.Replace)
	}
	if strings.TrimSpace(s) == "" || tool.CheckOnlySymbolText(strings.TrimSpace(s)) {
		return ""
	}
	return s
}

//Output the content into subtitle file(.srt)
func MakeSubtitleText(index int, startTime int64, endTime int64, text string, translateText string, bilingualSubtitleSwitch bool, bilingualAsc bool) string {
	var content bytes.Buffer
	content.WriteString(strconv.Itoa(index))
	content.WriteString("\r\n")
	content.WriteString(tool.SubtitleTimeMillisecond(startTime, true))
	content.WriteString(" --> ")
	content.WriteString(tool.SubtitleTimeMillisecond(endTime, true))
	content.WriteString("\r\n")

	//Output bilingual subtitles
	if bilingualSubtitleSwitch {
		if bilingualAsc {
			content.WriteString(text)
			content.WriteString("\r\n")
			content.WriteString(translateText)
		} else {
			content.WriteString(translateText)
			content.WriteString("\r\n")
			content.WriteString(text)
		}
	} else {
		content.WriteString(text)
	}

	content.WriteString("\r\n")
	content.WriteString("\r\n")
	return content.String()
}

//Output the content into text file(.txt)
func MakeText(index int, startTime int64, endTime int64, text string) string {
	var content bytes.Buffer
	content.WriteString(text)
	content.WriteString("\r\n")
	content.WriteString("\r\n")
	return content.String()
}

//Output the content into music lrc file(.lrc)
func MakeMusicLrcText(index int, startTime int64, endTime int64, text string) string {
	var content bytes.Buffer
	content.WriteString("[")
	content.WriteString(tool.MusicLrcTextMillisecond(startTime))
	content.WriteString("]")
	content.WriteString(text)
	content.WriteString("\r\n")
	return content.String()
}
