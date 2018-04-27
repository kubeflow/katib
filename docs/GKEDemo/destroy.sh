#/bin/bash
set -x
gcloud compute firewall-rules delete katibservice
gcloud container clusters delete katib
gcloud compute disks delete katib-pv
