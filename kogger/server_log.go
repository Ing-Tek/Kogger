package kogger

import (
	"bufio"
	"context"
	"sync"

	. "github.com/Ing-Tek/Kogger/koggerrpc"
	grpctoken "github.com/ZolaraProject/library/grpctoken"
	logger "github.com/ZolaraProject/library/logger"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func (*server) GetLogs(ctx context.Context, in *Void) (*Pods, error) {
	grpcToken := grpctoken.GetToken(ctx)
	logger.Debug(grpcToken, "GetLogs called with %v", in)

	// list all pods with k8s client
	kubeconfig, err := rest.InClusterConfig()
	if err != nil {
		logger.Err(grpcToken, "Failed to retrieve in-cluster Kubernetes config: %s", err)
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		logger.Err(grpcToken, "Failed to initialize Kubernetes client: %s", err)
		return nil, err
	}

	allpods, err := clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil {
		logger.Err(grpcToken, "Failed to list pods: %s", err)
		return nil, err
	}

	podLogsChan := make(chan *Pod, len(allpods.Items))
	errChan := make(chan error, len(allpods.Items))
	wg := &sync.WaitGroup{}
	wg.Add(len(allpods.Items))
	for _, pod := range allpods.Items {
		wg.Add(1)
		defer wg.Done()

		logger.Debug(grpcToken, "Pod found: %s in namespace %s", pod.Name, pod.Namespace)
		go func(pod v1.Pod) {
			req := clientset.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &v1.PodLogOptions{})
			podLogs, err := req.Stream(context.TODO())
			if err != nil {
				logger.Err(grpcToken, "Failed to get logs for pod %s in namespace %s: %s", pod.Name, pod.Namespace, err)
				errChan <- err
				return
			}
			defer podLogs.Close()

			logs := ""
			scanner := bufio.NewScanner(podLogs)
			for scanner.Scan() {
				logs += scanner.Text() + "\n"
			}

			if err := scanner.Err(); err != nil {
				logger.Err(grpcToken, "Error reading logs for pod %s in namespace %s: %s", pod.Name, pod.Namespace, err)
				errChan <- err
				return
			}

			podLogsChan <- &Pod{
				Name:      pod.Name,
				Namespace: pod.Namespace,
				Status:    string(pod.Status.Phase),
				NodeName:  pod.Spec.NodeName,
				Logs:      logs,
			}
		}(pod)
	}
	wg.Wait()

	close(errChan)
	close(podLogsChan)
	if len(errChan) > 0 {
		for err := range errChan {
			if err != nil {
				logger.Err(grpcToken, "Error occurred while fetching pod logs: %s", err)
				return nil, err
			}
		}
	}

	pods := &Pods{}
	for podLog := range podLogsChan {
		if podLog != nil {
			pods.Pods = append(pods.Pods, podLog)
		}
	}
	logger.Debug(grpcToken, "Returning %d pods with logs", len(pods.Pods))
	if len(pods.Pods) == 0 {
		logger.Debug(grpcToken, "No pods found with logs")
		return &Pods{}, nil
	}

	return pods, nil
}
