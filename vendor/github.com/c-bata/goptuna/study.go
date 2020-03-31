package goptuna

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

var errCreateNewTrial = errors.New("failed to create a new trial")

// FuncObjective is a type of objective function
type FuncObjective func(trial Trial) (float64, error)

// StudyDirection represents the direction of the optimization
type StudyDirection string

const (
	// StudyDirectionMaximize maximizes objective function value
	StudyDirectionMaximize StudyDirection = "maximize"
	// StudyDirectionMinimize minimizes objective function value
	StudyDirectionMinimize StudyDirection = "minimize"
)

// Study corresponds to an optimization task, i.e., a set of trials.
type Study struct {
	ID                int
	Storage           Storage
	Sampler           Sampler
	RelativeSampler   RelativeSampler
	Pruner            Pruner
	direction         StudyDirection
	logger            Logger
	ignoreErr         bool
	trialNotification chan FrozenTrial
	loadIfExists      bool
	mu                sync.RWMutex
	ctx               context.Context
}

// EnqueueTrial to enqueue a trial with given parameter values.
// You can fix the next sampling parameters which will be evaluated in your
// objective function.
//
// This is an EXPERIMENTAL API and may be changed in the future.
// Currently EnqueueTrial only accepts internal representations.
// This means that you need to encode categorical parameter to its index number.
// Furthermore, please caution that there is a concurrency problem on SQLite3.
func (s *Study) EnqueueTrial(internalParams map[string]float64) error {
	systemAttrs := make(map[string]string, 8)
	paramJSONBytes, err := json.Marshal(internalParams)
	if err != nil {
		return err
	}

	systemAttrs["fixed_params"] = string(paramJSONBytes)
	return s.appendTrial(
		0,
		nil,
		nil,
		nil,
		systemAttrs,
		nil,
		TrialStateWaiting,
		time.Now(),
		time.Time{},
	)
}

func (s *Study) popWaitingTrialID() (int, error) {
	trials, err := s.Storage.GetAllTrials(s.ID)
	if err != nil {
		return -1, err
	}

	// TODO(c-bata): Reduce database query counts for extracting waiting trials.
	for i := range trials {
		if trials[i].State != TrialStateWaiting {
			continue
		}

		err = s.Storage.SetTrialState(trials[i].ID, TrialStateRunning)
		if err == ErrTrialCannotBeUpdated {
			err = nil
			continue
		} else if err != nil {
			return -1, err
		}
		s.logger.Debug("trial is popped from the trial queue.",
			fmt.Sprintf("number=%d", trials[i].Number))
		return trials[i].ID, nil
	}
	return -1, nil
}

// AppendTrial to inject a trial into the Study.
func (s *Study) appendTrial(
	value float64,
	internalParams map[string]float64,
	distributions map[string]interface{},
	userAttrs map[string]string,
	systemAttrs map[string]string,
	intermediateValues map[int]float64,
	state TrialState,
	datetimeStart time.Time,
	datetimeComplete time.Time,
) error {
	params := make(map[string]interface{}, len(internalParams))
	for name := range internalParams {
		d, ok := distributions[name]
		if !ok {
			return fmt.Errorf("distribution '%s' is not found", name)
		}
		xr, err := ToExternalRepresentation(d, internalParams[name])
		if err != nil {
			return err
		}
		params[name] = xr
	}
	if state.IsFinished() && datetimeComplete.IsZero() {
		datetimeComplete = time.Now()
	}
	if distributions == nil {
		distributions = make(map[string]interface{})
	}
	if internalParams == nil {
		internalParams = make(map[string]float64)
	}
	if userAttrs == nil {
		userAttrs = make(map[string]string)
	}
	if systemAttrs == nil {
		systemAttrs = make(map[string]string)
	}
	if intermediateValues == nil {
		intermediateValues = make(map[int]float64)
	}
	trial := FrozenTrial{
		ID:                 -1, // dummy value
		StudyID:            s.ID,
		Number:             -1, // dummy value
		State:              state,
		Value:              value,
		IntermediateValues: intermediateValues,
		DatetimeStart:      datetimeStart,
		DatetimeComplete:   datetimeComplete,
		InternalParams:     internalParams,
		Params:             params,
		Distributions:      distributions,
		UserAttrs:          userAttrs,
		SystemAttrs:        systemAttrs,
	}
	err := trial.validate()
	if err != nil {
		return err
	}
	_, err = s.Storage.CloneTrial(s.ID, trial)
	return err
}

// GetTrials returns all trials in this study.
func (s *Study) GetTrials() ([]FrozenTrial, error) {
	return s.Storage.GetAllTrials(s.ID)
}

// Direction returns the direction of objective function value
func (s *Study) Direction() StudyDirection {
	return s.direction
}

// WithContext sets a context and it might cancel the execution of Optimize.
func (s *Study) WithContext(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ctx = ctx
}

func (s *Study) callRelativeSampler(trialID int) (
	map[string]interface{},
	map[string]float64,
	error,
) {
	if s.RelativeSampler == nil {
		return nil, nil, nil
	}

	frozen, err := s.Storage.GetTrial(trialID)
	if err != nil {
		return nil, nil, err
	}

	intersection, err := IntersectionSearchSpace(s)
	if err != nil {
		return nil, nil, err
	}
	if intersection == nil {
		return nil, nil, nil
	}

	searchSpace := make(map[string]interface{}, len(intersection))
	for paramName := range intersection {
		distribution := intersection[paramName]
		if yes, _ := DistributionIsSingle(distribution); yes {
			continue
		}
		searchSpace[paramName] = intersection[paramName]
	}
	relativeParams, err := s.RelativeSampler.SampleRelative(s, frozen, searchSpace)
	if err == ErrUnsupportedSearchSpace {
		s.logger.Warn("Your objective function contains unsupported search space for RelativeSampler.",
			fmt.Sprintf("trialID=%d", trialID),
			fmt.Sprintf("searchSpace=%#v", searchSpace))
		return nil, nil, nil
	} else if err != nil {
		return nil, nil, err
	}
	return searchSpace, relativeParams, nil
}

func (s *Study) runTrial(objective FuncObjective) (int, error) {
	trialID, err := s.popWaitingTrialID()
	if err != nil {
		s.logger.Error("failed to pop a waiting trial",
			fmt.Sprintf("err=%s", err))
		return -1, err
	}
	if trialID == -1 {
		trialID, err = s.Storage.CreateNewTrial(s.ID)
		if err != nil {
			s.logger.Error("failed to create a new trial",
				fmt.Sprintf("err=%s", err))
			return -1, errCreateNewTrial
		}
	}
	searchSpace, relativeParams, err := s.callRelativeSampler(trialID)
	if err != nil {
		s.logger.Error("failed to call relative sampler",
			fmt.Sprintf("err=%s", err))
		return -1, err
	}

	trial := Trial{
		Study:               s,
		ID:                  trialID,
		relativeParams:      relativeParams,
		relativeSearchSpace: searchSpace,
	}
	evaluation, objerr := objective(trial)
	var state TrialState
	if objerr == ErrTrialPruned {
		state = TrialStatePruned
		objerr = nil
	} else if objerr != nil {
		state = TrialStateFail
	} else {
		state = TrialStateComplete
	}

	if objerr != nil {
		s.logger.Error("Objective function returns error",
			fmt.Sprintf("trialID=%d", trialID),
			fmt.Sprintf("state=%s", state.String()),
			fmt.Sprintf("err=%s", objerr))
	} else {
		s.logger.Info("Trial finished",
			fmt.Sprintf("trialID=%d", trialID),
			fmt.Sprintf("state=%s", state.String()),
			fmt.Sprintf("evaluation=%f", evaluation))
	}

	if state == TrialStateComplete {
		// The trial.value of pruned trials are already set at trial.Report().
		err = s.Storage.SetTrialValue(trialID, evaluation)
		if err != nil {
			s.logger.Error("Failed to set trial value",
				fmt.Sprintf("trialID=%d", trialID),
				fmt.Sprintf("state=%s", state.String()),
				fmt.Sprintf("evaluation=%f", evaluation),
				fmt.Sprintf("err=%s", err))
			return trialID, err
		}
	} else if state == TrialStatePruned {
		// Register the last intermediate value if present as the value of the trial.
		trial, err := s.Storage.GetTrial(trialID)
		if err != nil {
			return -1, err
		}
		if lastStep, exists := trial.GetLatestStep(); exists {
			err = s.Storage.SetTrialValue(trialID, trial.IntermediateValues[lastStep])
			s.logger.Error("Failed to set trial value",
				fmt.Sprintf("trialID=%d", trialID),
				fmt.Sprintf("state=%s", state.String()),
				fmt.Sprintf("evaluation=%f", evaluation),
				fmt.Sprintf("err=%s", err))
			if err != nil {
				return -1, err
			}
		}
	}

	err = s.Storage.SetTrialState(trialID, state)
	if err != nil {
		s.logger.Error("Failed to set trial state",
			fmt.Sprintf("trialID=%d", trialID),
			fmt.Sprintf("state=%s", state.String()),
			fmt.Sprintf("evaluation=%f", evaluation),
			fmt.Sprintf("err=%s", err))
		return trialID, err
	}
	return trialID, objerr
}

// Optimize optimizes an objective function.
func (s *Study) Optimize(objective FuncObjective, evaluateMax int) error {
	evaluateCnt := 0
	for {
		if evaluateCnt >= evaluateMax {
			break
		}

		if s.ctx != nil {
			select {
			case <-s.ctx.Done():
				err := s.ctx.Err()
				s.logger.Debug("context is canceled", err)
				return err
			default:
				// do nothing
			}
		}
		// Evaluate an objective function
		trialID, err := s.runTrial(objective)
		if err == errCreateNewTrial {
			continue
		}
		evaluateCnt++

		// Send trial notification
		if s.trialNotification != nil {
			frozen, gerr := s.Storage.GetTrial(trialID)
			if gerr != nil {
				s.logger.Error("Failed to send trial notification",
					fmt.Sprintf("trialID=%d", trialID),
					fmt.Sprintf("err=%s", gerr))
				if !s.ignoreErr {
					return gerr
				}
			}
			s.trialNotification <- frozen
		}

		if !s.ignoreErr && err != nil {
			return err
		}
	}
	return nil
}

// GetBestValue return the best objective value
func (s *Study) GetBestValue() (float64, error) {
	trial, err := s.Storage.GetBestTrial(s.ID)
	if err != nil {
		return 0.0, err
	}
	return trial.Value, nil
}

// GetBestParams return parameters of the best trial
func (s *Study) GetBestParams() (map[string]interface{}, error) {
	trial, err := s.Storage.GetBestTrial(s.ID)
	if err != nil {
		return nil, err
	}
	return trial.Params, nil
}

// SetUserAttr to store the value for the user.
func (s *Study) SetUserAttr(key, value string) error {
	return s.Storage.SetStudyUserAttr(s.ID, key, value)
}

// SetSystemAttr to store the value for the system.
func (s *Study) SetSystemAttr(key, value string) error {
	return s.Storage.SetStudySystemAttr(s.ID, key, value)
}

// GetUserAttrs to store the value for the user.
func (s *Study) GetUserAttrs() (map[string]string, error) {
	return s.Storage.GetStudyUserAttrs(s.ID)
}

// GetSystemAttrs to store the value for the system.
func (s *Study) GetSystemAttrs() (map[string]string, error) {
	return s.Storage.GetStudySystemAttrs(s.ID)
}

// GetLogger returns logger object.
func (s *Study) GetLogger() Logger {
	return s.logger
}

// CreateStudy creates a new Study object.
func CreateStudy(
	name string,
	opts ...StudyOption,
) (*Study, error) {
	storage := NewInMemoryStorage()
	sampler := NewRandomSearchSampler()
	study := &Study{
		ID:              0,
		Storage:         storage,
		Sampler:         sampler,
		RelativeSampler: nil,
		Pruner:          nil,
		direction:       StudyDirectionMinimize,
		logger: &StdLogger{
			Logger: log.New(os.Stdout, "", log.LstdFlags),
			Level:  LoggerLevelDebug,
			Color:  true,
		},
		ignoreErr: false,
	}

	for _, opt := range opts {
		if err := opt(study); err != nil {
			return nil, err
		}
	}

	if study.loadIfExists {
		study, err := LoadStudy(name, opts...)
		if err == nil {
			return study, nil
		}
	}

	studyID, err := study.Storage.CreateNewStudy(name)
	if err != nil {
		return nil, err
	}
	err = study.Storage.SetStudyDirection(studyID, study.direction)
	if err != nil {
		return nil, err
	}
	study.ID = studyID
	return study, nil
}

// LoadStudy loads an existing study.
func LoadStudy(
	name string,
	opts ...StudyOption,
) (*Study, error) {
	storage := NewInMemoryStorage()
	sampler := NewRandomSearchSampler()
	study := &Study{
		ID:              0,
		Storage:         storage,
		Sampler:         sampler,
		RelativeSampler: nil,
		Pruner:          nil,
		direction:       "",
		logger: &StdLogger{
			Logger: log.New(os.Stdout, "", log.LstdFlags),
			Level:  LoggerLevelDebug,
			Color:  true,
		},
		ignoreErr: false,
	}

	for _, opt := range opts {
		if err := opt(study); err != nil {
			return nil, err
		}
	}

	studyID, err := study.Storage.GetStudyIDFromName(name)
	if err != nil {
		return nil, err
	}
	study.ID = studyID
	direction, err := study.Storage.GetStudyDirection(studyID)
	if err != nil {
		return nil, err
	}
	study.direction = direction
	return study, nil
}

// DeleteStudy delete a study object.
func DeleteStudy(
	name string,
	storage Storage,
) error {
	studyID, err := storage.GetStudyIDFromName(name)
	if err != nil {
		return err
	}
	return storage.DeleteStudy(studyID)
}

// StudyOption to pass the custom option
type StudyOption func(study *Study) error

// StudyOptionSetDirection change the direction of optimize
func StudyOptionSetDirection(direction StudyDirection) StudyOption {
	return func(s *Study) error {
		s.direction = direction
		return nil
	}
}

// StudyOptionLogger sets Logger.
func StudyOptionLogger(logger Logger) StudyOption {
	return func(s *Study) error {
		if logger == nil {
			s.logger = &StdLogger{Logger: nil}
		} else {
			s.logger = logger
		}
		return nil
	}
}

// StudyOptionStorage sets the storage object.
func StudyOptionStorage(storage Storage) StudyOption {
	return func(s *Study) error {
		s.Storage = storage
		return nil
	}
}

// StudyOptionSampler sets the sampler object.
func StudyOptionSampler(sampler Sampler) StudyOption {
	return func(s *Study) error {
		s.Sampler = sampler
		return nil
	}
}

// StudyOptionRelativeSampler sets the relative sampler object.
func StudyOptionRelativeSampler(sampler RelativeSampler) StudyOption {
	return func(s *Study) error {
		s.RelativeSampler = sampler
		return nil
	}
}

// StudyOptionPruner sets the pruner object.
func StudyOptionPruner(pruner Pruner) StudyOption {
	return func(s *Study) error {
		s.Pruner = pruner
		return nil
	}
}

// StudyOptionIgnoreError is an option to continue even if
// it receive error while running Optimize method.
func StudyOptionIgnoreError(ignore bool) StudyOption {
	return func(s *Study) error {
		s.ignoreErr = ignore
		return nil
	}
}

// StudyOptionSetTrialNotifyChannel to subscribe the finished trials.
func StudyOptionSetTrialNotifyChannel(notify chan FrozenTrial) StudyOption {
	return func(s *Study) error {
		s.trialNotification = notify
		return nil
	}
}

// StudyOptionLoadIfExists to load the study if exists.
func StudyOptionLoadIfExists(loadIfExists bool) StudyOption {
	return func(s *Study) error {
		s.loadIfExists = loadIfExists
		return nil
	}
}

// StudyOptionSetLogger sets Logger.
// Deprecated: please use StudyOptionLogger instead.
var StudyOptionSetLogger = StudyOptionLogger
