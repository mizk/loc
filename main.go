package main

import (
	"github.com/spf13/cobra"
	"loc/utils"
	"log"
	"os"
	"strings"
)

const PathSeparator = string(os.PathSeparator)

func main() {

	root := &cobra.Command{}
	addInitCommand(root)
	addPatchCommand(root)
	addRestoreCommand(root)
	root.Execute()
}

//restore
func addRestoreCommand(root *cobra.Command) {
	command := &cobra.Command{
		Use:   "restore {--lang=lang} {excel} {strings}",
		Short: "restore {--lang=lang} {excel} {strings} 使用EXCEL文件生成strings文件",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				cmd.Usage()
				return
			}
			input := args[0]
			if ok, err := utils.PathExists(input); err != nil || ok == false {
				log.Println("翻译文件不存在")
			}
			output := args[1]
			language := parseLangFlag(cmd)
			if len(language) == 0 {
				return
			}

			if !strings.HasSuffix(input, ".xlsx") {
				log.Println("翻译文件必须是xlsx类型")
				return
			}
			translate := utils.ReadExcel(input, language)
			err := utils.RestoreStrings(translate, language, output)
			if err != nil {
				log.Println(err)
			}

		},
	}
	command.PersistentFlags().String("lang", "zh_CN", "语言类型,lang只能是如下值:zh_CN,zh_Hans,en_US,ja_JP,ko_KR")
	root.AddCommand(command)
}

//patch
func addPatchCommand(root *cobra.Command) {
	command := &cobra.Command{
		Use:   "patch {--lang=lang} {patch} {translate}",
		Short: "patch {--lang=lang} {patch} {translate}  使用strings或现有的Excel文件更s新Excel",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				cmd.Usage()
				return
			}
			patch := args[0]
			if ok, err := utils.PathExists(patch); err != nil || ok == false {
				log.Println("缺少patch文件")
				return
			}
			input := args[1]
			if ok, err := utils.PathExists(input); err != nil || ok == false {
				log.Println("缺少翻译文件")
				return
			}
			language := parseLangFlag(cmd)
			if len(language) == 0 {
				return
			}
			translate := make(map[string]string)
			if strings.HasSuffix(patch, ".xlsx") {
				excel := utils.ReadExcel(patch, language)
				for k, v := range excel {
					translate[k] = v
				}
			} else if strings.HasSuffix(patch, ".strings") {
				r := utils.ReadStrings(patch)
				for _, record := range r {
					translate[record.Key] = record.Value
				}
			} else {
				log.Println("patch文件必须是xlsx或strings类型")
				return
			}
			rd := utils.ReadExcel(input, language)
			for key, _ := range rd {
				if value, ok := translate[key]; ok {
					rd[key] = value
				}
			}
			err := utils.UpdateExcel(input, language, rd)
			if err != nil {
				log.Println(err)
			}
		},
	}
	command.PersistentFlags().String("lang", "zh_CN", "语言类型,lang只能是如下值:zh_CN,zh_Hans,en_US,ja_JP,ko_KR")
	root.AddCommand(command)
}

//
func addInitCommand(root *cobra.Command) {
	command := &cobra.Command{
		Use:   "init {--lang=lang} {strings} {excel}",
		Short: "init {--lang=lang} {strings} {excel} 使用strings文件生成EXCEL文件",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				cmd.Usage()
				return
			}
			input := args[0]
			if ok, err := utils.PathExists(input); err != nil || ok == false {
				log.Println("缺少strings文件")
				return
			}
			records := utils.ReadStrings(input)
			language := parseLangFlag(cmd)
			if len(language) == 0 {
				return
			}
			output := args[1]
			title := utils.LoadTitle(language)
			utils.SaveRecords(output, language, records, title)
		},
	}
	command.PersistentFlags().String("lang", "zh_CN", "语言类型,lang只能是如下值:zh_CN,zh_Hans,en_US,ja_JP,ko_KR")
	root.AddCommand(command)
}

func parseLangFlag(command *cobra.Command) string {
	language := ""
	lang := command.Flag("lang")
	if lang == nil {
		log.Println("不受支持的语言,lang只能是如下值:zh_CN,zh_Hans,en_US,ja_JP,ko_KR")
		return ""
	}
	language = lang.Value.String()
	valid := language == "zh_CN" || language == "en_US" || language == "ja_JP" || language == "zh_Hans" || language == "ko_KR"
	if !valid {
		log.Println("不受支持的语言,lang只能是如下值:zh_CN,zh_Hans,en_US,ja_JP,ko_KR")
		return ""
	}
	return language
}
