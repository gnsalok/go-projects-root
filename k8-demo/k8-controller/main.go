package main

import (
	"flag"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
)

// MyController defines the structure of our custom controller.
type MyController struct {
	clientset  *kubernetes.Clientset
	cmInformer cache.SharedIndexInformer
	workqueue  workqueue.RateLimitingInterface
}

// NewMyController is a constructor for MyController.
func NewMyController(clientset *kubernetes.Clientset, informer cache.SharedIndexInformer) *MyController {
	// We’ll use a simple rate-limiting workqueue to handle reconciliation events.
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	controller := &MyController{
		clientset:  clientset,
		cmInformer: informer,
		workqueue:  queue,
	}

	// Register event handlers for the informer
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			klog.Info("ConfigMap Added")
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(key)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			klog.Info("ConfigMap Updated")
			key, err := cache.MetaNamespaceKeyFunc(newObj)
			if err == nil {
				queue.Add(key)
			}
		},
		DeleteFunc: func(obj interface{}) {
			klog.Info("ConfigMap Deleted")
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(key)
			}
		},
	})

	return controller
}

// Run starts the controller loop.
func (c *MyController) Run(stopCh <-chan struct{}) {
	// We start processing the informer in a separate goroutine.
	go c.cmInformer.Run(stopCh)

	// Wait until all caches have been synced, meaning
	// the informer has a complete view of the cluster state.
	if !cache.WaitForCacheSync(stopCh, c.cmInformer.HasSynced) {
		runtime.HandleError(fmt.Errorf("failed to sync informer caches"))
		return
	}

	// Start a worker to process items from the queue.
	go c.runWorker()

	// Block until a stop signal is received.
	<-stopCh
	klog.Info("Stopping MyController...")
}

func (c *MyController) runWorker() {
	for c.processNextItem() {
	}
}

// processNextItem pulls an item from the queue and processes it.
func (c *MyController) processNextItem() bool {
	// Pull the item off the queue
	key, quit := c.workqueue.Get()
	if quit {
		return false
	}
	// Always mark the item done once processed
	defer c.workqueue.Done(key)

	// The key is usually in "namespace/name" format.
	// We'll parse that to fetch the actual ConfigMap if needed.
	err := c.reconcile(key.(string))
	if err != nil {
		// Re-queue the item to retry later.
		c.workqueue.AddRateLimited(key)
		return true
	}

	// If no error, tell the queue we successfully processed it.
	c.workqueue.Forget(key)
	return true
}

// reconcile is where you’d implement the logic to ensure your desired state.
func (c *MyController) reconcile(key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	// Get the ConfigMap from the cache (or from the cluster via clientset).
	cmObj, err := c.cmInformer.GetIndexer().ByIndex(cache.NamespaceIndex, namespace)
	if err != nil {
		klog.Errorf("Error fetching ConfigMap from indexer: %v", err)
		return err
	}

	// Alternatively, to be thorough, you can do:
	// cm, err := c.clientset.CoreV1().ConfigMaps(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	// if err != nil {
	//     // If it’s a “NotFound” error, the ConfigMap was deleted.
	//     // Otherwise, handle other errors.
	// }

	// For demonstration, we’ll just log that we’re reconciling.
	klog.Infof("Reconciling ConfigMap: %s/%s", namespace, name)

	// Here you would put your business logic, e.g.:
	// - Compare the ConfigMap with a desired state
	// - Possibly create, update, or delete other resources
	// - Perform some external action, etc.

	return nil
}

func main() {
	// Accept a kubeconfig flag so we can run this controller outside the cluster as well.
	var kubeconfig string
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig file")
	flag.Parse()

	// Build the config from flags or assume in-cluster.
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		klog.Fatalf("Error building kubeconfig: %v", err)
	}

	// Create a clientset.
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatalf("Error creating clientset: %v", err)
	}

	// Create a shared informer factory to watch ConfigMaps across all namespaces.
	// You can restrict this to a specific namespace if desired.
	informerFactory := informers.NewSharedInformerFactory(clientset, 30*time.Second)

	// Create an informer for ConfigMaps.
	cmInformer := informerFactory.Core().V1().ConfigMaps().Informer()

	// Create our controller.
	controller := NewMyController(clientset, cmInformer)

	// Create a channel to handle stop signals (SIGINT, SIGTERM, etc.).
	stopCh := make(chan struct{})
	defer close(stopCh)

	// Start the controller.
	go controller.Run(stopCh)

	// Keep the main thread alive until a stop signal is received.
	// In a real program, you might listen for OS signals and close stopCh on SIGTERM.
	select {}
}
