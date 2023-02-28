Kluster

Group: vitu.dev
Version: v1alpha1

- generate deepcopy objects, clientset, informer, listers

- doc.go is used to specify the global tags for the api
    - tags are used to control the behavior of the code-generator
        - global tags are specified in doc.go
        - local tags
        
- code-generator gen command
    1. export GOPATH=~/go
    2. ~/go/src/k8s.io/code-generator //location
    3. cd ~/go/src/github.com/vitu1234/kluster
    4. execDir=~/go/src/k8s.io/code-generator
    5. "${execDir}"/generate-groups.sh all github.com/vitu1234/kluster/pkg/client github.com/vitu1234/kluster/pkg/apis vitu.dev:v1alpha1 --go-header-file "${execDir}"/hack/boilerplate.go.txt

        WHERE; project path: github.com/vitu1234/kluster

- controller-gen
    generate CRD for the type
    1. export GOPATH=~/go
    2. go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.6.0
    3. export PATH=$PATH:$GOPATH/bin
    4. source ~/.bashrc
    5. controller-gen --version
