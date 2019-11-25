package clir

import "testing"

func TestClir(t *testing.T) {
	var mockCli *Cli
	t.Run("Run NewCli()", func(t *testing.T) {
		mockCli = NewCli("name", "description", "version")
		t.Log(mockCli)
	})

	t.Run("Run defaultBannerFunction()", func(t *testing.T) {
		err := defaultBannerFunction(mockCli)
		t.Log(err)
	})
}
