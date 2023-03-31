# Set Namespace
load('ext://namespace', 'namespace_create', 'namespace_inject')
namespace_create('go-oauth2-server')

# Build Docker image
#   More info: https://docs.tilt.dev/api.html#api.docker_build
docker_build('k3d-registry.tashima.space:5345/tashima42/go-oauth2-server',
             context='.',
             live_update=[
                sync('./api', '/app'),
             ]
)

# Create database secret
load('ext://secret', 'secret_create_generic')
secret_create_generic('pgpassword', namespace="go-oauth2-server", from_file="PGPASSWORD=./secrets/.pgpassword", secret_type="generic")
secret_create_generic('jwtsecret', namespace="go-oauth2-server", from_file="JWTSECRET=./secrets/.jwtsecret", secret_type="generic")

# Apply Kubernetes manifests
#   More info: https://docs.tilt.dev/api.html#api.k8s_yaml
k8s_yaml([
  'k8s/database-persistent-volume-claim.yaml', 
  'k8s/database-deployment.yaml', 
  'k8s/database-cluster-ip-service.yaml', 
  'k8s/api-deployment.yaml', 
  'k8s/api-service.yaml',
  ])

k8s_resource('database-deployment', port_forwards=5432)
k8s_resource('api-deployment', port_forwards=8096)

load('ext://git_resource', 'git_checkout')
