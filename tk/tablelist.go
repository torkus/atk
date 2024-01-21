package tk

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// https://www.nemethi.de/tablelist/tablelistWidget.html#col_options
type TablelistColumn struct {
	Align string
	// background
	// changesnipside
	// changetitlesnipside
	// editable
	// editwindow
	// font
	// foreground
	// formatcommand
	// hide
	// labelalign
	// labelbackground
	// labelborderwidth
	// labelcommand
	// labelcommand2
	// labelfont
	// labelforeground
	// labelheight
	// labelpady
	// labelrelief
	// selectfiltercommand
	// labelimage
	// labelvalign
	// labelwindow
	// maxwidth
	// name
	// resizable
	// selectbackground
	// selectforeground
	// showarrow
	// showlinenumbers
	// sortcommand
	// sortmode
	// stretchable
	// stretchwindow
	// stripebackground
	// stripeforeground
	// text
	Title string
	// valign
	Width int
	// windowdestroy
	// windowupdate
	// wrap
}

type TablelistSelectMode string

const (
	TablelistSelectSingle   TablelistSelectMode = "single"
	TablelistSelectBrowse   TablelistSelectMode = "browse"
	TablelistSelectMultiple TablelistSelectMode = "multiple"
	TablelistSelectExtended TablelistSelectMode = "extended"
)

type Tablelist struct {
	BaseWidget
	xscrollcommand *CommandEx
	yscrollcommand *CommandEx
}

func NewTablelist(parent Widget, attributes ...*WidgetAttr) *Tablelist {
	theme := checkInitUseTheme(attributes)
	iid := makeNamedWidgetId(parent, "atk_tablelist")
	info := CreateWidgetInfo(iid, WidgetTypeTablelist, theme, attributes)
	if info == nil {
		return nil
	}
	w := &Tablelist{}
	w.id = iid
	w.info = info
	RegisterWidget(w)
	return w
}

// "Specifies a Tcl command to be invoked when mouse button 1 is pressed over one of the header labels and later released over the same label."
func (w *Tablelist) SetLabelCommand(cmd string) {
	eval(fmt.Sprintf("%v configure -labelcommand %s", w.id, cmd))
}

// "this command sorts the items based on the column whose index was passed to it as second argument."
func (w *Tablelist) SetLabelCommandSortByColumn() {
	w.SetLabelCommand("tablelist::sortByColumn")
}

// "Specifies a Tcl command to be invoked when mouse button 1 is pressed together with the Shift key over one of the header labels and later released over the same label."
func (w *Tablelist) SetLabelCommand2(cmd string) {
	eval(fmt.Sprintf("%v configure -labelcommand2 %s", w.id, cmd))
}

// "this command adds the column index passed to it as second argument to the list of sort columns and sorts the items based on the columns indicated by the modified list"
func (w *Tablelist) SetLabelCommand2AddToSortColumns() {
	w.SetLabelCommand2("tablelist::addToSortColumns")
}

func (w *Tablelist) Attach(id string) error {
	info, err := CheckWidgetInfo(id, WidgetTypeTablelist)
	if err != nil {
		return err
	}
	w.id = id
	w.info = info
	RegisterWidget(w)
	return nil
}

func (w *Tablelist) SetTakeFocus(takefocus bool) error {
	return eval(fmt.Sprintf("%v configure -takefocus {%v}", w.id, boolToInt(takefocus)))
}

func (w *Tablelist) IsTakeFocus() bool {
	r, _ := evalAsBool(fmt.Sprintf("%v cget -takefocus", w.id))
	return r
}

func (w *Tablelist) SetHeight(row int) error {
	return eval(fmt.Sprintf("%v configure -height {%v}", w.id, row))
}

func (w *Tablelist) Height() int {
	r, _ := evalAsInt(fmt.Sprintf("%v cget -height", w.id))
	return r
}

func (w *Tablelist) RowCget(idx int, option string) (string, error) {
	return evalAsString(fmt.Sprintf("%v rowcget %v %v", w.id, idx, option))
}

func (w *Tablelist) RowConfigure(idx string, key string, value string) error {
	return eval(fmt.Sprintf("%v rowconfigure %v %v %v", w.id, idx, key, value))
}

/*

func (w *Tablelist) SetPaddingN(padx int, pady int) error {
	if w.info.IsTtk {
		return eval(fmt.Sprintf("%v configure -padding {%v %v}", w.id, padx, pady))
	}
	return eval(fmt.Sprintf("%v configure -padx {%v} -pady {%v}", w.id, padx, pady))
}


	func (w *Tablelist) PaddingN() (int, int) {
		var r string
		var err error
		if w.info.IsTtk {
			r, err = evalAsString(fmt.Sprintf("%v cget -padding", w.id))
		} else {
			r1, _ := evalAsString(fmt.Sprintf("%v cget -padx", w.id))
			r2, _ := evalAsString(fmt.Sprintf("%v cget -pady", w.id))
			r = r1 + " " + r2
		}
		return parserPaddingResult(r, err)
	}

	func (w *Tablelist) SetPadding(pad Pad) error {
		return w.SetPaddingN(pad.X, pad.Y)
	}

	func (w *Tablelist) Padding() Pad {
		x, y := w.PaddingN()
		return Pad{x, y}
	}
*/

func (w *Tablelist) SetSelectMode(mode TablelistSelectMode) error {
	return eval(fmt.Sprintf("%v configure -selectmode {%v}", w.id, mode))
}

func (w *Tablelist) SelectMode() string { //TablelistSelectMode { // todo
	r, _ := evalAsString(fmt.Sprintf("%v cget -selectmode", w.id))
	//return parserTreeSebectModeResult(r, err)
	return r
}

/*
untested

	func (w *Tablelist) SetHeaderHidden(hide bool) error {
		var value string
		if hide {
			value = "tree"
		} else {
			value = "tree headings"
		}
		return eval(fmt.Sprintf("%v configure -show {%v}", w.id, value))
	}

	func (w *Tablelist) IsHeaderHidden() bool {
		r, _ := evalAsString(fmt.Sprintf("%v cget -show", w.id))
		return r == "tree"
	}
*/

func (w *Tablelist) DeleteAllColumns() error {
	return eval(fmt.Sprintf("%v deletecolumns 0 end", w.id))
}

// "Inserts the columns specified by the list columnList just before the column given by columnIndex"
func (w *Tablelist) InsertColumnList(index interface{}, columns []TablelistColumn) error {
	if len(columns) == 0 {
		return nil
	}
	var index_str string
	switch t := index.(type) {
	case int:
		// todo: enum valid vals
		index_str = strconv.Itoa(t)
	case string:
		index_str = t
	default:
		panic("programming error, InsertColumnList received unsupported type for 'index'. index must be an int or a string")
	}
	var tlc_strings []string
	for _, tlc := range columns {
		tlc_strings = append(tlc_strings, fmt.Sprintf(`%v %v %s`, tlc.Width, Quote(tlc.Title), tlc.Align))
	}
	return eval(fmt.Sprintf("%v insertcolumnlist %v {%v}", w.id, index_str, strings.Join(tlc_strings, " ")))
}

func (w *Tablelist) ColumnCount() int {
	num, _ := evalAsInt(fmt.Sprintf("%v columncount", w.id))
	return num

}

/*
	func (w *Tablelist) SetHeaderLabels(labels []string) error {
		for n, label := range labels {
			setObjText("atk_heading_label", label)
			err := eval(fmt.Sprintf("%v heading #%v -text $atk_heading_label", w.id, n))
			if err != nil {
				return err
			}
		}
		return nil
	}

	func (w *Tablelist) SetHeaderLabel(column int, label string) error {
		setObjText("atk_heading_label", label)
		return eval(fmt.Sprintf("%v heading #%v -text $atk_heading_label", w.id, column))
	}

	func (w *Tablelist) HeaderLabel(column int) string {
		r, _ := evalAsString(fmt.Sprintf("%v heading #%v -text", w.id, column))
		return r
	}

	func (w *Tablelist) SetHeaderImage(column int, img *Image) error {
		var iid string
		if img != nil {
			iid = img.Id()
		}
		return eval(fmt.Sprintf("%v heading #%v -image {%v}", w.id, column, iid))
	}

	func (w *Tablelist) HeaderImage(column int) *Image {
		r, err := evalAsString(fmt.Sprintf("%v heading #%v -image", w.id, column))
		return parserImageResult(r, err)
	}

	func (w *Tablelist) SetHeaderAnchor(column int, anchor Anchor) error {
		return eval(fmt.Sprintf("%v heading #%v -anchor %v", w.id, column, anchor))
	}

	func (w *Tablelist) HeaderAnchor(column int) Anchor {
		r, err := evalAsString(fmt.Sprintf("%v heading #%v -anchor", w.id, column))
		return parserAnchorResult(r, err)
	}
*/

/*
	func (w *Tablelist) SetColumnWidth(column int, width int) error {
		return eval(fmt.Sprintf("%v column #%v -width %v", w.id, column, width))
	}

	func (w *Tablelist) ColumnWidth(column int) int {
		r, _ := evalAsInt(fmt.Sprintf("%v column #%v -width", w.id, column))
		return r
	}

	func (w *Tablelist) SetColumnMinimumWidth(column int, width int) error {
		return eval(fmt.Sprintf("%v column #%v -minwidth %v", w.id, column, width))
	}

	func (w *Tablelist) ColumnMinimumWidth(column int) int {
		r, _ := evalAsInt(fmt.Sprintf("%v column #%v -minwidth", w.id, column))
		return r
	}

	func (w *Tablelist) SetColumnAnchor(column int, anchor Anchor) error {
		return eval(fmt.Sprintf("%v column #%v -anchor %v", w.id, column, anchor))
	}

	func (w *Tablelist) ColumnAnchor(column int) Anchor {
		r, err := evalAsString(fmt.Sprintf("%v column #%v -anchor", w.id, column))
		return parserAnchorResult(r, err)
	}

// default all column stretch 1

	func (w *Tablelist) SetColumnStretch(column int, stretch bool) error {
		return eval(fmt.Sprintf("%v column #%v -stretch %v", w.id, column, stretch))
	}

// default all column stretch 1

	func (w *Tablelist) ColumnStretch(column int) bool {
		r, _ := evalAsBool(fmt.Sprintf("%v column #%v -stretch", w.id, column))
		return r
	}
*/

func (w *Tablelist) IsValidItem(item *TablelistItem) bool {
	return item != nil && item.tablelist != nil && item.tablelist.id == w.id
}

/*
func (w *Tablelist) RootItem() *TablelistItem {
	return &TablelistItem{w, ""}
}

func (w *Tablelist) ToplevelItems() []*TablelistItem {
	return w.RootItem().Children()
}
*/
// "Inserts zero or more new items in the widget's internal list just before the item given by index"
func (w *Tablelist) Insert(index int, item_list [][]string) []*TablelistItem {
	item_strings := []string{}
	for _, val_list := range item_list {
		cell_list := []string{}
		for _, val := range val_list {
			cell_list = append(cell_list, fmt.Sprintf("%v", Quote(val)))
		}
		item_strings = append(item_strings, fmt.Sprintf(`{%v}`, strings.Join(cell_list, " ")))
	}
	id_list, err := evalAsStringList(fmt.Sprintf("%v insert 0 %v", w.id, strings.Join(item_strings, " ")))
	if err != nil {
		return nil
	}

	// todo: could this ever happen?
	if len(id_list) != len(item_list) {
		dumpError(errors.New("number of ids returned does not match the number of rows given"))
	}

	result := []*TablelistItem{}
	for _, id := range id_list {
		result = append(result, &TablelistItem{w, id})
	}

	return result
}

func (w *Tablelist) InsertSingle(index int, item_list []string) []*TablelistItem {
	return w.Insert(index, [][]string{item_list})
}

func (w *Tablelist) InsertChildren(pidx interface{}, cidx int, item_list [][]string) ([]string, error) {
	var pidx_str string

	switch t := pidx.(type) {
	case string:
		// todo: ensure valid string value, i.e. 'root', 'end', etc.
		pidx_str = t
	case int:
		pidx_str = strconv.Itoa(t)
	default:
		panic("programming error, InsertChildren received unsupported type for 'pidx'. pidx must be an int or a string")
	}

	child_list := []string{}
	for _, item := range item_list {
		child_row := []string{}
		for _, cell := range item {
			child_row = append(child_row, fmt.Sprintf("%v", Quote(cell)))
		}
		child_list = append(child_list, fmt.Sprintf("{%v}", strings.Join(child_row, " ")))
	}

	return evalAsStringList(fmt.Sprintf("%v insertchildren %v %v %v", w.id, pidx_str, cidx, strings.Join(child_list, " ")))
}

func (w *Tablelist) InsertChildList(pidx interface{}, cidx int, item_list [][]string) ([]string, error) {
	var pidx_str string

	switch t := pidx.(type) {
	case string:
		// todo: ensure valid string value, i.e. 'root', 'end', etc.
		pidx_str = t
	case int:
		pidx_str = strconv.Itoa(t)
	default:
		panic("programming error, InsertChildList received unsupported type for 'pidx'. pidx must be an int or a string")
	}

	item_strings := []string{}
	for _, val_list := range item_list {
		cell_list := []string{}
		for _, val := range val_list {
			cell_list = append(cell_list, fmt.Sprintf("%v", Quote(val)))
		}
		item_strings = append(item_strings, fmt.Sprintf(`{%v}`, strings.Join(cell_list, " ")))
	}
	return evalAsStringList(fmt.Sprintf("%v insertchildlist %v %v {%v}", w.id, pidx_str, cidx, strings.Join(item_strings, " ")))
}

/*
	func (w *Tablelist) InsertChildList(parent *TablelistItem, index int, text string, values []string) *TablelistItem {
		var pid string
		if parent != nil {
			if !w.IsValidItem(parent) {
				return nil
			}
			pid = parent.id
		}
		//setObjText("atk_tablelist_item", text)
		//setObjTextList("atk_tablelist_values", values)
		cid := makeTablelistItemId(w.id, pid)

		err := eval(fmt.Sprintf("%v insert 0 {%v} ", w.id, strings.Join(values, " ")))
		if err != nil {
			return nil
		}
		return &TablelistItem{w, cid}
	}
*/

func (w *Tablelist) DeleteItem(item *TablelistItem) error {
	if !w.IsValidItem(item) || item.IsRoot() {
		return ErrInvalid
	}
	return eval(fmt.Sprintf("%v delete {%v}", w.id, item.id))
}

func (w *Tablelist) DeleteAllItems() error {
	return eval(fmt.Sprintf("%v delete 0 end", w.id))
}

func (w *Tablelist) MovableColumns(b bool) error {
	return eval(fmt.Sprintf("%v configure -movablecolumns %v", w.id, b))
}

func (w *Tablelist) MoveItem(item *TablelistItem, parent *TablelistItem, index int) error {
	if !w.IsValidItem(item) || item.IsRoot() {
		return ErrInvalid
	}
	var pid string
	if parent != nil {
		if !w.IsValidItem(parent) {
			return ErrInvalid
		}
		pid = parent.id
	}
	return eval(fmt.Sprintf("%v move {%v} {%v} %v", w.id, item.id, pid, index))
}

func (w *Tablelist) RefreshSorting(parentNodeIndex int) error {
	return eval(fmt.Sprintf("%v refreshsorting %v", w.id, parentNodeIndex))
}

func (w *Tablelist) ScrollTo(item *TablelistItem) error {
	if !w.IsValidItem(item) || item.IsRoot() {
		return ErrInvalid
	}
	/*
		children := w.RootItem().Children()
		if len(children) == 0 {
			return ErrInvalid
		}
		//fix see bug: first scroll to root
		eval(fmt.Sprintf("%v see %v", w.id, children[0].id))
		return eval(fmt.Sprintf("%v see %v", w.id, item.id))
	*/
	panic("not implemented")
}

func (w *Tablelist) CurrentIndex() *TablelistItem {
	lst := w.SelectionList()
	if len(lst) == 0 {
		return nil
	}
	return lst[0]
}

func (w *Tablelist) SetCurrentIndex(item *TablelistItem) error {
	return w.SetSelections(item)
}

func (w *Tablelist) SelectionList() (lst []*TablelistItem) {
	ids, err := evalAsStringList(fmt.Sprintf("%v selection", w.id))
	if err != nil {
		return
	}
	for _, id := range ids {
		lst = append(lst, &TablelistItem{w, id})
	}
	return lst
}

func (w *Tablelist) SetSelections(items ...*TablelistItem) error {
	return w.SetSelectionList(items)
}

func (w *Tablelist) RemoveSelections(items ...*TablelistItem) error {
	return w.RemoveSelectionList(items)
}

func (w *Tablelist) AddSelections(items ...*TablelistItem) error {
	return w.AddSelectionList(items)
}

func (w *Tablelist) ToggleSelections(items ...*TablelistItem) error {
	return w.ToggleSelectionList(items)
}

func (w *Tablelist) SetSelectionList(items []*TablelistItem) error {
	var ids []string
	for _, item := range items {
		if w.IsValidItem(item) && !item.IsRoot() {
			ids = append(ids, item.id)
		}
	}
	if len(ids) == 0 {
		return ErrInvalid
	}
	return eval(fmt.Sprintf("%v selection set {%v}", w.id, strings.Join(ids, " ")))
}

func (w *Tablelist) RemoveSelectionList(items []*TablelistItem) error {
	var ids []string
	for _, item := range items {
		if w.IsValidItem(item) && !item.IsRoot() {
			ids = append(ids, item.id)
		}
	}
	if len(ids) == 0 {
		return ErrInvalid
	}
	return eval(fmt.Sprintf("%v selection remove {%v}", w.id, strings.Join(ids, " ")))
}

func (w *Tablelist) AddSelectionList(items []*TablelistItem) error {
	var ids []string
	for _, item := range items {
		if w.IsValidItem(item) && !item.IsRoot() {
			ids = append(ids, item.id)
		}
	}
	if len(ids) == 0 {
		return ErrInvalid
	}
	return eval(fmt.Sprintf("%v selection add {%v}", w.id, strings.Join(ids, " ")))
}

func (w *Tablelist) ToggleSelectionList(items []*TablelistItem) error {
	var ids []string
	for _, item := range items {
		if w.IsValidItem(item) && !item.IsRoot() {
			ids = append(ids, item.id)
		}
	}
	if len(ids) == 0 {
		return ErrInvalid
	}
	return eval(fmt.Sprintf("%v selection toggle {%v}", w.id, strings.Join(ids, " ")))
}

// "... expands all top-level rows of a tablelist used as a tree widget, i.e., makes all their children visible."
// "... the command will be performed recursively, i.e., all of the descendants of the top-level nodes will be displayed."
func (w *Tablelist) ExpandAll() error {
	return eval(fmt.Sprintf("%v expandall -fully", w.id))
}

// "... expands all top-level rows of a tablelist used as a tree widget, i.e., makes all their children visible."
// "... restricts the operation to just one hierarchy level, indicating that only the children of the top-level
// nodes will be displayed, without changing the expanded/collapsed state of the child nodes."
func (w *Tablelist) ExpandAllPartly() error {
	return eval(fmt.Sprintf("%v expandall -partly", w.id))
}

// "... collapses all top-level rows of a tablelist used as a tree widget, i.e., elides all their descendants."
// "... the command will be performed recursively, i.e., all of the descendants of the top-level nodes will be collapsed ..."
func (w *Tablelist) CollapseAll() error {
	return eval(fmt.Sprintf("%v collapseall -fully", w.id))
}

// "... collapses all top-level rows of a tablelist used as a tree widget, i.e., elides all their descendants."
// "... restricts the operation to just one hierarchy level ..."
func (w *Tablelist) CollapseAllPartly() error {
	return eval(fmt.Sprintf("%v collapseall -partly", w.id))
}

func (w *Tablelist) Expand(item *TablelistItem) error {
	return w.SetExpanded(item, true)
}

// [does not work]
// "Returns the list of full keys of the expanded items."
func (w *Tablelist) ExpandedKeys() []int {
	key_list, err := evalAsIntList(fmt.Sprintf("%v expandedkeys", w.id))
	dumpError(err)
	return key_list
}

// [does not work]
// "returns a list ... of items between firstIndex and lastIndex, inclusive."
// "The [mod] argument can be used to restrict the items when building the result list.
// The default value -all means: no restriction.
// The value -nonhidden filters out the items whose row has its -hide option set to true.
// Finally, the value -viewable restricts the items to the viewable ones.
func (w *Tablelist) GetKeys(fidx, lidx int, mod string) (string, error) {
	return evalAsString(fmt.Sprintf("%v getkeys %v %v %v", w.id, fidx, lidx, mod))
}

func (w *Tablelist) Collapse(item *TablelistItem) error {
	return w.SetExpanded(item, false)
}

func (w *Tablelist) SetExpanded(item *TablelistItem, expand bool) error {
	if !w.IsValidItem(item) || item.IsRoot() {
		return ErrInvalid
	}
	return item.SetExpanded(expand)
}

func (w *Tablelist) IsExpanded(item *TablelistItem) bool {
	if !w.IsValidItem(item) || item.IsRoot() {
		return false
	}
	return item.IsExpanded()
}

func (w *Tablelist) FocusItem() *TablelistItem {
	r, _ := evalAsString(fmt.Sprintf("%v focus", w.id))
	if r == "" {
		return nil
	}
	return &TablelistItem{w, r}
}

func (w *Tablelist) SetFocusItem(item *TablelistItem) error {
	if !w.IsValidItem(item) || item.IsRoot() {
		return ErrInvalid
	}
	return eval(fmt.Sprintf("%v focus %v", w.id, item.id))
}

func (w *Tablelist) OnSelectionChanged(fn func()) error {
	if fn == nil {
		return ErrInvalid
	}
	return w.BindEvent("<<TablelistSelect>>", func(e *Event) {
		fn()
	})
}

type TL_ROW_STATE string

const (
	TL_ROW_ALL       TL_ROW_STATE = "-all"
	TL_ROW_NONHIDDEN TL_ROW_STATE = "-nonhidden"
	TL_ROW_VIEWABLE  TL_ROW_STATE = "-viewable"
)

// "Returns a list containing the numerical indices of all of the items in the tablelist that contain at least one selected element."
func (w *Tablelist) CurSelectionWithState(state TL_ROW_STATE) []int {
	idx_list, _ := evalAsIntList(fmt.Sprintf("%v curselection %v", w.id, state))
	return idx_list
}

// "Returns a list containing the numerical indices of all of the items in the tablelist that contain at least one selected element."
func (w *Tablelist) CurSelection() []int {
	return w.CurSelectionWithState(TL_ROW_ALL)
}

/*
func (w *Tablelist) OnItemExpanded(fn func()) error {
	if fn == nil {
		return ErrInvalid
	}
	return w.BindEvent("<<TreeviewOpen>>", func(e *Event) {
		fn()
	})
}
*/

/*
func (w *Tablelist) OnItemCollapsed(fn func()) error {
	if fn == nil {
		return ErrInvalid
	}
	return w.BindEvent("<<TreeviewSelect>>", func(e *Event) {
		fn()
	})
}
*/

/*
func (w *Tablelist) ItemAt(x int, y int) *TablelistItem {
	id, err := evalAsString(fmt.Sprintf("%v identify item %v %v", w.id, x, y))
	if err != nil {
		return nil
	}
	if id == "" {
		return nil
	}
	return &TablelistItem{w, id}
}
*/

func (w *Tablelist) ItemAt(x int, y int) *TablelistItem {
	idx, err := evalAsString(fmt.Sprintf("%v nearestcell %v %v", w.id, x, y))
	if err != nil {
		return nil
	}
	if idx == "" {
		return nil
	}
	id, _ := evalAsString(fmt.Sprintf("%v get %v", w.id, idx))
	return &TablelistItem{w, id}
}

func (w *Tablelist) OnDoubleClickedItem(fn func(item *TablelistItem)) {
	if fn == nil {
		return
	}
	w.BindEvent("<Double-ButtonPress-1>", func(e *Event) {
		item := w.ItemAt(e.PosX, e.PosY)
		fn(item)
	})
}

func (w *Tablelist) SetXViewArgs(args []string) error {
	return eval(fmt.Sprintf("%v xview %v", w.id, strings.Join(args, " ")))
}

func (w *Tablelist) SetYViewArgs(args []string) error {
	return eval(fmt.Sprintf("%v yview %v", w.id, strings.Join(args, " ")))
}

func (w *Tablelist) OnXScrollEx(fn func([]string) error) error {
	if fn == nil {
		return ErrInvalid
	}
	if w.xscrollcommand == nil {
		w.xscrollcommand = &CommandEx{}
		bindCommandEx(w.id, "xscrollcommand", w.xscrollcommand.Invoke)
	}
	w.xscrollcommand.Bind(fn)
	return nil
}

func (w *Tablelist) OnYScrollEx(fn func([]string) error) error {
	if fn == nil {
		return ErrInvalid
	}
	if w.yscrollcommand == nil {
		w.yscrollcommand = &CommandEx{}
		bindCommandEx(w.id, "yscrollcommand", w.yscrollcommand.Invoke)
	}
	w.yscrollcommand.Bind(fn)
	return nil
}

func (w *Tablelist) BindXScrollBar(bar *ScrollBar) error {
	if !IsValidWidget(bar) {
		return ErrInvalid
	}
	w.OnXScrollEx(bar.SetScrollArgs)
	bar.OnCommandEx(w.SetXViewArgs)
	return nil
}

func (w *Tablelist) BindYScrollBar(bar *ScrollBar) error {
	if !IsValidWidget(bar) {
		return ErrInvalid
	}
	w.OnYScrollEx(bar.SetScrollArgs)
	bar.OnCommandEx(w.SetYViewArgs)
	return nil
}

type TablelistEx struct {
	*ScrollLayout
	*Tablelist
}

// a composite widget consisting of `Tablelist` and scrollbars bound to scroll events,
// wrapped in a `ScrollLayout`.
func NewTablelistEx(parent Widget, attributes ...*WidgetAttr) *TablelistEx {
	w := &TablelistEx{}
	w.ScrollLayout = NewScrollLayout(parent)
	w.Tablelist = NewTablelist(parent, attributes...)
	w.SetWidget(w.Tablelist)
	w.Tablelist.BindXScrollBar(w.XScrollBar)
	w.Tablelist.BindYScrollBar(w.YScrollBar)
	RegisterWidget(w)
	return w
}

func TablelistAttrTakeFocus(takefocus bool) *WidgetAttr {
	return &WidgetAttr{"takefocus", boolToInt(takefocus)}
}

func TablelistAttrHeight(row int) *WidgetAttr {
	return &WidgetAttr{"height", row}
}

func TablelistAttrPadding(padding Pad) *WidgetAttr {
	return &WidgetAttr{"padding", padding}
}

func TablelistAttrTreeSelectMode(mode TreeSelectMode) *WidgetAttr {
	return &WidgetAttr{"selectmode", mode}
}
