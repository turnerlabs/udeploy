package app

// Copy ...
func (a Application) Copy() Application {
	c := Application{
		ID:        a.ID,
		Name:      a.Name,
		Type:      a.Type,
		Repo:      a.Repo.Copy(),
		ProjectID: a.ProjectID,
		Instances: map[string]Instance{},
	}

	for name, inst := range a.Instances {
		c.Instances[name] = inst.Copy()
	}

	return c
}

// Copy ...
func (r Repository) Copy() Repository {
	return Repository{
		Org:          r.Org,
		Name:         r.Name,
		AccessToken:  r.AccessToken,
		CommitConfig: r.CommitConfig.Copy(),
	}
}

// Copy ...
func (r CommitConfig) Copy() CommitConfig {
	return CommitConfig{
		ExistingValue: r.ExistingValue,
		NewValue:      r.NewValue,
		Filter:        r.Filter,
		Limit:         r.Limit,
	}
}

// Copy ...
func (t Task) Copy() Task {
	n := Task{
		ImageTagEx:     t.ImageTagEx,
		CronEnabled:    t.CronEnabled,
		Family:         t.Family,
		Registry:       t.Registry,
		Revisions:      t.Revisions,
		CloneEnvVars:   t.CloneEnvVars,
		DesiredCount:   t.DesiredCount,
		CronExpression: t.CronExpression,
		Definition:     t.Definition.Copy(),
		TasksInfo:      []TaskInfo{},
	}

	for _, ti := range t.TasksInfo {
		n.TasksInfo = append(n.TasksInfo, ti)
	}

	return n
}

// Copy ...
func (i Instance) Copy() Instance {

	n := Instance{
		Cluster:          i.Cluster,
		Service:          i.Service,
		Role:             i.Role,
		RepositoryRole:   i.RepositoryRole,
		EventRule:        i.EventRule,
		Repository:       i.Repository,
		FunctionName:     i.FunctionName,
		FunctionAlias:    i.FunctionAlias,
		S3Bucket:         i.S3Bucket,
		S3Prefix:         i.S3Prefix,
		S3ConfigKey:      i.S3ConfigKey,
		S3RegistryBucket: i.S3RegistryBucket,
		S3RegistryPrefix: i.S3RegistryPrefix,
		DeployCode:       i.DeployCode,
		AutoPropagate:    i.AutoPropagate,
		AutoScale:        i.AutoScale,
		Order:            i.Order,
		Task:             i.Task.Copy(),
		CurrentState:     i.CurrentState,
		PreviousState:    i.PreviousState,
	}

	for _, l := range i.Links {
		n.Links = append(n.Links, l)
	}

	return n
}

// Copy ...
func (d Definition) Copy() Definition {
	n := Definition{
		ID:          d.ID,
		Version:     d.Version,
		Build:       d.Build,
		Revision:    d.Revision,
		Description: d.Description,
		Environment: map[string]string{},
		Secrets:     map[string]string{},
	}

	for k, v := range d.Environment {
		n.Environment[k] = v
	}

	for k, v := range d.Secrets {
		n.Secrets[k] = v
	}

	return n
}
