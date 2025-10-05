package cli

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TestRunEncryptWithStubbedPassphrase(t *testing.T) {
	origReadPassword := readPassword
	defer func() { readPassword = origReadPassword }()

	callCount := 0
	passphrase := []byte("TestingOneTwoThree")

	readPassword = func(int) ([]byte, error) {
		callCount++
		buf := make([]byte, len(passphrase))
		copy(buf, passphrase)
		return buf, nil
	}

	forceCompressed = false
	forceUncompressed = false

	cmd := &cobra.Command{Use: "encrypt"}

	origStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w

	runErr := runEncrypt(cmd, []string{"5HueCGU8rMjxEXxiPuD5BDku4MkFqeZyd4dZ1jvhTVqvbTLvyTJ"})

	w.Close()
	os.Stdout = origStdout
	outputBytes, _ := io.ReadAll(r)

	if runErr != nil {
		t.Fatalf("runEncrypt returned error: %v", runErr)
	}

	if callCount != 2 {
		t.Fatalf("expected readPassword to be called twice, got %d", callCount)
	}

	output := string(outputBytes)
	if !strings.Contains(output, "Encrypted key: 6P") {
		t.Fatalf("expected encrypted key output, got %q", output)
	}
}

func TestRunGenerateIntermediateRequiresLotAndSequence(t *testing.T) {
	origReadPassword := readPassword
	defer func() { readPassword = origReadPassword }()

	readPassword = func(int) ([]byte, error) {
		return []byte("pass"), nil
	}

	lotNumber = 0
	sequenceNumber = 0
	useLotSeq = false
	defer func() {
		lotNumber = 0
		sequenceNumber = 0
		useLotSeq = false
	}()

	cmd := &cobra.Command{Use: "intermediate-generate"}
	cmd.Flags().Uint32Var(&lotNumber, "lot", 0, "")
	cmd.Flags().Uint32Var(&sequenceNumber, "sequence", 0, "")
	cmd.Flags().BoolVar(&useLotSeq, "use-lot-sequence", false, "")

	if err := cmd.Flags().Set("lot", "123"); err != nil {
		t.Fatalf("failed to set lot flag: %v", err)
	}

	err := runGenerateIntermediate(cmd, nil)
	if err == nil {
		t.Fatal("expected error when only --lot is provided")
	}

	if !strings.Contains(err.Error(), "both --lot and --sequence") {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func TestIsVerboseUsesViperConfig(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)

	viper.Set("verbose", true)

	cmd := &cobra.Command{Use: "root"}

	if !isVerbose(cmd) {
		t.Fatal("expected verbose to be true when set via config")
	}

	viper.Set("verbose", false)
	cmd.Flags().Bool("verbose", false, "")
	if err := cmd.Flags().Set("verbose", "true"); err != nil {
		t.Fatalf("failed to set verbose flag: %v", err)
	}

	if !isVerbose(cmd) {
		t.Fatal("expected verbose flag to take precedence")
	}
}
