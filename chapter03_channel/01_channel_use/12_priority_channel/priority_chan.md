# select中实现channel优先级
## 需求背景
需求：我们有一个函数会持续不间断地从ch1和ch2中分别接收任务1和任务2，

```go
// kubernetes/pkg/controller/nodelifecycle/scheduler/taint_manager.go 
func (tc *NoExecuteTaintManager) worker(worker int, done func(), stopCh <-chan struct{}) {
	defer done()

	// 当处理具体事件的时候，我们会希望 Node 的更新操作优先于 Pod 的更新
	// 因为 NodeUpdates 与 NoExecuteTaintManager无关应该尽快处理
	// -- 我们不希望用户(或系统)等到PodUpdate队列被耗尽后，才开始从受污染的Node中清除pod。
	for {
		select {
		case <-stopCh:
			return
		case nodeUpdate := <-tc.nodeUpdateChannels[worker]:
			tc.handleNodeUpdate(nodeUpdate)
			tc.nodeUpdateQueue.Done(nodeUpdate)
		case podUpdate := <-tc.podUpdateChannels[worker]:
			// 如果我们发现了一个 Pod 需要更新，我么你需要先清空 Node 队列.
		priority:
			for {
				select {
				case nodeUpdate := <-tc.nodeUpdateChannels[worker]:
					tc.handleNodeUpdate(nodeUpdate)
					tc.nodeUpdateQueue.Done(nodeUpdate)
				default:
					break priority
				}
			}
			// 在 Node 队列清空后我们再处理 podUpdate.
			tc.handlePodUpdate(podUpdate)
			tc.podUpdateQueue.Done(podUpdate)
		}
	}
}
```