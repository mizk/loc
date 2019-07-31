package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"log"
	"os"
	"path"
	"sort"
	"strings"
)

type Title struct {
	Keyword string `json:"keyword"`
	Value   string `json:"value"`
}

type Record struct {
	Key   string
	Value string
}

func LoadTitle(lang string) *Title {
	title := `
	{
	  "base": {
		"keyword": "关键词",
		"value": "翻译"
	  },
	  "zh_CN": {
		"keyword": "关键词",
		"value": "翻译"
	  },
	  "zh_Hans": {
		"keyword": "關鍵詞",
		"value": "翻譯"
	  },
	  "en_US": {
		"keyword": "Keyword",
		"value": "Translation"
	  },
	  "ko_KR": {
		"keyword": "키워드",
		"value": "번역하다"
	  },
	  "ja_JP": {
		"keyword": "キーワード",
		"value": "翻訳"
	  }
	}
	`
	bytes := []byte(title)
	titles := make(map[string]*Title)
	err := json.Unmarshal(bytes, &titles)
	if err != nil {
		log.Println(err)
		return nil
	}
	return titles[lang]
}
func RestoreStrings(translate map[string]string, lang, output string) error {
	outputDir := path.Join(output, strings.ToLower(lang))
	if !IsDir(outputDir) {
		err := os.MkdirAll(outputDir, 0x666)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	target := path.Join(outputDir, "Localizable.strings")
	if ok, err := PathExists(target); ok && err == nil {
		err := os.Truncate(target, 0)
		if err != nil {
			return err
		}
	}
	file, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	keys := make([]string, 0, len(translate))
	for key := range translate {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		f := fmt.Sprintf("\"%s\"=\"%s\";\n", key, translate[key])
		file.WriteString(f)
	}
	err = file.Sync()
	return err
}
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
func UpdateExcel(path, lang string, translate map[string]string) error {
	file, err := excelize.OpenFile(path)
	if err != nil {
		return err
	}
	rows := file.GetRows(lang)
	for index, row := range rows {
		if index == 0 { //忽略标题
			continue
		}
		if len(row) >= 1 {
			key := row[0]
			if t, ok := translate[key]; ok {
				axis := fmt.Sprintf("B%d", index+1)
				file.SetCellStr(lang, axis, t)
			}
		}
	}
	err = file.Save()
	return err

}
func SaveRecords(path string, lang string, records []*Record, title *Title) error {
	var file *excelize.File = nil
	if ok, err := PathExists(path); ok {
		file, err = excelize.OpenFile(path)
		if err != nil {
			return err
		}
	} else {
		file = excelize.NewFile()
	}
	var num = 0
	if _, ok := file.Sheet[lang]; ok {
		num = file.GetSheetIndex(lang)
	} else {
		num = file.NewSheet(lang)
	}
	t1 := "关键词"
	t2 := "翻译"
	if title != nil {
		t1 = title.Keyword
		t2 = title.Value
	}
	file.SetCellStr(lang, "A1", t1)
	file.SetCellStr(lang, "B1", t2)
	from := 2
	for index, record := range records {
		key := fmt.Sprintf("A%d", index+from)
		value := fmt.Sprintf("B%d", index+from)
		file.SetCellStr(lang, key, record.Key)
		file.SetCellStr(lang, value, record.Value)
	}
	file.SetActiveSheet(num)
	err := file.SaveAs(path)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
func ReadExcel(path, lang string) map[string]string {
	records := make(map[string]string)
	file, err := excelize.OpenFile(path)
	if err != nil {
		log.Println(err)
		return records
	}
	rows := file.GetRows(lang)
	for index, row := range rows {
		if len(row) >= 2 {
			key := row[0]
			value := row[1]
			if index == 0 { //第一行为标题
				continue
			}
			if len(key) > 0 {
				records[key] = value
			}
		}
	}
	return records
}
func ReadStrings(path string) []*Record {
	records := make([]*Record, 0)
	file, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return records
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, "\"") && strings.HasSuffix(text, ";") {
			//log.Println(text)
			lines := strings.Split(text, "=")
			if len(lines) == 2 {
				key := strings.Trim(strings.TrimSpace(lines[0]), "\"")
				value := strings.Trim(strings.TrimSpace(lines[1]), "\";")
				record := &Record{Key: key, Value: value}
				records = append(records, record)
			}
		}
	}
	return records
}
