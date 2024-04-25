package mock

import (
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStart(t *testing.T) {
	c := &Command{}
	assert.NoError(t, c.Start())
}

func TestWait(t *testing.T) {
	t.Run("No Error", func(t *testing.T) {
		c := &Command{}
		assert.NoError(t, c.Wait())
	})
	t.Run("An Error", func(t *testing.T) {
		c := &Command{waitErr: errors.New("some error")}
		assert.ErrorContains(t, c.Wait(), "some error")
	})
}

func TestString(t *testing.T) {
	c := &Command{}
	assert.Empty(t, c.String())
}

func TestStdoutPipe(t *testing.T) {
	c := &Command{}
	result, err := c.StdoutPipe()
	require.NoError(t, err)

	b, err := io.ReadAll(result)
	assert.NoError(t, err)
	assert.Empty(t, b)
}
func TestStderrPipe(t *testing.T) {
	c := &Command{}
	result, err := c.StderrPipe()
	require.NoError(t, err)

	b, err := io.ReadAll(result)
	assert.NoError(t, err)
	assert.Empty(t, b)
}

func TestCommandBuilderWithError(t *testing.T) {
	result := CommandBuilderWithError(errors.New("some error"), nil, nil, nil)
	assert.ErrorContains(t, result("start").Wait(), "some error")
}

func TestThresholdError(t *testing.T) {
	th := &ThresholdError{}
	assert.Equal(t, th.ExitCode(), thresholdsBreachedExitCode)
	assert.Contains(t, th.Error(), "mock")
}
