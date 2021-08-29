package main

import (
	"testing"

	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_initApp(t *testing.T) {
	type args struct {
		configPath string
	}
	tests := []struct {
		name      string
		args      args
		want      selfApp
		assertion assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := initApp(tt.args.configPath)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_getTitleForNote(t *testing.T) {
	type args struct {
		content string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, getTitleForNote(tt.args.content))
		})
	}
}

func Test_getSubContentForNote(t *testing.T) {
	type args struct {
		n *note
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, getSubContentForNote(tt.args.n))
		})
	}
}

func Test_selfApp_newNote(t *testing.T) {
	type fields struct {
		notes        []*note
		noteList     *tview.List
		notesChanged bool
		app          *tview.Application
		config       appConfigs
	}
	tests := []struct {
		name      string
		fields    fields
		assertion assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			self := &selfApp{
				notes:        tt.fields.notes,
				noteList:     tt.fields.noteList,
				notesChanged: tt.fields.notesChanged,
				app:          tt.fields.app,
				config:       tt.fields.config,
			}
			tt.assertion(t, self.newNote())
		})
	}
}

func Test_selfApp_openNote(t *testing.T) {
	type fields struct {
		notes        []*note
		noteList     *tview.List
		notesChanged bool
		app          *tview.Application
		config       appConfigs
	}
	type args struct {
		index int
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			self := &selfApp{
				notes:        tt.fields.notes,
				noteList:     tt.fields.noteList,
				notesChanged: tt.fields.notesChanged,
				app:          tt.fields.app,
				config:       tt.fields.config,
			}
			tt.assertion(t, self.openNote(tt.args.index))
		})
	}
}

func Test_sortByModTime(t *testing.T) {
	type args struct {
		notes []*note
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sortByModTime(tt.args.notes)
		})
	}
}

func Test_selfApp_loadNotes(t *testing.T) {
	testapp, err := initApp(getFileFromConfig("config.ini"))
	require.Nil(t, err, "failed to init app")

	type fields struct {
		notes        []*note
		noteList     *tview.List
		notesChanged bool
		app          *tview.Application
		config       appConfigs
	}
	tests := []struct {
		name      string
		fields    fields
		assertion assert.ErrorAssertionFunc
	}{
		// TODO: these test dont test anything...
		{
			name: "basic_s3",
			fields: fields{
				notes:        []*note{},
				noteList:     nil,
				notesChanged: false,
				app:          nil,
				config:       testapp.config, // Use the current users config
			},
			assertion: assert.NoError,
		},
		{
			name: "basic_local",
			fields: fields{
				notes:        []*note{},
				noteList:     nil,
				notesChanged: false,
				app:          nil,
				config: appConfigs{
					s3Active: false,
				},
			},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			self := &selfApp{
				notes:        tt.fields.notes,
				noteList:     tt.fields.noteList,
				notesChanged: tt.fields.notesChanged,
				app:          tt.fields.app,
				config:       tt.fields.config,
			}
			tt.assertion(t, self.loadNotes())
		})
	}
}

func Test_selfApp_saveNotes(t *testing.T) {
	type fields struct {
		notes        []*note
		noteList     *tview.List
		notesChanged bool
		app          *tview.Application
		config       appConfigs
	}
	tests := []struct {
		name      string
		fields    fields
		assertion assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			self := &selfApp{
				notes:        tt.fields.notes,
				noteList:     tt.fields.noteList,
				notesChanged: tt.fields.notesChanged,
				app:          tt.fields.app,
				config:       tt.fields.config,
			}
			tt.assertion(t, self.saveNotes())
		})
	}
}
