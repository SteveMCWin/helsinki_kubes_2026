# 3.9: DBaaS vs DIY

## DBaaS pros
- Less database management overhead
- High availability, automated backups/point-in-time recovery, monitoring, and patching come out of the box
- Storage and compute scale without manual intervention
- Pay for what you provision, on demand

## DBaaS cons
- Higher cost per unit of compute/storage than the equivalent self-managed setup
- Networking gets more complicated (needs Cloud SQL Auth Proxy sidecar or private IP/VPC peering)
- Slight vendor lock-in (mainly around backup/restore tooling and auth, not the data itself)
- More infrastructure components living outside the cluster (separate console/IAM model)

## DIY pros
- More control over how the DB is set up and configured
- Kubernetes-native approach (same manifests, same GitOps flow, same cluster as everything else)
- Lower raw compute/storage cost
- No extra networking setup (plain in-cluster Service)
- You get a new friend (someone to manage the database)

## DIY cons
- You own backups, HA/failover, monitoring, and security patching yourself
- Storage scaling is manual (StorageClass/PVC resizing)
- Your new friend won't have time to play because they are managing backups, security patches and other stuff
