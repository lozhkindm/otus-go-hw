package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"testing"

	exts "github.com/lozhkindm/otus-go-hw/hw09_struct_validator/extractors"
	"github.com/stretchr/testify/require"
)

// Test the function on different structures and other types.
type (
	UserRole string

	User struct {
		ID         int             `json:"id" validate:"maximum:10"`      // no tag, undefined rule
		Hash       string          `validate:"len:10"`                    // len rule
		Name       string          `json:"name"`                          // no tag
		Age        int             `validate:"min:18|max:50"`             // min & max rules
		Experience int             `validate:"in:3,4,5"`                  // in rule for int (range)
		Email      string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"` // regexp rule
		Role       UserRole        `validate:"in:admin,stuff"`            // in rule for string (enum)
		Phones     []int           `validate:"min:12345"`                 // int slice min rule
		Nicknames  []string        `validate:"len:5"`                     // string slice len rule
		Token      Token           `json:"token" validate:"min:15"`       // undefined rule for nested struct
		LastLogin  Response        `validate:"nested"`                    // nested validation
		hobbies    []string        `validate:"len:1"`                     // private
		meta       json.RawMessage `validate:"no"`                        // private, undefined tag
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code        int      `validate:"in:200,404,500"` // in rule for int (range)
		Body        string   `json:"omitempty"`          // no tag
		Headers     []string `validate:"len:3"`          // string slice len rule
		Connections []int    `validate:"min:3|max:5"`    // int slice min & max rules
	}

	InvalidIn struct {
		Field int `validate:"in:s,d"`
	}

	EmptyInString struct {
		Field string `validate:"in:"`
	}

	EmptyInInt struct {
		Field int `validate:"in:"`
	}

	InvalidMin struct {
		Field int `validate:"min:s"`
	}

	InvalidMax struct {
		Field int `validate:"max:s"`
	}

	InvalidLen struct {
		Field string `validate:"len:s"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name string
		in   interface{}
		err  error
	}{
		{
			name: "success",
			in: User{
				ID:         15,
				Hash:       "abcdefghij",
				Name:       "Random guy",
				Age:        25,
				Experience: 3,
				Email:      "randomguy@gmail.com",
				Role:       "stuff",
				Phones:     []int{12345, 54321},
				Nicknames:  []string{"guy__", "__guy", "_guy_"},
				Token:      Token{},
				LastLogin: Response{
					Code:        200,
					Body:        "does not matter",
					Headers:     []string{"acc", "hed", "der"},
					Connections: []int{3, 4, 5},
				},
				hobbies: []string{"walking", "running"},
				meta:    nil,
			},
			err: nil,
		},
		{
			name: "validation error",
			in: User{
				ID:         15,
				Hash:       "abc",
				Name:       "does not matter",
				Age:        15,
				Experience: 1,
				Email:      "===@gmail.com",
				Role:       "random",
				Phones:     []int{4321, 123456},
				Nicknames:  []string{"random_guy", "_guy_"},
				Token:      Token{},
				LastLogin: Response{
					Code:        403,
					Body:        "does not matter",
					Headers:     []string{"hed", "application/json"},
					Connections: []int{3, 1},
				},
				hobbies: []string{"walking", "running"},
				meta:    nil,
			},
			err: ValidationErrors{
				{
					Field: "Hash",
					Err:   errors.New(`the field "Hash" must be 10 characters`),
				},
				{
					Field: "Age",
					Err:   errors.New(`the field "Age" must be at least 18`),
				},
				{
					Field: "Experience",
					Err:   errors.New(`the field "Experience" must have one of the following values: 3, 4, 5`),
				},
				{
					Field: "Email",
					Err:   errors.New(`the field "Email" format is invalid`),
				},
				{
					Field: "Role",
					Err:   errors.New(`the field "Role" must have one of the following values: admin, stuff`),
				},
				{
					Field: "Phones",
					Err:   errors.New(`the field "Phones" (4321) must be at least 12345`),
				},
				{
					Field: "Nicknames",
					Err:   errors.New(`the field "Nicknames" (random_guy) must be 5 characters`),
				},
				{
					Field: "Code",
					Err:   errors.New(`the field "LastLogin->Code" must have one of the following values: 200, 404, 500`),
				},
				{
					Field: "Headers",
					Err:   errors.New(`the field "LastLogin->Headers" (application/json) must be 3 characters`),
				},
				{
					Field: "Connections",
					Err:   errors.New(`the field "LastLogin->Connections" (1) must be at least 3`),
				},
			},
		},
		{
			name: "not struct",
			in:   "abc",
			err:  errNotStructure,
		},
		{
			name: "invalid in",
			in:   InvalidIn{},
			err:  exts.ErrInInvalidOptions,
		},
		{
			name: "invalid min",
			in:   InvalidMin{},
			err:  exts.ErrInvalidMinValue,
		},
		{
			name: "invalid max",
			in:   InvalidMax{},
			err:  exts.ErrInvalidMaxValue,
		},
		{
			name: "invalid len",
			in:   InvalidLen{},
			err:  exts.ErrInvalidLenValue,
		},
		{
			name: "empty in string",
			in:   EmptyInString{},
			err:  exts.ErrInNoOptions,
		},
		{
			name: "empty in int",
			in:   EmptyInInt{},
			err:  exts.ErrInNoOptions,
		},
	}

	var vErr ValidationErrors
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := Validate(tt.in, "")
			if tt.err == nil {
				require.NoError(t, err)
			}
			if errors.As(err, &vErr) {
				require.Equal(t, tt.err.Error(), err.Error())
			} else {
				require.ErrorIs(t, err, tt.err)
			}
		})
	}
}
