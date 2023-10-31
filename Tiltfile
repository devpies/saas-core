k8s_yaml([
    './manifests/nats.yaml',
    './manifests/configmap.yaml',
    './manifests/secrets.yaml',
    './manifests/db-admin.yaml',
    './manifests/db-project.yaml',
    './manifests/db-user.yaml',
    './manifests/db-billing.yaml',
    './manifests/traefik-crd.yaml',
    './manifests/traefik-depl.yaml',
    './manifests/traefik-svc.yaml',
    './manifests/traefik-headers.yaml',
    './manifests/mic-tenant.yaml',
    './manifests/mic-registration.yaml',
    './manifests/mic-admin.yaml',
    './manifests/mic-project.yaml',
    './manifests/mic-user.yaml',
    './manifests/mic-billing.yaml'
])

docker_build('billing:latest', '.', dockerfile = 'deploy/billing.dockerfile')
docker_build('tenant:latest', '.', dockerfile = 'deploy/tenant.dockerfile')
docker_build('registration:latest', '.' ,dockerfile = 'deploy/registration.dockerfile')
docker_build('admin:latest', '.', dockerfile = 'deploy/admin.dockerfile')
docker_build('project:latest', '.', dockerfile = 'deploy/project.dockerfile')
docker_build('user:latest', '.', dockerfile = 'deploy/user.dockerfile')