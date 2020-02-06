module github.com/lyft/flytepropeller

go 1.13

require (
	github.com/DiSiqueira/GoTree v1.0.1-0.20180907134536-53a8e837f295
	github.com/aws/aws-sdk-go v1.28.12 // indirect
	github.com/benlaurie/objecthash v0.0.0-20180202135721-d1e3d6079fc1
	github.com/coreos/go-oidc v2.2.1+incompatible // indirect
	github.com/fatih/color v1.9.0
	github.com/ghodss/yaml v1.0.0
	github.com/go-redis/redis v6.15.7+incompatible
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.3.3
	github.com/google/uuid v1.1.1
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.0
	github.com/imdario/mergo v0.3.8 // indirect
	github.com/lyft/datacatalog v0.1.2
	github.com/lyft/flyteidl v0.17.2
	github.com/lyft/flyteplugins v0.3.2
	github.com/lyft/flytestdlib v0.3.2
	github.com/magiconair/properties v1.8.1
	github.com/mitchellh/mapstructure v1.1.2
	github.com/ncw/swift v1.0.50 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.4.0
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.4.0
	golang.org/x/crypto v0.0.0-20200206161412-a0c6ece9d31a // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0
	google.golang.org/api v0.17.0 // indirect
	google.golang.org/grpc v1.27.1
	gopkg.in/ini.v1 v1.52.0 // indirect
	k8s.io/api v0.17.2
	k8s.io/apimachinery v0.17.2
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/klog v1.0.0
	k8s.io/kube-openapi v0.0.0-20200204173128-addea2498afe // indirect
	sigs.k8s.io/controller-runtime v0.4.0
)

// Pin the version of client-go to something that's compatible with katrogan's fork of api and apimachinery
// Type the following
//   replace k8s.io/client-go => k8s.io/client-go kubernetes-1.16.2
// and it will be replaced with the 'sha' variant of the version

replace (
	github.com/GoogleCloudPlatform/spark-on-k8s-operator => github.com/lyft/spark-on-k8s-operator v0.1.3
	github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.3.1
	gopkg.in/fsnotify.v1 => github.com/fsnotify/fsnotify v1.4.7
	k8s.io/api => github.com/lyft/api v0.0.0-20191031200350-b49a72c274e0
	k8s.io/apimachinery => github.com/lyft/apimachinery v0.0.0-20191031200210-047e3ea32d7f
	k8s.io/client-go => k8s.io/client-go v0.0.0-20191016111102-bec269661e48
)
