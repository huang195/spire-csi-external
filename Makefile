all:
	podman build -t haih/spiffe-csi-driver:external .

push:
	podman push haih/spiffe-csi-driver:external
