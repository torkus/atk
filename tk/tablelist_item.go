// Copyright 2018 visualfc. All rights reserved.

package tk

import (
	"fmt"
	"strings"
)

type TablelistItem struct {
	tablelist *Tablelist
	pid       string // parent id, if any. typically "root" in flat lists.
	id        string // the 'fullkey' given to the item upon insertion.
}

func (t *TablelistItem) Id() string {
	return t.id
}

func (t *TablelistItem) IsValid() bool {
	return t != nil && t.tablelist != nil
}

func (t *TablelistItem) IsRoot() bool {
	return t.id == "root"
}

// ---

func makeTablelistItemId(treeid string, pid string, id string) string {
	return id
	/*
		if pid != "" && pid != "root" {
			return pid + "." + id // widj.foo.bar
		}
		return treeid + "." + id // widj.foo
	*/
}

func NewTablelistItem(parent_full_key string, full_key string, tablelist *Tablelist) *TablelistItem {
	return &TablelistItem{
		tablelist: tablelist,
		id:        makeTablelistItemId(tablelist.id, parent_full_key, full_key),
	}
}

// ---

// convenience
func (t *TablelistItem) Children() []*TablelistItem {
	ids, _ := evalAsStringList(fmt.Sprintf("%v childkeys {%v}", t.tablelist.id, t.id))
	lst := []*TablelistItem{}
	for _, id := range ids {
		lst = append(lst, NewTablelistItem(t.id, id, t.tablelist))
	}
	return lst
}

// --- https://www.nemethi.de/tablelist/tablelistWidget.html#row_options

func (t *TablelistItem) Background() string {
	sl, _ := evalAsString(fmt.Sprintf("%v rowcget {%v} -background", t.tablelist.id, t.id))
	return sl
}

func (t *TablelistItem) SetBackground(background string) error {
	return eval(fmt.Sprintf("%v rowconfigure {%v} -background %v", t.tablelist.id, t.id, background))
}

func (t *TablelistItem) Font() string {
	sl, _ := evalAsString(fmt.Sprintf("%v rowcget {%v} -font", t.tablelist.id, t.id))
	return sl
}

func (t *TablelistItem) SetFont(font []string) error {
	return eval(fmt.Sprintf("%v rowconfigure {%v} -font %v", t.tablelist.id, t.id, font))
}

func (t *TablelistItem) Foreground() string {
	sl, _ := evalAsString(fmt.Sprintf("%v rowcget {%v} -foreground", t.tablelist.id, t.id))
	return sl
}

func (t *TablelistItem) SetForeground(font []string) error {
	return eval(fmt.Sprintf("%v rowconfigure {%v} -foreground %v", t.tablelist.id, t.id, font))
}

func (t *TablelistItem) Hide() bool {
	sl, _ := evalAsBool(fmt.Sprintf("%v rowcget {%v} -hide", t.tablelist.id, t.id))
	return sl
}

func (t *TablelistItem) SetHide(hide bool) error {
	return eval(fmt.Sprintf("%v rowconfigure {%v} -hide %v", t.tablelist.id, t.id, hide))
}

func (t *TablelistItem) Name() string {
	sl, _ := evalAsString(fmt.Sprintf("%v rowcget {%v} -name", t.tablelist.id, t.id))
	return sl
}

func (t *TablelistItem) SetName(name string) error {
	return eval(fmt.Sprintf("%v rowconfigure {%v} -name %v", t.tablelist.id, t.id, Quote(name)))
}

func (t *TablelistItem) Selectable() bool {
	sl, _ := evalAsBool(fmt.Sprintf("%v rowcget {%v} -selectable", t.tablelist.id, t.id))
	return sl
}

func (t *TablelistItem) SetSelectable(selectable bool) error {
	return eval(fmt.Sprintf("%v rowconfigure {%v} -name %v", t.tablelist.id, t.id, selectable))
}

func (t *TablelistItem) SelectBackground() string {
	sl, _ := evalAsString(fmt.Sprintf("%v rowcget {%v} -selectbackground", t.tablelist.id, t.id))
	return sl
}

func (t *TablelistItem) SetSelectBackground(select_background string) error {
	return eval(fmt.Sprintf("%v rowconfigure {%v} -selectbackground %v", t.tablelist.id, t.id, select_background))
}

func (t *TablelistItem) SelectForeground() string {
	sl, _ := evalAsString(fmt.Sprintf("%v rowcget {%v} -selectforeground", t.tablelist.id, t.id))
	return sl
}

func (t *TablelistItem) SetSelectForeground(select_foreground string) error {
	return eval(fmt.Sprintf("%v rowconfigure {%v} -selectforeground %v", t.tablelist.id, t.id, select_foreground))
}

func (t *TablelistItem) Text() []string {
	sl, _ := evalAsStringList(fmt.Sprintf("%v rowcget {%v} -text", t.tablelist.id, t.id))
	return sl
}

func (t *TablelistItem) SetText(text []string) error {
	return eval(fmt.Sprintf("%v rowconfigure {%v} -text %v", t.tablelist.id, t.id, strings.Join(QuoteAll(text), " ")))
}

// --- Tablelist-level commands that target TablelistItems

/*
   pathName activate index
   pathName bbox index
   pathName childcount nodeIndex
   pathName childindex index
   pathName childkeys nodeIndex
   (!) pathName collapse indexList ?-fully|-partly?
   pathName configrowlist {index option value index option value ...}
   pathName configrows ?index option value index option value ...?
   pathName delete firstIndex lastIndex
   pathName depth nodeIndex
   pathName descendantcount nodeIndex
   (!) pathName expand indexList ?-fully|-partly?
   pathName findrowname name ?-descend? ?-parent nodeIndex?
   pathName getformatted indexList
   pathName getfullkeys indexList
   pathName getkeys indexList
   pathName hasrowattrib index name
   pathName index index
   pathName insert index ?item item ...?
   pathName insertchildlist parentNodeIndex childIndex itemList
   pathName insertchild(ren) parentNodeIndex childIndex ?item item ...?
   pathName insertlist index itemList
   (!) pathName isexpanded index
   pathName isviewable index
   pathName move sourceIndex targetIndex
   pathName move sourceIndex targetParentNodeIndex targetChildIndex
   pathName noderow parentNodeIndex childIndex
   pathName parentkey nodeIndex
   pathName refreshsorting ?parentNodeIndex?
   pathName rowattrib index ?name? ?value name value ...?
   pathName rowcget index option
   pathName rowconfigure index ?option? ?value option value ...?
   pathName see index
   pathName togglerowhide indexList
   pathName toplevelkey index
   pathName unsetrowattrib index name
   pathName viewablerowcount ?firstIndex lastIndex?
*/

// returns true if the row is expanded.
func (t *TablelistItem) IsExpanded() bool {
	b, _ := evalAsBool(fmt.Sprintf("%v isexpanded {%v}", t.tablelist.id, t.id))
	return b
}

// expands row.
func (t *TablelistItem) Expand() error {
	return eval(fmt.Sprintf("%v expand {%v} -partly", t.tablelist.id, t.id))
}

// expands row and all chilren.
func (t *TablelistItem) ExpandAll() error {
	return eval(fmt.Sprintf("%v expand {%v} -fully", t.tablelist.id, t.id))
}

// returns true if the row is collapsed.
func (t *TablelistItem) IsCollapsed() bool {
	return !t.IsExpanded()
}

// collapses row ('partly').
func (t *TablelistItem) Collapse() error {
	return eval(fmt.Sprintf("%v collapse {%v} -partly", t.tablelist.id, t.id))
}

// collapses row and all children ('fully').
func (t *TablelistItem) CollapseAll() error {
	return eval(fmt.Sprintf("%v collapse {%v} -fully", t.tablelist.id, t.id))
}
