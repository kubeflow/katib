import pytest
from kubeflow.katib import search
from kubeflow.katib.utils import utils


@pytest.mark.parametrize(
    "pip_index_urls, expected",
    [
        (["https://pypi.org/simple"], "--index-url https://pypi.org/simple"),
        (
            ["https://pypi.org/simple", "https://private-repo.com/simple"],
            "--index-url https://pypi.org/simple --extra-index-url https://private-repo.com/simple",
        ),
        (
            [
                "https://pypi.org/simple",
                "https://private-repo.com/simple",
                "https://another-repo.com/simple",
            ],
            "--index-url https://pypi.org/simple --extra-index-url https://private-repo.com/simple "
            "--extra-index-url https://another-repo.com/simple",
        ),
    ],
)
def test_format_pip_index_urls(pip_index_urls, expected):
    assert utils.format_pip_index_urls(pip_index_urls) == expected


def test_reuse_search_parameters():
    bounds = search.double(min=0.1, max=0.9)
    parameters = {"alpha": bounds, "beta": bounds, "epochs": 5}

    for _ in range(2):
        experiment_parameters = []
        trial_parameters = []
        inputs = utils.get_trial_substitutions_from_dict(
            parameters, experiment_parameters, trial_parameters
        )

        assert [param.name for param in experiment_parameters] == ["alpha", "beta"]
        assert [param.reference for param in trial_parameters] == ["alpha", "beta"]
        assert inputs == {
            "alpha": "${trialParameters.alpha}",
            "beta": "${trialParameters.beta}",
            "epochs": 5,
        }

    assert parameters == {"alpha": bounds, "beta": bounds, "epochs": 5}
    assert bounds.name is None
