load('ext://helm_resource', 'helm_resource', 'helm_repo')
# Set Namespace
load('ext://namespace', 'namespace_create', 'namespace_inject')
namespace_create('go-oauth2-server-local')

# Local resource to build the API binary
local_resource(
  'api-compile',
  'CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o build/go-oauth2-server ./',
  dir='./api',
  deps=['./api/main.go', './api/go.mod', './api/go.sum', './api/api.go', './api/db', './api/handlers', './api/helpers' ],
  labels=["api"]
)

# Local resource to build the UI dist
local_resource(
  'ui-compile',
  'yarn run build',
  dir='./ui',
  deps=['./ui/package.json', './ui/package-lock.json', './ui/src'],
  labels=["ui"]
)

# Build API Docker image
#   More info: https://docs.tilt.dev/api.html#api.docker_build
load('ext://restart_process', 'docker_build_with_restart')
docker_build_with_restart(
  'k3d-registry.local.tashima.space/tashima42/go-oauth2-server/api',
  context='.',
  dockerfile='api/Dockerfile.dev',
  entrypoint="/app/build/go-oauth2-server",
  only=['./api/build', "./api/db/schema_migrations"],
  live_update=[
    sync('./api/build', '/app/build'),
    sync('./api/db/schema_migrations', '/app/db/schema_migrations'),
  ],
)

# Build UI Docker image
docker_build(
  'k3d-registry.local.tashima.space/tashima42/go-oauth2-server/ui',
  context='.',
  dockerfile='ui/Dockerfile.dev',
  only=['./ui/dist', './ui/nginx.conf', './ui/hack'],
  live_update=[
    sync('./ui/dist', '/app'),
    sync('./ui/nginx.conf', '/etc/nginx/nginx.conf'),
  ],
)

# Create database secret
load('ext://secret', 'secret_create_generic')
secret_create_generic('pgpassword', namespace="go-oauth2-server-local", from_file="PGPASSWORD=./.secrets/.pgpassword", secret_type="generic")
secret_create_generic('jwtsecret', namespace="go-oauth2-server-local", from_file="JWTSECRET=./.secrets/.jwtsecret", secret_type="generic")

# Apply Kubernetes manifests
#   More info: https://docs.tilt.dev/api.html#api.k8s_yaml
k8s_yaml([
  'k8s/database-persistent-volume-claim.yaml', 
  'k8s/database-deployment.yaml', 
  'k8s/database-cluster-ip-service.yaml', 
])
k8s_yaml([
  'k8s/api-deployment.yaml', 
  'k8s/api-service.yaml',
  'k8s/api-ingressroute.yaml',
])
k8s_yaml([
  'k8s/ui-deployment.yaml',
  'k8s/ui-service.yaml',
  'k8s/ui-ingressroute.yaml',
])

k8s_resource('database-deployment', port_forwards='5432:5432', labels=["db"])
k8s_resource('api-deployment',  labels=["api"])
k8s_resource('ui-deployment', labels=["ui"])
