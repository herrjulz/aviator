# MinGOak

A lightweight, easy-to-use, in-memory file tree. 

--> `mingoak` implements [os.FilePath](https://golang.org/pkg/os/#FileInfo)

```
$ go get github.com/JulzDiverse/mingoak
```

```
import github.com/JulzDiverse/mingoak  
```

## Usage

```go

  root := mingoak.MkRoot()

  root.MkDirAll("path/to/dir/")
  root.WriteFile("path/to/dir/file", []byte("test"))

  fileInfo, _ := root.ReadDir("path/to/dir")
  for _, v := fileInfo {
     pintln(v.IsDir()) //true or false
     println(v.Name())  //name of file/dir
     println(v.ModTime()) //modification time
     println(v.Size()) //file size
  }

  file, _ = root.ReadFile("path/to/dir/file")
  
  //Walk also works:
  files, _ := root.Walk("path")
  for _, v := files {
     fmt.Prinln(v) //prints the file path
  }
```


