/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8scopilot/utils"

	"github.com/spf13/cobra"
)

// eventCmd represents the event command
var eventCmd = &cobra.Command{
	Use:   "event",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		eventLog, err := getPodEnentsAndLogs()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		result, err := sendToChatGPT(eventLog)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(result)
	},
}

func sendToChatGPT(podInfo map[string][]string) (string, error) {
	client, err := utils.NewOpenAIClient()
	if err != nil {
		return "", err
	}
	combinedInfo := "The following are Pod warning events and logs: \n"
	for podName, info := range podInfo {
		combinedInfo += fmt.Sprintf("Pod: %s\n", podName)
		for _, i := range info {
			combinedInfo += fmt.Sprintf("%s\n", i)
		}
		combinedInfo += "\n"
	}
	fmt.Println(combinedInfo)
	// 构造 chatgpt 请求信息
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "You are a Kubernetes expert and you help users diagnose multiple Pod issues.",
			//Content: "您是一位 Kubernetes 专家，你要帮助用户诊断多个 Pod 问题。",
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: fmt.Sprintf("The following are multiple Pod Event events and their corresponding logs: \n%s. Please provide substantive and actionable suggestions according to Pod Log."),
			//Content: fmt.Sprintf("以下是多个 Pod Event 事件和对应的日志:\n%s\n请主要针对 Pod Log 给出实质性、可操作的建议", combinedInfo),
		},
	}
	resp, err := client.Client.CreateChatCompletion(context.TODO(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT4o,
			Messages: messages,
		})
	if err != nil {
		return "", err
	}
	responseText := resp.Choices[0].Message.Content
	return responseText, nil
}

func getPodEnentsAndLogs() (map[string][]string, error) {
	// podslist, event 事件
	clientGo, err := utils.NewClientGo(kubeconfig)
	result := make(map[string][]string)
	// 获取 Warning 级别的事件
	events, err := clientGo.ClientSet.CoreV1().Events("").List(context.TODO(), metav1.ListOptions{
		FieldSelector: "type=Warning",
	})
	if err != nil {
		return nil, fmt.Errorf("error getting events: %v", err)
	}
	for _, event := range events.Items {
		podName := event.InvolvedObject.Name
		namespace = event.InvolvedObject.Namespace
		message := event.Message

		// 获取 Pod 的日志
		if event.InvolvedObject.Kind == "Pod" {
			logOption := &corev1.PodLogOptions{}
			req := clientGo.ClientSet.CoreV1().Pods(namespace).GetLogs(podName, logOption)
			podLog, err := req.Stream(context.TODO())
			if err != nil {
				continue
			}
			defer podLog.Close()

			buf := new(bytes.Buffer)
			_, err = buf.ReadFrom(podLog)
			if err != nil {
				continue
			}
			// 如果有日志的话
			result[podName] = append(result[podName], fmt.Sprintf("Event Message: %s", message))
			result[podName] = append(result[podName], fmt.Sprintf("Namespace: %s", namespace))
			// 日志信息
			result[podName] = append(result[podName], fmt.Sprintf("Logs: %s", buf.String()))
		}
	}

	return result, nil
}

func init() {
	analyzeCmd.AddCommand(eventCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// eventCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// eventCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
