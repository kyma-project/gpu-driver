# syntax=docker/dockerfile:1
ARG REGISTRY_PATH=gardenlinux/kmodbuild

FROM debian:bookworm-slim AS packager
ARG TARGET_ARCH
ARG DRIVER_VERSION

COPY resources/scripts/* /opt/nvidia-installer/

RUN apt-get update && apt-get install --no-install-recommends -y \
    kmod \
    pciutils \
    ca-certificates \
    wget \
    xz-utils

 RUN rm -rf /var/lib/apt/lists/*

RUN /opt/nvidia-installer/download_fabricmanager.sh

# Remove several things that are not needed, some of which raise Black Duck scan vulnerabilities
RUN apt-get remove -y --autoremove --allow-remove-essential --ignore-hold \
      libgnutls30 apt openssl wget ncurses-base ncurses-bin

RUN rm -rf /var/lib/apt/lists/* /usr/bin/dpkg /sbin/start-stop-daemon /usr/lib/x86_64-linux-gnu/libsystemd.so* \
         /var/lib/dpkg/info/libdb5.3* /usr/lib/x86_64-linux-gnu/libdb-5.3.so* /usr/share/doc/libdb5.3 \
         /usr/bin/chfn /usr/bin/gpasswd

RUN mkdir -p /rootfs \
        && cp -ar /bin /boot /etc /home /lib /lib64 /media /mnt /opt /root /run /sbin /srv /tmp /usr /var /rootfs \
        && rm -rf /rootfs/opt/actions-runner

FROM scratch

COPY --from=packager /rootfs /

ENTRYPOINT ["/opt/nvidia-installer/entrypoint.sh"]
