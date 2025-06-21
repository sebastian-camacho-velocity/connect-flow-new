package tid

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"
)

// 0100011010110011011101101001111110000000000 00000 00000000 00000000 00000000 00000000
// 2
func TestXxx(t *testing.T) {
	count := 280
	fmt.Println(strconv.FormatInt(int64(os.Getpid()), 2))

	for i := 0; i < count; i++ {
		d := time.Now().AddDate(i, 0, 0)
		id := New(d)
		_ = id

		// fmt.Printf("%s: %s %08b\n",
		// 	d.Format("2006"),
		// 	id,
		// 	id,
		// )

		fmt.Println(id, id[0])

		// b := fmt.Sprintf("%08b", id)
		// b = strings.ReplaceAll(b, " ", "")
		// b = strings.ReplaceAll(b, "[", "")
		// b = strings.ReplaceAll(b, "]", "")

		// fmt.Println(b[:43], b[43:56], b[56:])
	}

}

func BenchmarkGenIds(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = New(time.Now())
	}

}
