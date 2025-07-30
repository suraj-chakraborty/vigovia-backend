package service

import (
	"github.com/jung-kurt/gofpdf"
)

func addFooter(pdf *gofpdf.Fpdf) {
	pdf.SetFooterFunc(func() {
		yBase := 287.0
		pdf.SetY(yBase)
		pdf.SetFont("Helvetica", "", 9)

		// Left section
		pdf.SetTextColor(0, 0, 0)
		pdf.SetXY(15, yBase-9)
		pdf.CellFormat(0, 4, "Vigovia Tech Pvt. Ltd", "", 0, "", false, 0, "")
		pdf.SetTextColor(68, 68, 68)
		pdf.SetXY(15, yBase-5)
		pdf.CellFormat(0, 4, "Registered Office: Hd-109 Cinnabar Hills,", "", 0, "", false, 0, "")
		pdf.SetXY(15, yBase-1)
		pdf.CellFormat(0, 4, "Links Business Park, Karnataka, India.", "", 0, "", false, 0, "")

		// middle section
		pdf.SetTextColor(0, 0, 0)
		pdf.SetXY(100, yBase-9)
		pdf.CellFormat(0, 4, "Phone: +91-99X9999999", "", 0, "", false, 0, "")
		pdf.SetXY(100, yBase-5)
		pdf.CellFormat(0, 4, "Email ID: Contact@Vigovia.Com", "", 0, "", false, 0, "")
		// right section
		logoPath := "./assets/logo.png"
		imgOptions := gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: false}
		pdf.ImageOptions(logoPath, 170, yBase-12, 30, 12, false, imgOptions, 0, "")
	})
}
