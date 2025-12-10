package views

import (
	"context"
	"net/http"

	"github.com/linn221/memory-sheets/app"
)

type ViewRenderer struct {
	ctx context.Context
	w   http.ResponseWriter
}

func (vr *ViewRenderer) ListSheets(sheets []*app.MemorySheet) error {
	return Index(sheets).Render(vr.ctx, vr.w)
}

func Handler(handle func(vr *ViewRenderer) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vr := ViewRenderer{
			ctx: r.Context(),
			w:   w,
		}
		err := handle(&vr)
		if err != nil {
			ErrorBox(err.Error()).Render(r.Context(), w)
		}
	}
}
