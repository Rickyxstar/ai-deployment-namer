package main

import (
	"net/http"
	"os"

	"github.com/rickyxstar/ai-deployment-namer/internal/repositories"
	"github.com/rickyxstar/ai-deployment-namer/internal/webhook"
	appsv1 "k8s.io/api/apps/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	setupLog = ctrl.Log.WithName("setup")
)

func main() {
	ctrl.SetLogger(zap.New())

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	var ng repositories.NameGenerator

	switch os.Getenv("NAME_GENERATOR") {
	case "ollama":
		ng = repositories.NewNameGeneratorOllama(
			os.Getenv("OLLAMA_HOST"),
			&http.Client{},
			os.Getenv("MODEL"),
		)
	default:
		ng = repositories.NewNameGeneratorChatGPT(
			os.Getenv("OPENAI_API_KEY"),
			&http.Client{},
			os.Getenv("MODEL"),
		)
	}

	dn := webhook.NewDeploymentNamer(ng)

	err = ctrl.NewWebhookManagedBy(mgr).
		For(&appsv1.Deployment{}).
		WithDefaulter(dn).
		Complete()
	if err != nil {
		setupLog.Error(err, "unable to create webhook")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
