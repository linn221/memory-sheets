package views

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
	"github.com/linn221/memory-sheets/models"
)

type ViewRenderer struct {
	ctx context.Context
	w   http.ResponseWriter
	r   *http.Request
}

func (vr *ViewRenderer) ResponseWriter() http.ResponseWriter {
	return vr.w
}

func (vr *ViewRenderer) Request() *http.Request {
	return vr.r
}

func (vr *ViewRenderer) render(component templ.Component) error {
	return component.Render(vr.ctx, vr.w)
}

func (vr *ViewRenderer) IndexPage(sheets []*models.MemorySheet, todaySheet *models.MemorySheet, navSheets []*models.NavSheet) error {
	return vr.render(Index(sheets, todaySheet, navSheets))
}

func (vr *ViewRenderer) SheetComponent(sheet *models.MemorySheet) error {
	return vr.render(SheetComponent(sheet))
}

func (vr *ViewRenderer) ShowEditSheet(date string, content string) error {
	return vr.render(EditSheetForm(date, content))
}

func (vr *ViewRenderer) ShowChangePattern(selectedDays map[int]bool) error {
	return vr.render(ChangePattern(selectedDays))
}

func (vr *ViewRenderer) SheetListingComponent(sheets []*models.MemorySheet) error {
	return vr.render(SheetListingComponent(sheets))
}

func (vr *ViewRenderer) NavSheetComponent(sheet *models.NavSheet) error {
	return vr.render(NavSheetComponent(sheet))
}

func (vr *ViewRenderer) NavSheetsComponent(navSheets []*models.NavSheet) error {
	return vr.render(NavSheetsComponent(navSheets, true))
}

func (vr *ViewRenderer) ShowCreateNavSheet() error {
	return vr.render(CreateNavSheetForm())
}

func (vr *ViewRenderer) ShowEditNavSheet(title string, content string) error {
	return vr.render(EditNavSheetForm(title, content))
}

func (vr *ViewRenderer) SearchResults(memorySheets []*models.MemorySheet, navSheets []*models.NavSheet) error {
	return vr.render(SearchResults(memorySheets, navSheets))
}

func Handler(handle func(vr *ViewRenderer) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vr := ViewRenderer{
			ctx: r.Context(),
			w:   w,
			r:   r,
		}
		err := handle(&vr)
		if err != nil {
			ErrorBox(err.Error()).Render(r.Context(), w)
		}
	}
}
