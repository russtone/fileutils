package fileutils_test

import (
	"testing"

	"github.com/russtone/fileutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLinesCount(t *testing.T) {
	tests := []struct {
		file  string
		count int
	}{
		{"lines3_no_newline", 3},
		{"lines3", 3},
		{"lines5", 5},
		{"lines10", 10},
		{"lines1000", 1000},
	}

	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			count, err := fileutils.LinesCount("test/" + tt.file)
			require.NoError(t, err)
			assert.Equal(t, tt.count, count)
		})
	}
}

func TestLinesCount_Error(t *testing.T) {
	_, err := fileutils.LinesCount("test/invalid")
	assert.Error(t, err)
}

func TestFirstLine(t *testing.T) {
	tests := []struct {
		file string
		line string
	}{
		{"lines3_no_newline", "1"},
		{"lines3", "1"},
		{"lines5", "1"},
		{"lines10", "1"},
		{"lines1000", "1"},
	}

	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			line, err := fileutils.FirstLine("test/" + tt.file)
			require.NoError(t, err)
			assert.Equal(t, tt.line, line)
		})
	}
}

func TestLinesIterator(t *testing.T) {

	tests := []struct {
		file  string
		count uint64
		res   []string
	}{
		{
			"lines3",
			3,
			[]string{
				"1",
				"2",
				"3",
			},
		},
		{
			"lines3_no_newline",
			3,
			[]string{
				"1",
				"2",
				"3",
			},
		},
		{
			"lines5",
			5,
			[]string{
				"1",
				"2",
				"3",
				"4",
				"5",
			},
		},
		{
			"lines10",
			10,
			[]string{
				"1",
				"2",
				"3",
				"4",
				"5",
				"6",
				"7",
				"8",
				"9",
				"10",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			it, err := fileutils.NewLinesIterator("test/" + tt.file)
			require.NoError(t, err)

			assert.Equal(t, tt.count, it.Count())

			var line string

			res := make([]string, 0)
			for it.Next(&line) {
				res = append(res, line)
			}

			assert.Equal(t, tt.res, res)

			it.Reset()

			res2 := make([]string, 0)
			for it.Next(&line) {
				res2 = append(res2, line)
			}

			assert.Equal(t, tt.res, res2)

			assert.NoError(t, it.Close())
		})
	}
}
