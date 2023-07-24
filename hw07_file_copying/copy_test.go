package main

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	temp, err := os.CreateTemp("/tmp", "hw07.")
	require.NoError(t, err)

	fromFile := "./testdata/input.txt"
	fromStat, err := os.Stat(fromFile)
	require.NoError(t, err)
	fromAllFile, err := os.ReadFile(fromFile)
	require.NoError(t, err)

	toFile := temp.Name()

	defer func(temp *os.File) {
		_ = os.Remove(temp.Name())
	}(temp)

	t.Run("Copy full file", func(t *testing.T) {
		err := Copy(fromFile, toFile, 0, 0)
		require.NoError(t, err)
		toStat, err := temp.Stat()
		require.NoError(t, err)

		require.Equal(t, fromStat.Size(), toStat.Size())
	})

	t.Run("Copy part file", func(t *testing.T) {
		err := Copy(fromFile, toFile, 0, 10)
		require.NoError(t, err)
		toStat, err := temp.Stat()
		require.NoError(t, err)

		toAllFile, err := os.ReadFile(toFile)
		require.NoError(t, err)

		require.Equal(t, int64(10), toStat.Size())
		require.Equal(t, fromAllFile[0:10], toAllFile)
	})

	t.Run("Copy file with offset", func(t *testing.T) {
		err := Copy(fromFile, toFile, 15, 0)
		require.NoError(t, err)
		toStat, err := temp.Stat()
		require.NoError(t, err)

		require.Equal(t, fromStat.Size()-15, toStat.Size())
	})

	t.Run("Open no such file", func(t *testing.T) {
		err := Copy("test12345", toFile, 0, 0)
		require.Error(t, err)
	})

	t.Run("Open directory", func(t *testing.T) {
		err := Copy("./testdata/", toFile, 0, 0)
		require.Error(t, err)
		require.Truef(t, errors.Is(err, ErrUnsupportedFile), "actual err - %v", err)
	})

	t.Run("Offset is more file size", func(t *testing.T) {
		err := Copy(fromFile, toFile, 10000, 0)
		require.Error(t, err)
		require.Truef(t, errors.Is(err, ErrOffsetExceedsFileSize), "actual err - %v", err)
	})

	t.Run("Check urandom file", func(t *testing.T) {
		err := Copy("/dev/urandom", toFile, 0, 0)
		require.Error(t, err)
		require.Truef(t, errors.Is(err, ErrUnsupportedFile), "actual err - %v", err)
	})
}
