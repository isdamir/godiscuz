package discuz

import (
	"fmt"
	"log"
	"net/url"
	"testing"
	"time"
)

func TestDiscuzAuthcode(t *testing.T) {
	Register("1", "asf", "url", "uc")
	en := DiscuzEncode("wqewqeqwewqewq")
	de := DiscuzDecode(en)
	if de != "wqewqeqwewqewq" {
		t.FailNow()
	}

}
