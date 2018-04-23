package nvml

import (
	"fmt"

	"github.com/shirou/gopsutil/process"
	log "github.com/sirupsen/logrus"
)

// Session manages the initialization and shutdown of NVML context.
type Session struct {
	active bool
}

// NewSession creates a new session.
func NewSession() (*Session, error) {
	return &Session{active: true}, nvmlInit()
}

// Close frees existing session and underlying NVML context.
func (s *Session) Close() {
	if !s.active {
		log.Fatal("Already closed.")
	}
	s.active = false
	err := nvmlShutdown()
	if err != nil {
		log.Fatal(err)
	}
}

// DeviceCount returns the number of devices.
func (s *Session) DeviceCount() (int, error) {
	if !s.active {
		return 0, fmt.Errorf("Already closed.")
	}
	return nvmlDeviceGetCount()
}

// GetDevice returns a specific device given its index.
func (s *Session) GetDevice(idx int) (*Device, error) {
	if !s.active {
		return nil, fmt.Errorf("Already closed.")
	}
	dev, err := nvmlDeviceGetHandleByIndex(idx)
	if err != nil {
		return nil, err
	}
	return &Device{handle: dev}, nil
}

// GetAllDevices returns all devices accessible.
func (s *Session) GetAllDevices() ([]Device, error) {
	if !s.active {
		return nil, fmt.Errorf("Already closed.")
	}
	count, err := s.DeviceCount()
	if err != nil {
		return nil, err
	}
	var ret []Device
	for i := 0; i < count; i++ {
		dev, err := s.GetDevice(i)
		if err != nil {
			return nil, err
		}
		ret = append(ret, *dev)
	}
	return ret, nil
}

// Device represents a single device.
type Device struct {
	handle DeviceHandle
}

// Vgpu represents a single virtual instance.
type Vgpu struct {
	handle VgpuHandle
}

// MemoryInfo holds memory consumption information for a device.
type MemoryInfo struct {
	Free  uint64
	Used  uint64
	Total uint64
}

// ProcessInfo holds process information on a device.
type ProcessInfo struct {
	PID        int32  `json:"pid"`
	UsedMemory uint64 `json:"usedMemory"`
	Username   string `json:"username"`
}

type VgpuProcessInfo struct {
	Name string `json:"name"`
	PID  uint32 `json:"pid"`
	GPUUtil uint32 `json:"gpuUtil"` //!< SM (3D/Compute) Util Value
	MemUtil uint32 `json:"memUtil"` //!< Frame Buffer Memory Util Value
	ENCUtil uint32 `json:"encUtil"` //!< Encoder Util Value
	DECUtil uint32 `json:"decUtil"` //!< Decoder Util Value
	TimeStamp uint64 `json:"timeStamp"`
	VGPU Vgpu
}

type UtilizationInfo struct {
    GPUUtil uint `json:"gpuUtil"`	//!< Percent of time over the past sample period during which one or more kernels was executing on the GPU
    MemUtil uint `json:"memUtil"`	//!< Percent of time over the past sample period during which global (device) memory was being read or written
}

// GetAllVgpus returns each Vgpu from a device.
func (d *Device) GetAllVgpus() ([]Vgpu, error) {
  vgpuHandles, err := nvmlDeviceGetActiveVgpus(d.handle)
  if err != nil {
		return nil, err
	}
	var result []Vgpu
	for i := 0; i < len(vgpuHandles); i++ {
		result = append(result, Vgpu{
			handle: vgpuHandles[i],
		})
	}
	return result, nil
}

// MemoryInfo returns memory consumption information from a device.
func (d *Device) MemoryInfo() (MemoryInfo, error) {
	return nvmlDeviceGetMemoryInfo(d.handle)
}

// Processes returns processes running on a device.
func (d *Device) Processes() ([]ProcessInfo, error) {
	processes, err := nvmlDeviceGetComputeRunningProcesses(d.handle)
	if err != nil {
		return nil, err
	}
	for idx, p := range processes {
		pp, err := process.NewProcess(p.PID)
		if err != nil {
			return nil, err
		}
		username, err := pp.Username()
		if err != nil {
			return nil, err
		}
		processes[idx].Username = username
	}
	return processes, nil
}

func (d *Device) GetUtilization() (UtilizationInfo, error) {
	return nvmlDeviceGetUtilizationRates(d.handle)
}

func (d *Device) GetEncoderUtilization() (uint, error) {
	return nvmlDeviceGetEncoderUtilization(d.handle)
}

func (d *Device) GetDecoderUtilization() (uint, error) {
	return nvmlDeviceGetDecoderUtilization(d.handle)
}

func (d *Device) GetVbiosVersion() (string, error) {
	vbios, err := nvmlDeviceGetVbiosVersion(d.handle)
	if err != nil {
	  return "", err
	}
	return vbios, nil
}

// Pascal
// func (d *Device) VgpuProcesses() ([]VgpuProcessInfo, error) {
// 	processes, err := nvmlDeviceGetVgpuProcessUtilization(d.handle)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return processes, nil
// }

// Pascal
// func (d *Device) GetEncodeCapacity() (uint, error) {
//   return nvmlDeviceGetEncoderCapacity(d.handle)
// }

func (d *Device) GetMaxInstances(v Vgpu) (uint, error) {
  vt, err := nvmlVgpuInstanceGetType(v.handle)
  if err != nil {
		return 0, err
	}
	max, err := nvmlVgpuTypeGetMaxInstances(d.handle, vt)
	if err != nil {
		return 0, err
	}
	return max, nil
}

func (d *Device) GetCurInstances() (uint, error) {
	vgpuHandles, err := nvmlDeviceGetActiveVgpus(d.handle)
	if err != nil {
		return 0, err
	}
	return uint(len(vgpuHandles)), nil
}


func (v *Vgpu) GetFrameRateLimit() (uint, error) {
	limit, err := nvmlVgpuInstanceGetFrameRateLimit(v.handle)
	if err != nil {
		return 0, err
	}
	return limit, nil
}

func (v *Vgpu) GetTypeName() (string, error) {
	vt, err := nvmlVgpuInstanceGetType(v.handle)
	if err != nil {
		  return "", err
	  }
	name, err := nvmlVgpuTypeGetName(vt)
	if err != nil {
	  return "", err
	}
	return name, nil
}

func (v *Vgpu) GetDriverVersion() (string, error) {
	vd, err := nvmlVgpuInstanceGetVmDriverVersion(v.handle)
	if err != nil {
	  return "", err
	}
	return vd, nil
}

func (v *Vgpu) GetVmId() (string, error) {
	return nvmlVgpuInstanceGetVmID(v.handle)
}

// Pascal
// func (v *Vgpu) GetEncodeCapacity() (int, error) {
//   return nvmlVgpuInstanceGetEncoderCapacity(v.handle)
// }

// Pascal
// func (v *Vgpu) SetEncodeCapacity(capacity uint) (error) {
//   return nvmlVgpuInstanceSetEncoderCapacity(v.handle, capacity)
// }

// This is a really stupid helper, but i guess it's the only..
// ...way to do a reverse lookup?
func (v *Vgpu) GetDevice() (Device, error) {
  
  s, err := NewSession()
  if err != nil {
    panic(err)
  }
  defer s.Close()
  
  devices, err := s.GetAllDevices()
  if err != nil {
    panic(err)
  }
  for _, d := range devices {
    vs, err := d.GetAllVgpus()
    if err != nil {
      panic(err)
    }
    for i := 0; i < len(vs); i++ {
      if vs[i] == *v { return d, nil }
    }
  }
  return Device{}, nil
}

// func VPFilter(vs []VgpuProcessInfo, f func(VgpuProcessInfo) bool) []VgpuProcessInfo {
//     vsf := make([]VgpuProcessInfo, 0)
//     for _, v := range vs {
//         if f(v) {
//             vsf = append(vsf, v)
//         }
//     }
//     return vsf
// }

