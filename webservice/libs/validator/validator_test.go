package validator_test

import (
	"testing"

	"webservice/libs/validator"
)

type testCase struct {
	name     string
	input    interface{}
	expected string
}

var tests = []*testCase{
	{
		name: `required`,
		input: struct {
			Title string `json:"title" validate:"required"`
		}{},
		expected: "O campo 'title' é obrigatório",
	},
	{
		name: `max`,
		input: struct {
			Course string `json:"course" validate:"max=7"`
		}{Course: "CS-0001."},
		expected: "O campo 'course' deve ter um tamanho de no máximo 7 caracteres",
	},
	{
		name: `datetime`,
		input: struct {
			Timestamp string `json:"timestamp" validate:"datetime"`
		}{Timestamp: "2025-10-01T09:35:14.00Z"},
		expected: "O campo 'timestamp' deve seguir o formato RFC3339 'YYYY-MM-DDThh:mm:ssZ'",
	},
	{
		name: `user_id`,
		input: struct {
			UserID string `json:"user_id" validate:"user_id"`
		}{UserID: "thisisnotright"},
		expected: "O campo 'user_id' deve seguir o padrão 'user_XXX'",
	},
	{
		name: `msg_id`,
		input: struct {
			ID string `json:"msg_id" validate:"msg_id"`
		}{ID: "thisIsAlsoNotRight"},
		expected: "O campo 'msg_id' deve seguir o padrão 'msg_XXX'",
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
