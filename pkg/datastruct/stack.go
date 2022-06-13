/*
 * Copyright © 2022 photowey (photowey@gmail.com)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package datastruct

const (
	EmptySlice = 0
)

type Element any

type Stack interface {
	Length() int
	IsEmpty() bool
	Push(element Element)
	PushList(elements ...Element)
	Pop() Element
	Peek() Element
	Clear()
}

type stack struct {
	list []Element
}

func NewStack() Stack {
	return &stack{list: make([]Element, 0)}
}

func (st *stack) Length() int {
	return len(st.list)
}

func (st *stack) IsEmpty() bool {
	return st.Length() == EmptySlice
}

func (st *stack) Push(element Element) {
	st.list = append(st.list, element)
}

func (st *stack) PushList(elements ...Element) {
	st.list = append(st.list, elements...)
}

func (st *stack) Pop() Element {
	if st.IsEmpty() {
		return nil
	} else {
		top := st.list[len(st.list)-1]
		st.list = st.list[:len(st.list)-1]

		return top
	}
}

func (st *stack) Peek() Element {
	if st.IsEmpty() {
		return nil
	} else {
		return st.list[len(st.list)-1]
	}
}

func (st *stack) Clear() {
	if len(st.list) == 0 {
		return
	}
	for i := 0; i < st.Length(); i++ {
		st.list[i] = nil
	}
	st.list = make([]Element, 0)
}

func ToElement(element any) Element {
	return Element(element)
}

func ToElements(elements ...any) []Element {
	elementz := make([]Element, 0, len(elements))
	for _, element := range elements {
		elementz = append(elementz, element)
	}

	return elementz
}
