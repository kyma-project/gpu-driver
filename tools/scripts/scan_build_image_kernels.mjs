#!/usr/bin/env zx

const tagMap = {}

let page = 1
while (page < 20) {
    const txt = await $`gh api "/orgs/gardenlinux/packages/container/gardenlinux%2Fkmodbuild/versions?per_page=100&page=${page}"`
    const data = JSON.parse(txt.text())
    for (let img of data) {
        for (let tag of (img.metadata.container.tags || [])) {
            if (typeof tag !== "string") {
                continue
            }
            if (tag.startsWith("amd64-")) {
                tagMap[tag] = true
            }
            if (tag === "amd64-1592.3") {
                break
            }
        }
    }
    page += 1
}
let arr = []
for (let tag in tagMap) {
    arr.push(tag)
}
arr.sort()

for (let TAG of arr) {
    console.log(TAG)
}

const IMAGE = "ghcr.io/gardenlinux/gardenlinux/kmodbuild"

for (let TAG of arr) {
    console.log(TAG)

    const cmd = `docker run --platform linux/amd64 --rm -v ${__dirname}/../../charts/gpu-driver/files/gardenlinux-nvidia-installer:/mnt/scripts ${IMAGE}:${TAG} /mnt/scripts/extract_kernel_name.sh cloud`
    console.log(cmd)
}

// gh api "/orgs/gardenlinux/packages/container/gardenlinux%2Fkmodbuild/versions?per_page=100&page=6" | jq '.[] | .metadata.container.tags[0]' | grep arm
