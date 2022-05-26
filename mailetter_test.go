package mailetter

import (
	"fmt"
	"testing"
)

func TestBorder (t *testing.T) {
	fmt.Println(border())
	t.Errorf("[border]")
}
