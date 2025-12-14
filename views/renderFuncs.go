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

func (vr *ViewRenderer) IndexPage(session *models.Session, sheets []*models.MemorySheet, todaySheet *models.MemorySheet) error {
	return vr.render(Index(session, sheets, todaySheet))
}

func (vr *ViewRenderer) SheetListingComponent(sheet *models.MemorySheet) error {
	return vr.render(ListingSheet(sheet))
}

func (vr *ViewRenderer) EditSheetPage(date string, content string) error {
	return vr.render(EditSheetPage(date, content))
}

func (vr *ViewRenderer) ChangePatternPage(selectedDays map[int]bool) error {
	return vr.render(ChangePatternPage(selectedDays))
}

func (vr *ViewRenderer) SheetListingComponents(sheets []*models.MemorySheet) error {
	return vr.render(SheetListingComponent(sheets))
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
