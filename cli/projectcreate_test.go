package cli_test

import (
	"testing"

	"github.com/coder/coder/cli/clitest"
	"github.com/coder/coder/coderd/coderdtest"
	"github.com/coder/coder/database"
	"github.com/coder/coder/provisioner/echo"
	"github.com/coder/coder/pty/ptytest"
	"github.com/stretchr/testify/require"
)

func TestProjectCreate(t *testing.T) {
	t.Parallel()
	t.Run("Create", func(t *testing.T) {
		t.Parallel()
		client := coderdtest.New(t, nil)
		coderdtest.CreateFirstUser(t, client)
		source := clitest.CreateProjectVersionSource(t, &echo.Responses{
			Parse:     echo.ParseComplete,
			Provision: echo.ProvisionComplete,
		})
		cmd, root := clitest.New(t, "projects", "create", "my-project", "--directory", source, "--provisioner", string(database.ProvisionerTypeEcho))
		clitest.SetupConfig(t, client, root)
		_ = coderdtest.NewProvisionerDaemon(t, client)
		doneChan := make(chan struct{})
		pty := ptytest.New(t)
		cmd.SetIn(pty.Input())
		cmd.SetOut(pty.Output())
		go func() {
			defer close(doneChan)
			err := cmd.Execute()
			require.NoError(t, err)
		}()
		matches := []string{
			"Create project?", "yes",
		}
		for i := 0; i < len(matches); i += 2 {
			match := matches[i]
			value := matches[i+1]
			pty.ExpectMatch(match)
			pty.WriteLine(value)
		}
		<-doneChan
	})
}
