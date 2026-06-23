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

import argparse
import json
import logging
from typing import Any, Optional

try:
    from transformers import TrainingArguments
except ImportError:  # pragma: no cover - exercised only when transformers is absent.

    class TrainingArguments:  # type: ignore[no-redef]
        def __init__(self, *args: Any, **kwargs: Any) -> None:
            raise ImportError(
                "transformers is required to construct HuggingFace TrainingArguments."
            )


logger = logging.getLogger(__name__)

DEFAULT_OUTPUT_DIR = "./output"


def parse_training_args(raw: Optional[str]) -> dict[str, Any]:
    """Parse a JSON string into a TrainingArguments configuration."""

    if raw is None:
        return {}

    if not isinstance(raw, str):
        raise ValueError(
            "training_parameters must be a JSON string or None; got "
            f"{type(raw).__name__}."
        )

    normalized = raw.strip()
    if not normalized:
        return {}

    try:
        parsed = json.loads(normalized)
    except json.JSONDecodeError as exc:
        raise ValueError(
            "Invalid JSON in training_parameters. Provide a JSON object string, for "
            f'example \'{{"output_dir": "./output"}}\'. Received: {raw!r}. '
            f"JSON error: {exc.msg} at line {exc.lineno}, column {exc.colno}."
        ) from exc

    if not isinstance(parsed, dict):
        raise ValueError(
            "training_parameters must decode to a JSON object. Received "
            f"{type(parsed).__name__}: {parsed!r}."
        )

    invalid_keys = [
        key for key in parsed.keys() if not isinstance(key, str) or not key.strip()
    ]
    if invalid_keys:
        raise ValueError(
            "training_parameters contains invalid keys. JSON object keys must be non-empty "
            f"strings. Invalid keys: {invalid_keys!r}."
        )

    return parsed


def build_training_arguments(raw: Optional[str]) -> TrainingArguments:
    logger.info("Raw training_parameters payload: %r", raw)
    parsed_config = parse_training_args(raw)

    if not parsed_config:
        logger.info(
            "training_parameters is empty or missing; using default "
            "TrainingArguments with output_dir=%s",
            DEFAULT_OUTPUT_DIR,
        )
        return TrainingArguments(output_dir=DEFAULT_OUTPUT_DIR)

    logger.info(
        "Parsed training_parameters config: %s",
        json.dumps(parsed_config, sort_keys=True),
    )
    try:
        return TrainingArguments(**parsed_config)
    except Exception as exc:
        logger.error(
            "Failed to create TrainingArguments from parsed training_parameters: %s",
            json.dumps(parsed_config, sort_keys=True),
            exc_info=True,
        )
        raise ValueError(
            "Failed to initialize TrainingArguments from training_parameters. "
            "Check the JSON keys and values, and ensure they match the HuggingFace "
            "TrainingArguments signature. Parsed config: "
            f"{json.dumps(parsed_config, sort_keys=True)}"
        ) from exc


def _build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(description="Run a HuggingFace training job.")
    parser.add_argument(
        "--training_parameters",
        type=str,
        default="{}",
        help="JSON object used to initialize HuggingFace TrainingArguments.",
    )
    return parser


def main() -> None:
    logging.basicConfig(
        level=logging.INFO, format="%(asctime)s %(levelname)s %(name)s %(message)s"
    )
    parser = _build_parser()
    args = parser.parse_args()

    training_args = build_training_arguments(args.training_parameters)
    logger.info("TrainingArguments initialized successfully: %s", training_args)

    # Replace this with the actual training workflow used by the example.
    logger.info("Trainer entrypoint completed parsing and initialization only.")


if __name__ == "__main__":
    main()
