package service

import (
	// "bytes"
	// "image"
	// "image/color"
	// "image/draw"
	// "image/png"
	"fmt"
	"itinerary/model"
	"sort"
	"time"

	"github.com/jung-kurt/gofpdf"
)

// // Clip srcImg to a circle and output PNG-encoded bytes
// func GetCircularImagePNG(srcImg image.Image, size int) ([]byte, error) {
// 	dst := image.NewRGBA(image.Rect(0, 0, size, size))

// 	// Fill background white
// 	draw.Draw(dst, dst.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

// 	// Draw circle mask
// 	for x := 0; x < size; x++ {
// 		for y := 0; y < size; y++ {
// 			dx := float64(x - size/2)
// 			dy := float64(y - size/2)
// 			if dx*dx+dy*dy <= float64((size/2)*(size/2)) {
// 				dst.Set(x, y, srcImg.At(x, y))
// 			}
// 		}
// 	}

// 	var buf bytes.Buffer
// 	if err := png.Encode(&buf, dst); err != nil {
// 		return nil, err
// 	}
// 	return buf.Bytes(), nil
// }

// custom table
func DrawTable(pdf *gofpdf.Fpdf, headings []string, rows [][]string, startX, startY float64, colWidths []float64, rowHeight float64) {
	gapWidth := 2.0
	x := startX
	y := startY

	// headings
	pdf.SetFont("Helvetica", "", 12)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFillColor(50, 30, 93)
	for i, head := range headings {
		w := colWidths[i]
		pdf.SetXY(x, y)
		pdf.RoundedRect(x, y, w, rowHeight, 5, "1,2", "F")
		pdf.CellFormat(w, rowHeight, head, "", 0, "CM", false, 0, "")
		x += w + gapWidth
	}
	pdf.Ln(rowHeight)
	y += rowHeight

	// rows
	pdf.SetFont("Helvetica", "", 12)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFillColor(245, 230, 255)
	for rowIndex, row := range rows {
		x = startX
		isLastRow := rowIndex == len(rows)-1

		// Calculate max height
		maxHeight := 0.0
		lineCounts := make([]int, len(row))
		for i, cell := range row {
			lines := pdf.SplitLines([]byte(cell), colWidths[i])
			lineCounts[i] = len(lines)
			height := float64(len(lines)) * rowHeight
			if height > maxHeight {
				maxHeight = height
			}
		}

		for i, cell := range row {
			w := colWidths[i]
			pdf.SetXY(x, y)

			// Draw background
			if isLastRow {
				corners := "3,4"

				if corners != "" {
					pdf.RoundedRect(x, y, w, maxHeight, 5, corners, "F")
				} else {
					pdf.Rect(x, y, w, maxHeight, "F")
				}
			} else {
				pdf.Rect(x, y, w, maxHeight, "F")
			}

			// Print text
			currX, currY := x, y
			pdf.SetXY(currX, currY)
			pdf.MultiCell(w, rowHeight, cell, "", "CM", false)

			x += w + gapWidth
		}

		y += maxHeight
		pdf.Ln(-1)
	}
}

func FormatCustomDate(dateStr string) string {
	// Parse the date in "yyyy-mm-dd" format
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return dateStr // fallback to original if parsing fails
	}
	// weekdays := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	months := []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
	// weekday := weekdays[t.Weekday()]
	day := t.Format("02")
	month := months[int(t.Month())-1]
	// yearShort := t.Format("06") // two-digit year
	return day + " " + month
}

// GroupActivitiesByDay organizes activities by unique dates and returns them as sorted day-wise plans
func GroupActivitiesByDay(activities []model.Activity) ([]model.Day, error) {
	// Create a map to group activities by their date string
	dateMap := make(map[string][]model.Activity)
	for _, act := range activities {
		dateMap[act.Date] = append(dateMap[act.Date], act)
	}

	// Extract unique dates from the map
	var uniqueDates []string
	for date := range dateMap {
		uniqueDates = append(uniqueDates, date)
	}

	// Sort the unique date strings in chronological order
	sort.Slice(uniqueDates, func(i, j int) bool {
		ti, _ := time.Parse("2006-01-02", uniqueDates[i])
		tj, _ := time.Parse("2006-01-02", uniqueDates[j])
		return ti.Before(tj)
	})

	// Build day-wise plans with formatted date strings and grouped activities
	var plans []model.Day
	for i, dateStr := range uniqueDates {
		formatted := FormatCustomDate(dateStr)
		plans = append(plans, model.Day{
			Day:           i + 1,
			FormattedDate: formatted,
			OriginalDate:  dateStr,
			Image:         nil,
			Activities:    dateMap[dateStr],
		})
	}
	return plans, nil
}

// getActivitiesForSlot filters the activities based on the time slot (Morning, Afternoon, Evening)
func getActivitiesForSlot(activities []model.Activity, slot string, date string) []model.Activity {
	var result []model.Activity
	for _, a := range activities {
		if a.Date != date {
			continue
		}

		parsedTime, err := time.Parse("3:04 PM", a.Time)
		if err != nil {
			fmt.Println("Failed to parse time:", a.Time, err)
			continue
		}

		hour := parsedTime.Hour()

		switch slot {
		case "Morning":
			if hour >= 6 && hour < 12 {
				result = append(result, a)
			}
		case "Afternoon":
			if hour >= 12 && hour < 17 {
				result = append(result, a)
			}
		case "Evening":
			if hour >= 17 && hour < 22 {
				result = append(result, a)
			}
		}
	}
	return result
}
