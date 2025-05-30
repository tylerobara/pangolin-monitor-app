APP_NAME=pangolin-monitor
IMAGE_REPO=your-docker-repo/$(APP_NAME)
IMAGE_TAG=latest
CHART_VERSION=0.1.0
CHART_PACKAGE=$(APP_NAME)-$(CHART_VERSION).tgz

.PHONY: all build push helm-package deploy test clean

all: build push helm-package deploy test

build:
	docker build -t $(IMAGE_REPO):$(IMAGE_TAG) .

push:
	docker push $(IMAGE_REPO):$(IMAGE_TAG)

helm-package:
	helm package ./pangolin-monitor --version $(CHART_VERSION)

deploy:
	helm upgrade --install $(APP_NAME) ./pangolin-monitor/$(CHART_PACKAGE)

test:
	helm test $(APP_NAME)

clean:
	rm -f pangolin-monitor/$(CHART_PACKAGE)
