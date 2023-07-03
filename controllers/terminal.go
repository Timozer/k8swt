package controllers

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/Timozer/k8swt/common"
	"github.com/Timozer/k8swt/k8s"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
)

func WsProcess(c *gin.Context) {
	logger := getLogger(c)

	pods, err := listPods(c)
	if err != nil {
		logger.Error().Err(err).Msg("list pods fail")
		c.JSON(500, gin.H{"msg": "internal error"})
		return
	}
	if len(pods) == 0 {
		c.JSON(404, gin.H{"msg": "not found valid pod"})
		return
	}
	if len(pods) > 1 {
		c.JSON(400, gin.H{"msg": fmt.Sprintf("found %d pods", len(pods))})
		return
	}

	pod := pods[0]

	wsConn, err := NewWsConn(c)
	if err != nil {
		logger.Error().Err(err).Msg("open websocket conn fail")
		return
	}
	// defer wsConn.Close()

	datas, _ := ioutil.ReadFile("conf/banner")
	tmp := strings.ReplaceAll(string(datas), "\n", "\r\n")
	wsConn.Write(websocket.TextMessage, []byte(tmp))

	shells := strings.Split(c.GetString(common.TERM_SHELLS), ":")
	for _, shell := range shells {
		err = RunTerminal(c.Request.Context(), wsConn, &pod, shell)
		if err != nil {
			logger.Error().Err(err).Msg("remotecommand executor fail")
			continue
		}
	}

}

func RunTerminal(ctx context.Context, conn *WsConn, pod *v1.Pod, shell string) error {
	req := k8s.GetClient().Client.CoreV1().RESTClient().Post().Resource("pods").
		Name(pod.Name).Namespace(pod.Namespace).SubResource("exec").
		VersionedParams(&v1.PodExecOptions{
			Container: pod.Status.ContainerStatuses[0].Name,
			Command:   []string{shell},
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		}, scheme.ParameterCodec)

	executor, err := remotecommand.NewSPDYExecutor(k8s.GetClient().Config, "POST", req.URL())
	if err != nil {
		return err
	}

	handler := NewStreamHandler(conn)
	return executor.StreamWithContext(ctx,
		remotecommand.StreamOptions{
			Stdin: handler, Stdout: handler, Stderr: handler,
			TerminalSizeQueue: handler, Tty: true,
		})
}
