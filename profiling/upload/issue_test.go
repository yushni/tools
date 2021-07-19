package upload

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f, err := os.Open("./file.txt")
		assert.NoError(b, err)

		err = NotFixed(uploader{}, f, "5f17f65591794efc048db9bea5132de5")
		assert.NoError(b, err)
	}
}
