package src

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"syscall"
)

func server() {
	read := func(writer http.ResponseWriter, request *http.Request) {
		// 尝试打开文件，并处理可能出现的错误
		f, err := os.Open("./hello.txt")
		if err != nil {
			// 如果打开文件失败，返回HTTP错误响应
			http.Error(writer, "Error opening file", http.StatusInternalServerError)
			return
		}
		defer f.Close() // 确保文件在函数返回后被关闭

		buf := make([]byte, 1024)
		// 尝试读取文件，并处理可能出现的错误
		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			// 如果读取文件失败，返回HTTP错误响应
			http.Error(writer, "Error reading file", http.StatusInternalServerError)
			return
		}
		// 将读取的内容写入HTTP响应
		writer.Write(buf[:n])
	}

	mmap := func(writer http.ResponseWriter, request *http.Request) {
		f, _ := os.Open("./hello.txt")
		data, err := syscall.Mmap(int(f.Fd()), 0, 12, syscall.PROT_READ, syscall.MAP_SHARED) // mmap返回了一个data的字节数组，这个字节数组的内容就是映射了文件内容
		/*
			syscall.Mmap(int(f.Fd()), 0, 12, syscall.PROT_READ, syscall.MAP_SHARED)
			mmap涉及的参数含义:
			第一个：要映射的文件描述符。
			第二、三个：映射的范围是从0个字节到第12个字节（也就是 Hello World!）。
			第四个: 代表映射的后的内存区域是只读的，类似的参数还有 syscall.PROT_WRITE表示内存区域可以被写入，syscall.PROT_NONE表示内存区域不可访问。
			第五个：映射的内存区域可以被多个进程共享，这样一个进程修改了这个内存区域的数据，对其他进程是可见的，并且修改后的内容会自动被操作系统同步到磁盘文件里。
		*/
		if err != nil {
			panic(err)
		}
		writer.Write(data) // 返回到网卡
	}

	http.HandleFunc("/hello", read)
	http.HandleFunc("/mmap", mmap)
	http.ListenAndServe(":8080", http.DefaultServeMux)

	// 启动HTTP服务器，并处理可能出现的错误
	err := http.ListenAndServe(":8088", http.DefaultServeMux)
	if err != nil {
		// 如果启动服务器失败，打印错误信息
		panic(err)
	}
}

func sendfile() {
	// 打开源文件
	sourceFile, err := os.Open("source.txt")
	if err != nil {
		fmt.Println("Error opening source file:", err)
		return
	}
	defer sourceFile.Close()

	// 获取源文件的文件描述符
	sourceFileFd := int(sourceFile.Fd())

	// 创建目标文件，这里以另一个文件为例，实际使用中可能是套接字
	// destFile, err := os.Open("dest.txt")
	destFile, err := os.OpenFile("dest.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error creating destination file:", err)
		return
	}
	defer destFile.Close()

	// 获取目标文件的文件描述符
	destFileFd := int(destFile.Fd())
	destInfo, _ := destFile.Stat()
	destSize := destInfo.Size()
	fileInfo, _ := sourceFile.Stat()
	fileSize := fileInfo.Size()
	remain := fileSize
	fmt.Println("Source size: ", fileSize)
	fmt.Println("Dest size: ", destSize)
	// 调用 sendfile 系统调用
	// 第三个参数是传输的字节数，0 表示直到文件末尾
	// 第四个参数是传输的偏移量，nil 表示从当前文件位置开始
	// 在这里，我们使用 0 作为传输的字节数，表示传输整个文件
	var written int64 = 0
	var tail int64 = destSize
	fmt.Println("Tail: ", tail)
	for remain > 0 {
		n := remain
		syscall.Sendfile(destFileFd, sourceFileFd, &tail, int(remain))
		if n > 0 {
			written += n
			remain -= n
		}

		fmt.Println("write ", n, " total: ", written, " remain ", remain)
	}
	fmt.Println("File sent successfully")
}

func main() {
	sendfile()
}
