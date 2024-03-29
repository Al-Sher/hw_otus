// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package hw10programoptimization

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson7b0a8cdfDecodeGithubComAlSherHwOtusHw10ProgramOptimization(in *jlexer.Lexer, out *User) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "Name":
			out.Name = string(in.String())
		case "Username":
			out.Username = string(in.String())
		case "Email":
			out.Email = string(in.String())
		case "Phone":
			out.Phone = string(in.String())
		case "Password":
			out.Password = string(in.String())
		case "Address":
			out.Address = string(in.String())
		case "ID":
			out.ID = int(in.Int())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson7b0a8cdfEncodeGithubComAlSherHwOtusHw10ProgramOptimization(out *jwriter.Writer, in User) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Name\":"
		out.RawString(prefix[1:])
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"Username\":"
		out.RawString(prefix)
		out.String(string(in.Username))
	}
	{
		const prefix string = ",\"Email\":"
		out.RawString(prefix)
		out.String(string(in.Email))
	}
	{
		const prefix string = ",\"Phone\":"
		out.RawString(prefix)
		out.String(string(in.Phone))
	}
	{
		const prefix string = ",\"Password\":"
		out.RawString(prefix)
		out.String(string(in.Password))
	}
	{
		const prefix string = ",\"Address\":"
		out.RawString(prefix)
		out.String(string(in.Address))
	}
	{
		const prefix string = ",\"ID\":"
		out.RawString(prefix)
		out.Int(int(in.ID))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v User) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson7b0a8cdfEncodeGithubComAlSherHwOtusHw10ProgramOptimization(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v User) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson7b0a8cdfEncodeGithubComAlSherHwOtusHw10ProgramOptimization(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *User) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson7b0a8cdfDecodeGithubComAlSherHwOtusHw10ProgramOptimization(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *User) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson7b0a8cdfDecodeGithubComAlSherHwOtusHw10ProgramOptimization(l, v)
}
