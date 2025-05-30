package container

import (
	"bytes"
	"context"
	"io"
	"log/slog"

	"github.com/docker/docker/api/types/container"
)

func (k *KeepalivedContainerCollectorHost) dockerExecCmd(cmd []string) (*bytes.Buffer, error) {
	rst, err := k.dockerCli.ContainerExecCreate(context.Background(), k.containerName, container.ExecOptions{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          cmd,
	})
	if err != nil {
		slog.Error("Error creating exec container", "CMD", cmd, "error", err)

		return nil, err
	}

	response, err := k.dockerCli.ContainerExecAttach(context.Background(), rst.ID, container.ExecStartOptions{})
	if err != nil {
		slog.Error("Error attaching a connection to an exec process", "CMD", cmd, "error", err)

		return nil, err
	}
	defer response.Close()

	data, err := io.ReadAll(response.Reader)
	if err != nil {
		slog.Error("Error reading response from docker command",
			"error", err,
			"CMD", cmd,
		)

		return nil, err
	}

	return bytes.NewBuffer(data), nil
}
