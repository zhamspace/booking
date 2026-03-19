package audit

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"

	"github.com/google/uuid"
	"github.com/zhamspace/booking/internal/domain/audit/model"
	"github.com/zhamspace/booking/internal/domain/common/session"
	"github.com/zhamspace/booking/internal/errs"
)

type Service struct {
	repoDb     RepoDbI
	serializer ObjectSerializer
	queue      chan *model.Job
	wg         sync.WaitGroup
	mu         sync.RWMutex
	startOnce  sync.Once
	stopOnce   sync.Once
	started    atomic.Bool
	stopped    atomic.Bool
}

func New(
	repoDb RepoDbI,
	serializer ObjectSerializer,
) *Service {
	return &Service{
		repoDb:     repoDb,
		serializer: serializer,
		queue:      make(chan *model.Job, 1024),
	}
}

func (s *Service) Start(_ context.Context) error {
	s.startOnce.Do(func() {
		s.started.Store(true)
		s.wg.Add(1)
		go s.run()
	})
	return nil
}

func (s *Service) Stop(ctx context.Context) error {
	if !s.started.Load() {
		return nil
	}

	s.stopOnce.Do(func() {
		s.stopped.Store(true)
		s.mu.Lock()
		close(s.queue)
		s.mu.Unlock()
	})
	return s.Wait(ctx)
}

func (s *Service) Wait(ctx context.Context) error {
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}

func (s *Service) run() {
	defer s.wg.Done()

	for job := range s.queue {
		if err := s.store(job); err != nil {
			slog.Error("audit.Service: failed to store job", "error", err)
		}
	}
}

func (s *Service) store(job *model.Job) error {
	if job == nil {
		return errs.New(errs.InvalidConfig, "audit job is required", nil)
	}
	if s.serializer == nil {
		return errs.New(errs.InvalidConfig, "audit serializer is required", nil)
	}

	data, err := s.serializer.Serialize(job.Resource, job.Object)
	if err != nil {
		return err
	}

	_, err = s.repoDb.Create(context.Background(), &model.Edit{
		Id:         new(uuid.NewString()),
		UserId:     &job.UserId,
		ChangeType: &job.ChangeType,
		Resource:   &job.Resource,
		ObjectId:   &job.ObjectId,
		Object:     &data,
	})
	return err
}

func (s *Service) Log(ctx context.Context, job *model.Job) {
	if job == nil {
		return
	}
	if s.stopped.Load() {
		return
	}

	ctx = context.WithoutCancel(ctx)

	select {
	case <-ctx.Done():
		return
	default:
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	ses := session.ExtractFromContext(ctx)
	if ses.Sub != "" {
		job.UserId = ses.Sub
	}

	select {
	case s.queue <- new(*job):
		return
	default:
		slog.Warn("audit.Service: queue is full, dropping job",
			"user_id", job.UserId,
			"resource", job.Resource,
			"change_type", job.ChangeType,
			"object", job.Object,
		)
		return
	}
}

func (s *Service) List(ctx context.Context, pars *model.ListReq) ([]*model.Main, int64, error) {
	items, totalCount, err := s.repoDb.List(ctx, pars)
	if err != nil {
		return nil, 0, err
	}
	return items, totalCount, nil
}

func (s *Service) Get(ctx context.Context, pars *model.GetReq, errNE bool) (*model.Main, bool, error) {
	item, found, err := s.repoDb.Get(ctx, pars)
	if err != nil {
		return nil, false, err
	}
	if !found {
		if errNE {
			return nil, false, errs.ObjectNotFound
		}
		return nil, false, nil
	}
	return item, found, nil
}
