import subprocess
from typing import List, Sequence

DEFAULT_TIMEOUT_SECONDS = 45
LOG_TIMEOUT_SECONDS = 90


def _print_section(title: str) -> None:
    print("")
    print("=" * 80)
    print(title)
    print("=" * 80, flush=True)


def _run_section(
    title: str,
    command: Sequence[str],
    timeout: int = DEFAULT_TIMEOUT_SECONDS,
) -> None:
    _print_section(title)
    print("+ " + " ".join(command), flush=True)
    try:
        result = subprocess.run(command, check=False, timeout=timeout)
    except FileNotFoundError as err:
        print(f"Failed to run command: {err}", flush=True)
        return
    except subprocess.TimeoutExpired:
        print(f"Command timed out after {timeout}s", flush=True)
        return

    if result.returncode != 0:
        print(f"Command exited with status {result.returncode}", flush=True)


def _get_pods(namespace: str) -> List[str]:
    try:
        result = subprocess.run(
            ["kubectl", "get", "pods", "-n", namespace, "-o", "name"],
            check=False,
            capture_output=True,
            text=True,
            timeout=DEFAULT_TIMEOUT_SECONDS,
        )
    except FileNotFoundError as err:
        print(f"Failed to list pods: {err}", flush=True)
        return []
    except subprocess.TimeoutExpired:
        print(
            f"Listing pods timed out after {DEFAULT_TIMEOUT_SECONDS}s",
            flush=True,
        )
        return []

    if result.returncode != 0:
        print("Failed to list pods for log collection.", flush=True)
        if result.stderr:
            print(result.stderr, flush=True)
        return []

    return [line.strip() for line in result.stdout.splitlines() if line.strip()]


def _dump_pod_logs(namespace: str) -> None:
    pods = _get_pods(namespace)
    if not pods:
        _print_section(f"No pods found for log collection in namespace {namespace}")
        return

    for pod in pods:
        _run_section(
            f"Logs for {namespace}/{pod}",
            [
                "kubectl",
                "logs",
                "-n",
                namespace,
                pod,
                "--all-containers",
                "--prefix",
                "--tail=300",
            ],
            timeout=LOG_TIMEOUT_SECONDS,
        )
        _run_section(
            f"Previous logs for {namespace}/{pod}",
            [
                "kubectl",
                "logs",
                "-n",
                namespace,
                pod,
                "--all-containers",
                "--prefix",
                "--tail=300",
                "--previous",
            ],
            timeout=LOG_TIMEOUT_SECONDS,
        )


def dump_e2e_diagnostics(exp_name: str, namespace: str) -> None:
    _print_section(f"Katib E2E diagnostics for {namespace}/{exp_name}")

    commands = [
        (
            "Experiment YAML",
            ["kubectl", "get", "experiment", exp_name, "-n", namespace, "-o", "yaml"],
        ),
        (
            "Experiment describe",
            ["kubectl", "describe", "experiment", exp_name, "-n", namespace],
        ),
        (
            "Suggestion YAML",
            ["kubectl", "get", "suggestion", exp_name, "-n", namespace, "-o", "yaml"],
        ),
        (
            "Suggestion describe",
            ["kubectl", "describe", "suggestion", exp_name, "-n", namespace],
        ),
        (
            "Experiment Trials",
            [
                "kubectl",
                "get",
                "trials",
                "-n",
                namespace,
                "-l",
                f"katib.kubeflow.org/experiment={exp_name}",
                "-o",
                "wide",
                "--show-labels",
            ],
        ),
        (
            "Experiment Trial describe",
            [
                "kubectl",
                "describe",
                "trials",
                "-n",
                namespace,
                "-l",
                f"katib.kubeflow.org/experiment={exp_name}",
            ],
        ),
        (
            "Namespace workload summary",
            [
                "kubectl",
                "get",
                "jobs,pods,deploy,svc,pvc",
                "-n",
                namespace,
                "-o",
                "wide",
                "--show-labels",
            ],
        ),
        (
            "Namespace jobs and pods YAML",
            ["kubectl", "get", "jobs,pods", "-n", namespace, "-o", "yaml"],
        ),
        (
            "Namespace PyTorchJobs YAML",
            ["kubectl", "get", "pytorchjobs.kubeflow.org", "-n", namespace, "-o", "yaml"],
        ),
        (
            "Namespace events",
            ["kubectl", "get", "events", "-n", namespace, "--sort-by=.lastTimestamp"],
        ),
        (
            "Namespace pod describe",
            ["kubectl", "describe", "pods", "-n", namespace],
        ),
    ]

    for title, command in commands:
        _run_section(title, command)

    _dump_pod_logs(namespace)

    _run_section(
        "Kubeflow Katib controller resources",
        [
            "kubectl",
            "-n",
            "kubeflow",
            "get",
            "deploy,svc,endpoints,endpointslice,pods",
            "-o",
            "wide",
        ],
    )
    _run_section(
        "Katib controller logs",
        [
            "kubectl",
            "-n",
            "kubeflow",
            "logs",
            "deploy/katib-controller",
            "--all-containers",
            "--tail=500",
        ],
        timeout=LOG_TIMEOUT_SECONDS,
    )
