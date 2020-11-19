# syntax=docker/dockerfile:1.0-experimental
ARG GOLANG_CONTAINER=golang:latest
ARG NGINX_VERSION
ARG BUILD_OS

######### Base image for Debian #########
FROM nginx:$NGINX_VERSION AS debian
ARG IC_VERSION

LABEL version=$NGINX_VERSION-debian-$IC_VERSION
LABEL maintainer="NGINX Docker Maintainers <docker-maint@nginx.com>"

RUN  apt-get update \
	&& apt-get install -y libcap2-bin


######### Base image for Alpine #########
FROM nginx:$NGINX_VERSION-alpine AS alpine
ARG IC_VERSION

LABEL version=$NGINX_VERSION-alpine-$IC_VERSION
LABEL maintainer="NGINX Docker Maintainers <docker-maint@nginx.com>"

RUN apk add --no-cache libcap


######### Base imager for Debian with NGINX Plus #########
FROM debian:buster-slim AS debian-plus
ARG IC_VERSION
ENV NGINX_PLUS_VERSION 22-1~buster

LABEL version=debian-nginx-plus-$NGINX_PLUS_VERSION-$IC_VERSION
LABEL maintainer="NGINX Docker Maintainers <docker-maint@nginx.com>"

# Install NGINX Plus
RUN --mount=type=secret,id=nginx-repo.crt,dst=/etc/ssl/nginx/nginx-repo.crt,mode=0644 \
	--mount=type=secret,id=nginx-repo.key,dst=/etc/ssl/nginx/nginx-repo.key,mode=0644 \
	set -x \
	&& apt-get update \
	&& apt-get install --no-install-recommends --no-install-suggests -y apt-transport-https ca-certificates gnupg1 libcap2-bin \
	&& \
	NGINX_GPGKEY=573BFD6B3D8FBC641079A6ABABF5BD827BD9BF62; \
	found=''; \
	for server in \
		ha.pool.sks-keyservers.net \
		hkp://keyserver.ubuntu.com:80 \
		hkp://p80.pool.sks-keyservers.net:80 \
		pgp.mit.edu \
	; do \
		echo "Fetching GPG key $NGINX_GPGKEY from $server"; \
		apt-key adv --keyserver "$server" --keyserver-options timeout=10 --recv-keys "$NGINX_GPGKEY" && found=yes && break; \
	done; \
	test -z "$found" && echo >&2 "error: failed to fetch GPG key $NGINX_GPGKEY" && exit 1; \
	echo "Acquire::https::plus-pkgs.nginx.com::Verify-Peer \"true\";" >> /etc/apt/apt.conf.d/90nginx \
	&& echo "Acquire::https::plus-pkgs.nginx.com::Verify-Host \"true\";" >> /etc/apt/apt.conf.d/90nginx \
	&& echo "Acquire::https::plus-pkgs.nginx.com::SslCert     \"/etc/ssl/nginx/nginx-repo.crt\";" >> /etc/apt/apt.conf.d/90nginx \
	&& echo "Acquire::https::plus-pkgs.nginx.com::SslKey      \"/etc/ssl/nginx/nginx-repo.key\";" >> /etc/apt/apt.conf.d/90nginx \
	&& echo "Acquire::https::plus-pkgs.nginx.com::User-Agent  \"k8s-ic-$IC_VERSION-apt\";" >> /etc/apt/apt.conf.d/90nginx \
	&& printf "deb https://plus-pkgs.nginx.com/debian buster nginx-plus\n" > /etc/apt/sources.list.d/nginx-plus.list \
	&& apt-get update && apt-get install -y nginx-plus=${NGINX_PLUS_VERSION}


######### Base image for Debian with NGINX Plus and App Protect #########
FROM debian-plus as debian-plus-ap
ARG IC_VERSION
LABEL version=debian-nginx-plus-ap-$IC_VERSION
LABEL maintainer="NGINX Docker Maintainers <docker-maint@nginx.com>"

ENV APPPROTECT_MODULE_VERSION 22+3.158.1-1~stretch
ENV APPPROTECT_PLUGIN_VERSION 3.158.1-1~stretch
ENV APPPROTECT_ENGINE_VERSION 4.23.3-1~stretch
ENV APPPROTECT_COMPILER_VERSION 4.23.3-1~stretch
ENV APPPROTECT_SIG_VERSION 2020.09.03-1~stretch
ENV APPPROTECT_THREAT_CAMPAIGNS_VERSION 2020.09.01-1~stretch
ENV NGINX_PLUS_RELEASE R22

RUN apt-get install --no-install-recommends --no-install-suggests -y wget

RUN --mount=type=secret,id=nginx-repo.crt,dst=/etc/ssl/nginx/nginx-repo.crt,mode=0644 \
	--mount=type=secret,id=nginx-repo.key,dst=/etc/ssl/nginx/nginx-repo.key,mode=0644 \
	&& sed -i '$ d' /etc/apt/apt.conf.d/90nginx \
	&& echo "Acquire::https::plus-pkgs.nginx.com::User-Agent  \"k8s-ic-$IC_VERSION-app-$APPPROTECT_VERSION-apt\";" >> /etc/apt/apt.conf.d/90nginx \
	&& echo "deb https://app-protect-sigs.nginx.com/debian/ stretch nginx-plus\n" | tee /etc/apt/sources.list.d/app-protect-sigs.list \
	&& wget https://nginx.org/keys/app-protect-sigs.key && apt-key add app-protect-sigs.key \
	&& echo "Acquire::https::app-protect-sigs.nginx.com::Verify-Peer \"true\";" >> /etc/apt/apt.conf.d/90app-protect-sigs \
	&& echo "Acquire::https::app-protect-sigs.nginx.com::Verify-Host \"true\";" >> /etc/apt/apt.conf.d/90app-protect-sigs \
	&& echo "Acquire::https::app-protect-sigs.nginx.com::SslCert     \"/etc/ssl/nginx/nginx-repo.crt\";" >> /etc/apt/apt.conf.d/90app-protect-sigs \
	&& echo "Acquire::https::app-protect-sigs.nginx.com::SslKey      \"/etc/ssl/nginx/nginx-repo.key\";" >> /etc/apt/apt.conf.d/90app-protect-sigs \
	&& apt-get update && \
	    apt-get install -y nginx-plus-module-appprotect=$APPPROTECT_MODULE_VERSION \
	    app-protect-plugin=$APPPROTECT_PLUGIN_VERSION \
	    app-protect-engine=$APPPROTECT_ENGINE_VERSION \
	    app-protect-compiler=$APPPROTECT_COMPILER_VERSION \
	    app-protect=$APPPROTECT_MODULE_VERSION \
	&& apt-get install -y app-protect-attack-signatures${APPPROTECT_SIG_VERSION:+=$APPPROTECT_SIG_VERSION} \
	&& apt-get install -y app-protect-threat-campaigns${APPPROTECT_THREAT_CAMPAIGNS_VERSION:+=$APPPROTECT_THREAT_CAMPAIGNS_VERSION}

# Uncomment the lines below if you want to install a custom CA certificate
# COPY build/appprotect/*.crt  /usr/local/share/ca-certificates/
# RUN update-ca-certificates


######### Base image for Openshift #########
FROM registry.access.redhat.com/ubi8/ubi:8.3 AS openshift-base

ARG NGINX_VERSION
LABEL name="NGINX Ingress Controller" \
      description="The Ingress controller is an application that runs in a cluster and configures an HTTP load balancer according to Ingress resources." \
      summary="The Ingress controller is an application that runs in a cluster and configures an HTTP load balancer according to Ingress resources." \
      io.openshift.tags="nginx,ingress-controller,ingress,controller,kubernetes,openshift" \
      maintainer="NGINX Docker Maintainers <docker-maint@nginx.com>" \
      vendor="NGINX Inc <kubernetes@nginx.com>"

RUN set -x \
	&& groupadd --system --gid 101 nginx \
	&& useradd --system --gid nginx --no-create-home --home-dir /nonexistent --comment "nginx user" --shell /bin/false --uid 101 nginx


RUN mkdir licenses
COPY LICENSE /licenses


######### Base image for Openshift OSS #########
FROM openshift-base AS openshift

RUN echo "[nginx]" >> /etc/yum.repos.d/nginx.repo \
	&& echo "name=nginx repo" >> /etc/yum.repos.d/nginx.repo \
	&& echo "baseurl=https://nginx.org/packages/mainline/rhel/8/\$basearch/" >> /etc/yum.repos.d/nginx.repo \
	&& echo "gpgcheck=0" >> /etc/yum.repos.d/nginx.repo \
	&& echo "enabled=1" >> /etc/yum.repos.d/nginx.repo \
	&& echo "module_hotfixes=true" >> /etc/yum.repos.d/nginx.repo \
	&& yum install -y nginx-${NGINX_VERSION}


######### Base image for Openshift with NGINX Plus #########
FROM openshift-base AS openshift-plus

ENV NGINX_PLUS_VERSION 22-1.el8.ngx
ARG IC_VERSION

# Install NGINX Plus
RUN --mount=type=secret,id=nginx-repo.crt,dst=/etc/ssl/nginx/nginx-repo.crt,mode=0644 \
	--mount=type=secret,id=nginx-repo.key,dst=/etc/ssl/nginx/nginx-repo.key,mode=0644 \
	&& yum install -y wget ca-certificates \
	&& wget https://nginx.org/keys/nginx_signing.key \
	&& \
	NGINX_GPGKEY=nginx_signing.key; \
	found=''; \
	for server in \
		ha.pool.sks-keyservers.net \
		hkp://p80.pool.sks-keyservers.net:80 \
		pgp.mit.edu \
	; do \
		echo "Fetching GPG key $NGINX_GPGKEY from $server"; \
		rpm --define="%_hkp_keyserver $server" --import $NGINX_GPGKEY && found=yes && break; \
	done; \
	test -z "$found" && echo >&2 "error: failed to fetch GPG key $NGINX_GPGKEY" && exit 1; \
	echo "[nginx-plus]" >> /etc/yum.repos.d/nginx-plus-8.repo \
	&& echo "name=nginx-plus repo" >> /etc/yum.repos.d/nginx-plus-8.repo \
	&& echo "baseurl=https://plus-pkgs.nginx.com/centos/8/\$basearch/" >> /etc/yum.repos.d/nginx-plus-8.repo \
	&& echo "sslclientcert=/etc/ssl/nginx/nginx-repo.crt" >> /etc/yum.repos.d/nginx-plus-8.repo \
	&& echo "sslclientkey=/etc/ssl/nginx/nginx-repo.key" >> /etc/yum.repos.d/nginx-plus-8.repo \
	&& echo "gpgcheck=1" >> /etc/yum.repos.d/nginx-plus-8.repo \
	&& echo "enabled=1" >> /etc/yum.repos.d/nginx-plus-8.repo \
	&& yum install -y nginx-plus-${NGINX_PLUS_VERSION}


######### Create common files, permissions and setcap #########
FROM $BUILD_OS as files

ARG PLUS

RUN mkdir -p /var/lib/nginx \
	&& mkdir -p /etc/nginx/secrets \
	&& mkdir -p /etc/nginx/stream-conf.d \
	&& setcap 'cap_net_bind_service=+ep' /usr/sbin/nginx \
	&& setcap 'cap_net_bind_service=+ep' /usr/sbin/nginx-debug \
	&& chown -R nginx:0 /etc/nginx \
	&& chown -R nginx:0 /var/cache/nginx \
	&& chown -R nginx:0 /var/lib/nginx \
	&& rm /etc/nginx/conf.d/*

COPY internal/configs/version1/nginx$PLUS.ingress.tmpl \
	internal/configs/version1/nginx$PLUS.tmpl \
	internal/configs/version2/nginx$PLUS.virtualserver.tmpl \
	internal/configs/version2/nginx$PLUS.transportserver.tmpl /

# Uncomment the line below if you would like to add the default.pem to the image
# and use it as a certificate and key for the default server
# ADD default.pem /etc/nginx/secrets/default

EXPOSE 80 443

ENTRYPOINT ["/nginx-ingress"]


######### Cleanup for Debian #########
FROM files as cleanup-debian
RUN apt-get remove --purge -y libcap2-bin \
	&& rm -rf /var/lib/apt/lists/*


######### Cleanup for Alpine #########
FROM files as cleanup-alpine
RUN apk del libcap \
	&& rm -rf /var/cache/apk/*


######### Cleanup for Debian wiht NGINX Plus #########
FROM files as cleanup-debian-plus
RUN apt-get remove --purge --auto-remove -y gnupg1 libcap2-bin \
	&& rm -rf /var/lib/apt/lists/* \
	&& rm -rf /etc/ssl/nginx \
	&& rm /etc/apt/apt.conf.d/90nginx /etc/apt/sources.list.d/nginx-plus.list


######### Cleanup for Debian wiht NGINX Plus and App Protect #########
FROM cleanup-debian-plus as cleanup-debian-plus-ap
RUN apt-get remove --purge --auto-remove -y wget \
	&& rm /etc/apt/apt.conf.d/90app-protect-sigs /etc/apt/sources.list.d/app-protect-sigs.list \
	&& mkdir -p /etc/nginx/waf \
	&& mkdir -p /etc/nginx/waf/nac-policies \
	&& mkdir -p /etc/nginx/waf/nac-logconfs \
	&& mkdir -p /var/log/app_protect \
	&& mkdir -p /opt/app_protect \
	&& chown -R nginx:0 /etc/app_protect \
	&& chown -R nginx:0 /usr/share/ts \
	&& chown -R nginx:0 /var/log/app_protect/ \
	&& chown -R nginx:0 /opt/app_protect/ \
	&& chown -R nginx:0 /var/log/nginx/

RUN printf "MODULE = ALL;\nLOG_LEVEL = TS_CRIT;\nFILE = 2;\n" > /etc/app_protect/bd/logger.cfg \
	&&  printf "[config_set_compiler]\nlog_level=fatal\n" >> /etc/app_protect/tools/asm_logging.conf \
	&& for v in \
		asm_config_server \
		lock_factory \
		bd_agent \
		import_export_policy \
		set_active \
	; do sed -i "/\[$v/a log_level=fatal" "/etc/app_protect/tools/asm_logging.conf" \
	; done

COPY --chown=nginx:0 build/log-default.json /etc/nginx


######### Cleanup for Openshift #########
FROM files as cleanup-openshift
RUN rm /etc/yum.repos.d/nginx.repo


######### Cleanup for Openshift with NGINX Plus #########
FROM files as cleanup-openshift-plus
RUN yum remove -y wget \
	&& rm -rf /etc/ssl/nginx \
	&& rm /etc/yum.repos.d/nginx-plus-8.repo \
	&& rm nginx_signing.key


######### Build nginx-ingress in golang container #########
FROM $GOLANG_CONTAINER AS builder
ARG VERSION
ARG GIT_COMMIT
WORKDIR /go/src/github.com/nginxinc/kubernetes-ingress/nginx-ingress/cmd/nginx-ingress
COPY . /go/src/github.com/nginxinc/kubernetes-ingress/nginx-ingress/
RUN CGO_ENABLED=0 GOFLAGS='-mod=vendor' \
	go build -installsuffix cgo -ldflags "-w -X main.version=${VERSION} -X main.gitCommit=${GIT_COMMIT}" -o /nginx-ingress


######### Create image with nginx-ingress built in container
FROM cleanup-$BUILD_OS AS container
COPY --from=builder /nginx-ingress /

USER nginx

######### Create image with nginx-ingress built locally #########
FROM cleanup-$BUILD_OS AS local
COPY nginx-ingress /

USER nginx
