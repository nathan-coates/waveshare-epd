package main

import (
	"github.com/mcules/waveshare-epd"
	"github.com/mcules/waveshare-epd/imageutil"

	"fmt"
	"os"

	"github.com/mcules/waveshare-epd/epd_config"
)

func main() {

	e := epd.Epd{
		Config: epd_config.EpdConfig{},
	}
	e.Setup_4Gray()
	e.Clear()

	img, err := imageutil.OpenImage("../imageutil/test/test_portrait.jpg")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	e.Display_4Gray(&img)

	e.Sleep()

}
