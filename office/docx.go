package office

import (
	"baliance.com/gooxml/common"
	"baliance.com/gooxml/document"
)

// ReadDocx function reads text from docx document.
func ReadDocx(p string) (string, error) {
	doc, err := document.Open(p)
	if err != nil {
		return "", err
	}
	t := ""
	paras := doc.Paragraphs()
	for _, para := range paras {
		runs := para.Runs()
		for _, run := range runs {
			t += run.Text()
		}
	}
	return t, nil
}

func addText(doc *document.Document, text string) {
	run := doc.AddParagraph().AddRun()
	run.AddText(text)
}

func addImage(doc *document.Document, pth string) error {
	img, err := common.ImageFromFile(pth)
	if err != nil {
		return err
	}

	imgRef, err := doc.AddImage(img)
	if err != nil {
		return err
	}

	run := doc.AddParagraph().AddRun()
	_, err = run.AddDrawingInline(imgRef)
	if err != nil {
		return err
	}

	return nil
}

type Para struct {
	Typ  string
	Text string
	Pth  string
}

func WriteDocx(doc *document.Document, paras []Para) {
	for _, para := range paras {
		if para.Typ == "text" {
			addText(doc, para.Text)
		}
		if para.Typ == "image" {
			err := addImage(doc, para.Pth)
			if err != nil {
				addText(doc, "Failed to insert image: "+para.Pth)
			}
		}
	}
}

func WriteDocxFile(pth string, paras []Para) {
	doc := document.New()
	WriteDocx(doc, paras)
	doc.SaveToFile(pth)
}
