// Copyright 2021 Authors of Cilium
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package exporter

import (
	"context"
	"encoding/json"
	"fmt"

	v1 "github.com/cilium/cilium/pkg/hubble/api/v1"
	"github.com/cilium/cilium/pkg/hubble/exporter/exporteroption"
	"github.com/cilium/cilium/pkg/hubble/observer"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// exporter is an implementation of OnDecodedEvent interface that writes Hubble events to a file.
type exporter struct {
	logger  logrus.FieldLogger
	encoder *json.Encoder
}

// NewExporter initializes an exporter.
func NewExporter(
	logger logrus.FieldLogger,
	options ...exporteroption.Option) (*exporter, error) {
	opts := exporteroption.Default // start with defaults
	for _, opt := range options {
		if err := opt(&opts); err != nil {
			return nil, fmt.Errorf("failed to apply option: %v", err)
		}
	}
	logger.WithField("options", opts).Info("Configuring Hubble event exporter")
	encoder := json.NewEncoder(&lumberjack.Logger{
		Filename:   opts.Path,
		MaxSize:    opts.MaxSizeMB,
		MaxBackups: opts.MaxBackups,
		Compress:   opts.Compress,
	})
	return newExporter(logger, encoder), nil
}

func newExporter(logger logrus.FieldLogger, encoder *json.Encoder) *exporter {
	return &exporter{
		logger:  logger,
		encoder: encoder,
	}
}

// Start calls GetFlows and writes responses to a file.
func (e *exporter) OnDecodedEvent(_ context.Context, ev *v1.Event) (bool, error) {
	res := observer.EventToGetFlowsResponse(ev)
	if res == nil {
		return false, nil
	}
	return false, e.encoder.Encode(res)
}
