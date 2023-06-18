// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package models

import (
	json "encoding/json"
	pgtype "github.com/jackc/pgx/pgtype"
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

func easyjson5a72dc82DecodeGithubComQqq4uTPDBMSTermProjectInternalModels(in *jlexer.Lexer, out *PostsList) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(PostsList, 0, 0)
			} else {
				*out = PostsList{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 Post
			easyjson5a72dc82DecodeGithubComQqq4uTPDBMSTermProjectInternalModels1(in, &v1)
			*out = append(*out, v1)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson5a72dc82EncodeGithubComQqq4uTPDBMSTermProjectInternalModels(out *jwriter.Writer, in PostsList) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v2, v3 := range in {
			if v2 > 0 {
				out.RawByte(',')
			}
			easyjson5a72dc82EncodeGithubComQqq4uTPDBMSTermProjectInternalModels1(out, v3)
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v PostsList) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson5a72dc82EncodeGithubComQqq4uTPDBMSTermProjectInternalModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v PostsList) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson5a72dc82EncodeGithubComQqq4uTPDBMSTermProjectInternalModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *PostsList) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson5a72dc82DecodeGithubComQqq4uTPDBMSTermProjectInternalModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *PostsList) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson5a72dc82DecodeGithubComQqq4uTPDBMSTermProjectInternalModels(l, v)
}
func easyjson5a72dc82DecodeGithubComQqq4uTPDBMSTermProjectInternalModels1(in *jlexer.Lexer, out *Post) {
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
		case "id":
			out.ID = int(in.Int())
		case "parent":
			out.Parent = int(in.Int())
		case "author":
			out.Author = string(in.String())
		case "message":
			out.Message = string(in.String())
		case "isEdited":
			out.IsEdited = bool(in.Bool())
		case "forum":
			out.Forum = string(in.String())
		case "thread":
			out.Thread = int(in.Int())
		case "created":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Created).UnmarshalJSON(data))
			}
		case "path":
			easyjson5a72dc82DecodeGithubComJackcPgxPgtype(in, &out.Path)
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
func easyjson5a72dc82EncodeGithubComQqq4uTPDBMSTermProjectInternalModels1(out *jwriter.Writer, in Post) {
	out.RawByte('{')
	first := true
	_ = first
	if in.ID != 0 {
		const prefix string = ",\"id\":"
		first = false
		out.RawString(prefix[1:])
		out.Int(int(in.ID))
	}
	if in.Parent != 0 {
		const prefix string = ",\"parent\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Parent))
	}
	{
		const prefix string = ",\"author\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Author))
	}
	{
		const prefix string = ",\"message\":"
		out.RawString(prefix)
		out.String(string(in.Message))
	}
	if in.IsEdited {
		const prefix string = ",\"isEdited\":"
		out.RawString(prefix)
		out.Bool(bool(in.IsEdited))
	}
	if in.Forum != "" {
		const prefix string = ",\"forum\":"
		out.RawString(prefix)
		out.String(string(in.Forum))
	}
	if in.Thread != 0 {
		const prefix string = ",\"thread\":"
		out.RawString(prefix)
		out.Int(int(in.Thread))
	}
	if true {
		const prefix string = ",\"created\":"
		out.RawString(prefix)
		out.Raw((in.Created).MarshalJSON())
	}
	if true {
		const prefix string = ",\"path\":"
		out.RawString(prefix)
		easyjson5a72dc82EncodeGithubComJackcPgxPgtype(out, in.Path)
	}
	out.RawByte('}')
}
func easyjson5a72dc82DecodeGithubComJackcPgxPgtype(in *jlexer.Lexer, out *pgtype.Int8Array) {
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
		case "Elements":
			if in.IsNull() {
				in.Skip()
				out.Elements = nil
			} else {
				in.Delim('[')
				if out.Elements == nil {
					if !in.IsDelim(']') {
						out.Elements = make([]pgtype.Int8, 0, 4)
					} else {
						out.Elements = []pgtype.Int8{}
					}
				} else {
					out.Elements = (out.Elements)[:0]
				}
				for !in.IsDelim(']') {
					var v4 pgtype.Int8
					if data := in.Raw(); in.Ok() {
						in.AddError((v4).UnmarshalJSON(data))
					}
					out.Elements = append(out.Elements, v4)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "Dimensions":
			if in.IsNull() {
				in.Skip()
				out.Dimensions = nil
			} else {
				in.Delim('[')
				if out.Dimensions == nil {
					if !in.IsDelim(']') {
						out.Dimensions = make([]pgtype.ArrayDimension, 0, 8)
					} else {
						out.Dimensions = []pgtype.ArrayDimension{}
					}
				} else {
					out.Dimensions = (out.Dimensions)[:0]
				}
				for !in.IsDelim(']') {
					var v5 pgtype.ArrayDimension
					easyjson5a72dc82DecodeGithubComJackcPgxPgtype1(in, &v5)
					out.Dimensions = append(out.Dimensions, v5)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "Status":
			out.Status = pgtype.Status(in.Uint8())
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
func easyjson5a72dc82EncodeGithubComJackcPgxPgtype(out *jwriter.Writer, in pgtype.Int8Array) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Elements\":"
		out.RawString(prefix[1:])
		if in.Elements == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v6, v7 := range in.Elements {
				if v6 > 0 {
					out.RawByte(',')
				}
				out.Raw((v7).MarshalJSON())
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"Dimensions\":"
		out.RawString(prefix)
		if in.Dimensions == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v8, v9 := range in.Dimensions {
				if v8 > 0 {
					out.RawByte(',')
				}
				easyjson5a72dc82EncodeGithubComJackcPgxPgtype1(out, v9)
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"Status\":"
		out.RawString(prefix)
		out.Uint8(uint8(in.Status))
	}
	out.RawByte('}')
}
func easyjson5a72dc82DecodeGithubComJackcPgxPgtype1(in *jlexer.Lexer, out *pgtype.ArrayDimension) {
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
		case "Length":
			out.Length = int32(in.Int32())
		case "LowerBound":
			out.LowerBound = int32(in.Int32())
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
func easyjson5a72dc82EncodeGithubComJackcPgxPgtype1(out *jwriter.Writer, in pgtype.ArrayDimension) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Length\":"
		out.RawString(prefix[1:])
		out.Int32(int32(in.Length))
	}
	{
		const prefix string = ",\"LowerBound\":"
		out.RawString(prefix)
		out.Int32(int32(in.LowerBound))
	}
	out.RawByte('}')
}
