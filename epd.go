package epd

import (
	"epd/epd_config"
	"epd/imageutil"
	"image"
	"image/color"
	"log"
	"time"

	"periph.io/x/conn/v3/gpio"
)

const (
	EPD_WIDTH int = 176
	EPD_HEIGHT int = 264
)

var (
	lut_vcom_dc = []byte{
		0x00, 0x00,
		0x00, 0x1A, 0x1A, 0x00, 0x00, 0x01,
		0x00, 0x0A, 0x0A, 0x00, 0x00, 0x08,
		0x00, 0x0E, 0x01, 0x0E, 0x01, 0x10,
		0x00, 0x0A, 0x0A, 0x00, 0x00, 0x08,
		0x00, 0x04, 0x10, 0x00, 0x00, 0x05,
		0x00, 0x03, 0x0E, 0x00, 0x00, 0x0A,
		0x00, 0x23, 0x00, 0x00, 0x00, 0x01,
	}

	lut_ww = []byte{
		0x90, 0x1A, 0x1A, 0x00, 0x00, 0x01,
		0x40, 0x0A, 0x0A, 0x00, 0x00, 0x08,
		0x84, 0x0E, 0x01, 0x0E, 0x01, 0x10,
		0x80, 0x0A, 0x0A, 0x00, 0x00, 0x08,
		0x00, 0x04, 0x10, 0x00, 0x00, 0x05,
		0x00, 0x03, 0x0E, 0x00, 0x00, 0x0A,
		0x00, 0x23, 0x00, 0x00, 0x00, 0x01,
	}

	// R22H    r
	lut_bw = []byte{
		0xA0, 0x1A, 0x1A, 0x00, 0x00, 0x01,
		0x00, 0x0A, 0x0A, 0x00, 0x00, 0x08,
		0x84, 0x0E, 0x01, 0x0E, 0x01, 0x10,
		0x90, 0x0A, 0x0A, 0x00, 0x00, 0x08,
		0xB0, 0x04, 0x10, 0x00, 0x00, 0x05,
		0xB0, 0x03, 0x0E, 0x00, 0x00, 0x0A,
		0xC0, 0x23, 0x00, 0x00, 0x00, 0x01,
	}

	// R23H    w
	lut_bb = []byte{
		0x90, 0x1A, 0x1A, 0x00, 0x00, 0x01,
		0x40, 0x0A, 0x0A, 0x00, 0x00, 0x08,
		0x84, 0x0E, 0x01, 0x0E, 0x01, 0x10,
		0x80, 0x0A, 0x0A, 0x00, 0x00, 0x08,
		0x00, 0x04, 0x10, 0x00, 0x00, 0x05,
		0x00, 0x03, 0x0E, 0x00, 0x00, 0x0A,
		0x00, 0x23, 0x00, 0x00, 0x00, 0x01,
	}
	// R24H    b
	lut_wb = []byte{
		0x90, 0x1A, 0x1A, 0x00, 0x00, 0x01,
		0x20, 0x0A, 0x0A, 0x00, 0x00, 0x08,
		0x84, 0x0E, 0x01, 0x0E, 0x01, 0x10,
		0x10, 0x0A, 0x0A, 0x00, 0x00, 0x08,
		0x00, 0x04, 0x10, 0x00, 0x00, 0x05,
		0x00, 0x03, 0x0E, 0x00, 0x00, 0x0A,
		0x00, 0x23, 0x00, 0x00, 0x00, 0x01,
	}
)

type Epd struct {
	Config epd_config.EpdConfig
}

func (e *Epd) Reset() {
	e.Config.Digital_writeRST(gpio.High)
	time.Sleep(200*time.Millisecond)
	e.Config.Digital_writeRST(gpio.Low)
	time.Sleep(5*time.Millisecond)
	e.Config.Digital_writeRST(gpio.High)
	time.Sleep(200*time.Millisecond)
}

func (e *Epd) Send_command(command byte){
	e.Config.Digital_writeDC(gpio.Low)
	e.Config.Digital_writeCS(gpio.Low)
	e.Config.WriteBytes([]byte{command})
	e.Config.Digital_writeCS(gpio.High)
}

func (e *Epd) Send_data(data byte){
	e.Config.Digital_writeDC(gpio.High)
	e.Config.Digital_writeCS(gpio.Low)
	e.Config.WriteBytes([]byte{data})
	e.Config.Digital_writeCS(gpio.High)
}

func (e *Epd) ReadBusy(){
	for e.Config.Digital_readBS() == gpio.Low { //low is idle; 1 is busy
		time.Sleep(200*time.Millisecond)
	}
}

func (e *Epd) Set_lut(){
	e.Send_command(0x20) //vcom
	for count:=0; count < 44; count++{
		e.Send_data(lut_vcom_dc[count])
	}
	e.Send_command(0x21) //ww --
	for count:=0; count < 42; count++{
		e.Send_data(lut_ww[count])
	}
	e.Send_command(0x22) //bw r
	for count:=0; count < 42; count++{
		e.Send_data(lut_bw[count])
	}
	e.Send_command(0x23) //wb w
	for count:=0; count < 42; count++{
		e.Send_data(lut_bb[count])
	}
	e.Send_command(0x24) //bb b
	for count:=0; count < 42; count++{
		e.Send_data(lut_wb[count])
	}
}

func (e *Epd) Setup(){
	e.Config.Setup()

	e.Reset()
	
	e.Send_command(0x01) // POWER_SETTING
    e.Send_data(0x03) // VDS_EN, VDG_EN
    e.Send_data(0x00) // VCOM_HV, VGHL_LV[1], VGHL_LV[0]
    e.Send_data(0x2b) // VDH
    e.Send_data(0x2b) // VDL
    e.Send_data(0x09) // VDHR
        
    e.Send_command(0x06) // BOOSTER_SOFT_START
    e.Send_data(0x07)
    e.Send_data(0x07)
    e.Send_data(0x17)
        
    // Power optimization
    e.Send_command(0xF8)
    e.Send_data(0x60)
    e.Send_data(0xA5)
        
    // Power optimization
    e.Send_command(0xF8)
    e.Send_data(0x89)
    e.Send_data(0xA5)
        
	// Power optimization
    e.Send_command(0xF8)
    e.Send_data(0x90)
    e.Send_data(0x00)
        
    // Power optimization
    e.Send_command(0xF8)
    e.Send_data(0x93)
    e.Send_data(0x2A)
        
    // Power optimization
    e.Send_command(0xF8)
    e.Send_data(0xA0)
    e.Send_data(0xA5)
        
    // Power optimization
    e.Send_command(0xF8)
    e.Send_data(0xA1)
    e.Send_data(0x00)
        
    // Power optimization
    e.Send_command(0xF8)
    e.Send_data(0x73)
    e.Send_data(0x41)
        
    e.Send_command(0x16) // PARTIAL_DISPLAY_REFRESH
    e.Send_data(0x00)
    e.Send_command(0x04) // POWER_ON
    e.ReadBusy()

    e.Send_command(0x00) // PANEL_SETTING
    e.Send_data(0xAF) // KW-BF   KWR-AF    BWROTP 0f
        
    e.Send_command(0x30) // PLL_CONTROL
    e.Send_data(0x3A)  // 3A 100HZ   29 150Hz 39 200HZ    31 171HZ
    
    e.Send_command(0X50) // VCOM AND DATA INTERVAL SETTING			
    e.Send_data(0x57)
        
    e.Send_command(0x82) // VCM_DC_SETTING_REGISTER
    e.Send_data(0x12)
    e.Set_lut()
}

func (e *Epd) Clear(){
	e.Send_command(0x10)
	for i:=0; i < EPD_WIDTH  * EPD_HEIGHT / 8; i++{
		e.Send_data(0xFF)
	}
	e.Send_command(0x13)
	for i:=0; i < EPD_WIDTH  * EPD_HEIGHT / 8; i++{
		e.Send_data(0xFF)
	}	
	e.Send_command(0x12)
	e.ReadBusy()

}

func (e *Epd) Sleep(){
	e.Send_command(0x50)
	e.Send_data(0xF7)
	e.Send_command(0x02)
	e.Send_command(0x07)
	e.Send_data(0xA5)

	time.Sleep(2000*time.Millisecond)
	e.Config.Destroy()
}

func (e *Epd) Display(img *image.Image){
	orientAndfittedImage := imageutil.OrientateAndFitImage(img, EPD_WIDTH, EPD_HEIGHT)
	
	monochromeTensor := ConvertImagetoMonochromeEPDTensor(&orientAndfittedImage)
	monochromeBslices:= GetEPDBuffer(monochromeTensor)
	e.Send_command(0x10)

	for i:=0; i < EPD_HEIGHT * EPD_WIDTH / 8; i++ {
		e.Send_data(0xFF)
	}
	e.Send_command(0x13)
	for i:=0; i < EPD_HEIGHT * EPD_WIDTH / 8; i++ {
		e.Send_data(monochromeBslices[i])
	}
	e.Send_command(0x12)
	e.ReadBusy()
	
}

func ConvertImagetoMonochromeEPDTensor(img *image.Image)(monochrome [][]uint8){
	p := imageutil.GetImageTensor(*img)

	//convert to greyscale tensor
	intermediateImg := imageutil.ConvertGreyScale(&p)

	p = imageutil.GetImageTensor(intermediateImg)

	for x:=0; x < len(p); x++ {
		col := []uint8{}
		for y:=0; y < len(p[0]); y++ {
			pix := p[x][y]
			originalColor, ok:= color.RGBAModel.Convert(pix).(color.RGBA)
			if !ok {
				log.Fatalf("color.color conversion went wrong")
			}
			c := originalColor.R
			if c > 200 {
				col = append(col, uint8(0))
			} else {
				col = append(col, uint8(1))
			}
		}
		monochrome = append(monochrome, col)
	}
	return 
}

func GetEPDBuffer(monochrome [][]uint8) []byte{
	imgWidth := len(monochrome)
	imgHeight := len(monochrome[0])

	buffLength := EPD_WIDTH / 8 * EPD_HEIGHT
	buf := make([]byte, buffLength)
	for i:=0; i < buffLength; i++{
		buf[i] = 0xff
	}

	if imgWidth == EPD_WIDTH && imgHeight == EPD_HEIGHT {
		//image is vertical
		for y:=0; y < imgHeight; y++{
			for x:=0; x < imgWidth; x++ {
				if monochrome[x][y] == 0 {
					index := (x + y * EPD_WIDTH) / 8
					buf[index] &= ^(0x80 >> (x % 8)) // x bit will be 0 while rest are 1s whick allow masking
				}
			}
		}
	} else if imgWidth == EPD_HEIGHT && imgHeight == EPD_WIDTH  {
		//image is horizontal
		for y:=0; y < imgHeight; y++ {
			for x:=0; x < imgWidth; x++{
				if monochrome[x][y] == 0 {
					newx := y
					newy := EPD_HEIGHT - newx - 1
					index := (newx + newy*EPD_WIDTH) / 8
					buf[index] &= ^(0x80 >> (newx % 8))

				}
			}
		}
	}

	return buf
}




