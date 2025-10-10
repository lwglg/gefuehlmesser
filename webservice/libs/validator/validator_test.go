package validator_test

import (
	"testing"

	"webservice/libs/validator"
	v "webservice/libs/validator"
)

type testCase struct {
	name     string
	input    interface{}
	expected v.ErrorObject
}

var tests = []*testCase{
	{
		name: `required`,
		input: struct {
			Title string `json:"title" validate:"required"`
		}{},
		expected: v.ErrorObject{
			Code:    "FIELD_REQUIRED",
			Message: "O campo 'title' é obrigatório",
		},
	},
	{
		name: `max`,
		input: struct {
			Course string `json:"course" validate:"max=7"`
		}{Course: "CS-0001."},
		expected: v.ErrorObject{
			Code:    "INVALID_SIZE",
			Message: "O campo 'course' deve ter um tamanho de no máximo 7 caracteres",
		},
	},
	{
		name: `datetime`,
		input: struct {
			Timestamp string `json:"timestamp" validate:"datetime"`
		}{Timestamp: "2025-10-01T09:35:14.00Z"},
		expected: v.ErrorObject{
			Code:    "INVALID_TIMESTAMP",
			Message: "O campo 'timestamp' deve seguir o formato RFC3339 'YYYY-MM-DDThh:mm:ssZ'",
		},
	},
	{
		name: `id_field`,
		input: struct {
			UserID string `json:"user_id" validate:"id_field"`
		}{UserID: "thisisnotright"},
		expected: v.ErrorObject{
			Code:    "INVALID_ID_FIELD",
			Message: "O campo 'user_id' deve seguir o padrão '[user|msg|perf]_XXX'",
		},
	},
	{
		name: `id_field`,
		input: struct {
			ID string `json:"msg_id" validate:"id_field"`
		}{ID: "thisIsAlsoNotRight"},
		expected: v.ErrorObject{
			Code:    "INVALID_ID_FIELD",
			Message: "O campo 'msg_id' deve seguir o padrão '[user|msg|perf]_XXX'",
		},
	},
}

func TestToErrResponse(t *testing.T) {
	vr := validator.New()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := vr.Struct(tc.input)
			errResp := validator.ToErrResponse(err)

			if errResp == nil || len(errResp.Errors) != 1 {
				t.Fatalf(`Expected: "{[%v]}", Got: "%v"`, tc.expected, errResp)
			} else if errResp.Errors[0] != tc.expected {
				t.Fatalf(`Expected: "%v", Got:"%v"`, tc.expected, errResp.Errors[0])
			}
		})
	}
}
