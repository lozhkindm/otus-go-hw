package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	test := []struct {
		name    string
		from    string
		to      string
		offset  int64
		limit   int64
		wantErr bool
		err     error
	}{
		{
			name:    "err unsupported file",
			from:    "/dev/urandom",
			to:      "./out.txt",
			offset:  0,
			limit:   0,
			wantErr: true,
			err:     ErrUnsupportedFile,
		},
		{
			name:    "err offset exceeds filesize",
			from:    "./testdata/input.txt",
			to:      "./out.txt",
			offset:  100500,
			limit:   0,
			wantErr: true,
			err:     ErrOffsetExceedsFileSize,
		},
		{name: "out_offset0_limit0.txt", from: "./testdata/input.txt", to: "./out.txt", offset: 0, limit: 0},
		{name: "out_offset0_limit10.txt", from: "./testdata/input.txt", to: "./out.txt", offset: 0, limit: 10},
		{name: "out_offset0_limit1000.txt", from: "./testdata/input.txt", to: "./out.txt", offset: 0, limit: 1000},
		{name: "out_offset0_limit10000.txt", from: "./testdata/input.txt", to: "./out.txt", offset: 0, limit: 10000},
		{name: "out_offset100_limit1000.txt", from: "./testdata/input.txt", to: "./out.txt", offset: 100, limit: 1000},
		{name: "out_offset6000_limit1000.txt", from: "./testdata/input.txt", to: "./out.txt", offset: 6000, limit: 1000},
	}

	for _, tc := range test {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := Copy(tc.from, tc.to, tc.offset, tc.limit)

			if tc.wantErr {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)

				expectedFile, err := os.Open(fmt.Sprintf("./testdata/%s", tc.name))
				require.NoError(t, err)
				defer expectedFile.Close()
				expectedFileStat, err := expectedFile.Stat()
				require.NoError(t, err)
				expectedFileContents := make([]byte, expectedFileStat.Size())
				_, err = expectedFile.ReadAt(expectedFileContents, 0)
				require.NoError(t, err)

				resultFile, err := os.Open(tc.to)
				require.NoError(t, err)
				defer func() {
					resultFile.Close()
					os.Remove(tc.to)
				}()
				resultFileStat, err := resultFile.Stat()
				require.NoError(t, err)
				resultFileContents := make([]byte, resultFileStat.Size())
				_, err = resultFile.ReadAt(resultFileContents, 0)
				require.NoError(t, err)

				require.Equal(t, expectedFileStat.Size(), resultFileStat.Size())
				require.Equal(t, expectedFileContents, resultFileContents)
			}
		})
	}
}
