package web

import (
	"fmt"
	"github.com/half-coconut/gocopilot/core-engine/pkg/jsonx"
	"testing"
	"time"
)

func TestAPI(t *testing.T) {
	jsonString := `{"jsonrpc": "2.0", "method": "eth_accounts", "params": [], "id": 1}`
	//jsonString := `{"Content-Type": "application/json","User-Agent": "PostmanRuntime/7.39.0"}`

	var result map[string]interface{}

	//var result map[string]string

	result = jsonx.JsonUnmarshal(jsonString, result)

	fmt.Println(result)
	fmt.Println(result["Content-Type"]) // 访问具体字段
}

func TestTimestemp(t *testing.T) {
	// Unix 时间戳
	timestamp := int64(1743144661)

	// 将时间戳转换为时间
	ts := time.Unix(timestamp, 0)

	// 打印时间
	fmt.Println("时间:", ts)
	fmt.Println("格式化时间:", ts.Format("2006-01-02 15:04:05"))

}

func TestAIds(t *testing.T) {
	AIds := []int64{1, 3, 5, 7, 9}
	jsonx.JsonMarshal(AIds)
	fmt.Println(AIds)

	var aids []int64
	a_ids := "[1, 3, 5, 7, 9]"
	aids = jsonx.JsonUnmarshal(a_ids, aids)
	for i := range aids {
		fmt.Println(aids[i])
	}
	//fmt.Println(aids)
}

func TestDuration(t *testing.T) {
	//durationStr := "10m0s"
	durationStr := "30s"
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		// 处理错误
	}
	fmt.Println("Received duration:", duration)
}

func TestWorkers(t *testing.T) {
	w := int64(5)
	mw := int64(10)
	Workers := uint64(w)
	MaxWorkers := uint64(mw)
	fmt.Println(Workers)
	fmt.Println(MaxWorkers)
}
