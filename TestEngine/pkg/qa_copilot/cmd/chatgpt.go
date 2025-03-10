/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"backend/pkg/qa_copilot/utils"
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/restmapper"
	"k8s.io/kubectl/pkg/scheme"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// chatgptCmd represents the chatgpt command
var chatgptCmd = &cobra.Command{
	Use:   "chatgpt",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		startChat()
	},
}

func startChat() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("I'm test cases generator copilot, what can I do for you?")

	for {
		fmt.Printf("> ")
		if scanner.Scan() {
			input := scanner.Text()
			if input == "exit" {
				fmt.Println("Bye!")
				break
			}
			if input == "" {
				continue
			}
			response := processInput(input)
			fmt.Println(response)
		}
	}
}

func processInput(input string) string {
	client, err := utils.NewOpenAIClient()
	if err != nil {
		return err.Error()
	}
	//response, err := client.SendMessage("Hi, you are a k8s assistant, you can produce yaml content, just output the content from yaml, do not put yaml in the code block", input)
	response := functionCalling(input, client)
	return response
}

func functionCalling(input string, client *utils.OpenAI) string {
	// 定义第一个函数，生成 k8s yaml 并部署资源
	f1 := openai.FunctionDefinition{
		Name:        "generateAndDeploymentResource",
		Description: "produce k8s yaml, and deploy",
		Parameters: jsonschema.Definition{
			Type: jsonschema.Object,
			Properties: map[string]jsonschema.Definition{
				"user_input": {
					Type:        jsonschema.String,
					Description: "raw content from user input must include resource type and image",
				},
			},
			Required: []string{"user_input"},
		},
	}

	t1 := openai.Tool{
		Type:     openai.ToolTypeFunction,
		Function: &f1,
	}

	// 定义第二个函数，查询k8s资源
	f2 := openai.FunctionDefinition{
		Name:        "queryResource",
		Description: "get k8s resource",
		Parameters: jsonschema.Definition{
			Type: jsonschema.Object,
			Properties: map[string]jsonschema.Definition{
				"namespace": {
					Type:        jsonschema.String,
					Description: "the namespace where the resource is located",
				},
				"resource_type": {
					Type:        jsonschema.String,
					Description: "resource type, such as Pod, Deployment, Service, etc.",
				},
			},
			Required: []string{"namespace", "resource_type"},
		},
	}

	t2 := openai.Tool{
		Type:     openai.ToolTypeFunction,
		Function: &f2,
	}
	//
	f3 := openai.FunctionDefinition{
		Name:        "deleteResource",
		Description: "get k8s resource",
		Parameters: jsonschema.Definition{
			Type: jsonschema.Object,
			Properties: map[string]jsonschema.Definition{
				"namespace": {
					Type:        jsonschema.String,
					Description: "the namespace where the resource is located",
				},
				"resource_type": {
					Type:        jsonschema.String,
					Description: "resource type, such as Pod, Deployment, Service, etc.",
				},
				"resource_name": {
					Type:        jsonschema.String,
					Description: "resource name",
				},
			},
			Required: []string{"namespace", "resource_type"},
		},
	}

	t3 := openai.Tool{
		Type:     openai.ToolTypeFunction,
		Function: &f3,
	}

	dialogue := []openai.ChatCompletionMessage{
		{Role: openai.ChatMessageRoleUser, Content: input},
	}

	resp, err := client.Client.CreateChatCompletion(context.TODO(), openai.ChatCompletionRequest{
		Model:    openai.GPT4o,
		Messages: dialogue,
		Tools:    []openai.Tool{t1, t2, t3},
	})
	if err != nil {
		return err.Error()
	}

	msg := resp.Choices[0].Message
	if len(msg.ToolCalls) != 1 {
		return fmt.Sprintf("can not find appropirate tools to call, %v", len(msg.ToolCalls))
	}
	// 组装历史对话
	dialogue = append(dialogue, msg)
	//return fmt.Sprintf("OpenAI hopes to request the function: %s, parameters: %s", msg.ToolCalls[0].Function.Name, msg.ToolCalls[0].Function.Arguments)
	result, err := callFunction(client, msg.ToolCalls[0].Function.Name, msg.ToolCalls[0].Function.Arguments)
	if err != nil {
		return fmt.Sprintf("Error calling function: %v\n", err)
	}
	return result
}

func callFunction(client *utils.OpenAI, name, arguments string) (string, error) {
	if name == "generateAndDeploymentResource" {
		params := struct {
			UserInput string `json:"user_input"`
		}{}
		if err := json.Unmarshal([]byte(arguments), &params); err != nil {
			return "", err
		}
		return generateAndDeploymentResource(client, params.UserInput)
	}
	if name == "queryResource" {
		params := struct {
			Namespace    string `json:"namespace"`
			ResourceType string `json:"resource_type"`
		}{}
		if err := json.Unmarshal([]byte(arguments), &params); err != nil {
			return "", err
		}
		return queryResource(params.Namespace, params.ResourceType)
	}
	if name == "deleteResource" {
		params := struct {
			Namespace    string `json:"namespace"`
			ResourceType string `json:"resource_type"`
			ResourceName string `json:"resource_name"`
		}{}
		if err := json.Unmarshal([]byte(arguments), &params); err != nil {
			return "", err
		}
		return deleteResource(client, params.Namespace, params.ResourceType, params.ResourceName)
	}
	return "", fmt.Errorf("unknown function: %s", name)
}

func generateAndDeploymentResource(client *utils.OpenAI, userInput string) (string, error) {
	yamlContent, err := client.SendMessage("You're a k8s resource generator, please produce k8s yaml according to user_input, just output content from yaml, do not put yaml in code block.", userInput)
	if err != nil {
		return "", err
	}
	//return yamlContent, nil
	// 调用 dynamic client 创建资源
	clientGo, err := utils.NewClientGo(kubeconfig)
	if err != nil {
		return "", err
	}
	resources, err := restmapper.GetAPIGroupResources(clientGo.DiscoveryClient)
	if err != nil {
		return "", err
	}
	// 把 yaml 转成 unstructured
	unstructuredObj := &unstructured.Unstructured{}
	_, _, err = scheme.Codecs.UniversalDeserializer().Decode([]byte(yamlContent), nil, unstructuredObj)
	if err != nil {
		return "", err
	}
	// 创建 mapper, GVK->GVR
	mapper := restmapper.NewDiscoveryRESTMapper(resources)
	// 从 ustructuredObj 里获取 GVK
	gvk := unstructuredObj.GroupVersionKind()
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return "", err
	}
	namespace = unstructuredObj.GetNamespace()
	if namespace == "" {
		namespace = "default"
	}
	// 使用 deployment 创建资源
	_, err = clientGo.DynamicClient.Resource(mapping.Resource).Namespace(namespace).Create(context.TODO(), unstructuredObj, metav1.CreateOptions{})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("YAML content:\n%s\n\nDeployment successful.", yamlContent), nil
}

func queryResource(namespace, resourceType string) (string, error) {
	clientGo, err := utils.NewClientGo(kubeconfig)
	resourceType = strings.ToLower(resourceType)
	var gvr schema.GroupVersionResource
	switch resourceType {
	case "deployment":
		gvr = schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	case "service":
		gvr = schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}
	case "pod":
		gvr = schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}
	default:
		return "", fmt.Errorf("unsupported resource type: %s", resourceType)
	}
	// 通过 dynamicClient 获取资源
	resourceList, err := clientGo.DynamicClient.Resource(gvr).Namespace(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to list resources: %w", err)
	}

	// Iterate over the resources and print their names (or handle them as needed)
	result := ""
	for _, item := range resourceList.Items {
		result += fmt.Sprintf("Found %s: %s\n", resourceType, item.GetName())
	}

	return result, nil
}

func deleteResource(client *utils.OpenAI, n string, resourceType string, name string) (string, error) {
	return "", nil
}

func init() {
	askCmd.AddCommand(chatgptCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// chatgptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// chatgptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
