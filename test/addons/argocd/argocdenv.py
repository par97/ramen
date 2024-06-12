#!/usr/bin/env python3

# SPDX-FileCopyrightText: The RamenDR authors
# SPDX-License-Identifier: Apache-2.0

import os

from drenv import kubectl


def get_argocd_env(tmpdir, cluster):
    # Create a temporary kubeconfig so we don't modify the shared config
    # used by other addons concurrently.
    kubeconfig = os.path.join(tmpdir, "kubeconfig")
    # print(f"Creating temporary kubeconfig {kubeconfig}")
    out = kubectl.config("view", "--flatten", "--output=yaml")
    with open(kubeconfig, "w") as f:
        f.write(out)
    kubectl.config("use-context", f"--kubeconfig={kubeconfig}", cluster)
    kubectl.config(
        "set-context",
        f"--kubeconfig={kubeconfig}",
        "--current",
        "--namespace=argocd",
    )

    # Create an environemnt so we can pass the kubeconfig to argocd.
    env = dict(os.environ)
    env["KUBECONFIG"] = kubeconfig
    return env
