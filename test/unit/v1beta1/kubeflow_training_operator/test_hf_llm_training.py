# Copyright 2024 The Kubeflow Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

from __future__ import annotations

import importlib.util
from pathlib import Path

import pytest

SCRIPT_PATH = (
    Path(__file__).resolve().parents[4]
    / "examples"
    / "v1beta1"
    / "kubeflow-training-operator"
    / "hf_llm_training.py"
)


class DummyTrainingArguments:
    def __init__(self, **kwargs):
        self.kwargs = kwargs


def load_module():
    spec = importlib.util.spec_from_file_location("hf_llm_training", SCRIPT_PATH)
    module = importlib.util.module_from_spec(spec)
    assert spec.loader is not None
    spec.loader.exec_module(module)
    return module


def test_parse_training_args_empty_string_returns_empty_dict():
    module = load_module()

    assert module.parse_training_args("") == {}


def test_parse_training_args_none_returns_empty_dict():
    module = load_module()

    assert module.parse_training_args(None) == {}


def test_parse_training_args_whitespace_returns_empty_dict():
    module = load_module()

    assert module.parse_training_args("   \n\t ") == {}


def test_parse_training_args_valid_json_returns_dict():
    module = load_module()

    assert module.parse_training_args(
        '{"output_dir": "./output", "learning_rate": 0.0001}'
    ) == {
        "output_dir": "./output",
        "learning_rate": 0.0001,
    }


def test_parse_training_args_invalid_json_raises_value_error():
    module = load_module()

    with pytest.raises(ValueError, match="Invalid JSON in training_parameters"):
        module.parse_training_args("{invalid-json")


def test_parse_training_args_malformed_keys_raises_value_error():
    module = load_module()

    with pytest.raises(ValueError, match="invalid keys"):
        module.parse_training_args('{"": 1, "  ": 2}')


def test_build_training_arguments_uses_default_when_empty(monkeypatch):
    module = load_module()
    monkeypatch.setattr(module, "TrainingArguments", DummyTrainingArguments)

    training_args = module.build_training_arguments("")

    assert training_args.kwargs == {"output_dir": "./output"}


def test_build_training_arguments_passes_valid_config(monkeypatch):
    module = load_module()
    monkeypatch.setattr(module, "TrainingArguments", DummyTrainingArguments)

    training_args = module.build_training_arguments(
        '{"output_dir": "./tmp", "num_train_epochs": 3}'
    )

    assert training_args.kwargs == {"output_dir": "./tmp", "num_train_epochs": 3}
