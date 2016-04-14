package lib

import (
	"fmt"
	"gopkg.in/gographics/imagick.v2/imagick"
	"math"
)

func crop(mw *imagick.MagickWand, x, y int, cols, rows uint) error {
	var result error
	result = nil

	imCols := mw.GetImageWidth()
	imRows := mw.GetImageHeight()

	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}

	if uint(x) >= imCols || uint(y) >= imRows {
		result = fmt.Errorf("x, y more than image cols, rows")
		return result
	}

	if cols == 0 || imCols < uint(x) + cols {
		cols = imCols - uint(x)
	}

	if rows == 0 || imRows < uint(y) + rows {
		rows = imRows - uint(y)
	}

	fmt.Print(fmt.Printf("wi_crop(im, %d, %d, %d, %d)\n", x, y, cols, rows))

	result = mw.CropImage(cols, rows, x, y)

	return result
}

func proportion(mw *imagick.MagickWand, proportion int, cols uint, rows uint) error {
	var result error
	result = nil

	imCols := mw.GetImageWidth()
	imRows := mw.GetImageHeight()

	if proportion == 0 {
		fmt.Printf("p=0, wi_scale(im, %d, %d)\n", cols, rows)
		result = mw.ResizeImage(cols, rows, imagick.FILTER_UNDEFINED, 1.0)

	} else if proportion == 1 {

		if cols == 0 || rows == 0 {
			if cols > 0 {
				rows = uint(round(float64((cols / imCols) * imRows)))
			} else {
				cols = uint(round(float64((rows / imRows) * imCols)))
			}
			fmt.Printf("p=1, wi_scale(im, %d, %d)\n", cols, rows)
			result = mw.ResizeImage(cols, rows, imagick.FILTER_UNDEFINED, 1.0)
		} else {
			var x, y, sCols, sRows uint
			x, y = 0, 0

			colsRate := float64(cols) / float64(imCols)
			rowsRate := float64(rows) / float64(imRows)

			if colsRate > rowsRate {
				sCols = cols
				sRows = uint(round(float64(colsRate * float64(imRows))))
				y = uint(math.Floor(float64((sRows - rows) / 2.0)))
			} else {
				sCols = uint(round(float64(rowsRate * float64(imCols))))
				sRows = rows
				x = uint(math.Floor(float64((sCols - cols) / 2.0)))
			}

			fmt.Printf("p=2, wi_scale(im, %d, %d)\n", sCols, sRows)
			//result = mw.ResizeImage(sCols, sRows, imagick.FILTER_UNDEFINED, 1.0)
			mw.SetGravity(imagick.GRAVITY_CENTER)
			result = mw.StripImage()
			result = mw.ThumbnailImage(sCols, sRows)
			fmt.Printf("p=2, wi_crop(im, %d, %d, %d, %d)\n", x, y, cols, rows)
			result = mw.CropImage(cols, rows, int(x), int(y))
		}

	} else if proportion == 2 {
		x := int(math.Floor(float64((imCols - cols) / 2.0)))
		y := int(math.Floor(float64((imRows - rows) / 2.0)))
		fmt.Printf("p=3, wi_crop(im, %d, %d, %d, %d)\n", x, y, cols, rows)
		result = mw.CropImage(cols, rows, x, y)

	} else if proportion == 3 {
		if cols == 0 || rows == 0 {
			var rate uint
			if cols > 0 {
				rate = cols
			} else {
				rate = rows
			}
			rows = uint(round(float64(imRows * rate / 100)))
			cols = uint(round(float64(imCols * rate / 100)))
			fmt.Printf("p=3, wi_scale(im, %d, %d)\n", cols, rows)
			//result = mw.ResizeImage(cols, rows, imagick.FILTER_UNDEFINED, 1.0)
			mw.SetGravity(imagick.GRAVITY_CENTER)
			result = mw.StripImage()
			result = mw.ThumbnailImage(cols, rows)
		} else {
			rows = uint(round(float64(imRows * rows / 100)))
			cols = uint(round(float64(imCols * cols / 100)))
			fmt.Printf("p=3, wi_scale(im, %d, %d)\n", cols, rows)
			//result = mw.ResizeImage(cols, rows, imagick.FILTER_UNDEFINED, 1.0)
			mw.SetGravity(imagick.GRAVITY_CENTER)
			result = mw.StripImage()
			result = mw.ThumbnailImage(cols, rows)
		}

	} else if proportion == 4 {
		var rate float64
		rate = 1.0
		if cols == 0 || rows == 0 {
			if cols > 0 {
				rate = float64(cols / imCols)
			} else {
				rate = float64(rows / imRows)
			}
		} else {
			rateCol := float64(cols) / float64(imCols)
			rateRow := float64(rows) / float64(imRows)
			if rateCol < rateRow {
				rate = rateCol
			} else {
				rate = rateRow
			}
		}

		cols = uint(round(float64(float64(imCols) * rate)))
		rows = uint(round(float64(float64(imRows) * rate)))
		fmt.Printf("p=4, wi_scale(im, %d, %d)\n", cols, rows)
		//result = mw.ResizeImage(cols, rows, imagick.FILTER_UNDEFINED, 1.0)
		mw.SetGravity(imagick.GRAVITY_CENTER)
		result = mw.StripImage()
		result = mw.ThumbnailImage(cols, rows)
	} else {
		fmt.Printf("p=%v\n", proportion)
	}

	return result

}

func convert(mw *imagick.MagickWand, conf ScaleConf) error {

	fmt.Println("call convert function......")

	var result error
	result = nil
	mw.ResetIterator()
	mw.SetImageOrientation(imagick.ORIENTATION_TOP_LEFT)

	x := conf.X
	y := conf.Y
	cols := uint(conf.Width)
	rows := uint(conf.Height)

	fmt.Printf("image cols %d, rows %d \n", cols, rows)

	if !(cols == 0 && rows == 0) {

		/* crop and scale */
		if x == -1 && y == -1 {
			fmt.Println("call proportion function......")
			fmt.Print(fmt.Printf("proportion(im, %d, %d, %d) \n", conf.Proportion, cols, rows))

			result = proportion(mw, conf.Proportion, cols, rows)
			if result != nil {
				return result
			}
		} else {
			fmt.Println("call crop function......")
			fmt.Print(fmt.Printf("crop(im, %d, %d, %d, %d) \n", x, y, cols, rows))

			result = crop(mw, x, y, cols, rows)
			if result != nil {
				return result
			}
		}
	}

	/* rotate image */
	if conf.Rotate != 0 {
		fmt.Print(fmt.Printf("wi_rotate(im, %d) \n", conf.Rotate))

		background := imagick.NewPixelWand()
		if background == nil {
			result = fmt.Errorf("init new pixelwand faile.")
			return result
		}
		defer background.Destroy()
		isOk := background.SetColor("#ffffff")
		if !isOk {
			result = fmt.Errorf("set background color faile.")
			return result
		}

		result = mw.RotateImage(background, float64(conf.Rotate))
		if result != nil {
			return result
		}
	}

	/* set gray */
	if conf.Gary == 1 {
		fmt.Print(fmt.Printf("wi_gray(im) \n"))
		result = mw.SetImageType(imagick.IMAGE_TYPE_GRAYSCALE)
		if result != nil {
			return result
		}
	}

	/* set quality */
	fmt.Print(fmt.Printf("wi_set_quality(im, %d) \n", conf.Quality))
	result = mw.SetImageCompressionQuality(uint(conf.Quality))
	if result != nil {
		return result
	}

	/* set format */
	if "none" != conf.Format {
		fmt.Print(fmt.Printf("wi_set_format(im, %s) \n", conf.Format))
		result = mw.SetImageFormat(conf.Format)
		if result != nil {
			return result
		}
	}

	fmt.Print(fmt.Printf("convert(im, req) %s \n", result))

	return result
}

func round(val float64) float64 {
	if val > 0.0 {
		return math.Floor(val + 0.5)
	} else {
		return math.Ceil(val - 0.5)
	}
}
