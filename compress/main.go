package main

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/pierrec/lz4"
)

var (
	origData = []byte("要压缩的数据内容")
	gzname   = "data.gz"
	oneKB    = make([]byte, 1024)
	oneMB    = make([]byte, 1024*1024)
)

func GenRandomData(buffer []byte) { // 纯随机数组
	// 使用 crypto/rand 生成随机字节
	_, err := rand.Read(buffer)
	if err != nil {
		log.Fatal(err) // 如果生成随机字节时发生错误，终止程序
	}
}

func GenRandomUser() []byte {
	users := make([]User, 1024)

	for i := range users {
		users[i] = generateUser(i)
	}

	jsonData, err := json.Marshal(users)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v\n", err)
	}
	return jsonData
}

type User struct {
	UserID    int    `json:"user_id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Status    string `json:"status"`
	LastLogin string `json:"last_login"`
}

// 随机生成用户名
func randomUsername(index int) string {
	return fmt.Sprintf("user%d", index)
}

// 随机生成状态
func randomStatus() string {
	statuses := []string{"active", "inactive", "suspended"}
	randomint, _ := rand.Int(rand.Reader, big.NewInt(3))
	return statuses[randomint.Int64()]
}

// 生成一个用户数据
func generateUser(index int) User {
	return User{
		UserID:    index,
		Username:  randomUsername(index),
		Email:     fmt.Sprintf("user%d@example.com", index),
		Status:    randomStatus(),
		LastLogin: time.Now().Add(-time.Duration(index) * time.Hour).Format(time.RFC3339),
	}
}

func TestGZIP() {
	// 创建一个用于写入的文件
	gzipFile, err := os.Create(gzname)
	if err != nil {
		panic(err)
	}
	defer gzipFile.Close()

	// 创建gzip.Writer，将文件作为底层写入器
	gw, err := gzip.NewWriterLevel(gzipFile, gzip.BestCompression) // with level
	// gw := gzip.NewWriter(gzipFile) // not with level
	defer gw.Close()

	// 写入数据
	_, err = gw.Write(origData)
	if err != nil {
		panic(err)
	}

	// 确保所有数据都已写入并压缩
	err = gw.Flush()
	if err != nil {
		panic(err)
	}
}

func TestUngzip() {
	// 打开gzip文件
	gzipFile, err := os.Open(gzname)
	if err != nil {
		panic(err)
	}
	defer gzipFile.Close()

	// 创建gzip.Reader
	gr, err := gzip.NewReader(gzipFile)
	if err != nil {
		panic(err)
	}
	defer gr.Close()

	// 读取解压后的数据
	decompressedData, err := ioutil.ReadAll(gr)
	if err != nil {
		panic(err)
	}

	// 打印解压后的数据
	println(string(decompressedData))
}

func main() {
	// GenRandomData(oneMB)
	data := GenRandomUser()
	// TestGZIP()
	// TestGZIPLevel()
	// TestUngzip()
	BenchMarkGZip(data)
	BenchMarkLz4(data)
	BenchMarkZlib(data)
}

// --------------------------- GZIP -------------------------------------//
// 压缩并返回压缩后的数据和压缩时间
func compressGzip(data []byte) ([]byte, time.Duration) {
	start := time.Now()
	var compressedData bytes.Buffer
	gzipWriter := gzip.NewWriter(&compressedData)
	_, err := gzipWriter.Write(data)
	if err != nil {
		log.Fatal(err)
	}
	gzipWriter.Close()
	duration := time.Since(start)
	return compressedData.Bytes(), duration
}

// 使用 gzip 压缩数据，并指定压缩等级
func compressGzipLevel(data []byte, level int) ([]byte, error) {
	var buf bytes.Buffer
	gzipWriter, err := gzip.NewWriterLevel(&buf, level)
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip writer: %v", err)
	}
	defer gzipWriter.Close()

	_, err = gzipWriter.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed to write data to gzip: %v", err)
	}

	return buf.Bytes(), nil
}

// 使用 gzip 解压数据
func decompressGzip(data []byte) ([]byte, time.Duration, error) {
	start := time.Now()
	gzipReader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create gzip reader: %v", err)
	}
	defer gzipReader.Close()
	ret, err := io.ReadAll(gzipReader)
	if err != nil {
		return nil, 0, err
	}
	duration := time.Since(start)
	return ret, duration, nil
}

// --------------------------- ZLIB ------------------------------------//
func compressZlib(data []byte) ([]byte, time.Duration) {
	start := time.Now()
	var compressedData bytes.Buffer
	zlibWriter := zlib.NewWriter(&compressedData)
	_, err := zlibWriter.Write(data)
	if err != nil {
		log.Fatal(err)
	}
	zlibWriter.Close()
	duration := time.Since(start)
	return compressedData.Bytes(), duration
}

func decompressZlib(data []byte) ([]byte, time.Duration, error) {
	start := time.Now()
	r, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create gzip reader: %v", err)
	}
	defer r.Close()
	duration := time.Since(start)
	ret, err := io.ReadAll(r)
	if err != nil {
		return nil, 0, err
	}
	return ret, duration, nil
}

// ---------------------------- LZ4 ------------------------------------//
func compressLz4(data []byte) ([]byte, time.Duration) {
	start := time.Now()
	var compressedData bytes.Buffer
	lz4Writer := lz4.NewWriter(&compressedData)
	_, err := lz4Writer.Write(data)
	if err != nil {
		log.Fatal(err)
	}
	lz4Writer.Close()
	duration := time.Since(start)
	return compressedData.Bytes(), duration
}

func decompressLZ4(data []byte) ([]byte, time.Duration, error) {
	start := time.Now()
	lz4Reader := lz4.NewReader(bytes.NewReader(data))
	data, err := io.ReadAll(lz4Reader)
	return data, time.Since(start), err
}

func BenchMarkGZip(data []byte) {
	// 测试 Gzip 压缩
	gzipData, gzipTime := compressGzip(data)

	_, gzipDeTime, _ := decompressGzip(gzipData)
	fmt.Printf("Gzip: CompressTime=%v, DeCompressTime=%v,  Compressed Size=%d bytes, Compression Ratio=%.4f\n",
		gzipTime, gzipDeTime, len(gzipData), float64(len(gzipData))/float64(len(data)))
}

func BenchMarkZlib(data []byte) {

	// 测试 Zlib 压缩
	zlibData, zlibTime := compressZlib(data)
	_, zlibDeTime, _ := decompressZlib(zlibData)
	fmt.Printf("Zlib: CompressTime=%v, DeCompressTime=%v, Compressed Size=%d bytes, Compression Ratio=%.4f\n",
		zlibTime, zlibDeTime, len(zlibData), float64(len(zlibData))/float64(len(data)))

}

func BenchMarkLz4(data []byte) {
	// 测试 LZ4 压缩
	compressData, lz4Time := compressLz4(data)
	_, deTime, _ := decompressZlib(compressData)
	fmt.Printf("LZ4: CompressTime=%v, DeCompressTime=%v, Compressed Size=%d bytes, Compression Ratio=%.4f\n",
		lz4Time, deTime, len(compressData), float64(len(compressData))/float64(len(data)))
}
