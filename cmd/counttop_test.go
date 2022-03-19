package cmd

import (
	"github.com/stevequadros/counttop/phrasecount"
	"io"
	"reflect"
	"testing"
)

func TestExecute(t *testing.T) {
	type args struct {
		readers []io.ReadCloser
		count   int
	}
	tests := []struct {
		name    string
		args    args
		want    []phrasecount.PhraseOutput
		wantErr bool
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Execute(tt.args.readers, tt.args.count)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Execute() got = %v, want %v", got, tt.want)
			}
		})
	}
}
