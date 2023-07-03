package controllers

import (
	"github.com/Timozer/k8swt/common"
	"github.com/Timozer/k8swt/k8s"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getLogger(c *gin.Context) *zerolog.Logger {
	loggerI, _ := c.Get(common.CTX_LOGGER)
	logger, _ := loggerI.(*zerolog.Logger)
	return logger
}

func listPods(c *gin.Context) ([]v1.Pod, error) {
	logger := getLogger(c)

	namespace := c.DefaultQuery("namespace", "")
	podName := c.DefaultQuery("podname", "")
	podIP := c.DefaultQuery("ip", "")

	logger.Debug().Str("Namespace", namespace).Str("PodName", podName).Str("IP", podIP).Msg("")

	selector := ""
	if len(podName) > 0 {
		selector = "metadata.name=" + podName
	} else if len(podIP) > 0 {
		selector = "status.podIP=" + podIP
	}

	return k8s.GetClient().ListPods(
		c.Request.Context(), namespace, metav1.ListOptions{FieldSelector: selector},
	)
}
