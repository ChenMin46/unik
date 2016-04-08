package state

import (
	"github.com/emc-advanced-dev/unik/pkg/types"
	"sync"
	"encoding/json"
	"github.com/layer-x/layerx-commons/lxerrors"
	"os"
	"path/filepath"
	"io/ioutil"
)

type State interface {
	GetImages() map[string]*types.Image
	GetInstances() map[string]*types.Instance
	GetVolumes() map[string]*types.Volume
	SetImages(map[string]*types.Image)
	SetInstances(map[string]*types.Instance)
	SetVolumes(map[string]*types.Volume)
	Save() error
	Load() error
}

type memoryState struct {
	lock      *sync.Mutex
	saveFile  string
	Images    map[string]*types.Image    `json:"Images"`
	Instances map[string]*types.Instance `json:"Instances"`
	Volumes   map[string]*types.Volume   `json:"Volumes"`
}

func NewMemoryState(saveFile string) *memoryState {
	return &memoryState{
		lock: &sync.Mutex{},
		saveFile: saveFile,
		Images: make(map[string]*types.Image),
		Instances: make(map[string]*types.Instance),
		Volumes: make(map[string]*types.Volume),
	}
}

func (s *memoryState) GetImages() map[string]*types.Image {
	imagesCopy := make(map[string]*types.Image)
	for id, image := range s.Images {
		deviceMappingsCopy := []*types.DeviceMapping{}
		for _, deviceMapping := range image.DeviceMappings {
			deviceMappingsCopy = append(deviceMappingsCopy, &types.DeviceMapping{
				MountPoint: deviceMapping.MountPoint,
				DeviceName: deviceMapping.DeviceName,
			})
		}

		imageCopy := &types.Image{
			Id: image.Id,
			Name: image.Name,
			DeviceMappings: deviceMappingsCopy,
			SizeMb: image.SizeMb,
			Infrastructure: image.Infrastructure,
		}
		imagesCopy[id] = imageCopy
	}
	return imagesCopy
}

func (s *memoryState) GetInstances() map[string]*types.Instance {
	instancesCopy := make(map[string]*types.Instance)
	for id, instance := range s.Instances {
		instanceCopy := &types.Instance{
			Id: instance.Id,
			ImageId: instance.ImageId,
			Infrastructure: instance.Infrastructure,
			Name: instance.Name,
			State: instance.State,
		}
		instancesCopy[id] = instanceCopy
	}
	return instancesCopy
}

func (s *memoryState) GetVolumes() map[string]*types.Volume {
	volumesCopy := make(map[string]*types.Volume)
	for id, volume := range s.Volumes {
		volumeCopy := &types.Volume{
			Id: volume.Id,
			Name: volume.Name,
			SizeMb: volume.SizeMb,
			Attachment: volume.Attachment,
			Infrastructure: volume.Infrastructure,
		}
		volumesCopy[id] = volumeCopy
	}
	return volumesCopy
}

func (s *memoryState) SetImages(images map[string]*types.Image) {
	s.Images = images
}

func (s *memoryState) SetInstances(instances map[string]*types.Instance) {
	s.Instances = instances
}

func (s *memoryState) SetVolumes(volumes map[string]*types.Volume) {
	s.Volumes = volumes
}

func (s *memoryState) Save() error {
	s.lock.Lock()
	defer s.lock.Unlock()
	data, err := json.Marshal(s)
	if err != nil {
		return lxerrors.New("failed to marshal memory state to json", err)
	}
	os.MkdirAll(filepath.Dir(s.saveFile), 0644)
	err = ioutil.WriteFile(s.saveFile, data, 0644)
	if err != nil {
		return lxerrors.New("writing save file "+s.saveFile, err)
	}
	return nil
}

func (s *memoryState) Load() error {
	data, err := ioutil.ReadFile(s.saveFile)
	if err != nil {
		return lxerrors.New("error reading save file "+s.saveFile, err)
	}
	var newS memoryState
	err = json.Unmarshal(data, &newS)
	if err != nil {
		return lxerrors.New("failed to unmarshal data "+string(data)+" to memory state", err)
	}
	s = newS
	return nil
}