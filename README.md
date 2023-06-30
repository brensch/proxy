# proxy

An easy way to get a bunch of free IP addresses.

from https://cloud.google.com/functions/docs/locations

## Tier 1 pricing:

- us-west1 (Oregon)
- us-central1 (Iowa)
- us-east1 (South Carolina)
- us-east4 (Northern Virginia)
- europe-west1 (Belgium)
- europe-west2 (London)
- asia-east1 (Taiwan)
- asia-east2 (Hong Kong)
- asia-northeast1 (Tokyo)
- asia-northeast2 (Osaka)

## Tier 2 pricing:

- us-west2 (Los Angeles)
- us-west3 (Salt Lake City)
- us-west4 (Las Vegas)
- northamerica-northeast1 (Montreal)
- southamerica-east1 (Sao Paulo)
- europe-west3 (Frankfurt)
- europe-west6 (Zurich)
- europe-central2 (Warsaw)
- australia-southeast1 (Sydney)
- asia-south1 (Mumbai)
- asia-southeast1 (Singapore)
- asia-southeast2 (Jakarta)
- asia-northeast3 (Seoul)


```
gcloud projects describe proxy-362608 --format='value(projectNumber)'

gcloud projects add-iam-policy-binding proxy-362608 \
  --member=serviceAccount:668769694702@cloudbuild.gserviceaccount.com \
  --role=roles/run.admin

gcloud projects add-iam-policy-binding proxy-362608 \
  --member=serviceAccount:668769694702@cloudbuild.gserviceaccount.com \
  --role=roles/iam.serviceAccountUser

gcloud builds submit --config cloudbuild.yaml --substitutions=_IMAGE_NAME=proxy,_PROJECT_ID=proxy-362608
```
