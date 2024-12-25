name: Deploy MultiK8s
on:
  push:
    branches:
      - main

env:
  SHA: $(git rev-parse HEAD)

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Test
        run: |-
          docker login -u ${{ secrets.DOCKER_USERNAME }} -p ${{ secrets.DOCKER_PASSWORD }}
          docker build -t rallycoding/react-test -f ./client/Dockerfile.dev ./client
          docker run -e CI=true rallycoding/react-test npm test

      - name: Set Service Key
        uses: 'google-github-actions/auth@v0'
        with:
          credentials_json: '${{ secrets.GKE_SA_KEY }}'

      - name: Set Project
        uses: google-github-actions/setup-gcloud@v0
        with:
          project_id: multi-338920

      - name: Auth
        run: |-
          gcloud --quiet auth configure-docker

      - name: Get Credentials
        uses: google-github-actions/get-gke-credentials@v0
        with:
          cluster_name: multi-cluster
          location: us-central1-c

      - name: Build
        run: |-
          docker build -t rallycoding/multi-client-k8s-gh:latest -t rallycoding/multi-client-k8s-gh:${{ env.SHA }} -f ./client/Dockerfile ./client
          docker build -t rallycoding/multi-server-k8s-gh:latest -t rallycoding/multi-server-k8s-gh:${{ env.SHA }} -f ./server/Dockerfile ./server
          docker build -t rallycoding/multi-worker-k8s-gh:latest -t rallycoding/multi-worker-k8s-gh:${{ env.SHA }} -f ./worker/Dockerfile ./worker

      - name: Push
        run: |-
          docker push rallycoding/multi-client-k8s-gh:latest
          docker push rallycoding/multi-server-k8s-gh:latest
          docker push rallycoding/multi-worker-k8s-gh:latest

          docker push rallycoding/multi-client-k8s-gh:${{ env.SHA }}
          docker push rallycoding/multi-server-k8s-gh:${{ env.SHA }}
          docker push rallycoding/multi-worker-k8s-gh:${{ env.SHA }}

      - name: Apply
        run: |-
          kubectl apply -f k8s
          kubectl set image deployments/server-deployment server=rallycoding/multi-server-k8s-gh:${{ env.SHA }}
          kubectl set image deployments/client-deployment client=rallycoding/multi-client-k8s-gh:${{ env.SHA }}
          kubectl set image deployments/worker-deployment worker=rallycoding/multi-worker-k8s-gh:${{ env.SHA }}
