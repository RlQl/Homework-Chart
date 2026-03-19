package edf

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
)

// EDF文件结构体
type EDFDocument struct {
	Header  *Header
	Signals []SignalHeader
	File    *os.File
}

// 前 256 字节
type Header struct {
	Version     string
	PatientID   string
	RecordingID string
	StartDate   string
	StartTime   string
	HeaderBytes int
	Reserved    string
	NumRecords  int
	DurationSec float64
	NumSignals  int
}

//通道头部

type SignalHeader struct {
	Label         string // 通道名称
	PhysicalDim   string // 物理单位
	PhysicalMin   float64
	PhysicalMax   float64
	DigitalMin    int
	DigitalMax    int
	SamplesPerRec int
}

func OpenEDF(filePath string) (*EDFDocument, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	//先读取Header
	headerData := make([]byte, 256)
	if _, err := io.ReadFull(file, headerData); err != nil {
		return nil, errors.New("无法读取完整EDF头部" + err.Error())
	}
	header := &Header{}
	header.Version = parseString(headerData[0:8])
	header.PatientID = parseString(headerData[8:88])
	header.RecordingID = parseString(headerData[88:168])
	header.StartDate = parseString(headerData[168:176])
	header.StartTime = parseString(headerData[176:184])
	header.HeaderBytes = parseInt(headerData[184:192])
	header.Reserved = parseString(headerData[192:236])
	header.NumRecords = parseInt(headerData[236:244])
	header.DurationSec = parseFloat(headerData[244:252])
	header.NumSignals = parseInt(headerData[252:256])

	//读取通道头
	signals, err := readSignalHeaders(file, header.NumSignals)
	if err != nil {
		file.Close()
		return nil, errors.New("无法读取通道头部: " + err.Error())
	}
	return &EDFDocument{
		Header:  header,
		Signals: signals,
		File:    file,
	}, nil
}

// 读取通道头的函数,将平铺的数据格式放入结构体中
func readSignalHeaders(file *os.File, numSignals int) ([]SignalHeader, error) {
	signals := make([]SignalHeader, numSignals)

	labels := make([]byte, 16*numSignals)
	transducers := make([]byte, 80*numSignals)
	phyDims := make([]byte, 8*numSignals)
	phyMins := make([]byte, 8*numSignals)
	phyMaxs := make([]byte, 8*numSignals)
	digMins := make([]byte, 8*numSignals)
	digMaxs := make([]byte, 8*numSignals)
	prefilters := make([]byte, 80*numSignals)
	samples := make([]byte, 8*numSignals)
	reserved := make([]byte, 32*numSignals)

	io.ReadFull(file, labels)
	io.ReadFull(file, transducers)
	io.ReadFull(file, phyDims)
	io.ReadFull(file, phyMins)
	io.ReadFull(file, phyMaxs)
	io.ReadFull(file, digMins)
	io.ReadFull(file, digMaxs)
	io.ReadFull(file, prefilters)
	io.ReadFull(file, samples)
	io.ReadFull(file, reserved)

	//写入结构体
	for i := 0; i < numSignals; i++ {
		signals[i].Label = parseString(labels[i*16 : (i+1)*16])
		signals[i].PhysicalDim = parseString(phyDims[i*8 : (i+1)*8])
		signals[i].PhysicalMin = parseFloat(phyMins[i*8 : (i+1)*8])
		signals[i].PhysicalMax = parseFloat(phyMaxs[i*8 : (i+1)*8])
		signals[i].DigitalMin = parseInt(digMins[i*8 : (i+1)*8])
		signals[i].DigitalMax = parseInt(digMaxs[i*8 : (i+1)*8])
		signals[i].SamplesPerRec = parseInt(samples[i*8 : (i+1)*8])

	}
	return signals, nil
}

// 辅助函数,处理ASCII字节到Go类别转换
func parseString(b []byte) string {
	return strings.TrimSpace(string(bytes.Trim(b, "\x00")))
}

func parseInt(b []byte) int {
	val, _ := strconv.Atoi(strings.TrimSpace(string(b)))
	return val
}
func parseFloat(b []byte) float64 {
	val, _ := strconv.ParseFloat(strings.TrimSpace(string(b)), 64)
	return val
}
