package tk

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type TABLELIST_SELECT_MODE string

const (
	TABLELIST_SELECT_MODE_SINGLE   TABLELIST_SELECT_MODE = "single"
	TABLELIST_SELECT_MODE_BROWSE   TABLELIST_SELECT_MODE = "browse"
	TABLELIST_SELECT_MODE_MULTIPLE TABLELIST_SELECT_MODE = "multiple"
	TABLELIST_SELECT_MODE_EXTENDED TABLELIST_SELECT_MODE = "extended"
)

var TABLELIST_SELECT_MODE_SET = map[TABLELIST_SELECT_MODE]bool{
	TABLELIST_SELECT_MODE_SINGLE:   true,
	TABLELIST_SELECT_MODE_BROWSE:   true,
	TABLELIST_SELECT_MODE_MULTIPLE: true,
	TABLELIST_SELECT_MODE_EXTENDED: true,
}

type TABLELIST_ROW_STATE string

const (
	TABLELIST_ROW_STATE_ALL       TABLELIST_ROW_STATE = "-all"
	TABLELIST_ROW_STATE_NONHIDDEN TABLELIST_ROW_STATE = "-nonhidden"
	TABLELIST_ROW_STATE_VIEWABLE  TABLELIST_ROW_STATE = "-viewable"
)

type TABLELIST_HALIGN string

const (
	TABLELIST_HALIGN_LEFT   TABLELIST_HALIGN = "left"
	TABLELIST_HALIGN_RIGHT  TABLELIST_HALIGN = "right"
	TABLELIST_HALIGN_CENTER TABLELIST_HALIGN = "center"
)

type TABLELIST_VALIGN string

const (
	TABLELIST_VALIGN_CENTER TABLELIST_VALIGN = "center"
	TABLELIST_VALIGN_TOP    TABLELIST_VALIGN = "top"
	TABLELIST_VALIGN_BOTTOM TABLELIST_VALIGN = "bottom"
)

type TABLELIST_RELIEF string
const (
	TABLELIST_RELIEF_RAISED TABLELIST_RELIEF = "raised"
	TABLELIST_RELIEF_SUNKEN TABLELIST_RELIEF = "sunken"
	TABLELIST_RELIEF_FLAT   TABLELIST_RELIEF = "flat"
	TABLELIST_RELIEF_RIDGE  TABLELIST_RELIEF = "ridge"
	TABLELIST_RELIEF_SOLID  TABLELIST_RELIEF = "solid"
	TABLELIST_RELIEF_GROOVE TABLELIST_RELIEF = "groove"
)

type TABLELIST_SORT_MODE string

const (
	TABLELIST_SORT_MODE_ASCII         TABLELIST_SORT_MODE = "ascii"
	TABLELIST_SORT_MODE_ASCII_NO_CASE TABLELIST_SORT_MODE = "asciinocase"
	TABLELIST_SORT_MODE_DICTIONARY    TABLELIST_SORT_MODE = "dictionary"
	TABLELIST_SORT_MODE_INTEGER       TABLELIST_SORT_MODE = "integer"
	TABLELIST_SORT_MODE_REAL          TABLELIST_SORT_MODE = "real"
	//TABLELIST_SORT_MODE_COMMAND     TABLELIST_SORT_MODE = "command"
)

// https://www.nemethi.de/tablelist/tablelistWidget.html#col_options
type TablelistColumn struct {
	Align               TABLELIST_HALIGN
	Background          string
	ChangeSnipSide      bool
	ChangeTitleSnipSide bool
	Editable            bool
	EditWindow          string
	Font                string
	Foreground          string
	//FormatCommand       string
	Hide             bool
	LabelAlign       TABLELIST_HALIGN
	LabelBackground  string
	LabelBorderWidth string // int? float?
	//LabelCommand        string
	//LabelCommand2       string
	LabelFont       string
	LabelForeground string
	LabelHeight     int
	LabelPadY       int
	LabelRelief     TABLELIST_RELIEF
	// selectfiltercommand
	// LabelImage string
	LabelVAlign      TABLELIST_VALIGN
	LabelWindow      string
	MaxWidth         int
	Name             string
	Resizable        bool
	SelectBackground string
	SelectForeground string
	ShowArrow        bool
	ShowLineNumbers  bool
	// sortcommand
	SortMode         TABLELIST_SORT_MODE
	Stretchable      bool
	StretchWindow    bool
	StripeBackground string
	StripeForeground string
	// text
	Title  string
	VAlign TABLELIST_VALIGN
	Width  int
	// windowdestroy
	// windowupdate
	Wrap bool
}

func NewTablelistColumn() *TablelistColumn {
	return &TablelistColumn{
		Align:       TABLELIST_HALIGN_LEFT,
		EditWindow:  "entry",
		LabelAlign:  TABLELIST_HALIGN_LEFT,
		LabelRelief: TABLELIST_RELIEF_RAISED, // todo: better default?
		LabelVAlign: TABLELIST_VALIGN_CENTER,
		MaxWidth:    0, // auto size. positive ints are measured in chars (???) and negative ints are pixels
		Resizable:   true,
		ShowArrow:   true,
		SortMode:    TABLELIST_SORT_MODE_ASCII,
		VAlign:      TABLELIST_VALIGN_CENTER,
		Width:       0, // auto size. positive ints are measured in chars (???) and negative ints are pixels
	}
}

type Tablelist struct {
	BaseWidget
	xscrollcommand *CommandEx
	yscrollcommand *CommandEx
}

// --- utils

func int_list_to_string_list(int_list []int) []string {
	int_str_list := []string{}
	for _, idx_str := range int_list {
		int_str_list = append(int_str_list, strconv.Itoa(idx_str))
	}
	return int_str_list
}

func column_triple_string(column_list ...*TablelistColumn) string {
	csl := []string{}
	for _, column := range column_list {
		csl = append(csl, fmt.Sprintf(`%v %v %s`, column.Width, Quote(column.Title), column.Align))
	}
	return strings.Join(csl, " ")
}

// ---

func NewTablelist(parent Widget, attributes ...*WidgetAttr) *Tablelist {

	// "Specifies a Tcl command to be invoked when expanding a row of a tablelist used as a tree widget"
	attributes = append(attributes, &WidgetAttr{Key: "expandcommand", Value: "tablelistExpandCmd"})

	attributes = append(attributes, &WidgetAttr{Key: "collapsecommand", Value: "tablelistCollapseCmd"})

	attributes = append(attributes, &WidgetAttr{Key: "populatecommand", Value: "tablelistPopulateCmd"})

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

	// generates virtual events 'TablelistRowExpand' and 'TablelistRowCollapse'.
	// use `Tablelist.OnItemExpanded` and `Tablelist.OnItemCollapsed` to handle these events.
	// row index is available as `Event.UserData[0]`
	// todo: get key for idx, pass it as second parameter
	// todo: populate command not working
	expand_collapse_cmd_proc := `
proc tablelistExpandCmd {tbl row} {
    event generate $tbl <<TablelistRowExpand>> -when now -data $row
}
proc tablelistCollapseCmd {tbl row} {
    event generate $tbl <<TablelistRowCollapse>> -when now -data $row
}
proc tablelistPopulateCmd {tbl row} {
    event generate $tbl <<TablelistRowPopulate>> -when now -data $row
}
`
	err := eval(expand_collapse_cmd_proc)
	dumpError(err)

	return w
}

// --- NewTablelistEx

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

// ---

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

// --- WIDGET-SPECIFIC OPTIONS

/*
   -acceptchildcommand command
   -acceptdropcommand command
   -activestyle frame|none|underline
   -arrowcolor color
   -arrowdisabledcolor color
   -arrowstyle flat7x4|flat7x7|flat8x4|flat8x5|flat9x5|flat10x5|
               flat11x6|flat12x6|flat13x7|flat14x7|flat15x8|flat16x8|
               flatAngle7x4|flatAngle7x5|flatAngle9x5|flatAngle9x6|
               flatAngle10x6|flatAngle11x6|flatAngle13x7|flatAngle15x8|
               photo0x0|photo7x4|photo7x7|photo9x5|photo11x6|
               photo13x7|photo15x8|sunken8x7|sunken10x9|sunken12x11
   -autofinishediting boolean
   -autoscan boolean
   -collapsecommand command
   -colorizecommand command
   -columns {width title ?left|right|center? width title ?left|right|center? ...}
   -columntitles {title title ...}
   -customdragsource boolean
   -displayondemand boolean
   -editendcommand command
   -editendonfocusout boolean
   -editendonmodclick boolean
   -editselectedonly boolean
   -editstartcommand command
   -expandcommand command
   -forceeditendcommand boolean
   -fullseparators boolean
   -height units
   -incrarrowtype up|down
   -instanttoggle boolean
   -itembackground color  or  -itembg color
   -labelactivebackground color
   -labelactiveforeground color
   -labelbackground color  or  -labelbg color
   -labelborderwidth screenDistance  or  -labelbd screenDistance
   (!) -labelcommand command
   (!) -labelcommand2 command
   -labeldisabledforeground color
   -labelfont font
   -labelforeground color  or  -labelfg color
   -labelheight lines
   -labelpady screenDistance
   -labelrelief raised|sunken|flat|ridge|solid|groove
   -listvariable variable
   -movablecolumns boolean
   -movablerows boolean
   -movecolumncursor cursor
   -movecursor cursor
   -populatecommand command
   -protecttitlecolumns boolean
   -resizablecolumns boolean
   -resizecursor cursor
   -selectfiltercommand command
   (!)-selectmode single|browse|multiple|extended
   -selecttype row|cell
   -setfocus boolean
   -showarrow boolean
   -showbusycursor boolean
   -showeditcursor boolean
   -showhorizseparator boolean
   -showlabels boolean
   -showseparators boolean
   -snipstring string
   -sortcommand command
   -spacing screenDistance
   -state normal|disabled
   -stretch all|columnIndexList
   -stripebackground color  or  -stripebg color
   -stripeforeground color  or  -stripefg color
   -stripeheight lines
   (!) -takefocus 0|1|""|command
   -targetcolor color
   -tight boolean
   -titlecolumns number
   -tooltipaddcommand command
   -tooltipdelcommand command
   -treecolumn columnIndex
   -treestyle adwaita|ambiance|aqua|aqua11|arc|baghira|bicolor|bicolor100|
              bicolor125|bicolor150|bicolor175|bicolor200|blueMenta|classic|
              classic100|classic125|classic150|classic175|classic200|dust|dustSand|
              gtk|klearlooks|mate|menta|mint|mint2|newWave|oxygen1|oxygen2|phase|
              plain|plain100|plain125|plain150|plain175|plain200|plastik|plastique|
              radiance|ubuntu|ubuntu2|ubuntu3|ubuntuMate|vistaAero|vistaClassic|
              white|white100|white125|white150|white175|white200|win7Aero|
              win7Classic|win10|winnative|winxpBlue|winxpOlive|winxpSilver|yuyo
   -width characters
   -xmousewheelwindow window
   -ymousewheelwindow window
*/

// "Specifies a Tcl command to be invoked when mouse button 1 is pressed over one of the header labels and later released over the same label."
func (w *Tablelist) LabelCommand(command string) error {
	return eval(fmt.Sprintf("%v configure -labelcommand %s", w.id, command))
}

// convenience
// "this command sorts the items based on the column whose index was passed to it as second argument."
func (w *Tablelist) LabelCommandSortByColumn() error {
	return w.LabelCommand("tablelist::sortByColumn")
}

// "Specifies a Tcl command to be invoked when mouse button 1 is pressed together with the Shift key over one of the header labels and later released over the same label."
func (w *Tablelist) LabelCommand2(command string) error {
	return eval(fmt.Sprintf("%v configure -labelcommand2 %s", w.id, command))
}

// convenience.
// "this command adds the column index passed to it as second argument to the list of sort columns and sorts the items based on the columns indicated by the modified list"
func (w *Tablelist) LabelCommand2AddToSortColumns() error {
	return w.LabelCommand2("tablelist::addToSortColumns")
}

func (w *Tablelist) SelectMode() TABLELIST_SELECT_MODE {
	mode, _ := evalAsString(fmt.Sprintf("%v cget -selectmode", w.id))
	mode_key := TABLELIST_SELECT_MODE(mode)
	_, present := TABLELIST_SELECT_MODE_SET[mode_key]
	if !present {
		dumpError(fmt.Errorf("unknown select mode returned: %s", mode))
	}
	return mode_key
}

func (w *Tablelist) SetSelectMode(mode TABLELIST_SELECT_MODE) error {
	return eval(fmt.Sprintf("%v configure -selectmode {%v}", w.id, mode))
}

func (w *Tablelist) TakeFocus() bool {
	r, _ := evalAsBool(fmt.Sprintf("%v cget -takefocus", w.id))
	return r
}

func (w *Tablelist) SetTakeFocus(takefocus string) error {
	return eval(fmt.Sprintf("%v configure -takefocus {%v}", w.id, takefocus))
}

// --- WIDGET COMMAND

/*

   pathName activate index
   pathName activatecell cellIndex
   pathName applysorting itemList
   pathName attrib ?name? ?value name value ...?
   pathName autoscrolltarget event x y
   pathName bbox index
   pathName bodypath
   pathName bodytag
   pathName canceledediting
   pathName cancelediting
   pathName cellattrib cellIndex ?name? ?value name value ...?
   pathName cellbbox cellIndex
   pathName cellcget cellIndex option
   pathName cellconfigure cellIndex ?option? ?value option value ...?
   pathName cellindex cellIndex
   pathName cellselection option args
   pathName cellselection anchor cellIndex
   pathName cellselection clear firstCell lastCell
   pathName cellselection clear cellIndexList
   pathName cellselection includes cellIndex
   pathName cellselection set firstCell lastCell
   pathName cellselection set cellIndexList
   pathName cget option
   pathName childcount nodeIndex
   pathName childindex index
   pathName childkeys nodeIndex
   pathName collapse indexList ?-fully|-partly?
   pathName collapseall ?-fully|-partly?
   pathName columnattrib columnIndex ?name? ?value name value ...?
   pathName columncget columnIndex option
   (!) pathName columnconfigure columnIndex ?option? ?value option value ...?
   (!) pathName columncount
   pathName columnindex columnIndex
   pathName columnwidth columnIndex ?-requested|-stretched|-total?
   pathName configcelllist {cellIndex option value cellIndex option value ...}
   pathName configcells ?cellIndex option value cellIndex option value ...?
   pathName configcolumnlist {columnIndex option value columnIndex option value ...}
   pathName configcolumns ?columnIndex option value columnIndex option value ...?
   pathName configrowlist {index option value index option value ...}
   pathName configrows ?index option value index option value ...?
   pathName configure ?option? ?value option value ...?
   pathName containing y
   pathName containingcell x y
   pathName containingcolumn x
   pathName cornerlabelpath
   pathName cornerpath ?-ne|-sw?
   pathName curcellselection ?-all|-nonhidden|-viewable?
   pathName curselection ?-all|-nonhidden|-viewable?
   (!) pathName delete firstIndex lastIndex
   (!) pathName delete indexList
   (!) pathName deletecolumns firstColumn lastColumn
   (!) pathName deletecolumns columnIndexList
   pathName depth nodeIndex
   pathName descendantcount nodeIndex
   pathName dicttoitem dictionary
   pathName dumptofile fileName
   pathName dumptostring
   pathName editcell cellIndex
   pathName editinfo
   pathName editwinpath
   pathName editwintag
   pathName embedcheckbutton cellIndex ?command?
   pathName embedcheckbuttons columnIndex ?command?
   pathName embedttkcheckbutton cellIndex ?command?
   pathName embedttkcheckbuttons columnIndex ?command?
   pathName entrypath
   pathName expand indexList ?-fully|-partly?
   pathName expandall ?-fully|-partly?
   pathName expandedkeys
   pathName fillcolumn columnIndex ?-text|-image|-window? value
   pathName findcolumnname name
   pathName findrowname name ?-descend? ?-parent nodeIndex?
   pathName finishediting
   pathName formatinfo
   pathName get firstIndex lastIndex ?-all|-nonhidden|-viewable?
   pathName get indexList
   pathName getcells firstCell lastCell ?-all|-nonhidden|-viewable?
   pathName getcells cellIndexList
   pathName getcolumns firstColumn lastColumn
   pathName getcolumns columnIndexList
   pathName getformatted firstIndex lastIndex ?-all|-nonhidden|-viewable?
   pathName getformatted indexList
   pathName getformattedcells firstCell lastCell ?-all|-nonhidden|-viewable?
   pathName getformattedcells cellIndexList
   pathName getformattedcolumns firstColumn lastColumn
   pathName getformattedcolumns columnIndexList
   pathName getfullkeys firstIndex lastIndex ?-all|-nonhidden|-viewable?
   pathName getfullkeys indexList
   pathName getkeys firstIndex lastIndex ?-all|-nonhidden|-viewable?
   pathName getkeys indexList
   pathName hasattrib name
   pathName hascellattrib cellIndex name
   pathName hascolumnattrib columnIndex name
   pathName hasrowattrib index name
   pathName header option ?arg arg ...?
   pathName headerpath
   pathName headertag
   pathName hidetargetmark
   pathName imagelabelpath cellIndex
   pathName index index
   pathName insert index ?item item ...?
   (!) pathName insertchildlist parentNodeIndex childIndex itemList
   (!) pathName insertchild(ren) parentNodeIndex childIndex ?item item ...?
   (!) pathName insertcolumnlist columnIndex {width title ?left|right|center? width title ?left|right|center? ...}
   (!) pathName insertcolumns columnIndex ?width title ?left|right|center? width title ?left|right|center? ...?
   pathName insertlist index itemList
   pathName iselemsnipped cellIndex fullTextName
   pathName isexpanded index
   pathName istitlesnipped columnIndex fullTextName
   pathName isviewable index
   pathName itemlistvar
   pathName itemtodict item ?excludedColumnIndexList?
   pathName labelpath columnIndex
   pathName labels
   pathName labeltag
   pathName labelwindowpath columnIndex
   pathName loadfromfile fileName ?-fully|-partly?
   pathName loadfromstring string ?-fully|-partly?
   pathName move sourceIndex targetIndex
   pathName move sourceIndex targetParentNodeIndex targetChildIndex
   pathName movecolumn sourceColumn targetColumn
   pathName nearest y
   pathName nearestcell x y
   pathName nearestcolumn x
   pathName noderow parentNodeIndex childIndex
   pathName parentkey nodeIndex
   pathName refreshsorting ?parentNodeIndex?
   pathName rejectinput
   pathName resetsortinfo
   pathName restorecursor
   pathName rowattrib index ?name? ?value name value ...?
   pathName rowcget index option
   pathName rowconfigure index ?option? ?value option value ...?
   pathName scan mark|dragto x y
   pathName searchcolumn columnIndex pattern ?options?
   pathName see index
   pathName seecell cellIndex
   pathName seecolumn columnIndex
   pathName selection option args
   pathName selection anchor index
   pathName selection clear firstIndex lastIndex
   pathName selection clear indexList
   pathName selection includes index
   pathName selection set firstIndex lastIndex
   pathName selection set indexList
   pathName separatorpath ?columnIndex?
   pathName separators
   pathName setbusycursor
   pathName showtargetmark before|inside index
   pathName size
   pathName sort ?-increasing|-decreasing?
   pathName sortbycolumn columnIndex ?-increasing|-decreasing?
   pathName sortbycolumnlist columnIndexList ?sortOrderList?
   pathName sortcolumn
   pathName sortcolumnlist
   pathName sortorder
   pathName sortorderlist
   pathName stopautoscroll
   pathName targetmarkpath
   pathName targetmarkpos y ?-any|-horizontal|-vertical?
   pathName togglecolumnhide firstColumn lastColumn
   pathName togglecolumnhide columnIndexList
   pathName togglerowhide firstIndex lastIndex
   pathName togglerowhide indexList
   pathName toplevelkey index
   pathName unsetattrib name
   pathName unsetcellattrib cellIndex name
   pathName unsetcolumnattrib columnIndex name
   pathName unsetrowattrib index name
   pathName viewablerowcount ?firstIndex lastIndex?
   pathName windowpath cellIndex
   pathName xview ?args?
   pathName xview
   pathName xview units
   pathName xview moveto fraction
   pathName xview scroll number units|pages
   pathName yview ?args?
   pathName yview
   pathName yview units
   pathName yview moveto fraction
   pathName yview scroll number units|pages


*/

func (w *Tablelist) ColumnCount() int {
	num, _ := evalAsInt(fmt.Sprintf("%v columncount", w.id))
	return num
}

func (w *Tablelist) ColumnConfigure(column_index string, attribute_list ...WidgetAttr) error {
	option_list := []string{}
	for _, attr := range attribute_list {
		option_list = append(option_list, fmt.Sprintf("-%v %v", attr.Key, attr.Value))
	}
	return eval(fmt.Sprintf("%v columnconfigure %v %v", w.id, column_index, strings.Join(option_list, " ")))
}

// similar to `ColumnConfigure` but introspects the fields and values on the given `column`,
// with special handling for certain fields and values,
// before issuing a `columnconfigure`.
func (w *Tablelist) ColumnConfigureEx(column_index string, column *TablelistColumn) error {
	option_list := []string{}
	column_as_value := reflect.ValueOf(*column)
	column_type := column_as_value.Type()
	for i := 0; i < column_as_value.NumField(); i++ {
		field := column_type.Field(i)
		field_value := column_as_value.Field(i).Interface()
		option := strings.ToLower(field.Name)

		if option == "labelheight" {
			// I'm getting: "bad option \"-labelheight\": must be -align, -background, -bg ... " ??
			continue
		}

		switch v := field_value.(type) {
		case string:
			if v == "" {
				// skip empty strings
				continue
			}
			if option == "title" {
				field_value = Quote(v)
			}
		}

		option_list = append(option_list, fmt.Sprintf("-%v %v", option, field_value))
	}
	return eval(fmt.Sprintf("%v columnconfigure %v %v", w.id, column_index, strings.Join(option_list, " ")))
}

func (w *Tablelist) DeleteColumns(first_column, last_column string) error {
	return eval(fmt.Sprintf("%v deletecolumns %v %v", w.id, first_column, last_column))
}

func (w *Tablelist) DeleteColumns2(column_index_list ...string) error {
	return eval(fmt.Sprintf("%v deletecolumns %v", w.id, strings.Join(column_index_list, " ")))
}

// convenience.
func (w *Tablelist) DeleteAllColumns() error {
	return eval(fmt.Sprintf("%v deletecolumns 0 end", w.id))
}

// "Inserts the columns specified by the list columnList just before the column given by columnIndex"
func (w *Tablelist) InsertColumnList(column_index string, column_list []*TablelistColumn) error {
	return eval(fmt.Sprintf("%v insertcolumnlist %v {%v}", w.id, column_index, column_triple_string(column_list...)))
}

func (w *Tablelist) InsertColumns(column_index string, column_list []*TablelistColumn) error {
	return eval(fmt.Sprintf("%v insertcolumns %v %v", w.id, column_index, column_triple_string(column_list...)))
}

// convenience. inserts the columns and then configures them.
func (w *Tablelist) InsertColumnsEx(column_index int, column_list []*TablelistColumn) {
	w.InsertColumns(strconv.Itoa(column_index), column_list)
	for offset, col := range column_list {
		w.ColumnConfigureEx(strconv.Itoa(column_index+offset), col)
	}
}

func (w *Tablelist) IsValidItem(item *TablelistItem) bool {
	return item != nil && item.tablelist != nil && item.tablelist.id == w.id
}

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
	full_key_list, _ := evalAsStringList(fmt.Sprintf("%v insert %v %v", w.id, index, strings.Join(item_strings, " ")))

	tablelist_item_list := []*TablelistItem{}
	for _, id := range full_key_list {
		tablelist_item_list = append(tablelist_item_list, NewTablelistItem("", id, w))
	}
	return tablelist_item_list
}

func (w *Tablelist) InsertSingle(index int, item_list []string) []*TablelistItem {
	return w.Insert(index, [][]string{item_list})
}

func (w *Tablelist) InsertChildren(parent_node_idx interface{}, child_index int, item_list [][]string) []string {
	var parent_node_idx_str string

	switch t := parent_node_idx.(type) {
	case string:
		// todo: ensure valid string value, i.e. 'root', 'end', etc.
		parent_node_idx_str = t
	case int:
		parent_node_idx_str = strconv.Itoa(t)
	default:
		panic("programming error, InsertChildren received unsupported type for 'parent_node_idx'. value must be an int or a string")
	}

	child_list := []string{}
	for _, item := range item_list {
		row_cell_list := []string{}
		for _, cell := range item {
			row_cell_list = append(row_cell_list, fmt.Sprintf("%v", Quote(cell)))
		}
		child_list = append(child_list, fmt.Sprintf("{%v}", strings.Join(row_cell_list, " ")))
	}
	full_key_list, _ := evalAsStringList(fmt.Sprintf("%v insertchildren %v %v %v", w.id, parent_node_idx_str, child_index, strings.Join(child_list, " ")))
	return full_key_list
}

// convenience. same as `InsertChildren` but returns a list of TablelistItem pointers.
func (w *Tablelist) InsertChildrenEx(parent_node_index interface{}, child_index int, item_list [][]string) ([]*TablelistItem, error) {
	full_key_list := w.InsertChildren(parent_node_index, child_index, item_list)
	tablelist_item_list := []*TablelistItem{}
	parent_full_key := "" // todo
	for _, full_key := range full_key_list {
		tablelist_item_list = append(tablelist_item_list, NewTablelistItem(parent_full_key, full_key, w))
	}
	return tablelist_item_list, nil

}

func (w *Tablelist) InsertChildList(parent_node_index interface{}, child_index int, item_list [][]string) []string {
	var pidx_str string

	switch t := parent_node_index.(type) {
	case string:
		// todo: ensure valid string value, i.e. 'root', 'end', etc.
		pidx_str = t
	case int:
		pidx_str = strconv.Itoa(t)
	default:
		panic("programming error, InsertChildList received unsupported type for 'pidx'. pidx must be an int or a string")
	}

	child_list := []string{}
	for _, item := range item_list {
		row_cell_list := []string{}
		for _, val := range item {
			row_cell_list = append(row_cell_list, fmt.Sprintf("%v", Quote(val)))
		}
		child_list = append(child_list, fmt.Sprintf(`{%v}`, strings.Join(row_cell_list, " ")))
	}
	key_list, _ := evalAsStringList(fmt.Sprintf("%v insertchildlist %v %v {%v}", w.id, pidx_str, child_index, strings.Join(child_list, " ")))
	return key_list
}

func (w *Tablelist) InsertChildListEx(pidx interface{}, cidx int, item_list [][]string) []*TablelistItem {
	key_list := w.InsertChildList(pidx, cidx, item_list)
	tablelist_item_list := []*TablelistItem{}
	parent_full_key := "" // todo
	for _, full_key := range key_list {
		tablelist_item_list = append(tablelist_item_list, NewTablelistItem(parent_full_key, full_key, w))
	}
	return tablelist_item_list
}

func (w *Tablelist) Delete(idx1, idx2 int) error {
	return eval(fmt.Sprintf("%v delete %v %v", w.id, idx1, idx2))
}

func (w *Tablelist) Delete2(idx ...int) error {
	return eval(fmt.Sprintf("%v delete %v", w.id, strings.Join(int_list_to_string_list(idx), " ")))
}

func (w *Tablelist) DeleteAllItems() error {
	return eval(fmt.Sprintf("%v delete 0 end", w.id))
}

func (w *Tablelist) MovableColumns(b bool) error {
	return eval(fmt.Sprintf("%v configure -movablecolumns %v", w.id, b))
}

func (w *Tablelist) RefreshSorting(parentNodeIndex int) error {
	return eval(fmt.Sprintf("%v refreshsorting %v", w.id, parentNodeIndex))
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

// [does not work]
// "Returns the list of full keys of the expanded items."
func (w *Tablelist) ExpandedKeys() []int {
	key_list, err := evalAsIntList(fmt.Sprintf("%v expandedkeys", w.id))
	dumpError(err)
	return key_list
}

func (w *Tablelist) GetKeys(idx1, idx2 int, option TABLELIST_ROW_STATE) []string {
	key_list, _ := evalAsStringList(fmt.Sprintf("%v getkeys %v %v %v", w.id, idx1, idx2, option))
	return key_list
}

func (w *Tablelist) GetKeys2(idx ...int) []string {
	idx_str_list := int_list_to_string_list(idx)
	key_list, _ := evalAsStringList(fmt.Sprintf("%v getkeys %v", w.id, strings.Join(idx_str_list, " ")))
	return key_list
}

func (w *Tablelist) RowCGet(idx int, option string) string {
	v, _ := evalAsString(fmt.Sprintf("%v rowcget %v %v", w.id, idx, option))
	return v
}

func (w *Tablelist) OnSelectionChanged(fn func()) error {
	if fn == nil {
		return ErrInvalid
	}
	return w.BindEvent("<<TablelistSelect>>", func(e *Event) {
		fn()
	})
}

// "... returns a list whose elements are all of the tablelist items (i.e., row contents) between firstIndex and lastIndex, inclusive."
func (w *Tablelist) Get(idx1, idx2 int, option TABLELIST_ROW_STATE) []string {
	rows, _ := evalAsStringList(fmt.Sprintf("%v get %v %v %v", w.id, idx1, idx2, option))
	return rows
}

func (w *Tablelist) Get2(idx ...int) []string {
	idx_str_list := int_list_to_string_list(idx)
	rows, _ := evalAsStringList(fmt.Sprintf("%v get %v", w.id, strings.Join(idx_str_list, " ")))
	return rows
}

// "Each item of a tablelist widget has a unique sequence number that remains unchanged until the item is deleted,
// thus acting as a key that uniquely identifies the item even if the latter's position (i.e., numerical row index) changes."
func (w *Tablelist) GetFullKeys(idx1, idx2 int, option TABLELIST_ROW_STATE) []string {
	key_list, _ := evalAsStringList(fmt.Sprintf("%v getfullkeys %v %v %v", w.id, idx1, idx2, option))
	return key_list
}

// "Each item of a tablelist widget has a unique sequence number that remains unchanged until the item is deleted,
// thus acting as a key that uniquely identifies the item even if the latter's position (i.e., numerical row index) changes."
func (w *Tablelist) GetFullKeys2(idx ...int) []string {
	idx_str_list := int_list_to_string_list(idx)
	key_list, _ := evalAsStringList(fmt.Sprintf("%v getfullkeys %v", w.id, strings.Join(idx_str_list, " ")))
	return key_list
}

// convenience.
// "Each item of a tablelist widget has a unique sequence number that remains unchanged until the item is deleted,
// thus acting as a key that uniquely identifies the item even if the latter's position (i.e., numerical row index) changes."
func (w *Tablelist) GetFullKey(idx int) (string, error) {
	key_list := w.GetFullKeys2(idx)
	if len(key_list) != 1 {
		return "", fmt.Errorf("key for row at index not found: %v", idx)
	}
	return key_list[0], nil
}

func (w *Tablelist) GetTablelistItemByIdx(idx int) *TablelistItem {
	full_key := w.GetFullKeys2(idx)
	return NewTablelistItem("", full_key[0], w)
}

func (w *Tablelist) OnItemExpanded(fn func(*TablelistItem)) error {
	event_fn := func(e *Event) {
		tli_idx, err := strconv.Atoi(e.UserData)
		dumpError(err)
		fn(w.GetTablelistItemByIdx(tli_idx))
	}
	return w.BindEvent("<<TablelistRowExpand>>", event_fn)
}

func (w *Tablelist) OnItemCollapsed(fn func(*TablelistItem)) error {
	event_fn := func(e *Event) {
		tli_idx, err := strconv.Atoi(e.UserData)
		dumpError(err)
		fn(w.GetTablelistItemByIdx(tli_idx))
	}
	return w.BindEvent("<<TablelistRowCollapse>>", event_fn)
}

// 'populate', singular.
func (w *Tablelist) OnItemPopulate(fn func(*Event)) error {
	return w.BindEvent("<<TablelistRowPopulate>>", fn)
}

// "Returns a list containing the numerical indices of all of the items in the tablelist that contain at least one selected element."
func (w *Tablelist) CurSelection(state TABLELIST_ROW_STATE) []int {
	idx_list, _ := evalAsIntList(fmt.Sprintf("%v curselection %v", w.id, state))
	return idx_list
}

// "Returns a list containing the numerical indices of all of the items in the tablelist that contain at least one selected element."
func (w *Tablelist) CurSelection2() []int {
	return w.CurSelection(TABLELIST_ROW_STATE_ALL)
}

func (w *Tablelist) ItemAt(x int, y int) *TablelistItem {
	idx, _ := evalAsString(fmt.Sprintf("%v nearestcell %v %v", w.id, x, y))
	if idx == "" {
		return nil
	}
	idx_int, err := strconv.Atoi(idx)
	dumpError(err)
	return w.GetTablelistItemByIdx(idx_int)
}

func (w *Tablelist) OnDoubleClickedItem(fn func(item *TablelistItem)) error {
	return w.BindEvent("<Double-ButtonPress-1>", func(e *Event) {
		item := w.ItemAt(e.PosX, e.PosY)
		fn(item)
	})
}
