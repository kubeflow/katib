from mock import MagicMock

from katib_suggestion.rpc_service import SuggestionService


def test_default_suggestion_service(rpc_request, study_config, observations):
    service = SuggestionService()
    service._parse_suggestion_parameters = MagicMock(return_value={})
    service._get_study_config = MagicMock(return_value=study_config)
    service._get_eval_history = MagicMock(return_value=(observations.parameters, observations.metrics))
    service._register_trials = MagicMock(return_value=[])
    service.GetSuggestions(request=rpc_request, context=None)


def test_boa_suggestion_service(rpc_request, study_config, observations):
    service = SuggestionService(search_algorithm="bayesian_optimization")
    service._parse_suggestion_parameters = MagicMock(return_value={})
    service._get_study_config = MagicMock(return_value=study_config)
    service._get_eval_history = MagicMock(return_value=(observations.parameters, observations.metrics))
    service._register_trials = MagicMock(return_value=[])
    service.GetSuggestions(request=rpc_request, context=None)

