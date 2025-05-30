package main

import (
	"bytes"
	"io"
	"log"
	"os"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
)

func test1() {
	logrus.SetReportCaller(true) // 设置log显示定位文件和行号

	logrus.SetLevel(logrus.TraceLevel)

	requestLogger := logrus.WithFields(logrus.Fields{
		"user_id": 10010,
		"ip":      "192.168.32.15",
	})

	requestLogger.Info("info msg")
	requestLogger.Error("error msg")
}

func test2() {
	writer1 := &bytes.Buffer{}
	writer2 := os.Stdout
	writer3, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Fatalf("create file log.txt failed: %v", err)
	}
	// 设置多个output
	logrus.SetOutput(io.MultiWriter(writer1, writer2, writer3))
	logrus.Info("info msg")
}

func test3() {
	log := logrus.New()
	log.SetLevel(logrus.InfoLevel)
	log.SetFormatter(&logrus.JSONFormatter{})
	// log条目以json形式
	log.Info("info msg")
}

func test4() {
	logrus.SetFormatter(&nested.Formatter{
		HideKeys:    true,
		FieldsOrder: []string{"component", "category"},
	})
	logrus.Info("info msg")
}

func main() {
	test4()
}
