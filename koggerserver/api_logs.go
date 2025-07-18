/*
 * API for kogger
 *
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package koggerserver

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	. "github.com/Ing-Tek/Kogger/kogger"
	"github.com/Ing-Tek/Kogger/koggerrpc"
	"github.com/Ing-Tek/Kogger/models"
	"github.com/ZolaraProject/library/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func GetLogs(w http.ResponseWriter, r *http.Request) {
	ctx, grpcToken := createContextFromHeader(r)

	conn, err := grpc.Dial(fmt.Sprintf("%v:%v", KoggerHost, KoggerPort), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logger.Err(grpcToken, "LogIn could not establish gRPC connection: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		writeStandardResponse(r, w, grpcToken, fmt.Sprintf("LogIn could not establish gRPC connection: %v", err))
		return
	}
	defer conn.Close()
	client := koggerrpc.NewKoggerClient(conn)

	res, err := client.GetLogs(ctx, &koggerrpc.Void{})
	if err != nil {
		logger.Err(grpcToken, "GetLogs failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		writeStandardResponse(r, w, grpcToken, fmt.Sprintf("GetLogs failed: %s", err))
		return
	}

	responseObj := make([]models.Logs, len(res.Pods))
	for i, pod := range res.Pods {
		logResponse := models.Logs{
			Pod:       pod.Name,
			Namespace: pod.Namespace,
			Status:    pod.Status,
			Node:      pod.NodeName,
			Logs:      pod.Logs,
		}

		responseObj[i] = logResponse
	}

	logResponse := models.LogResponse{
		Logs: responseObj,
	}

	response, _ := json.Marshal(logResponse)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

var (
	grpcTokenAlphabet = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ23456789")
)

func generateGrpcToken() string {
	tk := make([]byte, 16)
	for i := range tk {
		tk[i] = grpcTokenAlphabet[rand.Intn(len(grpcTokenAlphabet))]
	}
	return string(tk)
}

func createContextFromHeader(r *http.Request) (context.Context, string) {
	grpcToken := generateGrpcToken()

	ctx := metadata.AppendToOutgoingContext(r.Context(), "grpc-token", grpcToken)

	return ctx, grpcToken
}
