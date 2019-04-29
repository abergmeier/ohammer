
package main

import (
	_ "k8s.io/client-go/kubernetes/scheme"
)

/*
decode := scheme.Codecs.UniversalDeserializer().Decode
obj, _, _ := decode([]byte(filecontent), nil, nil)

pod := obj.(*v1beta1.Pod)

`kubectl apply -f - <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: kaniko
spec:
  containers:
  - name: kaniko
    image: gcr.io/kaniko-project/executor:latest
    args: ["--dockerfile=<path to Dockerfile within the build context>",
            "--context=gs://<GCS bucket>/<path to .tar.gz>",
            "--destination=<gcr.io/$PROJECT/$IMAGE:$TAG>"]
    volumeMounts:
      - name: kaniko-secret
        mountPath: /secret
    env:
      - name: GOOGLE_APPLICATION_CREDENTIALS
        value: /secret/kaniko-secret.json
  restartPolicy: Never
  volumes:
    - name: kaniko-secret
      secret:
        secretName: kaniko-secret
EOF`
*/
