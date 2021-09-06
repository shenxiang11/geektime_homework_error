package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type mysqlDockerCli struct {
	*client.Client
	context        context.Context
	mysqlContainer container.ContainerCreateCreatedBody
}

func (c *mysqlDockerCli) Start() error {
	return c.ContainerStart(c.context, c.mysqlContainer.ID, types.ContainerStartOptions{})
}

func (c *mysqlDockerCli) Remove() error {
	return c.ContainerRemove(c.context, c.mysqlContainer.ID, types.ContainerRemoveOptions{
		Force: true,
	})
}

func NewMysqlDockerCli() (*mysqlDockerCli, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	mysqlContainer, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image: "mysql:5.7",
			ExposedPorts: nat.PortSet{
				"3306/tcp": {},
			},
			Env: []string{
				"MYSQL_USER=docker",
				"MYSQL_PASSWORD=123456",
				"MYSQL_ROOT_PASSWORD=654321",
				"MYSQL_DATABASE=homework",
			},
		},
		&container.HostConfig{
			PortBindings: nat.PortMap{
				"3306/tcp": []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: "33061",
					},
				},
			},
		},
		nil,
		nil,
		"dbhomework",
	)
	if err != nil {
		return nil, err
	}

	return &mysqlDockerCli{
		Client:         cli,
		context:        ctx,
		mysqlContainer: mysqlContainer,
	}, nil
}
