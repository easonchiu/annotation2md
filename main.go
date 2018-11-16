package main

import (
  "annotation2md/engine"
  "annotation2md/feather"
  "flag"
  "io/ioutil"
  "os"
)

func main() {
  // 参数： --dir=目录 --title=文档标题 --outfile=abc
  dirname := ""
  flag.StringVar(&dirname, "dir", "", "解析目录")
  title := ""
  flag.StringVar(&title, "title", "", "文档标题")
  outfile := ""
  flag.StringVar(&outfile, "outfile", "doc", "导出的文件名")
  flag.Parse()

  filenames := feather.GetFileNamesFromDir(dirname)
  markdown := engine.Start(title, filenames)

  ioutil.WriteFile(outfile + ".md", []byte(markdown), os.ModePerm)
}
