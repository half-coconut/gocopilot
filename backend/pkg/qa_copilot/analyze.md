添加 analyze 命令

```shell
cobra-cli add analyze
cobra-cli add event -p 'analyzeCmd'

```

获取 Warning 级别的日志

```shell
chenchen@chenchendeMacBook-Pro-2 k8scopilot % go build -o k8scopilot    
chenchen@chenchendeMacBook-Pro-2 k8scopilot % ./k8scopilot analyze event
map[demo-1-6d9fbb798-4tjvs:[Event Message: Back-off restarting failed container demo-1 in pod demo-1-6d9fbb798-4tjvs_default(a4c4210d-b083-4bb7-a89d-20d3ee400a18) Namespace: default Logs: demo-1-6d9fbb798-4tjvs
nginx-deployment-96b9d695-8gqsm
nginx-deployment-96b9d695-g8sts
nginx-deployment-96b9d695-klw8v
nginx-pod
demo-1
nginx-deployment
]]

```

给出建议

```shell
chenchen@chenchendeMacBook-Pro-2 k8scopilot % go build -o k8scopilot    
chenchen@chenchendeMacBook-Pro-2 k8scopilot % ./k8scopilot analyze event
find the following Pod Warning Events and Logs: 
Pod: demo-1-6d9fbb798-4tjvs
Event Message: Back-off restarting failed container demo-1 in pod demo-1-6d9fbb798-4tjvs_default(a4c4210d-b083-4bb7-a89d-20d3ee400a18)
Namespace: default
Logs: demo-1-6d9fbb798-4tjvs
nginx-deployment-96b9d695-8gqsm
nginx-deployment-96b9d695-g8sts
nginx-deployment-96b9d695-klw8v
nginx-pod
demo-1
nginx-deployment



Certainly! However, it seems like there was an error in the way the information was presented ("%!s(MISSING)"). To provide practical and feasible suggestions, I would need a bit more detail on the specific pod events and logs you're seeing.

Here are some general troubleshooting steps and common scenarios for Kubernetes pod issues:

1. **Pod Pending:**
   - **Insufficient Resources:** Check if there are enough resources (CPU/memory) available in the cluster.
   - **Scheduling Issues:** Ensure there are nodes that satisfy the node selector, affinity rules, and taints/tolerations.
   - **Persistent Volume Claims:** Verify that any required persistent volumes are available and properly bound.

2. **Pod CrashLoopBackOff/Failed:**
   - **Application Errors:** Check the application logs for errors specific to the container application.
   - **Configuration Errors:** Double-check your environment variables, config maps, etc., for any misconfigurations.
   - **Dependencies:** Ensure all required services and endpoints are accessible from the pod's network.

3. **Image Pull Issues:**
   - **Image Not Found:** Verify the image name and tag. Ensure it's available in the specified registry.
   - **Authentication Issues:** Check if there’s a need for secret to pull images from private registries.
   - **Network/Firewall:** Ensure your cluster nodes can reach the image registry, possibly checking firewall rules.

4. **OOMKilled:**
   - **Resource Limits:** The container might be exceeding its memory limit. Consider increasing the memory limit or optimizing the application to use less memory.
   - **Memory Leaks:** Investigate the application for memory leaks or inefficient memory usage patterns.

5. **Networking Issues:**
   - **Service Discovery:** Ensure all services and endpoints are correctly configured and DNS resolution is working.
   - **Network Policies:** Check network policies that might be preventing connections.

6. **Pod Stuck in Terminating:**
   - **Finalizers:** Check if there are finalizers blocking the deletion and whether they should be removed.
   - **Hanging Processes:** Investigate if there are stuck processes within the container that prevent it from shutting down cleanly.

Please share the specific logs and events if available, and I can provide more targeted advice based on the actual outputs you're dealing with.

```

不同的提示词下，给的不同的建议：

```shell

chenchen@chenchendeMacBook-Pro-2 k8scopilot % go build -o k8scopilot    
chenchen@chenchendeMacBook-Pro-2 k8scopilot % ./k8scopilot analyze event
The following are Pod warning events and logs: 
Pod: demo-1-6d9fbb798-57tx9
Event Message: Back-off restarting failed container demo-1 in pod demo-1-6d9fbb798-57tx9_default(461a9f80-3b13-4f4b-b4f5-511a4d425e2e)
Namespace: default
Logs: deployments.apps is forbidden: User "system:serviceaccount:default:default" cannot list resource "deployments" in API group "apps" in the namespace "default": RBAC: role.rbac.authorization.k8s.io "demo" not found



根据 Pod 的日志信息，问题似乎与权限不足导致无法访问资源有关。具体来说，Pod 中运行的容器尝试列出命名空间 `default` 中的 `deployments` 资源，但是由于缺少适当的RBAC权限，被禁止执行此操作。

以下是一些可操作的步骤来解决这个问题：

1. **检查ServiceAccount**: Pod 使用了 `default` ServiceAccount。考虑为此应用创建一个专门的 ServiceAccount，以避免在 `default` ServiceAccount 中进行不必要的更改。

2. **创建或更新Role**:
   - 若 Role 不存在：创建一个允许 `list` 访问 `deployments` 资源的 Role。例如：
     ```yaml
     apiVersion: rbac.authorization.k8s.io/v1
     kind: Role
     metadata:
       namespace: default
       name: demo-role
     rules:
     - apiGroups: ["apps"]
       resources: ["deployments"]
       verbs: ["get", "list", "watch"]
     ```

   - 若 Role 已存在但缺少权限：确保它有 `list` 权限。

3. **绑定Role与ServiceAccount**:
   - 创建一个 RoleBinding，用于将上述 Role 绑定到特定的 ServiceAccount（假设新建的 ServiceAccount 为 `demo-sa`）。
     ```yaml
     apiVersion: rbac.authorization.k8s.io/v1
     kind: RoleBinding
     metadata:
       name: demo-rolebinding
       namespace: default
     subjects:
     - kind: ServiceAccount
       name: demo-sa
       namespace: default
     roleRef:
       kind: Role
       name: demo-role
       apiGroup: rbac.authorization.k8s.io
     ```

4. **更新Pod的ServiceAccount**:
   - 更新Pod的配置以使用新的 ServiceAccount `demo-sa`：
     ```yaml
     apiVersion: v1
     kind: Pod
     metadata:
       name: demo-1
       namespace: default
     spec:
       serviceAccountName: demo-sa
       containers:
       - name: demo-1
         image: your-image
     ```

5. **重新部署应用**:
   - 重新应用修改后的配置，确保每个 Pod 在使用新的 ServiceAccount。

完成以上步骤后，Pod 应该能够正确地列出 `deployments` 资源。如果问题仍然存在，请检查其他可能的配置错误或应用程序逻辑。

```