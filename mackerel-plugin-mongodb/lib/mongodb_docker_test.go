//go:build docker

package mpmongodb

import (
	"fmt"
	"strings"
	"testing"

	"github.com/fsouza/go-dockerclient"
)

// These tests do not check logger's output yet.
// You can see logger's output by `go test -v`
func TestFetchMetrics(t *testing.T) {
	port := "27017"

	for _, tag := range []string{"3.0", "3.2", "latest"} {
		opts := createContainerOptions(tag, port)

		err := invokeFunctionWithContainer(func() {
			var mongodb MongoDBPlugin = MongoDBPlugin{URL: "localhost:" + port}
			_, err := mongodb.FetchMetrics()
			if err != nil {
				t.Errorf("Error: %s", err.Error())
			}
		}, opts)

		if err != nil {
			t.Errorf("Error: %s", err.Error())
		}
	}
}

func invokeFunctionWithContainer(f func(), opts docker.CreateContainerOptions) (err error) {
	cl, err := docker.NewClientFromEnv()
	if err != nil {
		return err
	}

	// Pull Docker Image if not exists.
	_, err = cl.InspectImage(opts.Config.Image)
	image := strings.Split(opts.Config.Image, ":")
	if len(image) < 2 {
		image = append(image, "latest")
	}
	if err != nil {
		err = cl.PullImage(
			docker.PullImageOptions{
				Repository: image[0],
				Tag:        image[1],
			},
			docker.AuthConfiguration{},
		)
		if err != nil {
			return err
		}
	}

	c, err := cl.CreateContainer(opts)
	if err != nil {
		return err
	}

	defer func() {
		errRemoveContainer := cl.RemoveContainer(docker.RemoveContainerOptions{
			ID:    c.ID,
			Force: true,
		})
		if errRemoveContainer != nil {
			if err != nil {
				err = fmt.Errorf("%s; %s", err, errRemoveContainer)
			} else {
				err = errRemoveContainer
			}
		}
	}()

	if err = cl.StartContainer(c.ID, &docker.HostConfig{}); err != nil {
		return err
	}

	f()

	return nil
}

func createContainerOptions(tag, hostPort string) docker.CreateContainerOptions {
	opts := docker.CreateContainerOptions{
		Config: &docker.Config{
			Image: "mongo:" + tag,
			ExposedPorts: map[docker.Port]struct{}{
				"27017/tcp": struct{}{},
			},
		},
		HostConfig: &docker.HostConfig{
			PortBindings: map[docker.Port][]docker.PortBinding{
				"27017/tcp": {
					{HostIP: "", HostPort: hostPort},
				},
			},
		},
	}

	return opts
}
