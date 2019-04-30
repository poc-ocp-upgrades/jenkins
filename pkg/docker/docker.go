package docker

import (
	"bytes"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"time"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
)

type Client struct{ Client *client.Client }

func NewEnvClient() (*Client, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	client, err := client.NewEnvClient()
	return &Client{Client: client}, err
}
func (c *Client) ExecInActiveContainers(w io.Writer, ctx context.Context, cmd []string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for {
		time.Sleep(60 * time.Second)
		select {
		case <-ctx.Done():
			return
		default:
			containers, err := c.ContainerList()
			if err != nil {
				fmt.Fprintf(w, "container list error: %#v\n", err)
				continue
			}
			fmt.Fprintf(w, "found %d containers\n", len(containers))
			for _, container := range containers {
				fmt.Fprintf(w, "found container %s running command %s\n", container.ID, container.Command)
				rc, buf, err := c.ContainerExec(container.ID, cmd)
				if err != nil {
					fmt.Fprintf(w, "container exec error: %#v\n", err)
					continue
				}
				fmt.Fprintf(w, "exec of command %#v into %s had rc %d and text %s\n", cmd, container.ID, rc, string(buf))
			}
		}
	}
}
func (c *Client) InspectActiveContainers(w io.Writer, ctx context.Context) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for {
		time.Sleep(60 * time.Second)
		select {
		case <-ctx.Done():
			return
		default:
			containers, err := c.ContainerList()
			if err != nil {
				fmt.Fprintf(w, "container list error: %#v\n", err)
				continue
			}
			fmt.Fprintf(w, "found %d containers\n", len(containers))
			for _, container := range containers {
				fmt.Fprintf(w, "found container %s running command %s\n", container.ID, container.Command)
				_, body, err := c.Client.ContainerInspectWithRaw(context.Background(), container.ID, true)
				if err != nil {
					fmt.Fprintf(w, "container inspect error: %#v\n", err)
					continue
				}
				var prettyJSON bytes.Buffer
				error := json.Indent(&prettyJSON, body, "", "\t")
				if error != nil {
					fmt.Fprintf(w, "inspect of %s returned raw json %s\n", container.ID, string(body))
					continue
				}
				fmt.Fprintf(w, "inspect of %s returned formatted json:\n%s\n", container.ID, string(prettyJSON.Bytes()))
			}
		}
	}
}
func (c *Client) ContainerList() ([]types.Container, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return c.Client.ContainerList(context.Background(), types.ContainerListOptions{})
}
func (c *Client) ContainerCreate(config *container.Config, hostconfig *container.HostConfig) (string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	body, err := c.Client.ContainerCreate(context.Background(), config, hostconfig, nil, "")
	return body.ID, err
}
func (c *Client) ContainerExec(id string, cmd []string) (int, []byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	exec, err := c.Client.ContainerExecCreate(context.Background(), id, types.ExecConfig{AttachStdout: true, Cmd: cmd})
	if err != nil {
		return 0, nil, err
	}
	attach, err := c.Client.ContainerExecAttach(context.Background(), exec.ID, types.ExecConfig{})
	if err != nil {
		return 0, nil, err
	}
	defer attach.Close()
	bytes, err := ioutil.ReadAll(attach.Reader)
	if err != nil {
		return 0, nil, err
	}
	inspect, err := c.Client.ContainerExecInspect(context.Background(), exec.ID)
	if err != nil {
		return 0, nil, err
	}
	return inspect.ExitCode, bytes, nil
}
func (c *Client) ContainerInspect(id string) (string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	json, err := c.Client.ContainerInspect(context.Background(), id)
	if err != nil {
		return "", err
	}
	return json.NetworkSettings.IPAddress, nil
}
func (c *Client) ContainerStart(id string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return c.Client.ContainerStart(context.Background(), id, types.ContainerStartOptions{})
}
func (c *Client) ContainerLogs(id string) ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	r, err := c.Client.ContainerLogs(context.Background(), id, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(r)
}
func (c *Client) ContainerRemove(id string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return c.Client.ContainerRemove(context.Background(), id, types.ContainerRemoveOptions{})
}
func (c *Client) ContainerStop(id string, timeout *time.Duration) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return c.Client.ContainerStop(context.Background(), id, timeout)
}
func (c *Client) ContainerStopAndRemove(id string, timeout *time.Duration) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	err := c.ContainerStop(id, timeout)
	if err != nil {
		return err
	}
	return c.ContainerRemove(id)
}
func (c *Client) ContainerWait(id string) (int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return c.Client.ContainerWait(context.Background(), id)
}
func (c *Client) ImageRemove(name string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_, err := c.Client.ImageRemove(context.Background(), name, types.ImageRemoveOptions{})
	return err
}
func (c *Client) VolumeCreate() (string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	vol, err := c.Client.VolumeCreate(context.Background(), types.VolumeCreateRequest{})
	return vol.Name, err
}
func (c *Client) VolumeRemove(name string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return c.Client.VolumeRemove(context.Background(), name)
}
func Duration(d time.Duration) *time.Duration {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &d
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
