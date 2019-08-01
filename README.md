loc 是一个将ios的国际化文件(.strings)内容转换为EXCEL的小工具。

获取代码：
~~~
go get -u -v github.com/mizk/loc
~~~
编译方法:
~~~
1. cd $GOPATH/src/github.com/mizk/loc
2. export GO111MODULE=on
3. go mod tidy
4. go install
5. 在终端运行:
loc --help
~~~