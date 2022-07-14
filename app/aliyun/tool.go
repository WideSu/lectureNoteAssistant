package aliyun

import (
	"lectureNoteAssistant/app/tool"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/buger/jsonparser"
)

type AliyunAudioRecognitionResultBlock struct {
	AliyunAudioRecognitionResult
	Blocks           []int
	BlockEmptyTag    bool
	BlockEmptyHandle bool
}

//Alibaba Cloud Audio Recording File Recognition - Intelligent Segmentation Processing
func AliyunAudioResultWordHandle(result []byte, callback func(vresult *AliyunAudioRecognitionResult)) {
	var audioResult = make(map[int64][]*AliyunAudioRecognitionResultBlock)
	var wordResult = make(map[int64][]*AliyunAudioWord)
	var err error

	//Get the recording recognition dataset
	_, err = jsonparser.ArrayEach(result, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		text, _ := jsonparser.GetString(value, "Text")
		channelId, _ := jsonparser.GetInt(value, "ChannelId")
		beginTime, _ := jsonparser.GetInt(value, "BeginTime")
		endTime, _ := jsonparser.GetInt(value, "EndTime")
		silenceDuration, _ := jsonparser.GetInt(value, "SilenceDuration")
		speechRate, _ := jsonparser.GetInt(value, "SpeechRate")
		emotionValue, _ := jsonparser.GetInt(value, "EmotionValue")

		vresult := &AliyunAudioRecognitionResultBlock{}
		vresult.Text = text
		vresult.ChannelId = channelId
		vresult.BeginTime = beginTime
		vresult.EndTime = endTime
		vresult.SilenceDuration = silenceDuration
		vresult.SpeechRate = speechRate
		vresult.EmotionValue = emotionValue

		_, isPresent := audioResult[channelId]
		if isPresent {
			//append on the existing audio result
			audioResult[channelId] = append(audioResult[channelId], vresult)
		} else {
			//Init
			audioResult[channelId] = []*AliyunAudioRecognitionResultBlock{}
			audioResult[channelId] = append(audioResult[channelId], vresult)
		}
	}, "Result", "Sentences")
	if err != nil {
		panic(err)
	}

	//Get word dataset
	_, err = jsonparser.ArrayEach(result, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		word, _ := jsonparser.GetString(value, "Word")
		channelId, _ := jsonparser.GetInt(value, "ChannelId")
		beginTime, _ := jsonparser.GetInt(value, "BeginTime")
		endTime, _ := jsonparser.GetInt(value, "EndTime")
		vresult := &AliyunAudioWord{
			Word:      word,
			ChannelId: channelId,
			BeginTime: beginTime,
			EndTime:   endTime,
		}
		_, isPresent := wordResult[channelId]
		if isPresent {
			//Append on previous result
			wordResult[channelId] = append(wordResult[channelId], vresult)
		} else {
			//Init
			wordResult[channelId] = []*AliyunAudioWord{}
			wordResult[channelId] = append(wordResult[channelId], vresult)
		}
	}, "Result", "Words")
	if err != nil {
		panic(err)
	}

	var symbol = []string{"？", "。", "，", "！", "；", "、", "?", ".", ",", "!"}
	//Dataset processing
	for _, value := range audioResult {
		for _, data := range value {
			// filter
			data.Text = FilterText(data.Text)

			data.Blocks = GetTextBlock(data.Text)
			data.Text = ReplaceStrs(data.Text, symbol, "")

			if len(data.Blocks) == 0 {
				data.BlockEmptyTag = true
			}
		}
	}

	//Iterate over the output
	for _, value := range wordResult {

		var block string = ""
		var blockRune int = 0
		var lastBlock int = 0

		var beginTime int64 = 0
		var blockBool = false

		var ischinese = IsChineseWords(value) //Check if it is Chinese

		var chineseNumberWordIndexs []int
		var chineseNumberDiffLength int = 0

		for i, word := range value {
			if blockBool || i == 0 {
				beginTime = word.BeginTime
				blockBool = false
			}

			if ischinese && block == "" {
				chineseNumberWordIndexs = []int{}
				chineseNumberDiffLength = 0
			}

			if ischinese {
				block += word.Word
				if tool.CheckChineseNumber(word.Word) && FindSliceIntCount(chineseNumberWordIndexs, i) == 0 {
					cl := tool.ChineseNumberToLowercaseLength(word.Word) - utf8.RuneCountInString(word.Word)
					if cl > 0 {
						chineseNumberDiffLength += cl
						chineseNumberWordIndexs = append(chineseNumberWordIndexs, i)
					} else {
						// catch expection
						if i != 0 {
							newWord := value[i-1].Word + word.Word
							cl := tool.ChineseNumberToLowercaseLength(newWord) - utf8.RuneCountInString(newWord)
							if cl > 0 {
								chineseNumberDiffLength += cl
								chineseNumberWordIndexs = append(chineseNumberWordIndexs, i)
							}
						}
					}
				}
			} else {
				block += CompleSpace(word.Word) //Complete spaces
			}

			blockRune = utf8.RuneCountInString(block)

			//fmt.Println("chineseNumberDiffLength : " , chineseNumberWordIndexs , chineseNumberDiffLength , word.Word)

			for channel, p := range audioResult {
				if word.ChannelId != channel {
					continue
				}

				for windex, w := range p {

					if (word.BeginTime >= w.BeginTime && word.EndTime <= w.EndTime) || ((word.BeginTime < w.EndTime && word.EndTime > w.EndTime) && (FindSliceIntCount(w.Blocks, -1) != len(w.Blocks))) {
						flag := false
						early := false

						if !w.BlockEmptyTag {
							for t, B := range w.Blocks {
								//fmt.Println("blockRune : " , blockRune , B , word.Word)
								if ((blockRune >= B) || (blockRune+chineseNumberDiffLength >= B)) && B != -1 {
									flag = true

									//fmt.Println(w.Blocks)
									//fmt.Println(B , lastBlock , (B - lastBlock) , word.Word)
									//fmt.Println(w.Text)
									//fmt.Println(  block )
									//fmt.Println("\n")

									var thisText = ""
									// fault tolerance mechanism
									if t == (len(w.Blocks) - 1) {
										thisText = SubString(w.Text, lastBlock, 10000)
									} else {
										// next word ends early
										if i < len(value)-1 && value[i+1].BeginTime >= w.EndTime {
											thisText = SubString(w.Text, lastBlock, 10000)
											early = true
										} else {
											thisText = SubString(w.Text, lastBlock, (B - lastBlock))
										}
									}

									lastBlock = B
									if early == true {
										//All set to -1
										for vt, vb := range w.Blocks {
											if vb != -1 {
												w.Blocks[vt] = -1
											}
										}
									} else {
										w.Blocks[t] = -1
									}

									vresult := &AliyunAudioRecognitionResult{
										Text:            thisText,
										ChannelId:       channel,
										BeginTime:       beginTime,
										EndTime:         word.EndTime,
										SilenceDuration: w.SilenceDuration,
										SpeechRate:      w.SpeechRate,
										EmotionValue:    w.EmotionValue,
									}
									callback(vresult) //callback parameter

									blockBool = true
									break
								}
							}

							//fmt.Println("word.Word : " , word.Word)
							//fmt.Println(block)

							if FindSliceIntCount(w.Blocks, -1) == len(w.Blocks) {
								//All interception completed
								block = ""
								lastBlock = 0
							}
							// Fault tolerance mechanism
							if FindSliceIntCount(w.Blocks, -1) == (len(w.Blocks)-1) && flag == false {
								var thisText = SubString(w.Text, lastBlock, 10000)

								w.Blocks[len(w.Blocks)-1] = -1
								//vresult
								vresult := &AliyunAudioRecognitionResult{
									Text:            thisText,
									ChannelId:       channel,
									BeginTime:       beginTime,
									EndTime:         w.EndTime,
									SilenceDuration: w.SilenceDuration,
									SpeechRate:      w.SpeechRate,
									EmotionValue:    w.EmotionValue,
								}

								//fmt.Println(  thisText )
								//fmt.Println(  block )
								//fmt.Println(  word.Word , beginTime, w.EndTime , flag  , word.EndTime  )

								callback(vresult) //Sending back parameters using the callback function

								// Overwrite the timestamp of the next paragraph
								if windex < (len(p) - 1) {
									beginTime = p[windex+1].BeginTime
								} else {
									beginTime = w.EndTime
								}
								//Clear the parameters
								block = ""
								lastBlock = 0
							}
						} else {
							//Clear the parameters
							block = ""
							lastBlock = 0
							blockBool = true

							if w.BlockEmptyHandle == false {
								vresult := &AliyunAudioRecognitionResult{
									Text:            w.Text,
									ChannelId:       w.ChannelId,
									BeginTime:       w.BeginTime,
									EndTime:         w.EndTime,
									SilenceDuration: w.SilenceDuration,
									SpeechRate:      w.SpeechRate,
									EmotionValue:    w.EmotionValue,
								}
								callback(vresult) //Sending back parameters using the callback function
								w.BlockEmptyHandle = true
							}

						}

					}
				}
			}
		}
	}
}

func FindSliceIntCount(slice []int, target int) int {
	c := 0
	for _, v := range slice {
		if target == v {
			c++
		}
	}
	return c
}

//Batch replace multiple keyword texts
func ReplaceStrs(strs string, olds []string, s string) string {
	for _, word := range olds {
		strs = strings.Replace(strs, word, s, -1)
	}
	return strs
}

func StringIndex(strs string, word rune) int {
	strsRune := []rune(strs)
	for i, v := range strsRune {
		if v == word {
			return i
		}
	}
	return -1
}

//Fill right spaces
func CompleSpace(s string) string {
	s = strings.TrimLeft(s, " ")
	s = strings.TrimRight(s, " ")
	return s + " "
}

func IsChineseWords(words []*AliyunAudioWord) bool {
	for _, v := range words {
		if IsChineseChar(v.Word) {
			return true
		}
	}
	return false
}

func IsChineseChar(str string) bool {
	for _, r := range str {
		if unicode.Is(unicode.Scripts["Han"], r) || (regexp.MustCompile("[\u3002\uff1b\uff0c\uff1a\u201c\u201d\uff08\uff09\u3001\uff1f\u300a\u300b]").MatchString(string(r))) {
			return true
		}
	}
	return false
}

func IndexRunes(strs string, olds []rune) int {
	min := -1
	for i, word := range olds {
		index := StringIndex(strs, word)
		//println( "ts : " , index)
		if i == 0 {
			min = index
		} else {
			if min == -1 {
				min = index
			} else {
				if index < min && index != -1 {
					min = index
				}
			}
		}
	}
	return min
}

func GetTextBlock(strs string) []int {
	var symbol_zhcn = []rune{'？', '。', '，', '！', '；', '、', '?', '.', ',', '!'}
	//var symbol_en = []rune{'?','.',',','!'}
	strsRune := []rune(strs)

	blocks := []int{}
	for {
		index := IndexRunes(strs, symbol_zhcn)
		if index == -1 {
			break
		}
		strs = string(strsRune[0:index]) + string(strsRune[(index+1):])
		strsRune = []rune(strs)
		blocks = append(blocks, index)
	}
	return blocks
}

func SubString(str string, begin int, length int) (substr string) {
	// convert string to[]rune
	rs := []rune(str)
	lth := len(rs)

	// simple out-of-bounds judgment
	if begin < 0 {
		begin = 0
	}
	if begin >= lth {
		begin = lth
	}
	end := begin + length
	if end > lth {
		end = lth
	}
	// return substring
	return string(rs[begin:end])
}

//Filter Text
func FilterText(text string) string {
	//remove newlines
	re, _ := regexp.Compile("[\n|\r|\r\n]+")
	text = re.ReplaceAllString(text, "")
	return text
}
