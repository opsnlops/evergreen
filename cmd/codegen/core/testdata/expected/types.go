// Code generated by rest/model/codegen.go. DO NOT EDIT.

package model

import "github.com/evergreen-ci/evergreen/rest/model"

type APIMockTypes struct {
	BoolType       bool     `json:"bool_type"`
	BoolPtrType    *bool    `json:"bool_ptr_type"`
	IntType        int      `json:"int_type"`
	IntPtrType     *int     `json:"int_ptr_type"`
	StringType     string   `json:"string_type"`
	StringPtrType  *string  `json:"string_ptr_type"`
	Uint64Type     int      `json:"uint64_type"`
	Uint64PtrType  *int     `json:"uint64_ptr_type"`
	Float64Type    float64  `json:"float64_type"`
	Float64PtrType *float64 `json:"float64_ptr_type"`
	RuneType       int      `json:"rune_type"`
	RunePtrType    *int     `json:"rune_ptr_type"`
}

// APIMockTypesBuildFromService takes the model.MockTypes DB struct and
// returns the REST struct *APIMockTypes with the corresponding fields populated
func APIMockTypesBuildFromService(t model.MockTypes) *APIMockTypes {
	m := APIMockTypes{}
	m.BoolType = BoolBool(t.BoolType)
	m.BoolPtrType = BoolPtrBoolPtr(t.BoolPtrType)
	m.IntType = IntInt(t.IntType)
	m.IntPtrType = IntPtrIntPtr(t.IntPtrType)
	m.StringType = StringString(t.StringType)
	m.StringPtrType = StringPtrStringPtr(t.StringPtrType)
	m.Uint64Type = Uint64Int(t.Uint64Type)
	m.Uint64PtrType = Uint64PtrIntPtr(t.Uint64PtrType)
	m.Float64Type = Float64Float64(t.Float64Type)
	m.Float64PtrType = Float64PtrFloat64Ptr(t.Float64PtrType)
	m.RuneType = RuneInt(t.RuneType)
	m.RunePtrType = RunePtrIntPtr(t.RunePtrType)
	return &m
}

// APIMockTypesToService takes the APIMockTypes REST struct and returns the DB struct
// *model.MockTypes with the corresponding fields populated
func APIMockTypesToService(m APIMockTypes) *model.MockTypes {
	out := &model.MockTypes{}
	out.BoolType = BoolBool(m.BoolType)
	out.BoolPtrType = BoolPtrBoolPtr(m.BoolPtrType)
	out.IntType = IntInt(m.IntType)
	out.IntPtrType = IntPtrIntPtr(m.IntPtrType)
	out.StringType = StringString(m.StringType)
	out.StringPtrType = StringPtrStringPtr(m.StringPtrType)
	out.Uint64Type = IntUint64(m.Uint64Type)
	out.Uint64PtrType = IntPtrUint64Ptr(m.Uint64PtrType)
	out.Float64Type = Float64Float64(m.Float64Type)
	out.Float64PtrType = Float64PtrFloat64Ptr(m.Float64PtrType)
	out.RuneType = IntRune(m.RuneType)
	out.RunePtrType = IntPtrRunePtr(m.RunePtrType)
	return out
}