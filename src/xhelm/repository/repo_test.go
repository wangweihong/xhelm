package repository

import (
	"fmt"
	"testing"
)

const testTemplate = `
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: {{ template "fullname" . }}
  labels:
    app: {{ template "name" . }}
    chart: {{ .Chart.Name }}-{{ .Chart.Version | replace "+" "_" }}
    release: {{ .Release.Name }}
    heritage: {{ .Release.Service }}
spec:
  replicas: {{ .Values.replicaCount }}
  template:
    metadata:
      labels:
        app: {{ template "name" . }}
        release: {{ .Release.Name }}
    spec:
      containers:
        - name: {{ .Chart.Name }}
          image: "{{ .Values.image.repository }}:{{ .Values.image.tag }}"
          imagePullPolicy: {{ .Values.image.pullPolicy }}
          ports:
            - containerPort: {{ .Values.service.internalPort }}
          livenessProbe:
            httpGet:
              path: /
              port: {{ .Values.service.internalPort }}
          readinessProbe:
            httpGet:
              path: /
              port: {{ .Values.service.internalPort }}
          resources:
{{ toYaml .Values.resources | indent 12 }}
    {{- if .Values.nodeSelector }}
      nodeSelector:
{{ toYaml .Values.nodeSelector | indent 8 }}
    {{- end }}

`
const testValue = `
# Default values for haha.
# This is a YAML-formatted file.
# Declare variables to be passed into your templates.
replicaCount: 1
image:
  repository: nginx
  tag: stable
  pullPolicy: IfNotPresent
service:
  name: nginx
  type: ClusterIP
  internalPort: 80
ingress:
  enabled: false
  # Used to create an Ingress record.
  hosts:
    - chart-example.local
  annotations:
    # kubernetes.io/ingress.class: nginx
    # kubernetes.io/tls-acme: "true"
  tls:
    # Secrets must be manually created in the namespace.
    # - secretName: chart-example-tls
    #   hosts:
    #     - chart-example.local
resources: {}
  # We usually recommend not to specify default resources and to leave this as a conscious 
  # choice for the user. This also increases chances charts run on environments with little 
  # resources, such as Minikube. If you do want to specify resources, uncomment the following 
  # lines, adjust them as necessary, and remove the curly braces after 'resources:'.
  # limits:
  #  cpu: 100m
  #  memory: 128Mi
  #requests:
  #  cpu: 100m
  #  memory: 128Mi

`

/*
func Test_Repo(t *testing.T) {
	t.Log("创建仓库test")
	err := RM.CreateRepo("test")
	if err != nil {
		t.Error("创建仓库失败:", err)
		t.Fail()
	}

	repos := RM.ListRepos()
	for _, v := range repos {
		fmt.Printf("%v %v %v\n", v.Name, v.State, v.CreateTime)
	}

	err = RM.DeleteRepo("test")
	if err != nil {
		t.Error("删除仓库失败:", err)
		t.Fail()
	}
}
*/

func Test_AddRemote(t *testing.T) {
	/*
		opt := CreateOption{
			URL: "http://127.0.0.1:8879",
		}
	*/

	err := RM.AddRepo("localdddd", nil)
	if err != nil {
		t.Error("添加本地仓库失败:", err)
		t.Fail()
	}

	t.Log("添加仓库成功")
}

func Test_ListCharts(t *testing.T) {
	cs, err := RM.ListCharts("localdddd")
	if err != nil {
		t.Error("查看chart失败")
		t.Fail()
	}

	for _, v := range cs {
		//	fmt.Println(v)
		fmt.Println(v.Name)
	}

}

func Test_CreateChart(t *testing.T) {
	var opt ChartCreateOption
	opt.Version = "0.0.1"
	opt.Template = []byte(testTemplate)
	opt.DefaultValues = []byte(testValue)
	err := RM.CreateChart("localdddd", "haha", opt)
	if err != nil {
		t.Error("创建chart失败", err)
		t.Fail()
	}

	t.Log("创建chart成功")
}

func Test_GetChart(t *testing.T) {
	c, err := RM.GetChartVersion("localdddd", "haha", "0.0.1")
	if err != nil {
		t.Error("查看chart失败", err)
		t.Fail()
	}
	err = RM.UncompressData("localdddd", c)
	if err != nil {
		t.Error("解压chart失败", err)
		t.Fail()
	}

}

/*
func Test_DeleteRepo(t *testing.T) {
	err := RM.DeleteRepo("localdddd")
	if err != nil {
		t.Error("删除repo失败", err)
		t.Fail()

	}

}
*/
