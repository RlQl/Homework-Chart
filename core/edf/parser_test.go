package edf

import (
	"fmt"
	"testing"
)

func TestOpenEDF(t *testing.T) {
	filePath := "../../sample_edf_dataset/sample_dataset/SC4001E0-PSG.edf"

	doc, err := OpenEDF(filePath)
	if err != nil {
		t.Fatalf("解析EDF失败: %v", err)
	}
	defer doc.File.Close()
	fmt.Println("=== EDF 解析成功 ===")
	fmt.Printf("数据块总数: %d\n", doc.Header.NumRecords)
	fmt.Printf("单块时长: %.4f 秒\n", doc.Header.DurationSec)
	fmt.Println("--- 7 个通道详情 ---")

	for i, sig := range doc.Signals {
		fmt.Printf("通道%d [%s]: 采样点=%d, 物理范围=[%.2f %s ~ %.2f %s], 数字范围=[%d ~ %d]\n",
			i+1, sig.Label, sig.SamplesPerRec, sig.PhysicalMin, sig.PhysicalDim, sig.PhysicalMax, sig.PhysicalDim, sig.DigitalMin, sig.DigitalMax)
	}
	fmt.Println("====================")
}
