load('ext://restart_process', 'docker_build_with_restart')

local_resource(
  'reception-compile',
  'cd reception && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/app .',
  deps=['./reception/main.go', './reception/app.go', './reception/handlers.go', './reception/go.mod', './reception/go.sum']
)

docker_build_with_restart(
  'bakery/baker',
  'baker/build',
  dockerfile='tilt/Dockerfile',
  entrypoint=['./app'],
  only=['./app'],
  live_update=[
    sync('./baker/build/app', '/bakery/app')
  ]
)

local_resource(
  'baker-compile',
  'cd baker && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/app .',
  deps=['./baker/main.go', './baker/app.go', './baker/bake.go', './baker/go.mod', './baker/go.sum']
)

docker_build_with_restart(
  'bakery/reception',
  'reception/build',
  dockerfile='tilt/Dockerfile',
  entrypoint=['./app'],
  only=['./app'],
  live_update=[
    sync('./reception/build/app', '/bakery/app')
  ]
)

k8s_yaml('tilt/deploy.yaml')