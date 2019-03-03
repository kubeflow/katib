import pytest
from mock import MagicMock
from box import Box

from pkg.suggestion.bayesianoptimization.src.rpc_service import SuggestionService


@pytest.fixture
def rpc_request():
    return Box({"param_id": "test_param",
                "study_id": "test_study",
                "request_number": 2})


def test_suggestion_service(rpc_request, study_config, observations):
    service = SuggestionService()
    service._parse_suggestion_parameters = MagicMock(return_value={})
    service._get_study_config = MagicMock(return_value=study_config)
    service._get_eval_history = MagicMock(return_value=(observations.parameters, observations.metrics))
    service._register_trials = MagicMock(return_value=[])
    service.GetSuggestions(request=rpc_request, context=None)
