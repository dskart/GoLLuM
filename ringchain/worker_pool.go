package ringchain

import (
	"context"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type nodeJob struct {
	Args   map[string]any
	Node   Node
	Logger *zap.Logger
	Result chan<- map[string]any
}

type runningNode struct {
	NodeHash string
	Result   <-chan map[string]any
}

type workerPool struct {
	numWorkers   int
	queueSize    int
	jobQueue     chan nodeJob
	doneChan     chan bool
	runningNodes []runningNode
	errChan      chan error

	eg *errgroup.Group
}

func newWorkerPool(numWorkers int, queueSize int) *workerPool {
	return &workerPool{
		numWorkers:   numWorkers,
		queueSize:    queueSize,
		jobQueue:     make(chan nodeJob, queueSize),
		doneChan:     make(chan bool, numWorkers),
		runningNodes: make([]runningNode, 0, queueSize),
		errChan:      make(chan error),
	}
}

func (p *workerPool) Run(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)
	for w := 1; w <= p.numWorkers; w++ {
		eg.Go(func() error {
			return worker(ctx, p.doneChan, p.jobQueue, p.errChan)
		})
	}
	p.eg = eg
	return nil
}

func (p *workerPool) NumRunningNodes() int {
	return len(p.runningNodes)
}

func (p *workerPool) AppendRunningNode(runningNode runningNode) {
	p.runningNodes = append(p.runningNodes, runningNode)
}

func (p *workerPool) PopRunningNode() runningNode {
	runningNode := p.runningNodes[0]
	p.runningNodes = p.runningNodes[1:]
	return runningNode
}

func (p *workerPool) AddNodeJob(logger *zap.Logger, nodeHash string, node Node, args map[string]any) {
	resultChan := make(chan map[string]any)
	p.jobQueue <- nodeJob{
		Args:   args,
		Node:   node,
		Logger: logger,
		Result: resultChan,
	}
	p.runningNodes = append(p.runningNodes, runningNode{
		NodeHash: nodeHash,
		Result:   resultChan,
	})
}

func (p *workerPool) Stop() {
	for i := 0; i < p.numWorkers; i++ {
		p.doneChan <- true
	}
}

func (p *workerPool) HasErr() <-chan error {
	return p.errChan
}

func (p *workerPool) Wait() error {
	if p.eg == nil {
		return nil
	}
	if err := p.eg.Wait(); err != nil {
		return err
	}
	return nil
}

func worker(ctx context.Context, done <-chan bool, jobQueue <-chan nodeJob, errChan chan<- error) error {
	for {
		select {
		case <-done:
			return nil
		case job := <-jobQueue:
			node := job.Node
			res, err := node.Run(ctx, job.Logger, job.Args)
			if err != nil {
				errChan <- err
				return err
			}
			job.Result <- res

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
