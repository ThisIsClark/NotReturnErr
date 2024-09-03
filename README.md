# NotReturnErr

NotReturnErr is a program for checking for not retuning error.


## Use

Modify this line to set the dir to parse:
```
err := filepath.Walk("E:\\code\\liuyu\\data\\ingestor", func(path string, info os.FileInfo, err error) error {
```
Run the main

ChangeLog
* v1.0 实现基本功能
* v1.1 实现对源码目录下vendor目录的排除
* v1.2 实现对返回值不包含error的函数的排除
* v1.3 实现对返回值个数为0的函数的排除