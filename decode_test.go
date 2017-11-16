package fixedwidth

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

func ExampleUnmarshal() {
	// define the format
	var people []struct {
		ID        int     `fixed:"1,5"`
		FirstName string  `fixed:"6,15"`
		LastName  string  `fixed:"16,25"`
		Grade     float64 `fixed:"26,30"`
	}

	// define some fixed-with data to parse
	data := []byte("" +
		"1    Ian       Lopshire  99.50" + "\n" +
		"2    John      Doe       89.50" + "\n" +
		"3    Jane      Doe       79.50" + "\n")

	err := Unmarshal(data, &people)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", people[0])
	fmt.Printf("%+v\n", people[1])
	fmt.Printf("%+v\n", people[2])
	// Output:
	//{ID:1 FirstName:Ian LastName:Lopshire Grade:99.5}
	//{ID:2 FirstName:John LastName:Doe Grade:89.5}
	//{ID:3 FirstName:Jane LastName:Doe Grade:79.5}
}

func TestUnmarshal(t *testing.T) {
	// allTypes contains a field with all current supported types.
	type allTypes struct {
		String          string          `fixed:"1,5"`
		Int             int             `fixed:"6,10"`
		Float           float64         `fixed:"11,15"`
		TextUnmarshaler EncodableString `fixed:"16,20"` // test encoding.TextUnmarshaler functionality
	}
	for _, tt := range []struct {
		name      string
		rawValue  []byte
		target    interface{}
		expected  interface{}
		shouldErr bool
	}{
		{
			name:     "Basic Slice Case",
			rawValue: []byte("foo  123  1.2  bar" + "\n" + "bar  321  2.1  foo"),
			target:   &[]allTypes{},
			expected: &[]allTypes{
				{"foo", 123, 1.2, EncodableString{"bar", nil}},
				{"bar", 321, 2.1, EncodableString{"foo", nil}},
			},
			shouldErr: false,
		},
		{
			name:      "Basic Struct Case",
			rawValue:  []byte("foo  123  1.2  bar"),
			target:    &allTypes{},
			expected:  &allTypes{"foo", 123, 1.2, EncodableString{"bar", nil}},
			shouldErr: false,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			err := Unmarshal(tt.rawValue, tt.target)
			if tt.shouldErr != (err != nil) {
				t.Errorf("Unmarshal() err want %v, have %v (%v)", tt.shouldErr, err != nil, err)
			}
			if !reflect.DeepEqual(tt.target, tt.expected) {
				t.Errorf("Unmarshal() want %+v, have %+v", tt.target, tt.expected)
			}

		})
	}

	t.Run("Invalid Unmarshal Errors", func(t *testing.T) {
		for _, tt := range []struct {
			name      string
			v         interface{}
			shouldErr bool
		}{
			{"Invalid Unmarshal Nil", nil, true},
			{"Invalid Unmarshal Not Pointer 1", struct{}{}, true},
			{"Invalid Unmarshal Not Pointer 2", []struct{}{}, true},
			{"Valid Unmarshal slice", &[]struct{}{}, false},
			{"Valid Unmarshal struct", &struct{}{}, false},
		} {
			t.Run(tt.name, func(t *testing.T) {
				err := Unmarshal([]byte{}, tt.v)
				if tt.shouldErr != (err != nil) {
					t.Errorf("Unmarshal() err want %v, have %v (%v)", tt.shouldErr, err != nil, err)
				}
			})
		}
	})
}
