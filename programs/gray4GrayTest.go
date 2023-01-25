package main

import (
	"fmt"
	"os"

	"github.com/mcules/waveshare-epd/epd_config"
	"github.com/mcules/waveshare-epd/imageutil"

	"github.com/mcules/waveshare-epd"
)

func main() {

	e := epd.Epd{
		Config: epd_config.EpdConfig{},
	}
	e.Setup_4Gray()
	e.Clear()

	img, err := imageutil.OpenImage("../imageutil/test/test.jpg")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	e.Display_4Gray(&img)

	e.Sleep()

}
