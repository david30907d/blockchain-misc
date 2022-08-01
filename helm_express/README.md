<!-- minikube image load davidtnfsh/helm_express:0.1
docker build -t davidtnfsh/helm_express:0.1 . -->

minikube start
➜ helm_express eval $(minikube docker-env)  
➜ helm_express docker build -t davidtnfsh/helm_express:0.1 .
Sending build context to Docker daemon 2.439MB
minikube image load mongo:5.0.9
helm install express helm-express
minikube service express-helm-express --url
helm upgrade express helm-express/