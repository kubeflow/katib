from copy import deepcopy
from importlib import util
from pathlib import Path

import yaml


ROOT_DIR = Path(__file__).resolve().parents[4]
OVERRIDES_PATH = (
    ROOT_DIR
    / "test"
    / "e2e"
    / "v1beta1"
    / "scripts"
    / "gh-actions"
    / "experiment_overrides.py"
)


def _load_overrides_module():
    module_spec = util.spec_from_file_location("experiment_overrides", OVERRIDES_PATH)
    module = util.module_from_spec(module_spec)
    module_spec.loader.exec_module(module)
    return module


def _load_example(example_name):
    with open(ROOT_DIR / "examples" / "v1beta1" / "nas" / example_name, "r") as stream:
        return yaml.safe_load(stream)


def _algorithm_settings_by_name(example):
    return {
        setting["name"]: setting["value"]
        for setting in example["spec"]["algorithm"].get("algorithmSettings", [])
    }


def _trial_container_command(example):
    return example["spec"]["trialTemplate"]["trialSpec"]["spec"]["template"]["spec"][
        "containers"
    ][0]["command"]


def test_darts_cpu_e2e_override_uses_small_synthetic_dataset():
    overrides = _load_overrides_module()
    public_example = _load_example("darts-cpu.yaml")
    e2e_example = deepcopy(public_example)

    public_settings = _algorithm_settings_by_name(public_example)
    assert "use_synthetic_data" not in public_settings
    assert "batch_size" not in public_settings

    overrides.apply_e2e_experiment_overrides(e2e_example)

    e2e_settings = _algorithm_settings_by_name(e2e_example)
    assert e2e_settings["use_synthetic_data"] == "true"
    assert e2e_settings["number_of_examples"] == "64"
    assert e2e_settings["batch_size"] == "16"
    assert e2e_settings["num_workers"] == "0"


def test_enas_cpu_e2e_override_uses_small_synthetic_dataset():
    overrides = _load_overrides_module()
    public_example = _load_example("enas-cpu.yaml")
    e2e_example = deepcopy(public_example)

    public_command = _trial_container_command(public_example)
    assert "--use_synthetic_data" not in public_command
    assert "--number_of_examples=64" not in public_command

    overrides.apply_e2e_experiment_overrides(e2e_example)

    e2e_command = _trial_container_command(e2e_example)
    assert "--use_synthetic_data" in e2e_command
    assert "--number_of_examples=64" in e2e_command
