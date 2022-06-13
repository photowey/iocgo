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

import (
	"reflect"
	"testing"
)

func TestNewStack(t *testing.T) {
	tests := []struct {
		name string
		want Stack
	}{
		{
			name: "Test NewStack()",
			want: &stack{list: make([]Element, 0)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStack(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStack() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToElement(t *testing.T) {
	type args struct {
		element any
	}
	tests := []struct {
		name string
		args args
		want Element
	}{
		{
			name: "Test ToElement()",
			args: args{
				element: "Hello world!",
			},
			want: Element("Hello world!"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToElement(tt.args.element); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToElement() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToElements(t *testing.T) {
	type args struct {
		elements []any
	}
	tests := []struct {
		name string
		args args
		want []Element
	}{
		{
			name: "Test ToElement()",
			args: args{
				elements: []any{"Hello", "world"},
			},
			want: []Element{Element("Hello"), Element("world")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToElements(tt.args.elements...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToElements() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stack_Clear(t *testing.T) {
	type fields struct {
		list []Element
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Test stack.Clear()",
			fields: fields{
				list: []Element{Element("Hello"), Element("world")},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := &stack{
				list: tt.fields.list,
			}
			st.Clear()
		})
	}
}

func Test_stack_IsEmpty(t *testing.T) {
	type fields struct {
		list []Element
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Test stack.IsEmpty()-false",
			fields: fields{
				list: []Element{Element("Hello"), Element("world")},
			},
			want: false,
		},
		{
			name: "Test stack.IsEmpty()-true",
			fields: fields{
				list: []Element{},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := &stack{
				list: tt.fields.list,
			}
			if got := st.IsEmpty(); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stack_Length(t *testing.T) {
	type fields struct {
		list []Element
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "Test stack.Length()-2",
			fields: fields{
				list: []Element{Element("Hello"), Element("world")},
			},
			want: 2,
		},
		{
			name: "Test stack.Length()-0",
			fields: fields{
				list: []Element{},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := &stack{
				list: tt.fields.list,
			}
			if got := st.Length(); got != tt.want {
				t.Errorf("Length() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stack_Peek(t *testing.T) {
	type fields struct {
		list []Element
	}
	tests := []struct {
		name   string
		fields fields
		want   Element
	}{
		{
			name: "Test stack.Peek()",
			fields: fields{
				list: []Element{Element("Hello"), Element("world")},
			},
			want: Element("world"),
		},
		{
			name: "Test stack.Peek()-nil",
			fields: fields{
				list: []Element{},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := &stack{
				list: tt.fields.list,
			}
			if got := st.Peek(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Peek() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stack_Pop(t *testing.T) {
	type fields struct {
		list []Element
	}
	tests := []struct {
		name   string
		fields fields
		want   Element
	}{
		{
			name: "Test stack.Pop()",
			fields: fields{
				list: []Element{Element("Hello"), Element("world")},
			},
			want: Element("world"),
		},
		{
			name: "Test stack.Pop()-nil",
			fields: fields{
				list: []Element{},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := &stack{
				list: tt.fields.list,
			}
			if got := st.Pop(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pop() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stack_Push(t *testing.T) {
	type fields struct {
		list []Element
	}
	type args struct {
		element Element
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Test stack.Push()",
			fields: fields{
				list: []Element{Element("Hello")},
			},
			args: args{
				element: Element("world"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := &stack{
				list: tt.fields.list,
			}
			st.Push(tt.args.element)
		})
	}
}

func Test_stack_PushList(t *testing.T) {
	type fields struct {
		list []Element
	}
	type args struct {
		elements []Element
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Test stack.PushList()",
			fields: fields{
				list: []Element{},
			},
			args: args{
				elements: []Element{Element("Hello"), Element("Hello")},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := &stack{
				list: tt.fields.list,
			}
			st.PushList(tt.args.elements...)
		})
	}
}
