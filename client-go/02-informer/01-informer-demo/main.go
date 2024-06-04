package main

import (
	"flag"
	"fmt"
	"k8s.io/apimachinery/pkg/labels"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
	"time"
)

func main() {

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	kubeClient, err := kubernetes.NewForConfig(config)
	// 创建一个informer factory
	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, time.Second*30)
	// factory已经为所有k8s的内置资源对象提供了创建对应informer实例的方法，调用具体informer实例的Lister或Informer方法
	// 就完成了将informer注册到factory的过程
	//deploymentLister := kubeInformerFactory.Apps().V1().Deployments().Lister()

	// 这里使用 Informer 也行, 本质上 Informer 会调用  f.factory.InformerFor(&corev1.Pod{}, f.defaultInformer) 来初始化

	// deploymentLister := kubeInformerFactory.Apps().V1().Deployments().Informer()
	podLister := kubeInformerFactory.Core().V1().Pods().Lister()
	// 启动注册到factory的所有informer

	pl := podLister.Pods("default")
	i := kubeInformerFactory.Core().V1().Pods().Informer()
	e := &EventNil{}

	// 开始
	stopCh := make(chan struct{})
	kubeInformerFactory.Start(stopCh)
	kubeInformerFactory.WaitForCacheSync(stopCh)
	i.AddEventHandler(e)

	pods, err := pl.List(labels.Everything())
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(len(pods))

	select {}
}

type EventNil struct {
}

func (receiver *EventNil) OnAdd(obj interface{}, isInInitialList bool) {
	fmt.Println("add")
}
func (receiver *EventNil) OnUpdate(oldObj, newObj interface{}) {
}

func (receiver *EventNil) OnDelete(obj interface{}) {
	fmt.Println("delete")
}
