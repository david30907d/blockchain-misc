<!-- minikube image load davidtnfsh/helm_express:0.1
docker build -t davidtnfsh/helm_express:0.1 . -->

➜ helm_express eval $(minikube docker-env)  
➜ helm_express docker build -t davidtnfsh/helm_express:0.1 .
Sending build context to Docker daemon 2.439MB
