package nvml

// #cgo CFLAGS: -I/home/trevor/git/vgpumon/src/github.com/hotpxl/nvml
// #cgo LDFLAGS: -L/home/trevor/git/vgpumon/src/github.com/hotpxl/nvml/lib -lnvidia-ml
// #include <nvml.h>
import "C"
import (
	"fmt"
	"strconv"
)

type DeviceHandle struct {
	handle C.nvmlDevice_t
}

type VgpuHandle struct {
	handle C.nvmlVgpuInstance_t
}

type VgpuTypeHandle struct {
	handle C.nvmlVgpuTypeId_t
}

func handleError(ret C.nvmlReturn_t) error {
	if ret == C.NVML_SUCCESS {
		return nil
	}
	err := C.GoString(C.nvmlErrorString(ret))
	return fmt.Errorf("NVML error: %s.", strconv.QuoteToASCII(err))
}

func nvmlInit() error {
	return handleError(C.nvmlInit())
}

func nvmlShutdown() error {
	return handleError(C.nvmlShutdown())
}

func nvmlDeviceGetCount() (int, error) {
	var n C.uint
	ret := C.nvmlDeviceGetCount(&n)
	return int(n), handleError(ret)
}

//nvmlDeviceGetActiveVgpus ( nvmlDevice_t device, unsigned int* vgpuCount, nvmlVgpuInstance_t* vgpuInstances )
func nvmlDeviceGetActiveVgpus(h DeviceHandle) ([]VgpuHandle, error) {
	count := C.uint(64)
	var vgpus [64]C.nvmlVgpuInstance_t
	ret := C.nvmlDeviceGetActiveVgpus(h.handle, &count, &vgpus[0])
	err := handleError(ret)
	if err != nil {
		return nil, err
	}
	var result []VgpuHandle
	for i := 0; i < int(count); i++ {
		result = append(result, VgpuHandle{
			handle: vgpus[i],
		})
	}
	return result, handleError(ret)
}

// Pascal
//nvmlReturn_t nvmlDeviceGetVgpuProcessUtilization ( nvmlDevice_t device, unsigned long long lastSeenTimeStamp, unsigned int* vgpuProcessSamplesCount, nvmlVgpuProcessUtilizationSample_t* utilizationSamples )
// func nvmlDeviceGetVgpuProcessUtilization(h DeviceHandle) ([]VgpuProcessInfo, error) {
//   count := C.uint(128)
//   var ts C.ulonglong
//   var processes [128]C.nvmlVgpuProcessUtilizationSample_t
//   ret := C.nvmlDeviceGetVgpuProcessUtilization(h.handle, ts, &count, &processes[0])
//   err := handleError(ret)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var result []VgpuProcessInfo
// 	for i := 0; i < int(count); i++ {
// 		result = append(result, VgpuProcessInfo{
// 	  	Name:       C.GoString(&processes[i].processName[0]),
// 			PID:        uint32(processes[i].pid),
// 			MemUtil: uint32(processes[i].memUtil),
// 			GPUUtil: uint32(processes[i].smUtil),
// 	    ENCUtil: uint32(processes[i].encUtil),
// 	    DECUtil: uint32(processes[i].decUtil),
// 	    TimeStamp: uint64(processes[i].timeStamp),
// 	    VGPU: Vgpu{
// 	      handle: VgpuHandle{handle: processes[i].vgpuInstance},
// 	    },
// 		})
// 	}
// 	return result, nil
	
// }

// nvmlReturn_t nvmlVgpuInstanceGetVmID ( nvmlVgpuInstance_t vgpuInstance, char* vmId, unsigned int  size, nvmlVgpuVmIdType_t* vmIdType )
func nvmlVgpuInstanceGetVmID(h VgpuHandle) (string, error) {
  var Id [C.NVML_DEVICE_UUID_BUFFER_SIZE]C.char
  var Type C.nvmlVgpuVmIdType_t
  ret := C.nvmlVgpuInstanceGetVmID(h.handle, &Id[0], C.NVML_DEVICE_UUID_BUFFER_SIZE, &Type)
  err := handleError(ret)
	if err != nil {
		return "", err
	}
	vmId := C.GoString(&Id[0])
	return vmId, nil
}

func nvmlDeviceGetHandleByIndex(idx int) (DeviceHandle, error) {
	var dev DeviceHandle
	ret := C.nvmlDeviceGetHandleByIndex(C.uint(idx), &dev.handle)
	return dev, handleError(ret)
}

func nvmlDeviceGetMemoryInfo(h DeviceHandle) (MemoryInfo, error) {
	var mem C.nvmlMemory_t
	ret := C.nvmlDeviceGetMemoryInfo(h.handle, &mem)
	return MemoryInfo{Free: uint64(mem.free), Used: uint64(mem.used), Total: uint64(mem.total)}, handleError(ret)
}

func nvmlDeviceGetComputeRunningProcesses(h DeviceHandle) ([]ProcessInfo, error) {
	count := C.uint(64)
	var processes [64]C.nvmlProcessInfo_t
	ret := C.nvmlDeviceGetComputeRunningProcesses(h.handle, &count, &processes[0])
	err := handleError(ret)
	if err != nil {
		return nil, err
	}
	var result []ProcessInfo
	for i := 0; i < int(count); i++ {
		result = append(result, ProcessInfo{
			PID:        int32(processes[i].pid),
			UsedMemory: uint64(processes[i].usedGpuMemory),
		})
	}
	return result, nil
}

//nvmlVgpuTypeGetName(nvmlVgpuTypeId_t vgpuTypeId, char *vgpuTypeName, unsigned int *size); string
func nvmlVgpuTypeGetName(h VgpuTypeHandle) (string, error) {
	var name [C.NVML_DEVICE_NAME_BUFFER_SIZE]C.char
	var size C.uint
	size = C.uint(C.NVML_DEVICE_NAME_BUFFER_SIZE)
	ret := C.nvmlVgpuTypeGetName(h.handle, &name[0], &size)
	err := handleError(ret)
	if err != nil {
		return "", err
	}
	out := C.GoString(&name[0])
	return out, nil
}

// func nvmlVgpuInstanceGetEncoderCapacity(h VgpuHandle) (int, error) {
// 	var n C.uint
// 	ret := C.nvmlVgpuInstanceGetEncoderCapacity(h.handle, &n)
// 	return int(n), handleError(ret)
// }

func nvmlVgpuInstanceGetVmDriverVersion(h VgpuHandle) (string, error) {
  var version [C.NVML_SYSTEM_DRIVER_VERSION_BUFFER_SIZE]C.char
  ret := C.nvmlVgpuInstanceGetVmDriverVersion(h.handle, &version[0], C.NVML_SYSTEM_DRIVER_VERSION_BUFFER_SIZE)
  err := handleError(ret)
  if err != nil {
    return "", err
  }
  vmId := C.GoString(&version[0])
  return vmId, nil
}

// Pascal
// func nvmlDeviceGetEncoderCapacity(h DeviceHandle) (uint, error) {
// 	var n C.uint
// 	ret := C.nvmlDeviceGetEncoderCapacity(h.handle, C.NVML_ENCODER_QUERY_H264, &n)
// 	return uint(n), handleError(ret)
// }

// func nvmlVgpuInstanceSetEncoderCapacity(h VgpuHandle, capacity uint) (error) {
// 	ret := C.nvmlVgpuInstanceSetEncoderCapacity(h.handle, C.uint(capacity))
// 	err := handleError(ret)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func nvmlVgpuInstanceGetType(h VgpuHandle) (VgpuTypeHandle, error) {
  var vgpuTypeId C.nvmlVgpuTypeId_t
  ret := C.nvmlVgpuInstanceGetType(h.handle, &vgpuTypeId)
	err := handleError(ret)
	if err != nil {
		return VgpuTypeHandle{}, err
	}
	return VgpuTypeHandle{handle: vgpuTypeId}, nil
}

func nvmlVgpuTypeGetMaxInstances(h DeviceHandle, vgpuTypeId VgpuTypeHandle) (uint, error) {
	var n C.uint
	ret := C.nvmlVgpuTypeGetMaxInstances(h.handle, vgpuTypeId.handle, &n)
	return uint(n), handleError(ret)
}

//nvmlDeviceGetUtilizationRates(nvmlDevice_t device, nvmlUtilization_t *utilization); struct
func nvmlDeviceGetUtilizationRates(h DeviceHandle) (UtilizationInfo, error) {
	var util C.nvmlUtilization_t
	ret := C.nvmlDeviceGetUtilizationRates(h.handle, &util)
	return UtilizationInfo{GPUUtil: uint(util.gpu), MemUtil: uint(util.memory)}, handleError(ret)
}

//nvmlDeviceGetEncoderUtilization(nvmlDevice_t device, unsigned int *utilization, unsigned int *samplingPeriodUs); uint
func nvmlDeviceGetEncoderUtilization(h DeviceHandle) (uint, error) {
	var n C.uint
	var p C.uint
	ret := C.nvmlDeviceGetEncoderUtilization(h.handle, &n, &p)
	return uint(n), handleError(ret)
}

//nvmlDeviceGetDecoderUtilization(nvmlDevice_t device, unsigned int *utilization, unsigned int *samplingPeriodUs); uint
func nvmlDeviceGetDecoderUtilization(h DeviceHandle) (uint, error) {
	var n C.uint
	var p C.uint
	ret := C.nvmlDeviceGetDecoderUtilization(h.handle, &n, &p)
	return uint(n), handleError(ret)
}

//nvmlVgpuInstanceGetFrameRateLimit ( nvmlVgpuInstance_t vgpuInstance, unsigned int* frameRateLimit ) uint
func nvmlVgpuInstanceGetFrameRateLimit(h VgpuHandle) (uint, error) {
	var n C.uint
	ret := C.nvmlVgpuInstanceGetFrameRateLimit(h.handle, &n)
	return uint(n), handleError(ret)
}

//nvmlDeviceGetVbiosVersion(nvmlDevice_t device, char *version, unsigned int length); string
func nvmlDeviceGetVbiosVersion(h DeviceHandle) (string, error) {
	var version [C.NVML_DEVICE_VBIOS_VERSION_BUFFER_SIZE]C.char
	ret := C.nvmlDeviceGetVbiosVersion(h.handle, &version[0], C.NVML_DEVICE_VBIOS_VERSION_BUFFER_SIZE)
	err := handleError(ret)
	  if err != nil {
		  return "", err
	  }
	  vbios := C.GoString(&version[0])
	  return vbios, nil
  }


