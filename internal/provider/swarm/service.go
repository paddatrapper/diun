package swarm

import (
	"fmt"
	"reflect"

	"github.com/crazy-max/diun/internal/model"
	"github.com/crazy-max/diun/internal/provider"
	"github.com/crazy-max/diun/pkg/docker"
	"github.com/docker/docker/api/types/filters"
	"github.com/rs/zerolog/log"
)

func (c *Client) listServiceImage(elt model.PrdSwarm) []model.Image {
	sublog := log.With().
		Str("provider", fmt.Sprintf("swarm-%s", elt.ID)).
		Logger()

	cli, err := docker.NewClient(elt.Endpoint, elt.ApiVersion, elt.TLSCertsPath, elt.TLSVerify)
	if err != nil {
		sublog.Error().Err(err).Msg("Cannot create Docker client")
		return []model.Image{}
	}

	svcs, err := cli.ServiceList(filters.NewArgs())
	if err != nil {
		sublog.Error().Err(err).Msg("Cannot list Swarm services")
		return []model.Image{}
	}

	var list []model.Image
	for _, svc := range svcs {
		image, err := provider.ValidateContainerImage(svc.Spec.TaskTemplate.ContainerSpec.Image, svc.Spec.Labels, elt.WatchByDefault)
		if err != nil {
			sublog.Error().Err(err).Msgf("Cannot get image from service %s", svc.ID)
			continue
		} else if reflect.DeepEqual(image, model.Image{}) {
			sublog.Debug().Msgf("Watch disabled for service %s", svc.ID)
			continue
		}
		list = append(list, image)
	}

	return list
}
