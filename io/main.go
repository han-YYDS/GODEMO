package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// 20250224(DONE): CopyN
// -----------------------------------------------    TEST COPYN    ---------------------------------------------------------------------------------------------------------------------
// 将一个特定大小的数组进行拷贝
// 1. 两次Copy的结果是前后连续的
// 2. copy的参数size如果大于len, 会有EOF错误
func TestCopyN() {
	r := strings.NewReader("some io.Reader stream to be read\n")

	if _, err := io.CopyN(os.Stdout, r, 5); err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n")

	if _, err := io.CopyN(os.Stdout, r, 80); err != nil {
		log.Fatal(err)
	}

}

// -----------------------------------------------    TEST COPYN    ---------------------------------------------------------------------------------------------------------------------

// 实现一个reader来过滤掉非字母部分
type alphaReader struct {
	reader io.Reader
}

func newAlphaReader(input io.Reader) *alphaReader {
	return &alphaReader{reader: input}
}

func alpha(r byte) byte {
	if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
		return r
	}
	// return 0 即可过滤掉
	return 0
}

// reader interface
func (a *alphaReader) Read(p []byte) (int, error) {
	n, err := a.reader.Read(p)
	if err != nil {
		return 0, err
	}

	j := 0
	for i := 0; i < n; i++ {
		if char := alpha(p[i]); char != 0 {
			p[j] = char
			j++
		}
	}
	return n, nil
}

func TestAlphaReader() {
	reader := newAlphaReader(strings.NewReader("Hello World"))
	p := make([]byte, 4)
	for {
		n, err := reader.Read(p)
		if err == io.EOF {
			break
		}
		fmt.Print(string(p[:n]))
	}
	fmt.Println()
}

func TestReader() {
	// 这里如果buf在read时内不是空, 则read覆盖时不会clean掉未重合部分
	// 所以应该用 :n 来限定一下
	reader := strings.NewReader("0123456789")
	buf := make([]byte, 3)
	for {
		n, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("EOF")
				return
			} else {
				fmt.Println("ERROR: ", err)
			}
		}
		fmt.Println("Reader: ", n, " ", string(buf), string(buf[:n]))
	}
}

func TestReadFull() {
	// readall 会保留上一次读取的offset, 在下次读取时从offset开始读取
	// reader 读取到缓冲区中,是流式的
	reader := strings.NewReader("Hello World")
	buf := make([]byte, 3)
	n, err := io.ReadFull(reader, buf)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Read: ", n, "bytes ", string(buf))

	_, err = io.ReadFull(reader, buf)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Read: ", n, "bytes ", string(buf))

}

func TestSectionReader() {
	// 确保不会为了读取部分文件而加载所有文件
	// new reader
	data := strings.NewReader("0123456789")
	section := io.NewSectionReader(data, 2, 7)

	buffer := make([]byte, 5)
	for {
		n, err := section.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("ERROR: ", err)
		}
		fmt.Println("Read ", n, " ", string(buffer[:n]))
	}
}

func TestReadFile() {
	file, err := os.Open("demo")
	if err != nil {
		fmt.Println("ERROR: ", err)
		return
	}
	defer file.Close()

	buf := make([]byte, 3)
	for {
		n, err := file.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("EOF")
				return
			} else {
				fmt.Println("ERROR: ", err)
			}
		}
		fmt.Println("Reader: ", n, " ", string(buf), string(buf[:n]))
	}
}

func TestWriter() {
	// Buffer实际上支持了io.Writer接口
	var buf bytes.Buffer
	strs := []string{
		"demo",
		"demo1",
	}

	for _, str := range strs {
		buf.Write([]byte(str))
	}

	fmt.Println(buf.String())
}

type upperWriter struct {
	ch chan byte
}

func newUpperWriter() *upperWriter {
	return &upperWriter{make(chan byte, 1024)}
}

// <-chan 是只读channel的意思
func (w *upperWriter) Chan() <-chan byte {
	return w.ch
}

func (w *upperWriter) Write(b []byte) (int, error) {
	// 将byte写入channel中
	n := 0
	for _, bb := range b {
		w.ch <- bb
		n++
	}
	return n, nil
}

func (w *upperWriter) Close() error {
	close(w.ch)
	return nil
}

func TestUpperWriter() {
	writer := newUpperWriter()
	go func() {
		defer writer.Close()
		writer.Write([]byte("Hello"))
		writer.Write([]byte(" World"))
	}()

	for c := range writer.Chan() {
		fmt.Print(string(c))
	}
}

func TestStdOut() {
	proverbs := []string{
		"Channels orchestrate mutexes serialize\n",
		"Cgo is not Go\n",
		"Errors are values\n",
		"Don't panic\n",
	}
	for _, p := range proverbs {
		// _, err := os.Stdout.Write([]byte(p))
		_, err := os.Stderr.Write([]byte(p))
		if err != nil {
			fmt.Println("ERROR: ", err)
			return
		}
	}
}

func TestCopy() {
	str := "Hello World"
	bytess := []byte(str)
	buf := bytes.NewBuffer(bytess)
	file, err := os.Create("copy.txt")
	if err != nil {
		fmt.Println("ERROR: ", err)
		return
	}
	defer file.Close()

	if _, err := io.Copy(file, buf); err != nil {
		fmt.Println("ERROR: ", err)
		return
	}

	reader := make([]byte, 1024)
	n, err := file.Read(reader)
	if err != nil {
		if err != io.EOF {
			fmt.Println("ERROR: ", err)
			return
		} else {
			n -= 1
		}
	}
	fmt.Println(n)
	// fmt.Println(reader[:n])
}

func main() {
	// TestSectionReader()
	// TestReadFull()
	// TestReader()
	// TestReadFile()
	// TestAlphaReader()
	// TestWriter()
	// TestUpperWriter()
	// TestStdOut()
	// TestCopy()
	TestCopyN()
}
