package cron_jobs

type CronJob struct {
	job        cron.Cron
	paramToJob any
}

func NewCronJob(paramToJob any) *CronJob {
	return &CronJob{
		job:        cron.New(),
		paramToJob: paramToJob,
	}
}

func (c *CronJob) RegisterJob(cron string, handler func(any) error) {
	c.job.AddFunc(cron, handler)
}

func (c *CronJob) Run() {
	c.job.Start()
}
