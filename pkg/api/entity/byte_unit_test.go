package entity

import (
	"reflect"
	"testing"
)

func TestParseByteUnit(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name    string
		args    args
		want    *ByteUnit
		wantErr bool
	}{
		{
			name: "100G",
			args: args{"100G"},
			want: &ByteUnit{bytes: 100 * Gigabyte},
		},
		{
			name: "12M",
			args: args{"12M"},
			want: &ByteUnit{bytes: 12 * Megabyte},
		},
		{
			name: "64K",
			args: args{"64K"},
			want: &ByteUnit{bytes: 64 * Kilobyte},
		},
		{
			name: "1234",
			args: args{"1234"},
			want: &ByteUnit{bytes: 1234},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseByteUnit(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseByteUnit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseByteUnit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestByteUnit_String(t *testing.T) {
	type fields struct {
		bytes int64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "1G",
			fields: fields{Gigabyte},
			want:   "1G",
		},
		{
			name:   "1M",
			fields: fields{Megabyte},
			want:   "1M",
		},
		{
			name:   "1K",
			fields: fields{Kilobyte},
			want:   "1K",
		},
		{
			name:   "10",
			fields: fields{10},
			want:   "10",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := ByteUnit{
				bytes: tt.fields.bytes,
			}
			if got := s.String(); got != tt.want {
				t.Errorf("ByteUnit.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
