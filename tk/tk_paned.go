package tk

import "fmt"

// panedwindow
type TKPaned struct {
	BaseWidget
}

func NewTKPaned(parent Widget, orient Orient, attributes ...*WidgetAttr) *TKPaned {
	iid := makeNamedWidgetId(parent, "atk_paned")
	attributes = append(attributes, &WidgetAttr{"orient", orient})
	info := CreateWidgetInfo(iid, WidgetTypePaned, false, attributes)
	if info == nil {
		return nil
	}
	w := &TKPaned{}
	w.id = iid
	w.info = info
	RegisterWidget(w)
	return w
}

func (w *TKPaned) Attach(id string) error {
	info, err := CheckWidgetInfo(id, WidgetTypePaned)
	if err != nil {
		return err
	}
	w.id = id
	w.info = info
	RegisterWidget(w)
	return nil
}

func (w *TKPaned) SetWidth(width int) error {
	return eval(fmt.Sprintf("%v configure -width {%v}", w.id, width))
}

func (w *TKPaned) Width() int {
	r, _ := evalAsInt(fmt.Sprintf("%v cget -width", w.id))
	return r
}

func (w *TKPaned) SetHeight(height int) error {
	return eval(fmt.Sprintf("%v configure -height {%v}", w.id, height))
}

func (w *TKPaned) Height() int {
	r, _ := evalAsInt(fmt.Sprintf("%v cget -height", w.id))
	return r
}

func (w *TKPaned) AddWidget(widget Widget, attribute ...*WidgetAttr) error {
	if !IsValidWidget(widget) {
		return ErrInvalid
	}
	is_ttk := false
	extra := buildWidgetAttributeScript(w.info.MetaClass, is_ttk, attribute)
	return eval(fmt.Sprintf("%v add %v %v", w.id, widget.Id(), extra))
}

func (w *TKPaned) InsertWidget(pane int, widget Widget, weight int) error {
	if !IsValidWidget(widget) {
		return ErrInvalid
	}
	return eval(fmt.Sprintf("%v insert %v %v -weight %v", w.id, pane, widget.Id(), weight))
}

func (w *TKPaned) SetPane(pane int, weight int) error {
	return eval(fmt.Sprintf("%v pane %v -weight %v", w.id, pane, weight))
}

func (w *TKPaned) RemovePane(pane int) error {
	return eval(fmt.Sprintf("%v forget %v", w.id, pane))
}

func (w *TKPaned) HidePane(pane int, hide bool) error {
	pane_list, err := evalAsStringList(fmt.Sprintf("%v panes", w.id))
	if err != nil {
		return err
	}
	window := pane_list[pane] // todo: unsafe and out of keeping.
	return eval(fmt.Sprintf("%v paneconfigure %v -hide %v", w.id, window, hide))
}
