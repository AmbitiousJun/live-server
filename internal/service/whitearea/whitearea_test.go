package whitearea_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/AmbitiousJun/live-server/internal/service/whitearea"
)

func TestSet(t *testing.T) {
	if err := whitearea.Init(); err != nil {
		log.Panic(err)
	}
	fmt.Printf("whitearea.Passable(\"广东佛山南海\"): %v\n", whitearea.Passable("广东佛山南海"))
	whitearea.Set("湖南")
	fmt.Printf("whitearea.Passable(\"广东佛山南海\"): %v\n", whitearea.Passable("广东佛山南海"))
	whitearea.Set("广东")
	fmt.Printf("whitearea.Passable(\"广东佛山南海\"): %v\n", whitearea.Passable("广东佛山南海"))
	whitearea.Set("广东/广州")
	fmt.Printf("whitearea.Passable(\"广东佛山南海\"): %v\n", whitearea.Passable("广东佛山南海"))
	whitearea.Set("广东/佛山")
	fmt.Printf("whitearea.Passable(\"广东佛山南海\"): %v\n", whitearea.Passable("广东佛山南海"))
	whitearea.Set("广东/佛山/顺德")
	fmt.Printf("whitearea.Passable(\"广东佛山南海\"): %v\n", whitearea.Passable("广东佛山南海"))
	whitearea.Set("广东/佛山/南海")
	fmt.Printf("whitearea.Passable(\"广东佛山南海\"): %v\n", whitearea.Passable("广东佛山南海"))
}

func TestDel(t *testing.T) {
	if err := whitearea.Init(); err != nil {
		log.Panic(err)
	}
	whitearea.Del("广东")
}
