package service

import (
	// "bytes"
	"fmt"
	"itinerary/model"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/jung-kurt/gofpdf"
)

func getTextWidth(pdf *gofpdf.Fpdf, text string) float64 {
	textlen := pdf.GetStringWidth(text)
	return textlen
}

func GeneratePDF(data model.ItineraryData) (string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")

	addFooter(pdf)
	pdf.AddPage()
	// variables
	y := 05.0
	margin := 10.0
	boxW := 170.0
	boxH := 40.0
	pageW, _ := pdf.GetPageSize()

	var totalPrice float32

	// <--!---------------------------logo section---------------------->
	logoOpts := gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: false}
	pdf.ImageOptions("./assets/logo.png", (pageW-50)/2, y, 50, 20, false, logoOpts, 0, "")
	y += 25

	dep, _ := time.Parse("2006-01-02", data.DepartureDate)
	ret, _ := time.Parse("2006-01-02", data.ReturnDate)
	totalDays := int(ret.Sub(dep).Hours()/24) + 1
	nights := totalDays - 1

	// <--!---------------------------Header section---------------------->

	boxX := (pageW - boxW) / 2
	pdf.SetDrawColor(84, 28, 156)
	pdf.RoundedRect(boxX, y, boxW, boxH, 4, "1234", "D")
	pdf.ImageOptions("./assets/headerBg.png", boxX-8, y, boxW+10, boxH, false, logoOpts, 0, "")
	pdf.SetFont("Helvetica", "B", 22)
	pdf.SetTextColor(255, 255, 255)
	// nameWidth := pdf.GetStringWidth(data.Name)
	pdf.Text(pageW/2-getTextWidth(pdf, data.Name)/2-10, y+14, fmt.Sprintf("Hi, %s!", data.Name))
	pdf.SetFontSize(18)
	// nameWidth := pdf.GetStringWidth(data.DestinationCity)
	pdf.Text(pageW/2-getTextWidth(pdf, data.Name)/2-10, y+22, fmt.Sprintf("%s Itinerary", data.DestinationCity))
	pdf.SetFont("Helvetica", "", 14)
	pdf.Text(pageW/2-getTextWidth(pdf, data.Name)/2-10, y+30, fmt.Sprintf("%d Days %d Nights", totalDays, nights))
	y += 35

	// Icons
	iconPaths := []string{"./assets/icon1.png", "./assets/icon2.png", "./assets/icon3.png", "./assets/icon4.png", "./assets/icon5.png"}
	iconSize := 4.0
	spacing := 5.0
	totalIconsW := iconSize*float64(len(iconPaths)) + spacing*float64(len(iconPaths)-1)
	startX := (pageW - totalIconsW) / 2
	for i, ip := range iconPaths {
		pdf.ImageOptions(ip, startX+float64(i)*(iconSize+spacing), y, iconSize, iconSize, false, logoOpts, 0, "")
	}
	y += 20
	// <--!--------------------Trip destination section---------------------->
	labels := []string{"Departure From:", "Departure:", "Arrival:", "Destination:", "No. of Travellers:"}
	values := []string{data.DepartureCity, data.DepartureDate, data.ReturnDate, data.DestinationCity, fmt.Sprintf("%d", data.Travelers)}
	infoBoxW := 180.0
	infoBoxH := 10.0
	x0 := (pageW - infoBoxW) / 2
	pdf.SetDrawColor(84, 28, 156)
	pdf.RoundedRect(x0-2, y-10, infoBoxW, infoBoxH+10, 4, "1234", "D")
	colW := infoBoxW / float64(len(labels))
	for i := range labels {
		// cx := x0 + colW/float64(len(labels))
		pdf.SetFont("Helvetica", "B", 10)
		pdf.SetTextColor(0, 0, 0)
		pdf.Text(x0+3+colW*float64(i), y, labels[i])
		pdf.SetFont("Helvetica", "", 10)
		pdf.SetTextColor(0, 0, 0)
		pdf.Text(x0+3+colW*float64(i), y+5, values[i])
	}
	y += infoBoxH + 5

	// <--!---------------------------day section---------------------->
	days, err := GroupActivitiesByDay(data.Activity)
	if err != nil {
		return "Something went wrong", err
	}

	for _, day := range days {
		if y > 230 {
			pdf.AddPage()
			y = margin
		}
		// left
		pdf.SetFillColor(50, 30, 93)
		pdf.RoundedRect(15, y+3, 12, 50, 5, "1234", "F")
		pdf.SetTextColor(255, 255, 255)
		pdf.SetFont("Helvetica", "B", 15)
		pdf.TransformBegin()
		pdf.TransformRotate(90, 15+14, y+30)
		pdf.Text(25, y+23, fmt.Sprintf("Day %d", day.Day))
		pdf.TransformEnd()

		// middle
		pdf.Circle(pageW/3-13, y+10+15, 15, "D")
		pdf.ClipCircle(pageW/3-13, y+10+15, 15, true)
		pdf.Image("./assets/bg.jpg", pageW/3-30, y+1, 40, 47, false, "", 0, "")
		pdf.ClipEnd()

		pdf.SetFont("Helvetica", "B", 11)
		pdf.SetTextColor(0, 0, 0)
		pdf.Text(pageW/4-5, y+10+30+7, day.FormattedDate)

		pdf.SetFont("Helvetica", "", 9)
		pdf.SetTextColor(100, 100, 100)
		desc := fmt.Sprintf("Arrival in %s & City Exploration", data.DestinationCity)
		pdf.Text(pageW/4-15, y+10+30+13, desc)

		// right
		timelineX := pageW / 2.3
		pdf.SetDrawColor(47, 128, 237)
		pdf.SetLineWidth(0.5)
		pdf.Line(timelineX, y+10, timelineX, y+10+40)

		activitySlots := []string{"Morning", "Afternoon", "Evening"}
		for i, slot := range activitySlots {
			cy := y + 10 + float64(i)*20
			pdf.SetFillColor(255, 255, 255)
			pdf.Circle(timelineX, cy, 2, "DF")
			pdf.SetFont("Helvetica", "B", 10)
			pdf.SetTextColor(0, 0, 0)
			pdf.Text(timelineX+6, cy+0.5, slot+":")

			slotActivities := getActivitiesForSlot(data.Activity, slot, day.OriginalDate)
			// fmt.Println(slotActivities)
			textY := cy - 3
			for _, act := range slotActivities {
				pdf.SetFont("Helvetica", "", 9)
				pdf.SetTextColor(68, 68, 68)
				pdf.SetXY(timelineX+25, textY)
				pdf.MultiCell(100, 5, act.Activity, "", "L", false)
				textY += 10
			}
		}

		y += 60
		pdf.SetDrawColor(200, 200, 200)
		pdf.Line(20, y, pageW-20, y)
		y += 10
	}
	// // Draw activity summary table
	if y > 250 {
		pdf.AddPage()
		y = margin
	}

	y += infoBoxH + 5
	// <--!---------------------------flight section---------------------->
	if len(data.Flights) > 0 {
		y += 10

		pdf.SetTextColor(0, 0, 0)
		pdf.SetFont("Helvetica", "B", 14)
		pdf.Text(15, y, "Flight ")
		pdf.SetTextColor(104, 0, 153)
		pdf.SetFont("Helvetica", "B", 14)
		pdf.Text(pdf.GetStringWidth("Important ")+5, y, "Summary")

		y += 10
		for _, f := range data.Flights {
			if y > 250 {
				pdf.AddPage()
				y = 20
			}
			totalPrice += f.FlightPrice
			// fmt.Println("flightPrice", totalPrice)
			boxX := 15.0
			boxY := y
			boxWidth := 180.0
			boxHeight := 15.0
			// Box border
			pdf.SetDrawColor(180, 180, 180)
			pdf.SetLineWidth(0.5)
			pdf.RoundedRect(boxX, boxY, boxWidth, boxHeight-3, 2, "1234", "")

			// Date section - fill with light purple rectangle or image
			pdf.SetFillColor(237, 233, 254)
			pdf.Image("./assets/flirect.png", boxX, boxY, 40, boxHeight-3, false, "", 0, "")

			// Departure date
			pdf.SetTextColor(0, 0, 0)
			pdf.SetFont("Helvetica", "B", 12)
			departureDate := f.DepartureDate
			if departureDate == "" {
				departureDate = "N/A"
			}
			pdf.Text(boxX+5, boxY+8, FormatCustomDate(departureDate))

			// Flight info
			infoX := boxX + 45
			airline := f.FlightName
			fromCity := f.From
			toCity := f.To

			if airline == "" {
				airline = "N/A"
			}
			if fromCity == "" {
				fromCity = "N/A"
			}
			if toCity == "" {
				toCity = "N/A"
			}

			pdf.SetFont("Helvetica", "B", 11)
			pdf.SetTextColor(0, 0, 0)
			pdf.Text(infoX, boxY+8, airline)

			// Measure airline width to align route text
			airlineWidth := pdf.GetStringWidth(airline + " ")

			// Route info
			route := fmt.Sprintf("From %s  To %s",
				fromCity,
				toCity)

			pdf.SetFont("Helvetica", "", 10)
			pdf.SetTextColor(68, 68, 68)
			pdf.Text(infoX+airlineWidth, boxY+8, route)

			y += boxHeight
		}

		pdf.SetFont("Helvetica", "", 8)
		pdf.SetTextColor(128, 128, 128)
		pdf.Text(20, y+10, "Note: All Flights Include Meals, Seat Choice (Excluding XL), and 20kg/25kg Checked Baggage.")
		y += 5
		pdf.SetDrawColor(200, 200, 200)
		pdf.Line(20, y, 190, y)
		y += 10
	}

	// <--!-----------------------Hotel booking section---------------------->
	if len(data.Bookings) > 0 {
		y += 10
		if y > 250 {
			pdf.AddPage()
			y = 20
		}

		pdf.SetTextColor(0, 0, 0)
		pdf.SetFont("Helvetica", "B", 14)
		pdf.Text(15, y, "Hotel ")
		pdf.SetTextColor(104, 0, 153)
		pdf.SetFont("Helvetica", "B", 14)
		pdf.Text(pdf.GetStringWidth("Important ")+5, y, "Bookings")

		y += 10
		headings := []string{"city", "Check-in", "Check-out", "Nights", "Hotel Name"}

		var rows [][]string
		for _, b := range data.Bookings {
			rows = append(rows, []string{
				b.City,
				b.CheckIn,
				b.CheckOut,
				strconv.Itoa(b.Nights),
				b.HotelName,
			})
			// fmt.Println("+booking", b.HotelPrice)
			// fmt.Println("flight+", totalPrice)
			totalPrice += b.HotelPrice
		}
		// fmt.Println("flight+booking", totalPrice)

		// Define column widths, adjust as needed
		colWidths := []float64{30, 30, 30, 30, 50}

		DrawTable(pdf, headings, rows, boxX, y, colWidths, 10)
		y += boxH
		pdf.SetFont("Helvetica", "", 8)
		pdf.SetTextColor(128, 128, 128)
		pdf.Text(20, y, "1. All hotels are tentative and can be replaced with similar.")
		pdf.Text(20, y+3, "2. Breakfast included for all hotel stays.")
		pdf.Text(20, y+6, "3. All Hotels will be 4* and above category")
		pdf.Text(20, y+9, "4. A maximum occupancy of 2 people/room is allowed in most hotels.")
		y += 5
	}
	y += 20

	// <--!---------------------Important Notes section---------------------->
	headerNotes := []string{"Point", "Details"}
	notes := [][]string{
		{"Airlines Standard Policy", "In Case Of Visa Rejection, Visa Fees Or Any Other Non Cancellable Component Cannot Be Reimbursed At Any Cost."},
		{"Flight/Hotel Cancellation", "In Case Of Visa Rejection, Visa Fees Or Any Other Non Cancellable Component Cannot Be Reimbursed At Any Cost."},
		{"Trip Insurance", "In Case Of Visa Rejection, Visa Fees Or Any Other Non Cancellable Component Cannot Be Reimbursed At Any Cost."},
		{"Hotel Check-In & Check Out", "In Case Of Visa Rejection, Visa Fees Or Any Other Non Cancellable Component Cannot Be Reimbursed At Any Cost."},
		{"Visa Rejection", "In Case Of Visa Rejection, Visa Fees Or Any Other Non Cancellable Component Cannot Be Reimbursed At Any Cost."},
	}

	if y > 250 {
		pdf.AddPage()
		y = 20
	}
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Helvetica", "B", 14)
	pdf.Text(15, y, "Important ")
	pdf.SetTextColor(104, 0, 153)
	pdf.SetFont("Helvetica", "B", 14)
	pdf.Text(pdf.GetStringWidth("Important ")+15, y, "Notes")

	y += 10
	DrawTable(pdf, headerNotes, notes, boxX-5, y, []float64{60, 120}, 10)
	y += 20

	// <--!--------------Scope of Service section---------------------->
	y += 80
	serviceHeader := []string{"Service", "Details"}
	services := [][]string{
		{"Flight Tickets And Hotel Vouchers", "Delivered 3 Days Post Full Payment"},
		{"Web Check-In", "Boarding Pass Delivery Via Email/WhatsApp"},
		{"Support", "Chat Support – Response Time: 4 Hours"},
		{"Cancellation Support", "Provided"},
		{"Trip Support", "Response Time: 5 Minutes"},
	}

	y += 30
	if y > 230 {
		pdf.AddPage()
		y = 20
	}
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Helvetica", "B", 14)
	pdf.Text(15, y, "Scope Of ")
	pdf.SetTextColor(104, 0, 153)
	pdf.SetFont("Helvetica", "B", 14)
	pdf.Text(pdf.GetStringWidth("Important ")+15, y, "Service")
	y += 10
	DrawTable(pdf, serviceHeader, services, boxX-10, y, []float64{80, 100}, 10)
	y += 40

	// <--!--------------Inclusion Summary section---------------------->
	y += 40
	if y > 230 {
		pdf.AddPage()
		y = 20
	}
	summaryHeader := []string{"Category", "Count", "Details", "Status/Comments"}
	summaryData := [][]string{
		{"Flight", "2", "All flights mentioned", "Awaiting Confirmation"},
		{"Tourist Tax", "2", "Yotel (Singapore), Oakwood (Sydney), Mercure (Cairns), Novotel(Gold Coast), Holiday Inn(Melbourne)", "Awaiting Confirmation"},
		{"Hotel", "2", "Airport to Hotel - Hotel to Attractions - Day trips if any", "Included"},
	}

	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Helvetica", "B", 14)
	pdf.Text(15, y, "Inclusion")
	pdf.SetTextColor(104, 0, 153)
	pdf.SetFont("Helvetica", "B", 14)
	pdf.Text(pdf.GetStringWidth("Important ")+15, y, "Summary")
	y += 10
	DrawTable(pdf, summaryHeader, summaryData, boxX-5, y, []float64{30, 30, 75, 40}, 10)
	y += 10

	if y > 250 {
		pdf.AddPage()
		y = 20
	}

	// ----- Transfer Policy -----
	y += 80
	pdf.SetFont("Helvetica", "", 10)
	pdf.SetTextColor(0, 0, 0)
	pdf.Text(10, y, "Transfer Policy (Refundable Upon Claim)")
	pdf.SetTextColor(100, 100, 100)
	pdf.SetFont("Helvetica", "", 10)
	policyText := "If any transfer is delayed beyond 15 minutes, customers may book an app-based or radio taxi and claim a refund for that specific leg."
	lines := pdf.SplitLines([]byte(policyText), 180)
	for i, line := range lines {
		pdf.Text(10, y+float64(4+i*5), string(line))
	}
	y += 20

	// <--!---------------------------Activity section---------------------->
	if y > 200 {
		pdf.AddPage()
		y = 20
	}
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Helvetica", "B", 14)
	pdf.Text(15, y, "Activity ")
	pdf.SetTextColor(104, 0, 153)
	pdf.SetFont("Helvetica", "B", 14)
	pdf.Text(pdf.GetStringWidth("Activity")+15, y, "Table")
	y += 10
	AHeading := []string{"city", "Activity", "Date/Time", "Time Required"}
	var ARow [][]string
	for _, A := range data.Activity {
		ARow = append(ARow, []string{
			A.City,
			A.Activity,
			A.Date + "/" + A.Time,
			A.TimeRequired,
		})
		// fmt.Println(A.ActivityCost)
		totalPrice += A.ActivityCost
		// fmt.Println(totalPrice)
	}

	DrawTable(pdf, AHeading, ARow, boxX-5, y, []float64{30, 70, 50, 30}, 10)

	y += 25
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Helvetica", "B", 14)
	pdf.Text(15, y+30, "Terms and")
	pdf.SetTextColor(104, 0, 153)
	pdf.SetFont("Helvetica", "B", 14)
	pdf.Text(pdf.GetStringWidth("Terms and")+15, y+30, "Condition")
	y += 5
	url := "https://www.vigovia.com/"
	pdf.SetTextColor(47, 128, 237)
	pdf.SetFont("Helvetica", "", 12)
	pdf.SetXY(boxX-5, y+30)
	pdf.WriteLinkString(10, "View all terms and conditions", url)
	y += 10
	// <--!---------------------------payment section---------------------->
	y += 50
	if y > 250 {
		pdf.AddPage()
		y = 20
	}

	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Helvetica", "B", 14)
	pdf.Text(15, y, "payment ")
	pdf.SetTextColor(104, 0, 153)
	pdf.SetFont("Helvetica", "B", 14)
	pdf.Text(pdf.GetStringWidth("Important ")+15, y, "Plan")
	y += 10
	for i, p := range data.Payments {
		if p.Label == "Total Amount" {
			// Ensure slice is initialized
			if len(p.Value) == 0 {
				p.Value = make([]string, 1)
			}
			pdf.SetFont("Helvetica", "", 14)
			// Assign totalPrice to the "Total Amount" only
			p.Value[0] = "Rs " + strconv.FormatFloat(float64(totalPrice), 'f', 2, 64) + " (inclusive of GST)"
			data.Payments[i] = p
		}

		if y > 250 {
			pdf.AddPage()
			y = 20
		}

		boxX := 15.0
		boxY := y
		boxWidth := 180.0
		boxHeight := 15.0

		// Ensure p.Value has at least one element
		if len(p.Value) == 0 {
			p.Value = []string{"Not Collected"}
		}

		// Draw box
		pdf.SetDrawColor(180, 180, 180)
		pdf.SetLineWidth(0.5)
		pdf.RoundedRect(boxX, boxY, boxWidth, boxHeight-3, 2, "1234", "")
		pdf.SetFillColor(237, 233, 254)
		pdf.Image("./assets/flirect.png", boxX, boxY, 40, boxHeight-3, false, "", 0, "")

		// Set label (as date substitute)
		pdf.SetTextColor(0, 0, 0)
		pdf.SetFont("Helvetica", "B", 12)
		label := p.Label
		if label == "" {
			label = "N/A"
		}
		pdf.Text(boxX+5, boxY+8, FormatCustomDate(label))

		// Set value
		value := p.Value[0]
		if value == "" {
			value = "Not Collected"
		}
		pdf.SetFont("Helvetica", "", 14)
		pdf.SetTextColor(0, 0, 0)
		pdf.Text(boxX+45, boxY+8, value)

		y += boxHeight
	}

	y += 5

	Iheading := []string{"Installment", "Amount", "DueDate"}
	var Irow [][]string
	var paidAmount float64 = 0

	// First pass: calculate paid amount
	for _, I := range data.Installments {
		if I.Installment != "pending" && I.Amount != "" {
			f64, err := strconv.ParseFloat(I.Amount, 64)
			if err == nil {
				paidAmount += f64
			}
		}
	}

	pendingAmount := float64(totalPrice) - paidAmount

	for _, I := range data.Installments {
		amount := I.Amount
		if I.Installment == "pending" {
			amount = fmt.Sprintf("%.2f", pendingAmount)
		}
		Irow = append(Irow, []string{
			I.Installment,
			amount,
			I.DueDate,
		})
	}

	if y > 230 {
		pdf.AddPage()
		y = 20
	}
	DrawTable(pdf, Iheading, Irow, boxX-5, y, []float64{60, 60, 60}, 10)

	y += 20
	if y > 230 {
		pdf.AddPage()
		y = 20
	}
	y += 40
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Helvetica", "B", 14)
	pdf.Text(15, y, "Visa ")
	pdf.SetTextColor(104, 0, 153)
	pdf.SetFont("Helvetica", "B", 14)
	pdf.Text(pdf.GetStringWidth("Important ")+5, y, "Details")

	y += 20

	labels1 := []string{"Visa Type", "Validity", "Processing Date"}
	var values1 [][]string
	for _, pd := range data.PaymentDetails {
		values1 = append(values1, []string{
			pd.Visa,
			pd.Validity,
			pd.ProcessingDate,
		})
	}

	pdf.SetDrawColor(84, 28, 156)
	pdf.RoundedRect(x0-2, y-10, infoBoxW, infoBoxH+6, 4, "1234", "D")
	// Header Row
	for i, label := range labels1 {
		pdf.SetFont("Helvetica", "B", 10)
		pdf.SetTextColor(0, 0, 0)
		pdf.Text(x0+8+colW*float64(i)*1.8, y-3, label)
	}

	// Data Rows
	pdf.SetFont("Helvetica", "", 10)
	pdf.SetTextColor(0, 0, 0)

	for rowIdx, row := range values1 {
		yOffset := y + 6 + float64(rowIdx)*5
		for colIdx, cell := range row {
			pdf.Text(x0+8+colW*float64(colIdx)*1.8, yOffset-3, cell)
		}
	}
	y += infoBoxH + 5

	if y > 230 {
		pdf.AddPage()
		y = 20
	}
	y += 10
	pdf.SetFont("Helvetica", "", 16)
	pdf.SetTextColor(50, 0, 120)
	pdf.Text(pageW/2-getTextWidth(pdf, "PLAN.PACK.GO!")/2-10, y, "PLAN.PACK.GO!")
	pdf.ImageOptions("./assets/button.png", pageW/2-70/2-10, y+10, 70, 15, false, logoOpts, 0, "")
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Helvetica", "B", 18)
	pdf.Text(pageW/2-getTextWidth(pdf, "Book Now")/2-10, y+20, "Book Now")

	// <--!----------------------The End--------------------------------->

	// Ensure the "pdfs" directory exists
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	downloadsDir := filepath.Join(homeDir, "Downloads", "Output")

	// Ensure the directory exists
	err = os.MkdirAll(downloadsDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	filename := fmt.Sprintf("itinerary_%d.pdf", time.Now().Unix())
	fullPath := filepath.Join(downloadsDir, filename)

	// Save the PDF
	err = pdf.OutputFileAndClose(fullPath)
	if err != nil {
		return "", err
	}

	return fullPath, nil
}
