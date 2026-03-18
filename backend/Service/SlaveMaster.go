package service

import (
	pq "github.com/lib/pq"
	"database/sql"
	"fmt"
)

type SlaveMasterPair struct {
	// Slave is always an oracle credential
	Slave Credentials
	// Master is always a postgres credential
	Master Credentials
}

func (p *SlaveMasterPair) CreateSlave() error {
	return p.Slave.ConnectOracle()
}

func (p *SlaveMasterPair) CreateMaster() error {
	m := p.Master
	if m.db != nil {
		return nil
	}

	cfg := pq.Config{
		Host: p.Master.Server,
		Port: uint16(p.Master.Port),
		Database: p.Master.Database,
		Password: p.Master.Password,
		User: p.Master.User,
	}
	
	c, err := pq.NewConnectorConfig(cfg)
	if err != nil {
		return err
	}

	p.Master.db = sql.OpenDB(c)
	return nil
}

func (p *SlaveMasterPair) Ping() error {
	if p.Slave.db == nil {
		if err := p.CreateSlave(); err != nil {
			return fmt.Errorf("Failed to create slave: %v", err);
		}
	}

	if p.Master.db == nil {
		if err := p.CreateMaster(); err != nil {
			return fmt.Errorf("Failed to create master: %v", err);
		}
	}

	if err := p.Master.Ping(); err != nil {
		return err;
	}

	if err := p.Slave.Ping(); err != nil {
		return err;
	}

	return nil
}

func (p *SlaveMasterPair) Close() error {
	if err := p.Master.Close(); err != nil {
		return fmt.Errorf("Failed to close master: %v", err)
	}

	if err := p.Slave.Close(); err != nil {
		return fmt.Errorf("Failed to close slave: %v", err)
	}

	return nil
}

