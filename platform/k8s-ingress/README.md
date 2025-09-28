# Kubernetes Ingress Setup

This folder contains the manifests required to enable HTTPS ingress on the cluster.

## Components

- **ingress-nginx**  
  Provides the Ingress controller (reverse proxy) that routes external traffic into the cluster.

- **cert-manager**  
  Manages TLS certificates from Let's Encrypt.

- **cluster-issuer**  
  Configures cert-manager to automatically request certificates from Let's Encrypt for Ingress resources.

## Deployment

Apply all manifests in this folder:

kubectl apply -R -f platform/k8s-ingress

This will:
1. Install ingress-nginx (creates a LoadBalancer in AKS).  
2. Install cert-manager (with CRDs).  
3. Configure a ClusterIssuer (letsencrypt-prod) for issuing TLS certs.

## Usage

- Create an Ingress resource in your app namespace.  
- Annotate it with:

  annotations:
    cert-manager.io/cluster-issuer: "letsencrypt-prod"

- Add a tls: section with your hostname and a secret name.  
- cert-manager will request a certificate and ingress-nginx will serve your app over HTTPS.
