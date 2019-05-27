package main

import (
  "annotation2md/engine"
  "annotation2md/feather"
  "flag"
  "io/ioutil"
  "os"
)

// TODO
/*
1. 变量，在某个文件中定义变量可以引用(done)
例:
定义: $skip | int | 跳过
使用: ${skip} 时，将会编译成 skip | int | 跳过
     ${skip}? 时，将会编译成 skip? | int | 跳过
     list.${skip}? 时，将会编译成 list.skip? | int | 跳过

2. 关联，可以把另一个id下的@xxx的内容关联过来
例:
定义: @Router=2.1.1@Router 就是取Id为2.1.1这条文档的Router数据
*/

func main() {
  // 参数： --dir=目录 --title=文档标题 --outfile=导出的文件名 --vars=声明变量的文档 --json
  dirname := ""
  flag.StringVar(&dirname, "dir", "", "解析目录")
  title := ""
  flag.StringVar(&title, "title", "", "文档标题")
  outfile := "doc"
  flag.StringVar(&outfile, "outfile", "doc", "导出的文件名")
  vars := ""
  flag.StringVar(&vars, "vars", "", "声明变量的文档")
  jsonFile := false
  flag.BoolVar(&jsonFile, "json", false, "生成json格式的文档")
  flag.Parse()

  // dirname = "test"
  // title = "title"
  // vars = "test/.docvars"
  // outfile = "test/doc"
  // jsonFile = true

  filenames := feather.GetFileNamesFromDir(dirname)
  keyvars := feather.GetKeyVarsFromFile(vars)
  markdown, json := engine.Start(title, filenames, keyvars)

  _= ioutil.WriteFile(outfile+".md", []byte(markdown), os.ModePerm)
  if jsonFile {
    _= ioutil.WriteFile(outfile+".api.json", []byte(json), os.ModePerm)
  }
}
