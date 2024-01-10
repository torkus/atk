// Copyright 2018 visualfc. All rights reserved.

package tk

import "fmt"

type TablelistItem struct {
	tablelist *Tablelist
	id        string
	//Vals      []string // todo: improve
}

// ---

// todo: probably belongs in ids.go
func makeTablelistItemId(treeid string, pid string) string {
	if pid != "" {
		return makeNamedId(pid + ".I")
	}
	return makeNamedId(treeid + ".I")
}

// ---

func (t *TablelistItem) Id() string {
	return t.id
}

func (t *TablelistItem) IsValid() bool {
	return t != nil && t.tablelist != nil
}

/*
	func (t *TablelistItem) InsertItem(index int, text string, values []string) *TablelistItem {
		if !t.IsValid() {
			return nil
		}
		return t.tablelist.InsertItem(t, index, text, values)
	}
*/
func (t *TablelistItem) Index() int {
	if !t.IsValid() || t.IsRoot() {
		return -1
	}
	r, err := evalAsIntEx(fmt.Sprintf("%v index {%v}", t.tablelist.id, t.id), false)
	if err != nil {
		return -1
	}
	return r
}

func (t *TablelistItem) IsRoot() bool {
	return t.id == ""
}

func (t *TablelistItem) Parent() *TablelistItem {
	if !t.IsValid() || t.IsRoot() {
		return nil
	}
	r, err := evalAsStringEx(fmt.Sprintf("%v parent {%v}", t.tablelist.id, t.id), false)
	if err != nil {
		return nil
	}
	return &TablelistItem{t.tablelist, r}
}

func (t *TablelistItem) Next() *TablelistItem {
	if !t.IsValid() || t.IsRoot() {
		return nil
	}
	r, err := evalAsStringEx(fmt.Sprintf("%v next {%v}", t.tablelist.id, t.id), false)
	if err != nil || r == "" {
		return nil
	}
	return &TablelistItem{t.tablelist, r}
}

func (t *TablelistItem) Prev() *TablelistItem {
	if !t.IsValid() || t.IsRoot() {
		return nil
	}
	r, err := evalAsStringEx(fmt.Sprintf("%v prev {%v}", t.tablelist.id, t.id), false)
	if err != nil || r == "" {
		return nil
	}
	return &TablelistItem{t.tablelist, r}
}

/*
func (t *TablelistItem) Children() (lst []*TablelistItem) {
	if !t.IsValid() {
		return
	}
	ids, err := evalAsStringList(fmt.Sprintf("%v children {%v}", t.tree.id, t.id))
	if err != nil {
		return
	}
	for _, id := range ids {
		lst = append(lst, &TablelistItem{t.tree, id})
	}
	return
}
*/

func (t *TablelistItem) Children() (lst []*TablelistItem) {
	if !t.IsValid() {
		return
	}
	ids, err := evalAsStringList(fmt.Sprintf("%v childkeys {%v}", t.tablelist.id, t.id))
	if err != nil {
		return
	}
	for _, id := range ids {
		lst = append(lst, &TablelistItem{t.tablelist, id})
	}
	return
}

/*
	func (t *TablelistItem) SetExpanded(expand bool) error {
		if !t.IsValid() || t.IsRoot() {
			return ErrInvalid
		}
		return eval(fmt.Sprintf("%v item {%v} -open %v", t.tree.id, t.id, expand))
	}
*/
func (t *TablelistItem) SetExpanded(expand bool) error {
	if !t.IsValid() || t.IsRoot() {
		return ErrInvalid
	}
	if expand {
		return eval(fmt.Sprintf("%v expand {%v} -partly", t.tablelist.id, t.id))
	} else {
		return eval(fmt.Sprintf("%v collapse {%v} -fully", t.tablelist.id, t.id))
	}
}

/*
func (t *TablelistItem) IsExpanded() bool {
	if !t.IsValid() || t.IsRoot() {
		return false
	}
	r, _ := evalAsBool(fmt.Sprintf("%v item {%v} -open", t.tree.id, t.id))
	return r
}
*/

func (t *TablelistItem) IsExpanded() bool {
	if !t.IsValid() || t.IsRoot() {
		return false
	}
	ids, _ := evalAsStringList(fmt.Sprintf("%v expanded", t.tablelist.id))
	for _, id := range ids {
		if t.id == id {
			return true
		}
	}
	return false
}

/*
func (t *TablelistItem) expandAll(item *TablelistItem) error {
	for _, child := range item.Children() {
		child.SetExpanded(true)
		t.expandAll(child)
	}
	return nil
}

func (t *TablelistItem) ExpandAll() error {
	return t.expandAll(t)
}
*/

func (t *TablelistItem) ExpandAll() error {
	if !t.IsValid() || t.IsRoot() {
		return ErrInvalid
	}
	return eval(fmt.Sprintf("%v expand {%v} -fully", t.tablelist.id, t.id))
}

func (t *TablelistItem) collapseAll(item *TablelistItem) error {
	for _, child := range item.Children() {
		child.SetExpanded(false)
		t.collapseAll(child)
	}
	return nil
}

func (t *TablelistItem) CollapseAll() error {
	return t.collapseAll(t)
}

func (t *TablelistItem) Expand() error {
	return t.SetExpanded(true)
}

func (t *TablelistItem) Collapse() error {
	return t.SetExpanded(false)
}

func (t *TablelistItem) SetText(text string) error {
	if !t.IsValid() || t.IsRoot() {
		return ErrInvalid
	}
	setObjText("atk_tablelist_item", text)
	return eval(fmt.Sprintf("%v item {%v} -text $atk_tablelist_item", t.tablelist.id, t.id))
}

func (t *TablelistItem) Text() string {
	if !t.IsValid() || t.IsRoot() {
		return ""
	}
	r, _ := evalAsString(fmt.Sprintf("%v item {%v} -text", t.tablelist.id, t.id))
	return r
}

func (t *TablelistItem) SetValues(values []string) error {
	if !t.IsValid() || t.IsRoot() {
		return ErrInvalid
	}
	setObjTextList("atk_tablelist_values", values)
	return eval(fmt.Sprintf("%v item {%v} -values $atk_tablelist_values", t.tablelist.id, t.id))
}

func (t *TablelistItem) Values() []string {
	if !t.IsValid() || t.IsRoot() {
		return nil
	}
	r, _ := evalAsStringList(fmt.Sprintf("%v item {%v} -values", t.tablelist.id, t.id))
	return r
}

func (t *TablelistItem) SetImage(img *Image) error {
	if !t.IsValid() || t.IsRoot() {
		return ErrInvalid
	}
	var iid string
	if img != nil {
		iid = img.Id()
	}
	return eval(fmt.Sprintf("%v item {%v} -image {%v}", t.tablelist.id, t.id, iid))
}

func (t *TablelistItem) Image() *Image {
	if !t.IsValid() || t.IsRoot() {
		return nil
	}
	r, err := evalAsString(fmt.Sprintf("%v item {%v} -image", t.tablelist.id, t.id))
	return parserImageResult(r, err)
}

func (t *TablelistItem) SetColumnText(column int, text string) error {
	if column < 0 {
		return ErrInvalid
	} else if column == 0 {
		return t.SetText(text)
	}
	if !t.IsValid() || t.IsRoot() {
		return ErrInvalid
	}
	setObjText("atk_tablelist_column", text)
	return eval(fmt.Sprintf("%v set {%v} %v $atk_tablelist_column", t.tablelist.id, t.id, column-1))
}

func (t *TablelistItem) ColumnText(column int) string {
	if column < 0 {
		return ""
	} else if column == 0 {
		return t.Text()
	}
	if !t.IsValid() || t.IsRoot() {
		return ""
	}
	r, _ := evalAsString(fmt.Sprintf("%v set {%v} %v", t.tablelist.id, t.id, column-1))
	return r
}
